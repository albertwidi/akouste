package main

import (
	"context"
	"errors"
	"flag"
	"net/http"

	"github.com/albertwidi/akouste/downloader"
	"github.com/albertwidi/akouste/pkg/log"
	"github.com/albertwidi/akouste/pkg/storage"
	"github.com/albertwidi/akouste/pkg/storage/gcs"
	"github.com/albertwidi/akouste/pkg/storage/local"
	"github.com/gorilla/mux"
)

// appFlag contains app command-line flag
type appFlag struct {
	downloaderFlag
	storageProviderFlag

	logLevel string
}

type downloaderFlag struct {
	keepOldCount int
	destPath     string
}

type storageProviderFlag struct {
	bucketName  string
	bucketProto string
}

func main() {
	ctx := context.Background()

	appFlag := &appFlag{}
	flag.StringVar(&appFlag.logLevel, "logLevel", "info", "set the log level")
	flag.StringVar(&appFlag.bucketProto, "bucketProto", "", "the bucket provider/protocol ('gs', 'local', etc.)")
	flag.StringVar(&appFlag.bucketName, "bucketName", "", "the bucket name (dir path for 'local' bucketProto)")
	flag.StringVar(&appFlag.destPath, "downloadDIR", "", "download destination")
	flag.IntVar(&appFlag.keepOldCount, "keepOldCount", 5, "the number of downloaded versions to keep")
	flag.Parse()

	log.SetLevelString(appFlag.logLevel)

	storageProvider, err := newStorageProvider(ctx, appFlag.bucketProto, appFlag.bucketName)
	if err != nil {
		log.Fatalf("error initializing storage provider: %s", err.Error())
	}

	downloader, err := downloader.New(ctx, storageProvider, downloader.Config{
		DestPath:     appFlag.destPath,
		KeepOldCount: appFlag.keepOldCount,
	})
	if err != nil {
		log.Fatalf("error initializing downloader: %s\n", err.Error())
	}

	router := mux.NewRouter()
	handler := router.PathPrefix("/v1").Subrouter()
	handler.Methods("GET").Path("/ping").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PONG\n"))
	})
	handler.Methods("POST").Path("/download").HandlerFunc(downloader.HandlerDownload)

	log.Fatal(http.ListenAndServe(":9000", handler))
}

func newStorageProvider(ctx context.Context, bucketProto, bucketName string) (*storage.Storage, error) {
	switch bucketProto {
	case "gs":
		gcs, err := gcs.New(ctx, gcs.NewConfig(bucketName, ""))
		if err != nil {
			return nil, err
		}
		return storage.New(gcs), nil

	case "local":
		loc, err := local.New(local.NewConfig(bucketName))
		if err != nil {
			return nil, err
		}
		return storage.New(loc), nil

	default:
		return nil, errors.New("unknown bucket protocol")
	}
}
