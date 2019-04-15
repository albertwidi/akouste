package gcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSuccess(t *testing.T) {
	cfg := NewConfig("bucketname", "")
	err := cfg.Validate()
	assert.NoError(t, err)
}

func TestValidateFail(t *testing.T) {
	cfg := NewConfig("", "")
	err := cfg.Validate()
	assert.EqualError(t, ErrEmptyBucketName, err.Error())
}
