package engine

import (
	"io"
	"mime"
	"os"
	"path/filepath"

	"github.com/AliAlbhrani/StudentsArchive/env"
	"github.com/AliAlbhrani/StudentsArchive/helpers"
	"github.com/studio-b12/gowebdav"
)

type CSClient struct {
	Client   *gowebdav.Client
	BaseURL  string
	AppName  string
	Password string
	Path     string
}

var CS *CSClient
var _ = func() bool {
	CS = &CSClient{
		Client:   gowebdav.NewClient(env.STORAGE_BASE_URL, env.STORAGE_APP_NAME, env.STORAGE_PASSWORD),
		BaseURL:  env.STORAGE_BASE_URL,
		AppName:  env.STORAGE_APP_NAME,
		Password: env.STORAGE_PASSWORD,
		Path:     env.STORAGE_PATH,
	}
	err := CS.Client.Connect()
	if err != nil {
		panic(err)
	}
	return true
}()

func (c *CSClient) Upload(file io.Reader, originalFilename string) (string, error) {
	ext := filepath.Ext(originalFilename) // e.g. ".png"
	objectName := helpers.GenerateUniqueName() + ext
	objectPath := c.Path + "/" + objectName

	err := c.Client.WriteStream(objectPath, file, os.FileMode(0644))
	if err != nil {
		return "", err
	}
	return objectName, nil
}

type DownloadResult struct {
	Body        io.ReadCloser
	Size        int64
	ContentType string
}

func (c *CSClient) Download(objectName string) (*DownloadResult, error) {
	objectPath := c.Path + "/" + objectName

	stream, err := c.Client.ReadStream(objectPath)
	if err != nil {
		return nil, err
	}

	info, err := c.Client.Stat(objectPath)
	if err != nil {
		stream.Close()
		return nil, err
	}

	contentType := mime.TypeByExtension(filepath.Ext(objectName))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return &DownloadResult{
		Body:        stream,
		Size:        info.Size(),
		ContentType: contentType,
	}, nil
}
