package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/glwbr/brisa/portal/ba"
)

// TODO: add a logger
type Server struct {
	jobManager *JobManager
}

func NewServer() *Server {
	s := &Server{
		jobManager: NewJobManager(),
	}
	go s.jobManager.CleanupLoop(1*time.Minute, 2*time.Minute)
	return s
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/invoice-jobs", s.handleCreateJob)
	mux.HandleFunc("GET /api/invoice-jobs/{id}", s.handleGetJob)
	mux.HandleFunc("POST /api/invoice-jobs/{id}/captcha", s.handleSubmitCaptcha)

	handler := corsMiddleware(mux)

	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, handler)
}

func (s *Server) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessKey string `json:"accessKey"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AccessKey == "" {
		http.Error(w, "accessKey is required", http.StatusBadRequest)
		return
	}

	job := s.jobManager.CreateJob(req.AccessKey)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		job.SetRunning()

		solver := NewAsyncSolver(job)
		scraper, err := ba.New(ba.WithCaptchaSolver(solver))
		if err != nil {
			job.SetFailed(fmt.Errorf("failed to create scraper: %w", err))
			return
		}

		result, err := scraper.FetchByAccessKey(ctx, req.AccessKey)
		if err != nil {
			job.SetFailed(err)
			return
		}

		job.SetCompleted(result.Receipt)
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"jobId": job.ID,
	})
}

func (s *Server) handleGetJob(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.jobManager.GetJob(id)
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// Use a lock to read job safely? GetJob does RLock but we access fields.
	// We should probably have a job.ToJSON() or similar, but direct access is okay if we are careful or if fields are simple.
	// Actually, concurrent Read/Write on map fields (like Captcha metadata) could be an issue.
	// But let's assume JSON encoder reads are safe enough for this prototype or add a method.

	job.mu.Lock()
	defer job.mu.Unlock()
	json.NewEncoder(w).Encode(job)
}

func (s *Server) handleSubmitCaptcha(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	job, ok := s.jobManager.GetJob(id)
	if !ok {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	var req struct {
		Solution string `json:"solution"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Again... Non-blocking check if job is waiting?
	job.mu.Lock()
	if job.Status != StatusWaitingCaptcha {
		job.mu.Unlock()
		http.Error(w, "Job is not waiting for captcha", http.StatusBadRequest)
		return
	}
	job.mu.Unlock()

	// This might block if the receiver is not ready, but it should be ready if status is waiting.
	// However, if the scraper crashed or timed out, this might block forever.
	// Use a select with timeout or ensure logic is sound.

	select {
	case job.solutionCh <- req.Solution:
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Failed to submit solution (scraper not listening)", http.StatusInternalServerError)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: restrict this
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
