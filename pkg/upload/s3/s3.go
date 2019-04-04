package s3

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/google/go-cloud/blob"
	"github.com/google/go-cloud/blob/s3blob"
	"github.com/kosanapp/kosan-backend/pkg/upload"
)

// S3 struct
type S3 struct {
	storageBlobBucket *blob.Bucket
	config            Config
}

// Config of S3
type Config struct {
	Region         string `yaml:"region"`
	Endpoint       string `yaml:"endpoint"`
	Bucket         string `yaml:"bucket" json:"bucket"`
	ClientID       string `yaml:"client_id" json:"client_id"`
	ClientSecret   string `yaml:"client_secret" json:"client_secret"`
	DisableSSL     bool   `yaml:"disable_ssl" json:"disable_ssl"`
	ForcePathStyle bool   `yaml:"force_path_style" json:"force_path_style"`
	BucketProto    string `yaml:"bucket_proto" json:"bucket_proto"`
	BucketURL      string `yaml:"bucket_url" json:"bucket_url"`
}

// New S3 storage
func New(ctx context.Context, config Config) (*S3, error) {
	c := aws.Config{
		Region:           aws.String(config.Region),
		Credentials:      credentials.NewStaticCredentials(config.ClientID, config.ClientSecret, ""),
		DisableSSL:       aws.Bool(config.DisableSSL),
		S3ForcePathStyle: aws.Bool(config.ForcePathStyle),
	}
	// if we want to use digitalocean or another api that compatible to s3
	if config.Endpoint != "" {
		c.Endpoint = aws.String(config.Endpoint)
	}

	sess, err := session.NewSession(&c)
	if err != nil {
		return nil, err
	}

	bb, err := s3blob.OpenBucket(ctx, config.Bucket, sess, nil)
	if err != nil {
		return nil, err
	}

	s := S3{
		storageBlobBucket: bb,
		config:            config,
	}
	return &s, nil
}

// GetBlobBucket function
func (s3 *S3) GetBlobBucket() *blob.Bucket {
	return s3.storageBlobBucket
}

// Name of upload provider
func (s3 *S3) Name() string {
	return upload.StorageS3
}

// BucketName return name of bucket used for upload
func (s3 *S3) BucketName() string {
	return s3.config.Bucket
}

// BucketURL return full path for bucket
func (s3 *S3) BucketURL() string {
	return fmt.Sprintf("%s%s.%s", s3.config.BucketProto, s3.BucketName(), s3.config.BucketURL)
}
