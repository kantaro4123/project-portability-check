package report

import "github.com/kantaro4123/project-portability-check/internal/model"

// Summarize converts findings into a stable portability score.
func Summarize(files int, findings []model.Finding) model.Summary {
	s := model.Summary{FilesScanned: files, Score: 100}
	for _, f := range findings {
		switch f.Severity {
		case model.SeverityError:
			s.Errors++
			s.Score -= 15
		case model.SeverityWarning:
			s.Warnings++
			s.Score -= 5
		default:
			s.Info++
		}
	}
	if s.Score < 0 {
		s.Score = 0
	}
	return s
}
