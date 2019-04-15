package storage

import (
	"context"
	"io"
	"io/ioutil"
	"path"

	"gocloud.dev/blob"
)

// Provider interface
type Provider interface {
	GetBlobBucket() *blob.Bucket
	Name() string
	BucketName() string
	BucketURL() string
}

// Storage struct
type Storage struct {
	provider Provider
}

// New artifact
func New(provider Provider) *Storage {
	return &Storage{provider}
}

// Name of provider
func (s *Storage) Name() string {
	return s.provider.Name()
}

// BucketName of storage provider
func (s *Storage) BucketName() string {
	return s.provider.BucketName()
}

// Download file
func (s *Storage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.download(ctx, key)
}

func (s *Storage) download(ctx context.Context, key string) (io.ReadCloser, error) {
	blobBucket := s.provider.GetBlobBucket()
	r, err := blobBucket.NewReader(ctx, key, &blob.ReaderOptions{})

	return r, err
}

// Upload file from bytes
func (s *Storage) Upload(ctx context.Context, content []byte, destination string) (string, error) {
	return s.upload(ctx, content, destination)
}

// UploadFile file from source and destination
func (s *Storage) UploadFile(ctx context.Context, source, destination string) (string, error) {
	p, err := ioutil.ReadFile(source)
	if err != nil {
		return "", err
	}

	return s.upload(ctx, p, destination)
}

func (s *Storage) upload(ctx context.Context, content []byte, destination string) (string, error) {
	uploadPath := path.Join(s.provider.BucketURL(), destination)
	blobBucket := s.provider.GetBlobBucket()
	err := blobBucket.WriteAll(ctx, destination, content, nil)

	return uploadPath, err
}
