package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"s.mcquay.me/sm/trash"
)

var InternalServerError = []byte("Internal Server Error")
var NotFound = []byte("Page Not Found")
var cmdOpts Options

func main() {
	var err error

	_, err = cmdOpts.parseFlags()
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}

	server := http.Server{
		Handler: http.HandlerFunc(handle),
		Addr:    fmt.Sprintf("%s:%d", cmdOpts.IP, cmdOpts.Port),
	}

	log.Printf("server started at %s:%d\n", cmdOpts.IP, cmdOpts.Port)
	log.Println(server.ListenAndServe())
}

type TrashReader struct {
	data      []byte
	readIndex int64
}

func (r *TrashReader) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(len(r.data)) {
		err = io.EOF
		return
	}

	n = copy(p, r.data[r.readIndex:])
	r.readIndex += int64(n)
	return
}

func setHeaders(w http.ResponseWriter, fileName string, fileSize int64) {
	//Represents binary file
	w.Header().Set("Content-Type", "application/octet-stream")
	//Tells client what filename should be used.
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	//The length of the data.
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	//No cache headers.
	w.Header().Set("Cache-Control", "private")
	//No cache headers.
	w.Header().Set("Pragma", "private")
	//No cache headers.
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func handleMegaFile(w http.ResponseWriter, r *http.Request, pseudoSize int64) {
	//Set the headers
	setHeaders(w, cmdOpts.Filename, pseudoSize)
	w.WriteHeader(http.StatusOK)

	//Copy without loading everything in memory
	n, err := io.CopyN(w, trash.LoHi, pseudoSize)
	if err != nil {
		log.Errorf("error writing: %v", err)
		return
	}
	log.Printf("LARGE :: Written : %d", n)
}

func handle(w http.ResponseWriter, r *http.Request) {

	log.Printf("requested path: %s", r.URL.Path)

	if cmdOpts.Uri == "/" || r.URL.Path == cmdOpts.Uri {
		handleMegaFile(w, r, 1024*1024*cmdOpts.Size)
	} else {
		notFound(w)
	}

}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write(NotFound)
}

// func handleError(err error, w http.ResponseWriter) {
// 	w.WriteHeader(http.StatusInternalServerError)
// 	w.Write(InternalServerError)
// }
