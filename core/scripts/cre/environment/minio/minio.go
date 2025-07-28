package minio

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

var MinioCommand = &cobra.Command{
	Use:   "minio",
	Short: "interact with MinIO storage",
	Long:  `Commands to upload and manage files in MinIO object storage`,
}

var (
	configPath string
	bucketName string
	objectName string
)

var UploadCmd = &cobra.Command{
	Use:   "upload [files...]",
	Short: "Upload files to MinIO storage",
	Long:  `Upload specified files to MinIO object storage using configuration from cre.yaml`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse the config file
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			return
		}

		var config crecli.Profiles
		if err := yaml.Unmarshal(configData, &config); err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			return
		}

		// Get MinIO config
		minioConfig := config.Test.WorkflowStorage.Minio // TODO @george-dorin: make it dynamic don't just use test

		// Initialize MinIO client
		minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
			Secure: minioConfig.UseSSL,
			Region: minioConfig.Region,
		})
		if err != nil {
			fmt.Printf("Error initializing MinIO client: %v\n", err)
			return
		}

		// Create bucket if it doesn't exist
		exists, err := minioClient.BucketExists(context.Background(), bucketName)
		if err != nil {
			fmt.Printf("Error checking bucket: %v\n", err)
			return
		}
		if !exists {
			err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{
				Region: minioConfig.Region,
			})
			if err != nil {
				fmt.Printf("Error creating bucket: %v\n", err)
				return
			}
			fmt.Printf("Created bucket %s\n", bucketName)
		}

		// Process each file
		for _, filePath := range args {
			// Determine object name for this file
			currentObjectName := objectName
			if currentObjectName == "" {
				currentObjectName = filepath.Base(filePath)
			}

			// Get file info and open file
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Printf("Error getting file info for %s: %v\n", filePath, err)
				continue
			}
			file, err := os.Open(filePath)
			if err != nil {
				fmt.Printf("Error opening file %s: %v\n", filePath, err)
				continue
			}

			// Upload the file
			info, err := minioClient.PutObject(context.Background(), bucketName, currentObjectName, file,
				fileInfo.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})

			// Close the file after upload attempt
			file.Close()

			if err != nil {
				fmt.Printf("Error uploading file %s: %v\n", filePath, err)
				continue
			}

			fmt.Printf("Successfully uploaded %s to %s/%s\n", filePath, bucketName, currentObjectName)
			fmt.Printf("ETag: %s, Size: %d bytes\n", info.ETag, info.Size)
		}
	},
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List objects in a MinIO bucket",
	Long:  `List all objects stored in a specified MinIO bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse the config file
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			return
		}

		var config crecli.Profiles
		if err := yaml.Unmarshal(configData, &config); err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			return
		}

		// Get MinIO config
		minioConfig := config.Test.WorkflowStorage.Minio

		// Initialize MinIO client
		minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
			Secure: minioConfig.UseSSL,
			Region: minioConfig.Region,
		})
		if err != nil {
			fmt.Printf("Error initializing MinIO client: %v\n", err)
			return
		}

		// List all objects in bucket
		fmt.Printf("Objects in bucket '%s':\n", bucketName)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		objectCh := minioClient.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
			Recursive: true,
		})

		for object := range objectCh {
			if object.Err != nil {
				fmt.Printf("Error: %v\n", object.Err)
				continue
			}
			fmt.Printf("- %s (size: %d bytes, last modified: %s)\n",
				object.Key, object.Size, object.LastModified)
		}
	},
}

var DownloadCmd = &cobra.Command{
	Use:   "download [object-name]",
	Short: "Download an object from MinIO",
	Long:  `Download an object from MinIO bucket to local filesystem`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objectToDownload := args[0]
		outputPath := filepath.Base(objectToDownload)

		// Read and parse the config file
		configData, err := ioutil.ReadFile(configPath)
		if err != nil {
			fmt.Printf("Error reading config file: %v\n", err)
			return
		}

		var config crecli.Profiles
		if err := yaml.Unmarshal(configData, &config); err != nil {
			fmt.Printf("Error parsing config file: %v\n", err)
			return
		}

		// Get MinIO config
		minioConfig := config.Test.WorkflowStorage.Minio

		// Initialize MinIO client
		minioClient, err := minio.New(minioConfig.Endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
			Secure: minioConfig.UseSSL,
			Region: minioConfig.Region,
		})
		if err != nil {
			fmt.Printf("Error initializing MinIO client: %v\n", err)
			return
		}

		// Download the object
		err = minioClient.FGetObject(context.Background(), bucketName, objectToDownload,
			outputPath, minio.GetObjectOptions{})
		if err != nil {
			fmt.Printf("Error downloading object: %v\n", err)
			return
		}

		fmt.Printf("Successfully downloaded %s/%s to %s\n", bucketName, objectToDownload, outputPath)
	},
}

func init() {
	UploadCmd.Flags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")
	UploadCmd.Flags().StringVar(&bucketName, "bucket", "test-bucket", "Bucket name")
	UploadCmd.Flags().StringVar(&objectName, "name", "", "Object name (defaults to filename)")

	ListCmd.Flags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")
	ListCmd.Flags().StringVar(&bucketName, "bucket", "test-bucket", "Bucket name")

	DownloadCmd.Flags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")
	DownloadCmd.Flags().StringVar(&bucketName, "bucket", "test-bucket", "Bucket name")

	MinioCommand.AddCommand(UploadCmd, ListCmd, DownloadCmd)
}
