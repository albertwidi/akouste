package local

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Local storage struct
type Local struct {
	config Config
}

// New local storage
func New(config Config) (*Local, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	l := Local{
		config: config,
	}
	return &l, nil
}

// Upload function
func (l *Local) Upload(ctx context.Context, filePath, fileName string) error {
	from := path.Join(filePath, fileName)
	ffrom, err := ioutil.ReadFile(from)
	if err != nil {
		return err
	}

	copypath := l.bucketPath(fileName)
	// might need to check error for a more verbose error handling
	if err := ioutil.WriteFile(copypath, ffrom, os.ModePerm); err != nil {
		return fmt.Errorf("Error when copying %s to %s in local storage. With error: %s", from, copypath, err.Error())
	}
	return nil
}

func (l *Local) bucketPath(filepath string) string {
	return path.Join(l.config.Bucket, filepath)
}

// Config for local storage
type Config struct {
	Bucket string
}

// Validate config
func (c *Config) Validate() error {
	if c.Bucket == "" {
		return errors.New("Empty bucket path")
	}
	return nil
}

// NewConfig function
func NewConfig(bucket string) Config {
	c := Config{
		Bucket: bucket,
	}
	return c
}
