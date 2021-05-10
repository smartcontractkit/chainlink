package web

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
	"github.com/smartcontractkit/chainlink/core/logger"
)

const (
	acceptEncodingHeader  = "Accept-Encoding"
	contentEncodingHeader = "Content-Encoding"
	contentLengthHeader   = "Content-Length"
	rangeHeader           = "Range"
	varyHeader            = "Vary"
)

// ServeFileSystem wraps a http.FileSystem with an additional file existence check
type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

// BoxFileSystem implements ServeFileSystem with a packr box
type BoxFileSystem struct {
	packr.Box
}

// Exists implements the ServeFileSystem interface
func (b *BoxFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		return b.Has(p)
	}

	return false
}

// gzipFileHandler implements a http.Handler which can serve either the base
// file or the gzipped file depending on the Accept-Content header and the
// existence of the file
type gzipFileHandler struct {
	root ServeFileSystem
}

// GzipFileServer is a drop-in replacement for Go's standard http.FileServer
// which adds support for static resources precompressed with gzip, at
// the cost of removing the support for directory browsing.
func GzipFileServer(root ServeFileSystem) http.Handler {
	return &gzipFileHandler{root}
}

func (f *gzipFileHandler) openAndStat(path string) (http.File, os.FileInfo, error) {
	file, err := f.root.Open(path)
	var info os.FileInfo
	// This slightly weird variable reuse is so we can get 100% test coverage
	// without having to come up with a test file that can be opened, yet
	// fails to stat.
	if err == nil {
		info, err = file.Stat()
	}
	if err != nil {
		return file, nil, err
	}
	if info.IsDir() {
		return file, nil, fmt.Errorf("%s is directory", path)
	}
	return file, info, nil
}

// List of encodings we would prefer to use, in order of preference, best first.
// We only support gzip for now
var preferredEncodings = []string{"gzip"}

// File extension to use for different encodings.
func extensionForEncoding(encname string) string {
	switch encname {
	case "gzip":
		return ".gz"
	}
	return ""
}

// Find the best file to serve based on the client's Accept-Encoding, and which
// files actually exist on the filesystem. If no file was found that can satisfy
// the request, the error field will be non-nil.
func (f *gzipFileHandler) findBestFile(w http.ResponseWriter, r *http.Request, fpath string) (http.File, os.FileInfo, error) {
	ae := r.Header.Get(acceptEncodingHeader)
	// Send the base file if no AcceptEncoding header is provided
	if ae == "" {
		return f.openAndStat(fpath)
	}

	// Got an accept header? See what possible encodings we can send by looking for files
	var available []string
	for _, posenc := range preferredEncodings {
		ext := extensionForEncoding(posenc)
		fname := fpath + ext

		if f.root.Exists("/", fname) {
			available = append(available, posenc)
		}
	}

	// Negotiate the best content encoding to use
	negenc := negotiateContentEncoding(r, available)
	if negenc == "" {
		// If we fail to negotiate anything try the base file
		return f.openAndStat(fpath)
	}

	ext := extensionForEncoding(negenc)
	if file, info, err := f.openAndStat(fpath + ext); err == nil {
		wHeader := w.Header()
		wHeader[contentEncodingHeader] = []string{negenc}
		wHeader.Add(varyHeader, acceptEncodingHeader)

		if len(r.Header[rangeHeader]) == 0 {
			// If not a range request then we can easily set the content length which the
			// Go standard library does not do if "Content-Encoding" is set.
			wHeader[contentLengthHeader] = []string{strconv.FormatInt(info.Size(), 10)}
		}
		return file, info, nil
	}

	// If all else failed, fall back to base file once again
	return f.openAndStat(fpath)
}

// Determines the best encoding to use
func negotiateContentEncoding(r *http.Request, available []string) string {
	values := strings.Split(r.Header.Get(acceptEncodingHeader), ",")
	aes := []string{}

	// Clean the values
	for _, v := range values {
		aes = append(aes, strings.TrimSpace(v))
	}

	for _, a := range available {
		for _, acceptEnc := range aes {
			if acceptEnc == a {
				return a
			}
		}
	}

	return ""
}

// Implements http.Handler
func (f *gzipFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	fpath := path.Clean(upath)
	if strings.HasSuffix(fpath, "/") {
		http.NotFound(w, r)
		return
	}

	// Find the best acceptable file, including trying uncompressed
	if file, info, err := f.findBestFile(w, r, fpath); err == nil {
		http.ServeContent(w, r, fpath, info.ModTime(), file)
		logger.ErrorIfCalling(file.Close)
		return
	}

	http.NotFound(w, r)
}

// Static returns a middleware handler that serves static files in the given directory.
func ServeGzippedAssets(urlPrefix string, fs ServeFileSystem) gin.HandlerFunc {
	fileserver := GzipFileServer(fs)
	if urlPrefix != "" {
		fileserver = http.StripPrefix(urlPrefix, fileserver)
	}
	return func(c *gin.Context) {
		if fs.Exists(urlPrefix, c.Request.URL.Path) {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
	}
}
