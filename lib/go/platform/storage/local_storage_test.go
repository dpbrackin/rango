package storage_test

import (
	"bytes"
	"context"
	"os"
	"rango/core"
	"rango/platform/storage"
	"strings"
	"testing"
)

func TestDownloadUploadedFileToBuffer(t *testing.T) {
	ctx := context.TODO()
	cwd, _ := os.Getwd()
	storage := storage.NewDiscStorage(storage.NewDiscStorageParams{
		BasePath: cwd,
	})

	fileName := "TEST.txt"
	content := "HELLO TEST"
	testReader := strings.NewReader(content)

	source, err := storage.Upload(ctx, core.UploadParams{
		Reader: testReader,
		Name:   fileName,
	})

	downloadedContent := bytes.NewBuffer(nil)
	err = storage.Download(ctx, source, downloadedContent)

	if err != nil {
		t.Errorf("Download had an error: %s", err.Error())
	}

	if downloadedContent.String() != content {
		t.Errorf("Expected: %s. Got: %s", content, downloadedContent.String())
	}

	// clean up
	os.Remove(source)
}
