package gcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSuccess(t *testing.T) {
	cfg := Config{Bucket: "bucketname"}
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestValidateFail(t *testing.T) {
	cfg := Config{Bucket: ""}
	err := cfg.Validate()
	assert.EqualError(t, ErrEmptyBucketName, err.Error())
}
