package downloader

import (
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync/atomic"
)

// DiskDownloader can download files to disk
type DiskDownloader struct {
	Path        string // Directory to save resources
	LoadedSpace uint64 // Current space load status
}

// NewDisckDownloader constructor for DiskDownloader
func NewDisckDownloader(path string) *DiskDownloader {
	os.MkdirAll(path, os.ModePerm)
	return &DiskDownloader{Path: path}
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
