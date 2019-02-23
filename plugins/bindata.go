// Code generated by go-bindata.
// sources:
// ../ui/dist/index.html
// ../ui/dist/main.js
// DO NOT EDIT!

package plugins

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _indexHtml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x3c\x8f\xb1\x8e\xc2\x30\x0c\x86\x5f\xc5\x97\xfd\x2e\xeb\x0d\x8e\x97\x3b\xd8\x10\x0c\x65\x60\x0c\x89\x45\x53\xb5\x69\x15\x9b\xa8\x7d\x7b\x54\x5a\x31\xd9\xdf\x27\xd9\xbf\x7e\xfc\xfa\x3f\xff\x35\xb7\xcb\x01\x5a\x1d\x7a\x02\x5c\x07\xf4\x3e\x3f\x1c\xe7\x15\xd9\x47\x02\x1c\x58\x3d\x84\xd6\x17\x61\x75\xd7\xe6\xf8\xfd\x4b\x80\x9a\xb4\x67\x3a\x8d\x99\x67\x0e\x68\x37\x04\xb4\xfb\xcd\x7d\x8c\x0b\x01\xc6\x54\x21\x45\x27\xcf\x89\x4b\x4d\x32\x16\x21\xb4\x31\x55\x02\x94\x50\xd2\xa4\xa0\xcb\xc4\xce\x28\xcf\x6a\x3b\x5f\xfd\x66\x0d\x48\x09\xce\x0c\x3e\xe5\x9f\x4e\x0c\xa1\xdd\x3c\xa1\xdd\x1f\x7f\x78\x5f\xd6\xe4\x77\x87\x57\x00\x00\x00\xff\xff\x65\x5a\x0d\x38\xd4\x00\x00\x00")

func indexHtmlBytes() ([]byte, error) {
	return bindataRead(
		_indexHtml,
		"index.html",
	)
}

func indexHtml() (*asset, error) {
	bytes, err := indexHtmlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "index.html", size: 212, mode: os.FileMode(420), modTime: time.Unix(1530190485, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _mainJs = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x3b\xed\x72\xdb\x38\x92\xaf\x22\xf3\x66\x39\x40\x05\xe6\xc8\xd9\xd9\xab\x39\x6a\x10\x5f\x26\xe3\x4c\x9c\xf5\x24\xb3\xf9\xd8\xdd\x2a\x9f\xcb\x43\x91\x2d\x09\x36\x05\x28\x20\x24\xcb\x23\xf3\x7e\xdf\x53\x5c\xdd\xb3\xdc\xa3\xdc\x93\x5c\x35\x40\x90\xa0\x3e\xf2\x71\x53\xb5\x75\xa9\x94\x49\x02\x8d\x06\xfa\x13\xdd\x40\xeb\x68\xb2\x94\xb9\x11\x4a\x12\xa0\x9b\x55\xa6\x07\x86\x6f\xea\x91\x6f\x1c\x48\x22\xe8\x46\x4c\x88\xb9\x14\x57\x54\x83\x59\x6a\x39\xc0\xf7\x04\xd6\x0b\xa5\x4d\x35\xc2\x21\x9a\x63\x13\xdf\x88\x54\xb0\x32\x3d\x3a\x61\x4d\x67\xba\xa9\xeb\x51\x33\x08\x70\x50\x9e\x95\x25\xd1\x7e\x2c\xd3\xac\x7b\x97\x94\xe9\xa4\xe4\x47\xc3\xae\xad\x96\xc9\x9c\x03\x93\x49\xce\x0d\x93\x49\xc1\xbb\xa5\x32\xc3\x04\xdd\xc8\x44\xe1\x2b\x7d\x78\x78\x3d\xbe\x81\xdc\x24\x05\x4c\x84\x84\x5f\xb4\x5a\x80\x36\xf7\x16\x6c\x03\x72\x39\x07\x9d\x8d\x4b\x48\x8f\x86\x6c\x0a\x26\x15\x35\xad\x99\x4c\x34\x0f\x49\x8f\x96\xd2\x8d\x2e\xa2\x23\x6e\xee\x17\xa0\x26\x83\xb7\xf7\xf3\xb1\x2a\xe3\xd8\x3d\x13\xa3\xde\x1a\x2d\xe4\xf4\x5d\x36\x8d\xe3\x43\x33\xee\xc2\xb2\xcd\x2a\x2b\x97\x90\x46\x3f\xab\x62\x59\x42\x54\x53\x76\x68\x70\x74\x7d\x0d\x55\x03\xe6\x87\x1d\x0d\xdd\x72\x4d\x8f\x7c\x2b\x94\x93\xd8\xc4\x31\x01\x8e\x04\x50\xf6\x5d\x6c\xbc\x84\x60\x24\x26\xe4\x5b\xec\x8d\x94\x9d\x2a\xe2\x9e\x26\x88\x63\xfc\x9f\x74\x33\x75\x83\x50\x96\x82\x37\x8b\xcb\x35\x64\x06\x88\x5c\x96\x25\x45\x74\x32\xd1\x44\x1c\x5a\xba\x60\x51\x01\x93\x6c\x59\x9a\x68\x9b\xe3\x8e\x0a\xa8\x29\x7b\x6c\x17\x54\x59\xbe\x74\x4c\x06\x3a\x51\x9a\x58\x35\x1a\x08\x39\x00\x2a\x93\x82\x08\xa6\x59\x4b\xae\xa1\x9b\x56\x89\xcc\x55\x9d\x8c\x85\x2c\xec\xba\x98\xa6\xd4\xeb\x97\x40\x1e\x49\xbe\xab\xcd\x5b\xd4\x9e\xb6\x10\x1d\xd6\xa4\x59\x7b\x9d\xee\xe9\x6c\x35\x18\xd7\x65\x58\x94\x45\xcc\x50\x66\x70\x3a\xb5\x25\x92\x06\xb0\x61\xd1\x42\x2b\xa3\x90\xc8\x64\x96\x55\xaf\xef\xa4\x67\x96\xb3\x02\x1c\x80\x38\x16\x3c\x8a\x98\x24\x32\xa9\xf8\x90\xd6\xe4\xb2\xa7\xe3\x12\xf5\xb2\x82\x01\xf2\x2c\x37\x51\x67\x96\x82\xd0\x4d\xdd\x7e\x69\x37\xbd\xe7\xa3\x44\x3e\x1a\x0a\x97\xf2\x8a\x9b\x4b\x79\xd5\x9a\x60\x37\x42\x1d\x1e\x71\xb2\x07\xbc\x72\xe0\x26\xc9\x16\x0b\x90\xc5\xb3\x99\x28\x0b\x02\xb4\x03\xc8\xfc\x72\x4d\x22\x64\x05\xda\xfc\x00\x13\xa5\x81\x00\x93\x01\x54\x8e\x52\x81\x64\x91\x69\x90\xe6\x95\x2a\x20\xd1\x30\x57\x2b\xd8\xc5\x57\x6e\xad\x8f\x0f\x47\xf2\x7b\x48\x4a\x90\x53\x33\x1b\xc9\x47\xfc\xc4\x2e\x36\x8e\xf1\x2f\xca\x25\x18\x3b\xc1\x59\x1a\x1a\x0a\x95\x2f\xe7\x20\xbd\x36\x9f\x95\x80\x5f\xbd\xa9\x96\x87\xc1\xdf\xc1\xda\x2e\xb3\x07\x5f\x90\x43\xe0\xcf\xd4\xdc\x62\x8f\xa2\x00\x7c\xe6\x39\x03\x49\x56\x14\x67\x2b\x90\xe6\x42\x54\x06\x24\x68\x62\x98\x64\x47\x27\x01\xf0\xb4\x03\x76\x9c\xf9\x04\xfc\xbc\x83\xaf\xcc\x7d\x09\x49\x05\xa6\xb5\x49\xec\xa8\xd1\x6a\x0d\xb5\x96\xbd\xe0\x1b\xbd\x94\x52\xc8\x29\xba\x68\xa3\x33\x59\x09\xc4\x52\xa5\x97\x57\x6c\xac\x96\xb2\x48\xad\x51\x59\x4c\xd5\x0c\xc0\xb8\xef\x2c\x37\x62\x05\x6f\x96\x25\xa0\x43\x67\x0b\xad\xe6\xa2\x82\xa6\xaf\x40\xb9\x6d\xcc\x4c\x54\x49\x80\x31\x59\x2c\xab\x19\x01\xca\x6c\x47\x33\xeb\xc3\x03\x09\x3f\xad\xab\x87\x0f\x4b\xa8\xcc\x53\x29\xe6\x19\x0e\x7c\xae\xb3\x39\x38\x28\xbb\x20\x3f\xc4\x7e\x70\xfb\x2a\x61\x6d\x9c\x07\xc0\x4f\x4a\x29\xad\x71\x15\xb8\xbc\xd6\x2f\x1e\x59\xc8\x8e\x0e\xba\xc9\x95\xac\xcc\x00\xf8\x84\x44\xb6\x39\xa2\xa3\x56\x78\x33\xc8\x8a\x2d\xc5\x66\x8b\x60\x34\x87\xc4\x3e\x6b\x8b\x35\xe0\xc6\xa5\xb9\xf2\x0b\xec\xb7\x22\x69\x5b\x4b\x68\xcc\xc2\x2e\xf3\xd7\x7f\xbd\x85\xfb\x09\x92\x5a\x0d\xbe\xda\x98\x7a\xf0\xd5\x06\xea\x5f\x77\x46\xe4\x55\x65\x11\x36\x5a\x8f\x84\x22\xf1\xa4\x61\x77\xcb\xc6\x93\x91\xa7\xee\x4e\xc8\x42\xdd\x25\x0b\xd0\x13\xa5\xe7\x99\xcc\x21\x91\xea\x8e\xd0\x51\x09\x66\x60\xf8\x8e\x94\x1a\x73\x42\x2b\x1b\x99\xe3\xe3\x91\xe7\x93\xdc\x01\xbd\x34\x57\x23\x89\x1e\x6d\xaa\xb3\x79\x1c\xc3\x13\xde\x7e\x25\x20\x8b\x38\x96\x49\xa1\x24\x10\x8a\x1e\x0d\x64\x21\xe4\xd4\x43\xb9\xaf\xa4\x32\x99\x36\x08\x67\x5f\x48\xdb\x81\x23\x1a\x52\x4e\x89\x4c\x96\x8b\x02\x77\x9d\x2d\xd5\xe1\x47\x43\x9a\xb6\x43\x1e\x1e\x76\x28\xa9\x16\xa5\xc8\x81\x18\x76\x42\x6b\x8c\x56\x82\xb1\xf4\x53\x4a\x46\x47\x50\x56\x30\xf0\xc3\x42\xb5\x41\xbe\x01\xff\x84\x60\x1c\xff\x00\xf9\xb7\x0d\x59\x40\x09\x06\x9c\x6e\xd2\xd1\xb6\xa6\xf0\x4d\x5d\xd7\x2c\x84\x41\xfd\xf5\xa6\x9c\xf9\xf5\xf2\x9d\x16\x4b\xae\x21\x11\x1b\x44\x34\x99\x88\xd2\x80\x26\xc0\x9f\x40\x1c\x1f\x9f\x70\xce\x21\x11\xb2\x80\xf5\xeb\x09\x31\x94\x26\x37\x4a\x48\x07\x5a\xb3\xa9\x56\xcb\xc5\xeb\xa5\xd1\xaa\xf2\x6a\xa4\xec\x17\xdf\x68\x98\x67\xc2\x3a\x87\x21\xc3\x0d\x6a\x9c\xe5\xb7\xe8\x19\xea\x9a\xdd\x65\xc2\xa4\x84\xf2\x27\x64\x91\x34\xd6\xff\xf0\xd0\xbd\xf3\x5f\xdc\x33\xd1\x50\xa9\x72\x85\x4a\xd0\xf6\x25\x66\x06\x92\xe0\xd8\x4d\x07\x8f\x9e\xa3\xa6\x01\x14\x0d\x02\xce\x95\xe3\x82\xf3\xfa\x26\x31\xea\x16\x64\x2f\x22\x15\x04\x98\x60\x8a\x55\x2e\x2e\x75\x10\x47\x9c\xcb\x26\x8e\x19\x19\xbf\x90\x82\xab\x38\xde\x5c\xaa\xab\xb4\xaa\x1b\x1b\xc9\xb8\x26\x9a\x6c\x6a\x66\x92\xdc\xac\x29\xeb\x60\x29\xcb\x31\x54\x20\x26\xc9\x97\x1a\xf7\x27\x0e\x14\x3f\xd4\x7c\xa1\x24\x48\xc3\x32\x3a\x32\xc9\xb8\x54\xf9\xad\x85\xb2\x6f\xd5\xa9\x7f\x49\x26\x4a\x9f\x65\xf9\x8c\xd8\x2d\x8f\x3f\xd9\xe0\x92\x84\x8d\xb5\xc8\x22\xe9\xf1\x9d\x41\xa2\x1c\x4b\x20\x29\xc8\x09\x2e\xc2\xe1\xc0\xdd\xb7\xe1\x4d\x4d\xd3\xa6\xd5\x81\xe4\x49\x4e\x28\xcb\x2f\xf3\x44\x9c\x46\x22\x4a\xa3\x79\x74\x45\x4c\x32\x57\x4b\x69\x08\x62\xc8\x64\x3e\x53\x1a\xdf\xda\x15\x27\x5a\x29\x83\x3b\x02\xd9\x20\xb7\x1b\x7c\x3c\x6f\xe7\x0b\xe8\xc0\xf0\x3d\xb7\xa6\x13\x06\x50\x3e\x10\x88\xe3\xc8\x37\x07\x81\xa4\x95\x6d\x8d\xa1\x27\x0a\xa2\x11\x35\xf0\x27\x1b\x81\x52\x99\x81\x64\x27\xcc\x24\x36\xfc\x63\x40\x6b\xe6\xbb\xf2\xcc\xe4\x33\xf6\x98\x99\x04\xb4\x56\x1a\xfb\xec\xba\x1d\xdb\x8f\x38\x37\xad\x6f\xf0\xb1\x1d\x69\x9b\xd8\x90\xb2\xa3\x61\x8d\x16\xeb\xc4\x1f\x0e\xc3\x49\xc3\x31\x3b\x8b\x60\x47\xc3\x50\x3d\x36\x97\x4d\xd7\x55\x0a\x75\xdd\x6d\xab\xe3\x76\x4f\x2b\xa0\x32\x5a\xdd\x73\xe1\xdc\xd1\x44\x68\x20\x51\xd3\x18\x35\x3e\xaa\x02\xe3\xfb\xaf\x27\x3a\x9b\xda\x4d\xa5\x20\x47\x27\x47\x9c\x7b\x3f\xd6\x76\x58\x01\x37\x6d\x95\xc9\x0c\xa0\x13\x68\x27\xbe\xee\xc5\x91\x70\xc4\xe1\xd4\x70\x6e\x52\x40\xf2\x1e\x1e\x60\x5f\x3c\xff\xf0\xb0\x4f\x38\x1d\xce\x9b\xd0\xa0\xc0\x06\x7a\x76\xf6\x59\x26\x8b\x12\x74\x15\xc7\xfd\xef\x4b\xb8\x4a\x2a\xeb\x4d\x5d\xd8\xdf\x86\xe7\x82\x0f\x47\xe2\x7b\xe9\xbd\x9e\xc0\x20\x6c\xe3\xd2\x3f\x79\x29\xae\x46\x3a\xb9\xbe\x46\xcf\xe1\xb6\xfa\xe0\xcb\x25\x75\x36\xea\xc5\xa9\x30\x7c\xee\xf5\x9e\xd0\x80\x05\xf7\x5d\x78\x15\x70\xa9\xeb\xbf\xf5\x5e\xb2\x5d\xf0\xbe\x84\x85\x41\x72\x8d\x41\x02\x37\xee\x89\x36\xb7\xb0\x5b\x05\x37\x0c\xac\x65\x70\x67\x20\x0f\x0f\xc0\xd0\xc1\x2a\x0d\xdc\xb8\xe7\xc3\x03\x34\xb6\x83\x5f\xdd\xd4\xeb\x9e\x6b\xda\x66\x9a\x0f\x05\xc2\x36\x7e\x79\x45\xbb\xf4\xc1\x06\x45\x86\xb2\x4d\x8e\x5b\x73\x19\x26\x1b\x88\x13\xb8\x0c\x1c\xf7\xe8\xdf\xc1\xee\x96\x6e\x63\x03\xdc\xd8\x02\x2e\xbd\x6a\x35\xf4\x1a\x0d\xdc\xba\x34\xcc\x01\xdd\xc6\x87\x6b\xbf\x46\xbb\x6e\x43\xae\xb6\x05\x65\xf1\x3a\x6c\x1c\xdb\x40\xdd\xf1\x8e\xf6\xbb\x94\xdc\xdb\x9c\x4d\x0c\x68\xdf\xb3\x83\x3e\x0c\x4f\xef\xba\x2c\x2c\x90\x25\x43\x7f\xce\x04\x06\x2e\x5e\xb5\x94\xcb\xfc\x1c\x50\x21\x26\x13\xd0\x15\x81\x4b\x75\xc5\xcc\xa5\xba\xa2\x71\x4c\xe4\xa5\xba\xe2\x02\xe3\x80\x91\x40\xbf\x15\x18\x90\xf7\xe8\x94\xb5\xd6\xa6\x01\xbd\xe0\x12\xb5\x21\xb4\x35\xdf\x8d\x1a\xe1\xd5\xde\x25\x93\xfb\xa0\xbc\xc9\xfa\xe9\x9c\xf1\x5b\x90\x88\x6d\xf2\x59\x26\xa7\x50\xa4\x92\x35\x1e\x28\x0d\x49\x5c\x68\x58\x09\xb5\xac\x52\x53\x6f\xa3\x4b\x16\x7b\xa7\x73\xd8\x5d\xf8\xf3\x45\xe8\x69\xc0\xf0\xd7\xc8\x70\x17\x8e\x60\xd6\xdb\x58\x2a\xc5\x00\x56\x4c\x0c\xa1\x24\x80\x7d\xd3\xe4\x74\xbd\xc5\x5d\x6e\xad\x35\xd8\x6d\x80\x99\x87\x07\x6b\x5c\x35\x8a\xec\x2d\xdf\x34\x7e\x30\x1d\xdb\x53\x95\x7b\x86\x24\xa4\x37\x4c\xc9\x74\xcd\x2a\x30\xe9\x2b\xd6\x09\x22\x15\x0c\x35\x35\xbd\x63\xd7\x76\xd3\x4a\xdf\x30\x2f\xe8\xf4\xba\xb6\xf9\xc9\x33\xbe\x99\x80\xc9\x67\x6f\x91\x44\xd2\xc4\xa2\xf8\x77\x22\xa6\x29\xb0\x71\x56\x41\x6a\x6a\xa7\x49\x53\x30\xa4\xb5\x2c\x3b\x8a\xfc\x8a\xa1\xf4\x37\x42\x56\x06\xad\xeb\x9b\x5f\x1f\x81\xcc\x55\x01\xef\xdf\x9c\x3f\xf3\x3b\x22\x81\xe4\x55\x36\x07\xca\x36\x73\x30\x33\x55\xa4\xd1\x14\x4c\x54\xd3\x76\xeb\xfa\x76\xf8\xad\x8d\x9f\x90\xcb\xcb\xea\x34\xc8\x97\xea\x14\x92\x9b\x0a\x6d\xb5\x83\xde\x78\xf7\x4f\x5a\x40\xf0\x51\x67\x4d\x6b\xea\x76\x3a\x62\xf8\x13\x4b\x8b\x2a\xc1\xed\x78\xcd\x32\x98\xa1\x35\xad\x99\x8b\x87\x3d\xb9\x2d\x22\xd6\x10\x6e\x1c\xe1\xb2\x47\x38\xe6\x38\xed\xb9\xcd\xee\x2a\x5c\x10\xc1\x5a\xbe\xc8\xfa\x9b\x6a\xb9\x00\xbd\x12\x95\xd2\xfb\x39\x63\xb6\x39\xb3\x50\x55\x8f\x35\xb8\xdb\x3e\x1e\x0e\x03\xf6\xf4\xe6\x0f\x25\x37\xf2\xd4\xde\x65\x5a\x12\x0f\xce\xfc\x0b\x66\xd5\x1d\x77\x60\x87\x3b\x6e\x29\x36\x2c\x40\xee\xa8\xc5\x97\x33\xe7\xcb\x78\xf3\x71\x9d\xf9\x4c\xce\x7c\xfb\x8f\xe6\x8c\x51\x6a\x5a\x82\x0f\xe3\x9d\xfb\xc8\x4b\x91\xdf\x42\x11\xd1\x3a\x08\x98\xdf\x21\xff\x4a\xc8\xf4\x3b\x31\x07\xb5\x34\xce\x9b\x19\xf7\x11\x38\x84\xb3\x60\x7b\x6b\x4c\x6b\x93\x13\xba\x21\x12\xd3\xe6\x42\xac\x22\x4a\x93\xbc\xcc\xaa\x0a\x57\xc1\xa3\x5c\xe8\xbc\x84\xc1\x78\x7a\x3c\xd5\x70\x3f\x18\x97\x42\xde\x0a\x39\x1d\x54\x2b\x28\x0d\x1c\x9f\x3c\x3e\xc9\x26\x37\xb3\xa8\x66\x73\x87\x38\x23\x92\xb9\x53\x2f\x7b\x66\x00\x71\x9c\x13\xd9\xdb\xd5\x9e\xff\x9e\x15\x80\xfc\x1d\x33\x9f\xff\x9f\x67\xd6\x50\xfc\x8e\x79\x2f\x0e\xcc\x2b\xf9\x92\x44\xf9\x0c\x72\xe4\x68\x92\x24\x98\xbc\x7d\x36\xd2\x5f\x42\xa4\x41\xde\xd4\x24\x97\x68\x50\xf5\x0e\x8d\x19\x52\x28\xa4\x04\xfd\xe2\xdd\xcf\x17\xfc\xeb\xef\x0b\xb1\x1a\x58\x8a\x79\x84\x63\xb6\x89\x7c\xf2\xfd\x37\x85\x58\x3d\xf9\x37\x39\xd8\xf3\x0f\x07\x7c\xcd\x66\x44\x32\xa7\x92\x11\x13\x98\xe7\x07\x0c\x94\x22\x87\xb1\xf9\x52\x91\xb1\x69\x0f\x65\x48\xf4\x57\x1f\x27\xda\xfa\xd8\x2f\xa3\x7a\x51\x66\xf7\x5f\x48\x75\xa6\xcd\x3f\x98\xec\x9f\x03\xb2\x99\x60\x9a\x29\x8e\x29\x20\x7a\xc7\xe4\x6f\x4a\xdf\xfe\x28\x74\x5f\xaf\xbc\x3a\x33\x81\x2a\xf6\xdf\xff\x95\x99\x74\x10\x51\xa6\xf9\x92\x28\xca\xe6\x38\xd1\x44\x49\x73\xec\x8e\xc9\x58\x24\x4c\x56\x8a\x3c\x0a\xbb\x26\xd9\x5c\x94\xf7\x11\x8b\xe6\x4a\xaa\x6a\x91\xe5\xe0\xbb\x73\x55\x2a\x1d\xb1\x68\xaa\xb3\x7b\xdf\x36\x56\xba\x00\x7d\x6c\xd4\x22\x62\x51\x29\xa6\x33\x83\xbd\x83\x42\x19\x03\xc5\xe0\x64\xb1\xf6\x80\x8b\xac\xc0\x0c\xaf\x81\x1c\x26\x7f\x82\xf9\x76\x57\x09\x13\x13\xb1\xe8\xbb\xdd\x41\x1a\x11\xf7\xbb\xe6\x99\x9e\x0a\xb9\x0f\x5d\xd3\xd3\x60\x3b\xde\x1d\xe3\xb1\xb9\xae\x1d\xe9\xb0\x8a\x08\x26\xf1\xa1\x99\xa4\x35\x5b\xf8\x9c\xc4\x31\x3e\x8e\xd5\x11\xe7\x64\x57\x12\x18\xc8\xea\xa4\xc8\x4c\xc6\xd5\x47\xec\xf7\xfd\x8e\x4c\x59\xc9\x0a\x36\xe3\x26\x91\xe8\xff\xa7\xbc\x49\x5b\x77\x24\x6b\xb4\x15\x2c\xbe\xcd\x1a\xa1\xce\x28\x53\x28\xe9\xfd\x2a\x1b\x51\x56\x5a\xf0\x22\xa2\xac\xe0\x4b\x32\x45\x46\x94\x2c\xba\x53\xba\x38\x1e\x6b\xc8\x6e\x23\x16\xd9\xe7\x71\x56\x96\x9f\xe0\x85\xc0\x87\x72\x5f\xa5\x7b\x14\xac\xec\x18\x44\x20\x71\xb9\x1a\x66\x58\x8e\x35\x34\x8e\x67\xc8\x2c\x4f\x5c\xc0\xa2\x19\x65\xfb\x07\x4c\x71\x40\xcb\x04\x1c\x51\xb8\x11\xd3\x8f\x30\xf5\xe9\xd6\x8d\x02\x13\xdc\x34\xc8\x13\x90\x46\x0b\xa8\x48\x2b\xaf\x33\xb9\x12\x5a\x49\x0c\x81\x91\x8b\x97\x57\x4c\xf1\xe1\x48\x7d\x2f\x7c\x48\xad\x30\xf9\xd5\x98\x93\xbc\x27\x43\xf6\x03\x31\x4c\x30\xd5\x5e\x46\x05\x02\xc9\xc6\xf6\xb0\xd9\x4f\x0b\x7c\x38\x82\xef\xb5\x47\x03\x0e\x0d\x5c\x25\x39\xa1\xa3\x9e\xa7\x00\xb9\xfa\xa4\x97\x18\xf5\x53\x73\xdd\x4b\xcd\xf5\xa5\xb8\x4a\x50\xa9\x5d\xcc\xee\x65\x60\x8f\x6a\x76\x98\xba\xf9\x5c\x6e\x74\x29\xdb\x1e\x86\x34\xc7\xc7\x15\xf7\x0c\x19\x21\x8b\x4e\xf1\x4f\x82\xd3\x57\x34\x25\x2d\xd3\x2a\xca\x6c\x47\x4e\x9a\x97\x76\xad\xb4\xb6\x29\x8c\xea\x08\x6a\xd9\x6d\x4f\xc3\x46\xbe\x9d\xfb\xe9\xeb\x2d\x7f\x59\x12\x7b\xac\x14\x88\xff\xc7\x2d\x9b\xda\x91\x54\xb5\xc8\x64\x44\x19\x11\xed\xb6\x60\x60\x6d\x9e\x29\x69\x40\x1a\x1e\x15\xea\x4e\x96\x2a\x2b\x06\xa5\x9a\x0e\x26\x02\x5d\xa3\x48\x66\x1a\x26\x5c\x73\x93\x60\x14\xfa\x28\x0a\x83\xec\xe8\x51\xcb\x3e\x94\xe7\xa3\xe8\x9b\x52\x4d\x77\x25\xd8\x58\x4f\x68\x21\x88\xab\xa7\xee\x1a\xd5\xfd\x4b\x66\x41\x8b\xf0\x8b\xfb\x88\x45\xfc\x6d\xbf\x9b\x99\xb3\x05\x5b\xb1\x31\xbb\x66\x37\xec\x9e\xdd\xb2\x35\x7b\xc5\xee\xd8\x6b\xf6\x86\xbd\xe5\xf6\x78\xb9\x9d\xf2\x0d\xd8\x1d\xef\x34\xfa\x9f\xff\xf8\xcf\x28\xdd\x6e\x66\xcf\x38\xd9\x6e\x6b\x22\xd0\x6f\x4e\xe0\x5f\x68\x62\xd4\x73\xb1\x86\x82\x3c\xa6\xec\x5d\x00\xfa\xd6\xa8\xc5\x21\xb8\xb3\x1d\x97\x1a\xc7\x3f\x93\xa1\x35\x85\x2e\x88\xb4\xae\xd8\x87\xc7\x36\x59\x3d\xe7\xfb\x94\x39\x8e\x9f\xda\xb1\xec\xa2\xeb\xbe\x50\xd3\xe7\xa2\x84\x38\xfe\xd1\xa1\x3d\xb4\x77\xb6\xaf\x67\x71\x7c\xe6\x74\x78\xd7\xcb\x46\x94\x9d\xc7\xf1\xb9\xed\xde\xe3\x84\xbd\xf3\x75\x88\x8a\x40\x0b\xe7\x08\xac\x1d\xcf\xec\xf6\xbc\xe0\x4b\xf2\x96\xb2\xd5\x7e\x4f\x1e\x51\x36\x0e\x46\x5f\x07\xa3\x07\x4d\x9c\x6f\xb1\xdc\xf0\x25\x79\x46\xd9\x3d\xf6\x57\x11\x65\xb7\x07\xd1\xad\x03\x74\xaf\x2c\x38\x86\x80\x21\xae\x3b\xbe\x24\xef\x28\x7b\xed\x71\xbd\x39\x88\xeb\x22\x8e\x2f\x2c\x0b\x44\xe8\xdf\xf4\x72\x27\x0a\xea\x45\x4e\xcf\xb7\x22\xa7\x79\x26\xa4\x3d\xab\xde\x19\x55\x84\x60\x2e\x11\xd0\x2e\xe8\xd8\x01\x1d\x87\xa0\xe3\x72\x09\x07\x21\xd7\x21\xa4\xd2\x99\x9c\x1e\x86\x2d\x43\xd8\x89\x52\x66\x17\xe6\xe0\xb6\x69\xd5\x67\x8e\xef\xf6\xd8\x32\x43\xc7\x85\xdd\x56\x6f\xdc\x18\x96\x11\xc5\x9a\x97\xd2\x0f\xc6\x7d\x95\x55\x64\xce\x0a\x7c\x2c\xdc\x63\xe5\x1a\xc7\xee\x71\xcd\xc6\xf8\xb8\x71\x8f\x7b\xf7\xb8\x75\x7d\x6b\xf7\x78\xc5\xd6\xf8\xb8\x73\x8f\xd7\xee\xf1\x06\xfb\xac\xd0\x30\x12\xe8\xef\x1e\xdb\xe6\x77\x7a\x76\x7a\x96\xb8\xbe\x94\x90\x33\xee\x8c\x91\x5a\x71\x07\x74\xd1\xf4\x2c\x8e\xc9\x99\xbb\xcb\x38\xe3\x8e\xd6\x7d\x26\x79\x7a\x7e\x7a\xde\xe1\x3b\xe7\x4f\x03\x7c\xc8\x0f\x15\x14\x0f\xe0\x86\x9b\x9e\xc7\x31\x39\x77\x78\xcf\xb9\x3f\xfa\xf5\x61\xd8\x5b\xf4\x9c\x5f\xe6\xb5\xa8\xbd\xaf\xb1\xd1\xc4\xdb\x10\xd5\x33\x44\xf5\xf9\x1e\x0d\xd1\xdc\x38\x34\xcf\x42\x34\xef\x10\xcd\x67\x7a\x3b\xc4\x71\xe7\x70\xbc\x0b\xb8\xd5\x78\xa8\xd3\x8b\xd3\x8b\x8e\x53\x17\xfc\xc7\x80\x53\x81\xe4\x68\x7a\x11\xc7\xe4\xc2\x71\xe8\x82\x37\xd2\xec\x6d\x96\x56\x05\x0b\xd2\xcf\x32\x9e\x53\x66\xfb\x75\xa3\x8a\xf6\x56\x1c\xe2\x98\xe4\x98\x26\xe4\xa4\xa4\x8d\x8e\x14\xa4\xb7\xa9\xfc\xe0\x2b\x12\x5c\x28\xb0\x5d\x4d\x04\x5d\xbd\x8e\x0d\xf7\x6c\x7d\xca\xe5\xf0\x8a\x09\x17\xcc\xb9\xef\x13\xfc\x86\x2c\x9f\x5d\x37\x8d\xfe\xd3\x1e\x8f\x73\xc9\x44\x37\xe1\x07\x24\xe5\xd6\x5d\x2c\xb4\x47\xb7\xc1\xc1\xb0\x9b\x3f\x75\x8f\x9a\xf5\xce\x81\xdc\xe9\x51\x94\x24\x11\x2b\xc0\x64\xa2\x84\x22\x3d\x3a\xa9\x51\x60\xc8\x76\x8f\x4d\x48\xa3\x55\x7b\xa7\xdf\x9e\xef\xb7\x97\x43\x97\xef\xae\x58\x5b\x4d\x30\x05\xd3\x54\x99\xfc\x70\x7f\x5e\x90\xa8\xef\x08\x9a\xb4\x8a\x3e\x3c\xec\x1c\xfe\x07\xa5\x09\x90\x88\x82\xef\x1f\xc9\xa0\x17\x97\x7c\x9d\x58\x87\x97\xf4\x61\x37\x36\x0f\x4b\x6d\x57\x9d\x38\xef\xb5\x1f\x44\x43\x51\x27\xe8\x08\xf7\x77\x63\x4f\x9d\xf8\xf3\x95\x6d\x98\x71\x96\xdf\x4e\xb5\x5a\xca\xe2\xd8\x81\x97\x62\x0e\xcd\xa4\xee\x64\xe4\x93\x23\xdc\xfc\xee\x04\xe9\x33\xd0\x37\xc9\x63\x9d\xb8\xf3\x97\xed\x11\x85\xa8\x30\x69\x4f\x85\x2c\x85\x84\x63\xbb\x5d\x8c\xee\x44\x61\x66\xe9\x09\xcc\x47\x33\x40\x04\xf6\xb5\x49\x4a\x75\x56\x88\x65\x95\xda\xe4\x70\xd4\xe4\x7e\x63\x65\x8c\x9a\xa7\x8f\x17\xeb\x3a\xd1\xcb\x1d\xa2\x1b\x74\xc3\xe1\x1f\x46\x4d\xee\x99\x7e\xb7\x58\x6f\x21\xfc\x13\x0e\x6e\x72\xfd\x43\x8b\x9c\x94\xb0\x1e\xe1\x9f\xe3\x42\x68\xb0\xba\x90\x6a\x75\x37\xca\x4a\x31\x95\xc7\xc2\xc0\xbc\x4a\x73\x90\x06\xf4\x21\x54\xe9\x4c\xad\x40\x07\x7c\x4a\xef\x66\xc2\x40\x35\x57\xb7\x50\x27\xed\x86\xf9\x85\x0b\xa8\x93\x66\x9f\xdb\x91\x47\x43\xa2\x65\x62\x91\xe9\x5b\x97\xc7\x67\xd5\xcc\xe5\xf1\xa3\x5e\x2e\x8e\x6c\xae\x13\x90\xab\x8f\x30\x10\x35\xf9\xd8\xd2\x9b\x62\x2e\x3e\x0a\x4e\x17\xd2\xf6\x6c\x01\xf5\xd3\x9d\x2c\x6e\x63\x6a\xeb\x21\xd2\x2d\x53\x69\x8f\x22\x1f\x57\x03\xd4\x04\xcc\x90\xe4\x44\x48\x61\xa0\x4e\xf2\x4c\xef\xea\xa5\x5a\x1f\x57\xb3\xac\x50\x77\xe9\x70\xf0\xed\x62\x3d\xf8\x6e\xb1\x1e\x0c\x07\x7a\x3a\xce\xc8\x90\x0d\x9a\xff\xc9\x63\x3a\xea\x8a\x4e\xd2\x61\xf2\xc7\x6a\x57\xee\x3d\x22\xa2\x8b\x65\x2e\x8a\x6c\xf0\x36\x93\xd5\xe0\xbd\x14\xb9\x2a\x20\x62\x03\xdf\xfc\x93\xce\xa4\x6d\xa8\x32\x59\x1d\x57\xa0\xc5\xa4\xd1\x42\x8b\x69\xaf\x3a\xcf\xb3\xf5\xb1\x63\xe1\x1f\x87\x43\x54\x33\xef\xb9\xb6\x69\xea\x00\x91\xd7\x83\x23\x31\x5f\x28\x6d\x32\x69\xf6\xb2\xc0\xeb\x52\xc8\x08\x64\xc2\xc9\x3f\x1f\xe0\x44\x8d\xfb\x91\xc9\x84\xdc\x55\x94\xc0\x2e\xec\xf8\xda\x96\x58\xed\xc2\x7d\xb9\x25\x8c\x6e\x96\x95\x11\x93\xfb\xe3\xdc\xf9\xbf\xce\x40\xb2\xf9\x8e\x37\xf0\xcb\x78\xec\x85\x52\x89\xdf\xc0\x5a\xbf\xfd\xba\x73\xde\x60\xac\xca\xc2\x29\xa2\x15\xed\x44\xe9\x79\xba\x5c\x2c\x40\xe7\x59\x05\x75\x82\x0b\xdc\xaf\xc1\x43\xef\x4f\x86\x5e\x09\xac\x87\x4e\x2b\x55\x8a\xc2\x37\xb5\xce\x67\x30\x1c\xe0\xdf\xc7\x9d\xef\x71\x3e\xcd\x4e\xea\x62\x9a\xc1\xa1\xf7\x7f\x1a\xda\x7f\x87\xf4\xc1\xea\x8b\x75\x57\x18\xa4\xef\x5f\xec\xe3\xce\xfd\xe1\xeb\x5e\x4c\x6e\x5d\xe9\xb8\xcc\x30\xc8\x46\x2a\x06\xc8\xba\x60\x82\xa0\x6e\xed\x80\xc5\x6d\xfe\x34\xfc\xc3\x46\x2d\xb2\x5c\x98\xfb\x74\x58\xd7\x5f\xb3\x8a\x00\xeb\x95\xda\xd1\xda\x96\xc0\xb8\xbb\xf7\xe6\x22\xd7\xdf\x34\xf3\xcb\xab\x9d\x82\x89\x7e\x89\xed\x76\xb6\x7a\x38\x57\xed\x65\xc6\xec\x59\xf7\xfd\x4c\xcd\xe7\x99\x2c\xd8\xbb\xae\xe9\xa9\x9e\x56\x4d\xa1\xd4\x20\x0a\x72\xc9\x9f\x83\x2a\x18\x0c\x18\xec\x1d\x8f\xaf\x56\x3b\xeb\x6e\xfb\x4e\x9f\xa7\xe7\x36\xd1\x7c\xcf\x7f\x26\x86\xb2\xa7\xdc\x9d\xe1\x8d\x82\xf3\x87\x7e\x56\x1a\x06\x4c\x87\xa7\xb8\x08\xa6\xf8\x25\xfd\xca\x4e\xf1\x81\xff\x80\x53\xbc\xe0\x1f\xb6\xa6\xf8\xed\xd0\x14\x3f\xed\xc9\x88\x5f\x72\xd3\x7a\x8e\x38\x76\x47\x01\x9f\x4e\x76\x75\xf7\xaa\x82\x0c\xb1\x74\xd9\x69\xb1\x37\xc5\x7d\x6a\x43\xd2\x45\xbf\x2f\xc2\x5c\xb6\xc5\x35\xee\x5e\x5f\x58\xe8\xeb\xbd\x98\x6e\x10\x6c\x11\x35\x09\xec\x4e\xf7\xad\xed\xd6\x60\xd3\x57\x9b\xe8\xda\xcc\x75\x3b\x59\xdd\x19\xf7\x32\x8e\x5f\xba\x14\xbd\x77\x52\x9f\xcd\x61\x4f\x82\xaa\xdb\x08\xf9\x47\xca\x74\x38\xc0\xf9\xb8\x3d\x43\x6e\xda\x21\xbf\x51\x36\x23\xb7\xed\xe7\x4f\x94\xcd\xf1\xd3\x59\x58\xc4\xa2\x21\xb2\xa5\x97\xa2\xaa\x62\xe7\x0a\xa2\x9f\x46\xb7\x4e\x78\x07\x2c\xcc\x9e\xdf\xf0\x08\x7d\xfe\x20\x7a\x44\x3a\xb1\x9f\x46\xfe\x2d\x4a\xa3\x88\x3e\x8a\x3e\x3b\x69\x0d\xce\x7a\xb5\x3b\xeb\x55\x2e\x27\xd5\x28\xee\x39\xd1\x4d\x36\x8b\x89\xa9\x70\x89\xa9\x70\x89\xe9\x0a\x25\x3c\x27\xe3\x16\xe0\x1a\x9b\x30\x49\x5d\xb9\x24\x75\xe5\x92\xd4\x95\x4b\x52\x6f\x5d\x92\x7a\xeb\x92\xd4\x5b\x97\xa4\xae\x1a\x91\xcd\xc9\xaa\xcb\x4e\xed\xef\x41\xcc\x76\xe6\x27\x42\x1f\x80\xf9\x54\xd9\xa6\x76\xef\x11\x02\xcd\x55\xd8\x44\xeb\xa9\xcb\x90\x88\xb3\x5c\xd1\x24\x52\x01\x31\x94\x7d\xc0\x11\x68\x7d\x6e\xc4\x8b\x66\x84\x33\x44\x3f\x22\xa0\xae\xcb\xd9\xda\xf4\x51\x6c\xf9\x20\xc4\xb3\x6e\xd3\x44\xb3\x9d\x26\x8a\xfd\xfe\xa9\x9f\x18\x8a\x4e\xa0\x2f\x4f\x5f\x26\x8e\x15\x29\x21\x2f\xf9\xdf\x82\x75\x05\xec\xa2\xe9\xcb\x38\x26\x2f\xdd\xf2\x5f\xb6\xc9\x78\xe7\x0d\xde\xe0\xec\x81\xca\x88\xcf\x57\x19\x5b\xaa\x13\x2a\xde\x76\xb2\xf9\xb4\xc9\x34\x7b\x76\xf4\xa2\x69\xec\x59\xca\xb4\x6f\x29\x56\xe4\x2e\xd9\x6c\xaa\xc8\x76\x4a\x68\xfa\x55\x4b\xae\xdc\xca\xd6\x59\xf6\x8a\x1a\x5d\x2d\xb4\x6d\x0f\x2f\xc2\x79\x05\x3e\x99\x27\x5d\x67\x78\x79\x9f\x4c\x84\xcc\xca\xf2\xbe\xe9\x26\xf6\x16\xfe\x31\xfc\x91\xd6\xa3\x5d\xd8\xb6\xee\xd5\x10\xea\xae\xf4\x9b\xea\xb7\x8f\x97\xfb\x28\x5b\xc9\x14\x52\xb6\xa7\xfc\xa7\xa6\x36\x45\x35\x99\x9e\x42\x5b\x9d\xd4\xd5\xec\xa0\xb8\x5d\x93\x2b\x0f\xf5\x90\x0c\xda\x2a\xd1\xd7\x5b\xfb\x2e\xa5\xb5\x26\x1f\xba\x5f\xac\xb0\xb7\x94\xf5\x1b\x9e\xb9\xdf\x10\xbc\xe0\x1f\xec\xf3\x37\xbe\x51\xf2\xdc\xc0\xfc\x2d\x94\x90\x1b\xe8\x7e\x05\x60\x8b\x2d\xaa\xa6\x35\xc5\xb5\x06\xa5\x08\x3f\x7d\xf4\x9a\x7d\xb1\x73\xea\x5e\xaa\xac\x68\x2e\xc0\xf7\x5c\x91\x2e\x52\x71\xf0\x70\xfb\x65\x58\xab\xe7\x8b\x88\xbc\x81\x35\xb5\x44\xf6\x40\xbd\xcb\xfc\x4d\xe2\x97\xdd\x9d\x16\xd5\x4c\x70\x09\x77\x83\x17\x64\x83\xda\x95\xba\xb8\x85\xd9\xc2\xc0\xb4\x29\x1b\x64\x68\x89\xa9\xac\x83\xa3\x0d\x25\xbb\x5a\x8c\xee\x77\x4c\xf6\x77\x1a\x5b\x6c\xf3\x13\xa1\x48\x2d\x2b\x44\x5f\x92\x2d\xd9\xa2\x15\xa7\xa3\xdd\xd6\x19\x6f\x0c\x97\xcd\x2f\xf0\x36\xf5\xc8\x1f\x33\x55\xf6\x0e\xcd\xbd\xb7\x94\xa0\xc6\x20\xbd\xb6\x0f\x5f\x9a\x0b\x05\x7b\xc9\xe6\x09\xef\xee\x1d\xaa\xe6\x22\xae\x61\x0e\xdf\xc7\x1c\xf4\x3d\xae\xfc\xd0\x5b\xb9\xf0\x67\x22\xa4\x7f\xfb\xf2\xe7\x7d\x97\x6f\xcd\x44\x9f\xb8\x5f\x43\x41\xfe\x75\xdf\xfd\xda\xe7\x5d\xa9\xf1\x22\xe0\xe1\x67\x5c\x97\xd9\x38\x28\x50\x31\x62\x98\x6a\x0a\xcd\x9b\xf5\x3e\x3c\x98\xe6\x8e\xa6\xe3\x09\xdd\x08\xae\x3c\x40\x7b\x45\x56\xf1\xe1\xa8\xea\x68\xaa\x82\x2b\xb2\x8c\xff\x95\x28\x26\x58\x45\x47\xfa\xb2\xba\x3a\xc5\x3f\xd6\x75\x67\xf6\x8a\xac\x72\x74\x67\x94\xd9\x0e\x77\x45\x56\xd9\x2b\xb2\xf0\x00\x54\xfa\xab\xb2\xaa\x23\xa6\x72\xc4\x54\x9f\xba\x2a\x73\x57\x64\x6c\x8f\xe9\xfc\x25\x2c\x42\xde\xe4\xa9\x60\xf3\x54\x38\x63\x4b\x45\x00\xf7\xd7\xcf\x3f\xea\xf3\xba\x78\x29\xf7\x9c\xec\xb9\xce\x7d\x67\x7b\x7f\xff\xd8\xd9\x5e\x63\xd4\xee\xd7\x4b\xed\x39\x5e\xeb\x79\x6c\x9d\xd7\xa1\x83\xbc\xc3\xe9\x87\xaf\x10\x0d\xaa\x62\xbb\xd6\xa0\xf2\xf5\xb3\x33\x15\xfb\xdb\x1b\x85\x2e\xa8\x29\x28\x4b\x81\xe5\x66\x9d\x9a\xd6\xb3\xdb\x93\xc8\xa6\xbc\x3d\xfd\x89\xe1\xce\x91\xfe\x99\xd9\x1a\xb0\xf4\x2f\xcd\xcf\x27\xa3\x86\xd8\x88\xd9\x4a\xb0\x34\xc2\x41\x51\xfb\xd3\xc4\x15\x09\x0d\x4a\x35\xae\xc4\x2a\x3f\x53\xcd\xef\x08\x42\x67\xd2\x06\x74\xbe\x13\xdb\x55\xb3\x3d\x70\xdb\x6e\x9d\x0d\xb7\x5b\x58\xa8\x73\x3d\xd7\x83\x11\xb3\x59\x73\xc3\xda\xe5\x09\xfb\x6b\x01\x81\x01\x44\xb0\x22\x1a\xc7\x2b\x82\x06\xfc\xf0\xe0\x27\x44\x24\x5d\xe1\xae\xea\x7e\x86\xb1\x1d\x30\xa8\xf6\x67\x10\x80\x49\x87\x15\xeb\xef\x0a\x01\x7c\x31\xa1\xd5\x19\xe8\x55\x0e\xb6\x45\x81\x10\x16\x4c\x56\xbf\x1e\xac\x17\xdd\xaa\x09\x35\xfc\x89\xff\x2d\x8c\xdd\xee\x93\x79\xb6\xc0\xc6\xfd\x78\x0f\x94\x1b\x1e\xae\x4e\x6d\x67\xa3\xfb\x8b\x50\xbd\x41\x40\x50\x7c\x1a\x46\x3e\x5b\x45\xa8\xb4\xfe\xff\x17\x98\xb8\xee\xad\x42\xf5\xfd\x35\xea\xbb\xe5\xe9\x7b\x2a\xd3\xbb\xa2\x74\x4d\xfe\xbe\x1d\xdf\x84\x0d\xbf\xd1\x11\xee\xf3\x7f\x27\x1b\xb7\xa6\xf4\xf0\x81\x7e\xa7\x18\x11\x75\x7b\xff\xa6\xae\x69\x7d\x45\x47\xff\x1b\x00\x00\xff\xff\x9a\x7a\x62\x6b\x4c\x3f\x00\x00")

func mainJsBytes() ([]byte, error) {
	return bindataRead(
		_mainJs,
		"main.js",
	)
}

func mainJs() (*asset, error) {
	bytes, err := mainJsBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "main.js", size: 16204, mode: os.FileMode(420), modTime: time.Unix(1550941822, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"index.html": indexHtml,
	"main.js": mainJs,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"index.html": &bintree{indexHtml, map[string]*bintree{}},
	"main.js": &bintree{mainJs, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

