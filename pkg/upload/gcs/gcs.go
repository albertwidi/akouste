package gcs

import (
	"context"
	"fmt"

	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/gcsblob"
	"github.com/google/go-cloud/gcp"
	"github.com/kosanapp/kosan-backend/pkg/upload"
	"golang.org/x/oauth2/google"
)

// Config of Google Cloud Storage
type Config struct {
	AccessJSON  string `yaml:"access_json" json:"access_json"`
	Bucket      string `yaml:"bucket" json:"bucket"`
	BucketProto string `yaml:"bucket_proto" json:"bucket_proto"`
	BucketURL   string `yaml:"bucket_url" json:"bucket_url"`
}

// GCS struct
type GCS struct {
	config            Config
	bucket            string
	baseURL           string
	storageBlobBucket *blob.Bucket
	httpClient        *gcp.HTTPClient
}

// New gcs storage
// gcs storage expect user already authenticated using 'gcs' command
// json path is a path to json service account path
func New(ctx context.Context, config Config) (*GCS, error) {
	var (
		creds *google.Credentials
		err   error
	)

	if config.AccessJSON != "" {
		jsonKey := []byte(config.AccessJSON)
		creds, err = google.CredentialsFromJSON(ctx, jsonKey, "https://www.googleapis.com/auth/cloud-platform")
	} else {
		// expecting the gcloud command is invoked
		creds, err = gcp.DefaultCredentials(ctx)
	}

	// check the credentials result
	if err != nil {
		return nil, err
	}

	client, err := gcp.NewHTTPClient(gcp.DefaultTransport(), gcp.CredentialsTokenSource(creds))
	if err != nil {
		return nil, err
	}

	bb, err := gcsblob.OpenBucket(ctx, config.Bucket, client, nil)
	if err != nil {
		return nil, err
	}

	gcs := GCS{
		config:            config,
		storageBlobBucket: bb,
		httpClient:        client,
	}
	return &gcs, nil
}

// GetBlobBucket function
func (gcs *GCS) GetBlobBucket() *blob.Bucket {
	return gcs.storageBlobBucket
}

// Name of upload provider
func (gcs *GCS) Name() string {
	return upload.StorageGCS
}

// BucketName return name of bucket used for upload
func (gcs *GCS) BucketName() string {
	return gcs.config.Bucket
}

// BucketURL return full path for bucket
func (gcs *GCS) BucketURL() string {
	return fmt.Sprintf("%s%s.%s", gcs.config.BucketProto, gcs.BucketName(), gcs.config.BucketURL)
}
