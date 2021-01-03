package downloader

// Loader interface can load data from links
type Loader interface {
	Save(url string) error
	Free(url string) error
	Get(url string) string
}

// Downloader can download data from links
type Downloader struct {
	Loader
}

// NewDownloader constructor for Downloader
func NewDownloader(path string, space uint64) *Downloader {
	return &Downloader{Loader: NewDisckDownloader(path, space)}
}
