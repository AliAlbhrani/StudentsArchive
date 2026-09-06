package engine

import (
	"bytes"
	"io"
	"os"

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

func (c *CSClient) Upload(file io.Reader) (string, error) {
	objectName := helpers.GenerateUniqueName()
	objectPath := c.Path + "/" + objectName
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = c.Client.Write(objectPath, fileBytes, os.FileMode(0644))
	if err != nil {
		return "", err
	}
	return objectName, nil
}

func (c *CSClient) Download(objectName string) (io.Reader, error) {
	objectPath := c.Path + "/" + objectName
	fileBytes, err := c.Client.Read(objectPath)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(fileBytes)), nil
}
