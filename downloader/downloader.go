package downloader

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/albertwidi/akouste/pkg/archive"
	"github.com/albertwidi/akouste/pkg/log"
	"github.com/albertwidi/akouste/pkg/storage"
)

// Config struct for downloader package
type Config struct {
	// Where to store downloads
	DestPath string

	// Number of downloads to keep
	KeepOldCount int
}

// Downloader contains necessary downloader dependencies
type Downloader struct {
	config  Config
	storage *storage.Storage
}

// New returns initialized downloader client
func New(ctx context.Context, storage *storage.Storage, config Config) (*Downloader, error) {
	// Make sure the downloads directory exists
	if _, err := os.Stat(config.DestPath); err != nil {
		if err := os.MkdirAll(config.DestPath, 0755); err != nil {
			return nil, err
		}
	}

	return &Downloader{
		config:  config,
		storage: storage,
	}, nil
}

// HandlerDownload handles downloads and optionally decompresses the specified archive
// Accepted POST form fields:
// - uri       : filepath in the bucket
// - unarchive : whether to unarchive downloaded file (true/false)
//
// e.g. curl -X POST -d "uri=config-1.tar.gz&unarchive=true" localhost:9000/v1/download
func (d Downloader) HandlerDownload(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	err := r.ParseForm()
	if err != nil {
		msg := fmt.Sprintf("parse form failed: %s", err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	downloadFrom := r.PostForm.Get("uri")
	if downloadFrom == "" {
		http.Error(w, "empty uri field", http.StatusBadRequest)
		return
	}
	reader, err := d.storage.Download(ctx, downloadFrom)
	if err != nil {
		msg := fmt.Sprintf("error downloading %s: %s", downloadFrom, err.Error())
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	defer reader.Close()

	destinationFile := filepath.Join(d.config.DestPath, filepath.Base(downloadFrom))
	err = writeToFile(destinationFile, reader)
	if err != nil {
		log.Warnf("write file error: %s", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	unarchive := r.PostForm.Get("unarchive")
	if strings.ToLower(unarchive) == "true" {
		defer func() {
			// Delete the downloaded archive
			err = os.RemoveAll(destinationFile)
			if err != nil {
				log.Warnf("error delete: %s", err.Error())
			}

			// Ensures only 'keepOldCount' number of files are in the downloads directory
			err = deleteFilesExceedingN(d.config.DestPath, d.config.KeepOldCount)
			if err != nil {
				log.Warnf("error delete: %s", err.Error())
			}
		}()

		unarchiveDir := filepath.Join(d.config.DestPath, folderNameFromFileName(destinationFile))
		err = archive.Unarchive(destinationFile, unarchiveDir)
		if err != nil {
			log.Warnf("error unarchive: %s\n", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	log.Debugf("download success")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("download success\n"))
}

// writeToFile reads from an io.Reader into filepath
func writeToFile(filepath string, r io.Reader) error {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

// folderNameFromFileName returns a name for a folder
// which will be stripped off of its extensions.
func folderNameFromFileName(filename string) string {
	base := filepath.Base(filename)
	firstDot := strings.Index(base, ".")
	if firstDot > -1 {
		return base[:firstDot]
	}

	return base
}

// deleteFilesExceedingN deletes files that exceed N
func deleteFilesExceedingN(dir string, n int) error {
	var err error

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	sort.Slice(files, func(left, right int) bool {
		// sort by last modified time, newest to oldest
		return files[left].ModTime().Unix() > files[right].ModTime().Unix()
	})

	// This essentially removes files[n:]
	for i := n; i < len(files); i++ {
		path := filepath.Join(dir, files[i].Name())
		err = os.RemoveAll(path)
	}

	return err
}
