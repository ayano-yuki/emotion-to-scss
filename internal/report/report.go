package report

import (
	"encoding/json"
	"os"
)

type Summary struct {
	Files     int `json:"files"`
	Styles    int `json:"styles"`
	Converted int `json:"converted"`
	Verified  int `json:"verified"`
	Failed    int `json:"failed"`
}

type Report struct {
	Summary Summary      `json:"summary"`
	Files   []FileReport `json:"files"`
}

type FileReport struct {
	Path   string        `json:"path"`
	Styles []StyleReport `json:"styles"`
}

type StyleReport struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	Line         int    `json:"line"`
	ClassName    string `json:"className"`
	Verification string `json:"verification"`
	Reason       string `json:"reason,omitempty"`
}

func Write(path string, r Report) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
