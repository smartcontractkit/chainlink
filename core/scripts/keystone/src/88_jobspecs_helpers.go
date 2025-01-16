package src

import (
	"fmt"
)

type OCRSpec struct {
	ContractID string
}

type BootSpec struct {
	ContractID string
}

type WorkflowSpec struct {
	WorkflowID string
}

type JobSpec struct {
	ID                           string
	Name                         string
	BootstrapSpec                BootSpec
	OffChainReporting2OracleSpec OCRSpec
	WorkflowSpec                 WorkflowSpec
}

func upsertJob(api *nodeAPI, jobSpecName string, jobSpecStr string) {
	jobsResp := api.mustExec(api.methods.ListJobs)
	jobs := mustJSON[[]JobSpec](jobsResp)
	for _, job := range *jobs {
		if job.Name == jobSpecName {
			fmt.Printf("Job already exists: %s, replacing..\n", jobSpecName)
			api.withArg(job.ID).mustExec(api.methods.DeleteJob)
			break
		}
	}

	fmt.Printf("Deploying jobspec: %s\n", jobSpecName)
	_, err := api.withArg(jobSpecStr).exec(api.methods.CreateJob)
	if err != nil {
		panic(fmt.Sprintf("Failed to deploy job spec: %s Error: %s", jobSpecStr, err))
	}
}

func clearJobs(api *nodeAPI) {
	jobsResp := api.mustExec(api.methods.ListJobs)
	jobs := mustJSON[[]JobSpec](jobsResp)
	for _, job := range *jobs {
		fmt.Printf("Deleting job: %s\n", job.Name)
		api.withArg(job.ID).mustExec(api.methods.DeleteJob)
	}
	fmt.Println("All jobs have been deleted.")
}
