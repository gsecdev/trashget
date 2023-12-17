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

var NotFound = []byte("Page Not Found")
var cmdOpts Options

func main() {
	var err error

	err = cmdOpts.parseFlags()
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cmdOpts.IP, cmdOpts.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if cmdOpts.DoesThrottle() {
		ln = bwlimit.NewListener(ln, bwlimit.Byte(cmdOpts.Throttle)*1000*100, -1)
		log.Infof("throtteling connection at %d Mbit/s", cmdOpts.Throttle)
	}

	if cmdOpts.DoesAbort() {
		log.Infof("will abort downloads after %d%%", cmdOpts.AbortAfter)
	}

	http.Handle("/", http.HandlerFunc(handle))

	srv := &http.Server{}
	log.Printf("serving file '%s' with a virtual size of %d MB", cmdOpts.Filename, cmdOpts.Size)
	log.Printf("server started at %s:%d\n", cmdOpts.IP, cmdOpts.Port)
	log.Fatalf("failed to serve: %v", srv.Serve(ln))
}

type TrashReader struct {
	size       int64
	readIndex  int64
	abortAfter int64
}

func NewTrashReader(size int64, abortAfterPercent int) TrashReader {
	if abortAfterPercent == -1 {
		abortAfterPercent = 101
	}

	abortAfter := int64(size * int64(abortAfterPercent) / 100.0)
	return TrashReader{size: size, abortAfter: abortAfter, readIndex: 0}
}

func (r *TrashReader) Read(p []byte) (n int, err error) {
	if r.readIndex >= int64(r.size) {
		err = io.EOF
		return
	}

	n, err = rand.Reader.Read(p)
	if err != nil {
		err = io.ErrUnexpectedEOF
		return
	}

	r.readIndex += int64(n)
	if int64(n)+r.readIndex > r.abortAfter {
		err = fmt.Errorf("forcefully aborted after meeting threshold")
	}
	return
}

func setHeaders(w http.ResponseWriter, fileName string, fileSize int64) {
	//binary file
	w.Header().Set("Content-Type", "application/octet-stream")
	//filename
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, fileName))
	//length of the data
	w.Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	//disable caches
	w.Header().Set("Cache-Control", "private")
	w.Header().Set("Pragma", "private")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func handleMegaFile(w http.ResponseWriter, r *http.Request, pseudoSize int64) {
	setHeaders(w, cmdOpts.Filename, pseudoSize)
	w.WriteHeader(http.StatusOK)

	trashReader := NewTrashReader(pseudoSize, cmdOpts.AbortAfter)

	n, err := io.Copy(w, &trashReader)
	if err != nil {
		log.Warnf("aborted writing after %d bytes: %v", n, err)
		return
	}
	log.Info("Sent file completely")
}

func handle(w http.ResponseWriter, r *http.Request) {

	log.Infof("%v requested path: %s", r.RemoteAddr, r.URL.Path)

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
