package downloader

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync/atomic"
)

// DiskDownloader can download files to disk
type DiskDownloader struct {
	Path        string // Directory to save resources
	MaxSpace    uint64 // Max size of resource files is Bytes
	LoadedSpace uint64 // Current space load status
}

// NewDisckDownloader constructor for DiskDownloader
func NewDisckDownloader(path string, space uint64) *DiskDownloader {
	os.MkdirAll(path, os.ModePerm)
	return &DiskDownloader{Path: path, MaxSpace: space}
}

// Save file with given url
// Format: https://addr.tmp/a/b/c/res.data -> Path + addrtmpabcres.data
func (d *DiskDownloader) Save(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	re := regexp.MustCompile("//")
	res := re.Split(url, -1)

	bty, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	size := len(bty)

	if uint64(size) >= d.MaxSpace {
		return errors.New("File is too large")
	}

	out, err := os.Create(path.Join(d.Path, normalizeURL(res[len(res)-1])))
	if err != nil {
		return err
	}
	defer out.Close()

	var success bool
	for !success {
		lastSpace := atomic.LoadUint64(&d.LoadedSpace)
		currentSpace := uint64(size) + lastSpace
		for !(currentSpace < d.MaxSpace) {
			lastSpace = atomic.LoadUint64(&d.LoadedSpace)
			currentSpace = uint64(size) + lastSpace
		}
		success = atomic.CompareAndSwapUint64(&d.LoadedSpace, lastSpace, currentSpace)
	}

	_, err = io.Copy(out, bytes.NewReader(bty))
	return err
}

// Free data from disk of given file
func (d *DiskDownloader) Free(url string) error {
	if !strings.HasSuffix(url, ".webm") {
		return d.delete(url)
	}
	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	fi, err := os.Stat(path.Join(d.Path, normalizeURL(res[len(res)-1])))
	if err != nil {
		return err
	}
	log.Printf("Free space: %d\n", atomic.AddUint64(&d.LoadedSpace, ^uint64(fi.Size()-1)))
	return d.delete(url)
}

// Deletes file from disk
func (d *DiskDownloader) delete(url string) error {
	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	return os.Remove(path.Join(d.Path, normalizeURL(res[len(res)-1])))
}

// Get path of file on disk with given url
func (d *DiskDownloader) Get(url string) string {
	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	return path.Join(d.Path, normalizeURL(res[len(res)-1]))
}

// Converts https://addr.tmp/a/b/c/res.data -> Path + addrtmpabcres.data
func normalizeURL(url string) string {
	return strings.Replace(strings.Replace(url, "/", "", -1), ".", "", 1)
}
