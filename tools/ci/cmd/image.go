package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/ghaction"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/image"
)

func newImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Commands for resolving and managing container images",
	}

	cmd.AddCommand(newImageResolveCmd())

	return cmd
}

func newImageResolveCmd() *cobra.Command {
	var (
		ecrType        string
		repositoryPath string
		imageTag       string
		awsAccount     string
		awsRegion      string
		jsonOutput     bool
	)

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a Chainlink Docker image URI based on ECR type and environment variables",
		RunE: func(cmd *cobra.Command, args []string) error {
			act := ghaction.NewAction(cmd.OutOrStdout())
			if ecrType == "" {
				ecrType = act.GetInput("ecr_type")
			}
			if ecrType == "" {
				ecrType = act.Getenv("ECR_TYPE")
			}
			if repositoryPath == "" {
				repositoryPath = act.GetInput("repo_path")
			}
			if repositoryPath == "" {
				repositoryPath = act.Getenv("CHAINLINK_IMAGE_REPO_PATH")
			}
			if imageTag == "" {
				imageTag = act.GetInput("tag")
			}
			if imageTag == "" {
				imageTag = act.Getenv("CHAINLINK_IMAGE_TAG")
			}
			if awsAccount == "" {
				awsAccount = act.GetInput("aws_account")
			}
			if awsAccount == "" {
				awsAccount = act.Getenv("AWS_ACCOUNT_NUMBER")
			}
			if awsRegion == "" {
				awsRegion = act.GetInput("aws_region")
			}
			if awsRegion == "" {
				awsRegion = act.Getenv("AWS_REGION")
			}

			resolved, err := image.Resolve(image.ResolveOptions{
				ECRType:          ecrType,
				RepositoryPath:   repositoryPath,
				ImageTag:         imageTag,
				AWSAccountNumber: awsAccount,
				AWSRegion:        awsRegion,
			})
			if err != nil {
				return err
			}

			if act.Getenv("GITHUB_OUTPUT") != "" {
				if err := act.SetOutput("resolved_image", resolved); err != nil {
					return err
				}
			}

			if jsonOutput {
				payload := map[string]string{
					"resolved_image": resolved,
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}

			fmt.Fprintln(cmd.OutOrStdout(), resolved)
			return nil
		},
	}

	cmd.Flags().StringVar(&ecrType, "ecr-type", "", "ECR type ('sdlc' or 'public') (env: ECR_TYPE)")
	cmd.Flags().StringVar(&repositoryPath, "repo-path", "", "Repository path (env: CHAINLINK_IMAGE_REPO_PATH)")
	cmd.Flags().StringVar(&imageTag, "tag", "", "Image tag (env: CHAINLINK_IMAGE_TAG)")
	cmd.Flags().StringVar(&awsAccount, "aws-account", "", "AWS Account Number (env: AWS_ACCOUNT_NUMBER)")
	cmd.Flags().StringVar(&awsRegion, "aws-region", "", "AWS Region (env: AWS_REGION)")
	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output resolved image in JSON format")

	return cmd
}
