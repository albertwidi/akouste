package storage

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/albertwidi/akouste/pkg/storage/local"
	"github.com/stretchr/testify/assert"
)

var storage *Storage

func TestNew(t *testing.T) {
	localProvider, err := local.New(local.NewConfig("."))
	assert.NoError(t, err)

	storage = New(localProvider)
}

func TestUpload(t *testing.T) {
	b := []byte("testing")
	target := "testfile.txt"
	targetAttrs := "testfile.txt.attrs"

	uploadedPath, err := storage.Upload(context.TODO(), b, target)
	assert.NoError(t, err)
	assert.Equal(t, uploadedPath, target)
	assert.FileExists(t, target)

	_ = os.RemoveAll(target)
	_ = os.RemoveAll(targetAttrs)
}

func TestDownload(t *testing.T) {
	testByte := []byte("hello")
	testfile := "testfile.txt"
	testByteBuffer := bytes.NewBuffer(testByte)

	f, err := os.OpenFile(testfile, os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)
	_, err = testByteBuffer.WriteTo(f)
	assert.NoError(t, err)
	defer os.Remove(testfile)

	readcloser, err := storage.Download(context.TODO(), testfile)
	downloadedBuf := new(bytes.Buffer)
	downloadedBuf.ReadFrom(readcloser)
	defer readcloser.Close()

	assert.Equal(t, testByte, downloadedBuf.Bytes())
}
