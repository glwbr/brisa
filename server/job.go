package server

import (
	"sync"
	"time"

	"github.com/glwbr/brisa/invoice"
	"github.com/glwbr/brisa/scraper"
)

type JobStatus string

const (
	StatusCreated        JobStatus = "created"
	StatusRunning        JobStatus = "running"
	StatusWaitingCaptcha JobStatus = "waiting_captcha"
	StatusCompleted      JobStatus = "completed"
	StatusFailed         JobStatus = "failed"
)

type Job struct {
	ID        string           `json:"id"`
	Status    JobStatus        `json:"status"`
	AccessKey string           `json:"accessKey"`
	Result    *invoice.Receipt `json:"result,omitempty"`
	Error     string           `json:"error,omitempty"`
	CreatedAt time.Time        `json:"createdAt"`

	Captcha *scraper.CaptchaChallenge `json:"captcha,omitempty"`

	solutionCh chan string
	mu         sync.Mutex
}

func NewJob(id, accessKey string) *Job {
	return &Job{
		ID:         id,
		Status:     StatusCreated,
		AccessKey:  accessKey,
		CreatedAt:  time.Now(),
		solutionCh: make(chan string),
	}
}

type JobManager struct {
	jobs map[string]*Job
	mu   sync.RWMutex
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]*Job),
	}
}

func (m *JobManager) CreateJob(accessKey string) *Job {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := generateID()
	job := NewJob(id, accessKey)
	m.jobs[id] = job
	return job
}

func (m *JobManager) GetJob(id string) (*Job, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[id]
	return job, ok
}

func (m *JobManager) CleanupLoop(interval time.Duration, maxAge time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		m.mu.Lock()
		for id, job := range m.jobs {
			if time.Since(job.CreatedAt) > maxAge {
				delete(m.jobs, id)
			}
		}
		m.mu.Unlock()
	}
}

func generateID() string {
	return time.Now().Format("20060102150405")
}

func (j *Job) SetWaitingCaptcha(challenge *scraper.CaptchaChallenge) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Status = StatusWaitingCaptcha
	j.Captcha = challenge
}

func (j *Job) SetRunning() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Status = StatusRunning
	j.Captcha = nil
}

func (j *Job) SetCompleted(result *invoice.Receipt) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Status = StatusCompleted
	j.Result = result
	j.Captcha = nil
}

func (j *Job) SetFailed(err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Status = StatusFailed
	j.Error = err.Error()
	j.Captcha = nil
}

func (j *Job) SubmitCaptcha(solution string) {
	// Non-blocking send or blocking send? Its th question
	// The solver is waiting on this channel.
	j.solutionCh <- solution
}
