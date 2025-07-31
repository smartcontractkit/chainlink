// Package minio provides commands to interact with MinIO object storage
// through a CLI interface. It supports uploading, downloading, and listing
// objects in MinIO buckets.
//
// The package relies on configuration provided in a cre.yaml file which
// contains credentials and connection details for the MinIO server.
package minio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/system-tests/lib/crecli"
)

// MinioCommand is the root command for all MinIO-related operations.
// It serves as a container for subcommands that perform specific actions
// with MinIO storage.
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

// UploadCmd handles uploading files to a MinIO bucket.
// It accepts one or more file paths as arguments and uploads each file
// to the specified bucket. If no object name is provided, the file's
// basename is used as the object name.
//
// Flags:
//
//	--config string   Path to cre.yaml config file (default "cre.yaml")
//	--bucket string   Bucket name (default "default")
//	--name string     Object name (defaults to filename)
var UploadCmd = &cobra.Command{
	Use:   "upload [files...]",
	Short: "Upload files to MinIO storage",
	Long:  `Upload specified files to MinIO object storage using configuration from cre.yaml`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse the config file
		minioClient, minioConfig, err := setupMinioClient(configPath)
		if err != nil {
			fmt.Printf("Error initializing MinIO client: %v\n", err)
			return
		}

		// Create bucket if it doesn't exist
		err = ensureBucketExists(context.Background(), minioClient, bucketName, minioConfig.Region)
		if err != nil {
			fmt.Printf("Error checking bucket: %v\n", err)
			return
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

// ListCmd displays all objects in a MinIO bucket.
// It lists each object along with its size and last modified date.
//
// Flags:
//
//	--config string   Path to cre.yaml config file (default "cre.yaml")
//	--bucket string   Bucket name (default "default")
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List objects in a MinIO bucket",
	Long:  `List all objects stored in a specified MinIO bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read and parse the config file
		minioClient, _, err := setupMinioClient(configPath)
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

// DownloadCmd retrieves an object from a MinIO bucket and saves it locally.
// It requires the object name as an argument and saves the downloaded
// object using its basename as the filename.
//
// Flags:
//
//	--config string   Path to cre.yaml config file (default "cre.yaml")
//	--bucket string   Bucket name (default "default")
var DownloadCmd = &cobra.Command{
	Use:   "download [object-name]",
	Short: "Download an object from MinIO",
	Long:  `Download an object from MinIO bucket to local filesystem`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objectToDownload := args[0]
		outputPath := filepath.Base(objectToDownload)

		minioClient, _, err := setupMinioClient(configPath)
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

// init registers all subcommands and their flags with the root MinIO command.
// It sets up default values for the config path, bucket name, and other options
// for each subcommand.
func init() {
	UploadCmd.Flags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")
	UploadCmd.Flags().StringVar(&bucketName, "bucket", "default", "Bucket name")
	UploadCmd.Flags().StringVar(&objectName, "name", "", "Object name (defaults to filename)")

	ListCmd.Flags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")
	ListCmd.Flags().StringVar(&bucketName, "bucket", "default", "Bucket name")

	DownloadCmd.Flags().StringVar(&configPath, "config", "cre.yaml", "Path to cre.yaml config file")
	DownloadCmd.Flags().StringVar(&bucketName, "bucket", "default", "Bucket name")

	MinioCommand.AddCommand(UploadCmd, ListCmd, DownloadCmd)
}

// setupMinioClient initializes and returns a MinIO client from config
func setupMinioClient(configPath string) (*minio.Client, crecli.MinioStorageSettings, error) {
	// Read and parse the config file
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, crecli.MinioStorageSettings{}, fmt.Errorf("error reading config file: %w", err)
	}

	var config crecli.Profiles
	if err = yaml.Unmarshal(configData, &config); err != nil {
		return nil, crecli.MinioStorageSettings{}, fmt.Errorf("error parsing config file: %w", err)
	}

	// Get MinIO config
	minioConfig := config.Test.WorkflowStorage.Minio

	// Initialize MinIO client
	client, err := minio.New(minioConfig.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioConfig.AccessKeyID, minioConfig.SecretAccessKey, ""),
		Secure: minioConfig.UseSSL,
		Region: minioConfig.Region,
	})

	return client, minioConfig, err
}

// ensureBucketExists checks if a bucket exists and creates it if needed
func ensureBucketExists(ctx context.Context, client *minio.Client, bucketName, region string) error {
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking bucket: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
		fmt.Printf("Created bucket %s\n", bucketName)
	}
	return nil
}
