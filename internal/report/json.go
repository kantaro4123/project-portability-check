package report

import (
	"encoding/json"
	"io"

	"github.com/kantaro4123/project-portability-check/internal/model"
)

func WriteJSON(w io.Writer, report model.Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
