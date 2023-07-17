package downloader

// Loader interface can load data from links
type Loader interface {
	Free(url string) error
	Get(url string) string
}

// Downloader can download data from links
type Downloader struct {
	Loader
}

// NewDownloader constructor for Downloader
func NewDownloader(path string) *Downloader {
	return &Downloader{Loader: NewDisckDownloader(path)}
}
