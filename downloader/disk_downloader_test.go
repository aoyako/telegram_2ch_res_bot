package downloader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskDownloader(t *testing.T) {
	assert := assert.New(t)

	type fields struct {
		Path     string
		MaxSpace uint64
	}
	type args struct {
		url string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   error
		wantSpace uint64
		wantPath  string
	}{
		{
			name: "Download and save",
			fields: fields{
				Path:     "src",
				MaxSpace: 100000000,
			},
			args: args{
				url: "https://s1.webmshare.com/DBj7M.webm",
			},
			wantErr:   nil,
			wantSpace: 3126038,
			wantPath:  "src/s1webmshare.comDBj7M.webm",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := NewDisckDownloader(tt.fields.Path, tt.fields.MaxSpace)
			err := d.Save(tt.args.url)
			assert.Equal(tt.wantErr, err)
			assert.Equal(tt.wantSpace, d.LoadedSpace)

			_, err = os.Stat(tt.wantPath)
			assert.Nil(err)

			path := d.Get(tt.args.url)
			assert.Equal(tt.wantPath, path)

			err = d.Free(tt.args.url)
			assert.Nil(err)

			_, err = os.Stat(tt.wantPath)
			assert.True(os.IsNotExist(err))

			assert.Equal(uint64(0), d.LoadedSpace)
		})
	}
}

func Test_normalizeURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeURL(tt.args.url); got != tt.want {
				t.Errorf("normalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
