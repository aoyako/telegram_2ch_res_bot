package downloader

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
	"sync/atomic"
)

type DiskDownloader struct {
	Path        string
	MaxSpace    uint64
	LoadedSpace uint64
}

func NewDisckDownloader(path string, space uint64) *DiskDownloader {
	return &DiskDownloader{Path: path, MaxSpace: space}
}

func (d *DiskDownloader) Save(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	out, err := os.Create(path.Join(d.Path, normalizeURL(res[len(res)-1])))
	if err != nil {
		return err
	}
	defer out.Close()

	// size, err := io.Copy(ioutil.Discard, resp.Body)
	bty, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	size := len(bty)

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

func (d *DiskDownloader) Free(url string) error {
	if !strings.HasSuffix(url, ".webm") {
		return d.Delete(url)
	}
	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	fi, err := os.Stat(path.Join(d.Path, normalizeURL(res[len(res)-1])))
	if err != nil {
		return err
	}
	fmt.Println(atomic.AddUint64(&d.LoadedSpace, ^uint64(fi.Size()-1)))
	return d.Delete(url)
}

func (d *DiskDownloader) Delete(url string) error {
	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	return os.Remove(path.Join(d.Path, normalizeURL(res[len(res)-1])))
}

func (d *DiskDownloader) Get(url string) string {
	re := regexp.MustCompile("//")
	res := re.Split(url, -1)
	return path.Join(d.Path, normalizeURL(res[len(res)-1]))
}

func normalizeURL(url string) string {
	return strings.Replace(strings.Replace(url, "/", "", -1), ".", "", 1)
}
