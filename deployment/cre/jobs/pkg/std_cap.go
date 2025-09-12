package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink/deployment/cre/jobs/pkg/templates"
)

const (
	ErrorEmptyJobName = "job name cannot be empty"
)

type StandardCapabilityJob struct {
	JobName string // Must be alphanumeric, with _, -, ., no spaces.
	Command string
	Config  string

	// If not provided, ExternalJobID is automatically filled in by calling `externalJobIDHashFunc`
	ExternalJobID string
	// OracleFactory is the configuration for the Oracle Factory job.
	OracleFactory OracleFactory
}

func (s *StandardCapabilityJob) Validate() error {
	if s.JobName == "" {
		return errors.New(ErrorEmptyJobName)
	}

	return nil
}

func (s *StandardCapabilityJob) Resolve() (string, error) {
	if s.ExternalJobID == "" {
		// We expect there to only be 1 instance of a standard capability per node
		// This is because adding duplicate capabilities to the registry will typically fail due to an ID clash.
		// Some capabilities, such as contract read and write, are unique per their config
		externalJobID, err := externalJobIDHashFunc([]byte(s.Command), []byte(s.Config))
		if err != nil {
			return "", fmt.Errorf("failed to create external job id: %w", err)
		}
		s.ExternalJobID = externalJobID.String()
	}

	t, err := template.New("s").ParseFS(templates.FS, "stdcap.tmpl")
	if err != nil {
		return "", fmt.Errorf("failed to parse stdcap.tmpl: %w", err)
	}

	b := &bytes.Buffer{}
	err = t.ExecuteTemplate(b, "stdcap.tmpl", s)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return b.String(), nil
}

func externalJobIDHashFunc(command, config []byte) (uuid.UUID, error) {
	var externalJobID uuid.UUID
	if len(config) > 0 {
		externalJobID = uuid.NewSHA1(uuid.Nil, append(command, config...))
		return externalJobID, nil
	}
	externalJobID = uuid.NewSHA1(uuid.Nil, command)
	return externalJobID, nil
}
