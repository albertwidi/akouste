package gcs

import (
	"context"
	"errors"
	"fmt"

	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"golang.org/x/oauth2/google"
)

// Error variables
var (
	ErrEmptyBucketName = errors.New("empty bucket name")
)

// Config of Google Cloud Storage
type Config struct {
	AccessJSON string `yaml:"access_json" json:"access_json"`
	Bucket     string `yaml:"bucket" json:"bucket"`
}

// Validate validates configuration
func (c Config) Validate() error {
	if c.Bucket == "" {
		return ErrEmptyBucketName
	}
	return nil
}

// GCS struct
type GCS struct {
	config            Config
	storageBlobBucket *blob.Bucket
}

// New gcs storage
// gcs storage expect user already authenticated using 'gcs' command
// json path is a path to json service account path
func New(ctx context.Context, config Config) (*GCS, error) {
	var (
		creds *google.Credentials
		err   error
	)

	if err = config.Validate(); err != nil {
		return nil, err
	}

	if config.AccessJSON != "" {
		jsonKey := []byte(config.AccessJSON)
		creds, err = google.CredentialsFromJSON(ctx, jsonKey, "https://www.googleapis.com/auth/cloud-platform")
	} else {
		// Expects JSON filepath in GOOGLE_APPLICATION_CREDENTIALS env
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

	bb, err := gcsblob.OpenBucket(ctx, client, config.Bucket, nil)
	if err != nil {
		return nil, err
	}

	gcs := GCS{
		config:            config,
		storageBlobBucket: bb,
	}
	return &gcs, nil
}

// GetBlobBucket function
func (gcs *GCS) GetBlobBucket() *blob.Bucket {
	return gcs.storageBlobBucket
}

// Name of upload provider
func (gcs *GCS) Name() string {
	return "google-cloud-storage"
}

// BucketName return name of bucket used for upload
func (gcs *GCS) BucketName() string {
	return gcs.config.Bucket
}

// BucketURL return full path for bucket
func (gcs *GCS) BucketURL() string {
	return fmt.Sprintf("gs://%s", gcs.BucketName())
}
