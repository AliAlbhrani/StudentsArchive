package main

import (
	"github.com/AliAlbhrani/StudentsArchive/api"
	_ "github.com/AliAlbhrani/StudentsArchive/engine"
	_ "github.com/AliAlbhrani/StudentsArchive/env"
	_ "github.com/studio-b12/gowebdav"
)

func main() {
	api.Init()
}
