package tst

import (
	"io/ioutil"
	"log"
)

func LoadJson(path string) []byte {
	log.Printf("Load json from file:%v", path)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	log.Printf("file: %s", file)

	file = ParseByte(file)

	log.Printf("parsed file: %s", file)

	return file
}
