package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

func flushMemory() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	mediaFolder := filepath.Join(wd, "media")
	files, err := ioutil.ReadDir(mediaFolder)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// we don't want to delete our ``.keepdir`` because of Git.
		if file.Name() == ".keepdir" {
			continue
		}
		err = os.RemoveAll(path.Join(mediaFolder, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
	}
}
