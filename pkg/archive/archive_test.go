package archive

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArchive(t *testing.T) {
	file := "testfile/unarchived.txt"
	target := "archived.tar.gz"

	err := Archive([]string{file}, target)
	assert.NoError(t, err)
	assert.FileExists(t, target)

	_ = os.RemoveAll(target)
}

func TestUnarchive(t *testing.T) {
	file := "testfile/archived.tar.gz"
	targetDIR := "./"
	targetFile := filepath.Join(targetDIR, "unarchived.txt")

	err := Unarchive(file, targetDIR)
	assert.NoError(t, err)
	assert.FileExists(t, targetFile)

	_ = os.RemoveAll(targetFile)
}
