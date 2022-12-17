package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const romsDir = "roms/"

func main() {
	library := readLibrary()

	fmt.Println("starting webserver on 8080")

	fs := http.FileServer(http.Dir("build"))
	http.Handle("/static", fs)

	http.HandleFunc("/library", func(writer http.ResponseWriter, request *http.Request) {
		handleLibraryReq(writer, library)
	})

	panic(http.ListenAndServe(":8080", nil))
}

func handleLibraryReq(w http.ResponseWriter, library map[string]string) {
	bytes, err := json.Marshal(library)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(bytes)
}

func readLibrary() map[string]string {
	fis, err := ioutil.ReadDir(romsDir)
	if err != nil {
		panic(err)
	}

	lib := map[string]string{}
	for _, fi := range fis {
		rom, err := os.Open(romsDir + fi.Name())
		if err != nil {
			panic(err)
		}

		contents, err := ioutil.ReadAll(rom)
		if err != nil {
			panic(err)
		}

		rom.Close()

		lib[fi.Name()] = string(contents)
	}

	return lib
}
