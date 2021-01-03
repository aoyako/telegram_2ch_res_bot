package downloader

type Loader interface {
	Save(url string) error
	Free(url string) error
	Get(url string) string
}

type Downloader struct {
	Loader
}

func NewDownloader(path string, space uint64) *Downloader {
	return &Downloader{Loader: NewDisckDownloader(path, space)}
}
