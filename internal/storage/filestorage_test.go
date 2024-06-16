package storage

import (
	"context"
	"os"
	"testing"

	"github.com/impr0ver/gophKeeper/internal/serverconfig"
	"github.com/impr0ver/gophKeeper/internal/userdata"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var filesPath = serverconfig.NewServerConfig().FilesStore

func TestNewFileStorage(t *testing.T) {
	assert.NotPanics(t, func() {
		newFileStorage(filesPath)
	})
}

func TestFileStorage_CreateRecord(t *testing.T) {
	storage := newFileStorage(filesPath)

	tc := []struct {
		name    string
		prepare func()
		valid   func()
	}{
		{
			"Create file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), userdata.Record{
					ID:   "1",
					Type: userdata.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				assert.DirExists(t, filesPath)
				assert.FileExists(t, filesPath+"/1")
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.prepare()
		test.valid()
	}

	assert.NoError(t, os.RemoveAll(filesPath))
}

func TestFileStorage_GetRecord(t *testing.T) {
	storage := newFileStorage(filesPath)

	tc := []struct {
		name    string
		prepare func()
		valid   func()
	}{
		{
			"Get existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), userdata.Record{
					ID:   "1",
					Type: userdata.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				md := metadata.Pairs("recordMetadata", "file.txt")
				ctx := metadata.NewIncomingContext(context.Background(), md)
				record, err := storage.GetRecord(ctx, "1")
				assert.NoError(t, err)

				assert.Equal(t, userdata.Record{
					ID:       "1",
					Metadata: "file.txt",
					Type:     userdata.TypeFile,
					Data:     []byte("text"),
				}, record)

			},
		},
		{
			"Get existed file record, but don't provide record metadata",
			func() {
				id, err := storage.CreateRecord(context.Background(), userdata.Record{
					ID:   "1",
					Type: userdata.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				record, err := storage.GetRecord(context.Background(), "1")
				assert.Equal(t, ErrUnknown, err)
				assert.Empty(t, record)
			},
		},
		{
			"Get non existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), userdata.Record{
					ID:   "1",
					Type: userdata.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				md := metadata.Pairs("recordMetadata", "file.txt")
				ctx := metadata.NewIncomingContext(context.Background(), md)

				record, err := storage.GetRecord(ctx, "2")
				assert.Equal(t, ErrNotFound, err)
				assert.Empty(t, record)
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.prepare()
		test.valid()
	}

	assert.NoError(t, os.RemoveAll(filesPath))
}

func TestFileStorage_DeleteRecord(t *testing.T) {
	storage := newFileStorage(filesPath)

	tc := []struct {
		name    string
		prepare func()
		valid   func()
	}{
		{
			"Delete existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), userdata.Record{
					ID:   "1",
					Type: userdata.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				err := storage.DeleteRecord(context.Background(), "1")
				assert.NoError(t, err)
				assert.NoFileExists(t, filesPath+"/1")
			},
		},
		{
			"Delete non existed file record",
			func() {
				id, err := storage.CreateRecord(context.Background(), userdata.Record{
					ID:   "1",
					Type: userdata.TypeFile,
					Data: []byte("text"),
				})
				assert.NoError(t, err)
				assert.Equal(t, "1", id)
			},
			func() {
				err := storage.DeleteRecord(context.Background(), "2")
				assert.Equal(t, ErrNotFound, err)
				assert.NoFileExists(t, filesPath+"/2")
			},
		},
	}

	for _, test := range tc {
		t.Log(test.name)
		test.prepare()
		test.valid()
	}

	assert.NoError(t, os.RemoveAll(filesPath))
}
