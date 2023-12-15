package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/conduitio/bwlimit"
	log "github.com/sirupsen/logrus"
)

var InternalServerError = []byte("Internal Server Error")
var NotFound = []byte("Page Not Found")
var cmdOpts Options

func main() {
	var err error

	writeHelp, err := cmdOpts.parseFlags()
	if err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}

	if cmdOpts.Filename == "" {
		writeHelp()
		log.Errorf("empty filename specified")
	}

	if cmdOpts.Port < 0 || cmdOpts.Port > 65535 {
		writeHelp()
		log.Errorf("port needs to be in rage 0-65535")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cmdOpts.IP, cmdOpts.Port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if cmdOpts.Throttle > 0 {
		ln = bwlimit.NewListener(ln, bwlimit.Byte(cmdOpts.Throttle)*1000*100, -1)
		log.Infof("throtteling connection at %d Mbit/s", cmdOpts.Throttle)
	}

	http.Handle("/", http.HandlerFunc(handle))

	srv := &http.Server{}
	log.Printf("serving file '%s' with a virtual size of %d MB", cmdOpts.Filename, cmdOpts.Size)
	log.Printf("server started at %s:%d\n", cmdOpts.IP, cmdOpts.Port)
	log.Fatalf("Failed to serve: %v", srv.Serve(ln))

	// log.Println(server.ListenAndServe())
}

type TrashReader struct {
	size      int64
	readIndex int64
}

func (r *TrashReader) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(r.size) {
		err = io.EOF
		return
	}

	// rand.Reader

	// n = copy(p, r.data[r.readIndex:])
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
	w.Header().Set("Pragma", "private")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func handleMegaFile(w http.ResponseWriter, r *http.Request, pseudoSize int64) {
	setHeaders(w, cmdOpts.Filename, pseudoSize)
	w.WriteHeader(http.StatusOK)

	n, err := io.CopyN(w, rand.Reader, pseudoSize)
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
