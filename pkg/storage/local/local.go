package local

import (
	"errors"
	"fmt"

	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
)

// Config for local storage
type Config struct {
	// Base local filesystem path
	// e.g. "/tmp" or "./test-download-location"
	Bucket string
}

// Validate config
func (c *Config) Validate() error {
	if c.Bucket == "" {
		return errors.New("empty bucket path")
	}
	return nil
}

// Local storage struct
type Local struct {
	config            Config
	storageBlobBucket *blob.Bucket
}

// New local storage
func New(config Config) (*Local, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	bb, err := fileblob.OpenBucket(config.Bucket, &fileblob.Options{})
	if err != nil {
		return nil, err
	}

	return &Local{
		config:            config,
		storageBlobBucket: bb,
	}, nil
}

// GetBlobBucket function
func (l *Local) GetBlobBucket() *blob.Bucket {
	return l.storageBlobBucket
}

// Name of upload provider
func (l *Local) Name() string {
	return "local-file"
}

// BucketName return name of bucket used for upload
func (l *Local) BucketName() string {
	return l.config.Bucket
}

// BucketURL return full path for bucket
func (l *Local) BucketURL() string {
	return fmt.Sprintf("%s.%s", "file://", l.BucketName())
}
