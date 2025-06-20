package mockcapability

import (
	"bytes"
	"text/template"

	"github.com/google/uuid"

	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
)

func MockCapabilitiesJob(nodeID, binaryPath string, mocks []*MockCapabilities) *jobv1.ProposeJobRequest {
	jobTemplate := `type = "standardcapabilities"
			schemaVersion = 1
			externalJobID = "{{ .JobID }}"
			name = "mock-capability"
			forwardingAllowed = false
			command = "{{ .BinaryPath }}"
			config = """
				port=7777
		{{ range $index, $m := .Mocks }}
 		  [[DefaultMocks]]
				id="{{ $m.ID }}"
				description="{{ $m.Description }}"
				type="{{ $m.Type }}"
 		{{- end }}
			"""`
	tmpl, err := template.New("mock-job").Parse(jobTemplate)

	if err != nil {
		panic(err)
	}
	mockJobsData := make([]map[string]string, 0)
	for _, m := range mocks {
		mockJobsData = append(mockJobsData, map[string]string{
			"ID":          m.Name + "@" + m.Version,
			"Description": m.Description,
			"Type":        m.Type,
		})
	}

	jobUUID := uuid.NewString()
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, map[string]interface{}{
		"JobID":      jobUUID,
		"ShortID":    jobUUID[0:8],
		"BinaryPath": binaryPath,
		"Mocks":      mockJobsData,
	})
	if err != nil {
		panic(err)
	}

	return &jobv1.ProposeJobRequest{
		NodeId: nodeID,
		Spec:   renderedTemplate.String(),
	}
}
