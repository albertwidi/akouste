package upload

import (
	"context"
	"io/ioutil"
	"path"

	"github.com/google/go-cloud/blob"
)

// storage type of artifact
const (
	// local storage for testing
	StorageLocal = "local"
	// google cloud storage
	StorageGCS = "gcs"
	// amazon s3 storage
	StorageS3 = "s3"
	// minio storage
	StorageMinio = "minio"
	// digital ocean storage
	// same API as S3
	StorageDO = "do"
)

// StorageProvider interface
type StorageProvider interface {
	GetBlobBucket() *blob.Bucket
	Name() string
	BucketName() string
	BucketURL() string
}

// Upload struct
type Upload struct {
	storage StorageProvider
}

// New artifact
func New(storage StorageProvider) *Upload {
	return &Upload{storage}
}

// Upload file from bytes
func (u *Upload) Upload(ctx context.Context, content []byte, destination string) (string, error) {
	return u.upload(ctx, content, destination)
}

// UploadFile file from source and destination
func (u *Upload) UploadFile(ctx context.Context, source, destination string) (string, error) {
	p, err := ioutil.ReadFile(source)
	if err != nil {
		return "", err
	}

	return u.upload(ctx, p, destination)
}

func (u *Upload) upload(ctx context.Context, content []byte, destination string) (string, error) {
	uploadPath := path.Join(u.storage.BucketURL(), destination)
	blobBucket := u.storage.GetBlobBucket()
	err := blobBucket.WriteAll(ctx, destination, content, nil)

	return uploadPath, err
}

// Name of provider
func (u *Upload) Name() string {
	return u.storage.Name()
}

// BucketName of storage provider
func (u *Upload) BucketName() string {
	return u.storage.BucketName()
}
