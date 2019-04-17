package downloader

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/albertwidi/akouste/pkg/storage"
	"github.com/albertwidi/akouste/pkg/storage/local"
	"github.com/stretchr/testify/assert"
)

var cfg Config
var downloader *Downloader

func TestNew(t *testing.T) {
	localProvider, err := local.New(local.Config{Bucket: "."})
	assert.NoError(t, err)
	strg := storage.New(localProvider)

	downloader, err = New(context.TODO(), strg, Config{
		DestPath:     ".",
		KeepOldCount: 2,
	})
	assert.NoError(t, err)
}

func TestHandlerDownload(t *testing.T) {
	testfile := "testfile.txt"
	testunarchive := "false"
	testcontent := "hello"

	// Can I use internal function of the package I'm testing? o.O
	err := writeToFile(testfile, strings.NewReader(testcontent))
	assert.NoError(t, err)
	defer os.Remove(testfile)

	// Post form values
	form := url.Values{}
	form.Add("uri", testfile)
	form.Add("unarchive", testunarchive)
	request := httptest.NewRequest("POST", "/", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Test HandlerDownload
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(downloader.HandlerDownload)
	handler.ServeHTTP(rr, request)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check downloaded testfile content
	f, err := os.Open(testfile)
	assert.NoError(t, err)
	downloadedcontent, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, string(downloadedcontent), testcontent)

	// TODO test KeepOldCount behaviour too
}

func TestDeleteFilesExceedingN(t *testing.T) {
	testfilelist := []string{}
	testfiledir := "testfile"

	err := os.Mkdir(testfiledir, os.ModePerm)
	assert.NoError(t, err)
	defer os.RemoveAll(testfiledir)

	// Generate 3 dummy files
	for i := 0; i < 3; i++ {
		filepath := filepath.Join(testfiledir, fmt.Sprintf("test-%d", i))
		f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0755)
		assert.NoError(t, err)

		w := strings.NewReader("test")
		_, err = w.WriteTo(f)
		assert.NoError(t, err)
		f.Close()

		testfilelist = append(testfilelist, filepath)

		// Since the test is related to file's last modified time,
		// we wait 0.1 second between files generation
		// because the sorting is done in the nanosecond precision level anyway
		time.Sleep(100 * time.Millisecond)
	}

	n := 2
	// We know it's already ordered from oldest to newest, get last n elements
	expect := testfilelist[len(testfilelist)-n:]
	err = deleteFilesExceedingN(testfiledir, n)
	assert.NoError(t, err)

	// Check if generated files left in the directory matches expect
	generatedtestfiles, err := ioutil.ReadDir(testfiledir)
	assert.NoError(t, err)
	generatedtestfilelist := []string{}
	for i := 0; i < len(generatedtestfiles); i++ {
		generatedtestfilelist = append(generatedtestfilelist,
			filepath.Join(testfiledir, generatedtestfiles[i].Name()))
	}
	assert.Equal(t, expect, generatedtestfilelist)
}

func TestFolderNameFromFileName(t *testing.T) {
	testfilename := "path/to/file/filename.tar.gz"
	expect := "filename"
	result := folderNameFromFileName(testfilename)
	assert.Equal(t, expect, result)
}
