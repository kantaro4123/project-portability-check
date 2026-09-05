package analyzer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

type barrierDetector struct {
	id      string
	ready   chan<- string
	release <-chan struct{}
	finding model.Finding
	err     error
}

func (d barrierDetector) ID() string { return d.id }

func (d barrierDetector) Detect(ctx context.Context, _ Project) ([]model.Finding, error) {
	d.ready <- d.id
	select {
	case <-d.release:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	if d.err != nil {
		return nil, d.err
	}
	return []model.Finding{d.finding}, nil
}

func TestAnalyzeRunsDetectorsConcurrentlyAndSortsResults(t *testing.T) {
	ready := make(chan string, 2)
	release := make(chan struct{})
	engine := New(
		barrierDetector{id: "one", ready: ready, release: release, finding: model.Finding{RuleID: "z.rule", Severity: model.SeverityWarning, Path: "b.txt", Line: 2}},
		barrierDetector{id: "two", ready: ready, release: release, finding: model.Finding{RuleID: "a.rule", Severity: model.SeverityError, Path: "a.txt", Line: 1}},
	)

	done := make(chan struct{})
	var findings []model.Finding
	var err error
	go func() {
		findings, err = engine.Analyze(context.Background(), Project{})
		close(done)
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-ready:
		case <-time.After(time.Second):
			t.Fatal("detectors did not reach the barrier concurrently")
		}
	}
	close(release)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("analysis did not finish")
	}
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 2 || findings[0].RuleID != "a.rule" || findings[1].RuleID != "z.rule" {
		t.Fatalf("unexpected deterministic ordering: %+v", findings)
	}
}

func TestAnalyzeReturnsDetectorError(t *testing.T) {
	ready := make(chan string, 1)
	release := make(chan struct{})
	close(release)
	want := errors.New("boom")
	engine := New(barrierDetector{id: "bad", ready: ready, release: release, err: want})
	_, err := engine.Analyze(context.Background(), Project{})
	if !errors.Is(err, want) {
		t.Fatalf("err=%v, want %v", err, want)
	}
}
