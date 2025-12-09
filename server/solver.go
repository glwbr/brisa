package server

import (
	"context"

	"github.com/glwbr/brisa/scraper"
)

type AsyncSolver struct {
	Job *Job
}

func NewAsyncSolver(job *Job) *AsyncSolver {
	return &AsyncSolver{Job: job}
}

func (s *AsyncSolver) Solve(ctx context.Context, challenge *scraper.CaptchaChallenge) (*scraper.CaptchaSolution, error) {
	s.Job.SetWaitingCaptcha(challenge)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case solution := <-s.Job.solutionCh:
		s.Job.SetRunning()
		return &scraper.CaptchaSolution{
			Text:        solution,
			ChallengeID: challenge.ID,
		}, nil
	}
}
