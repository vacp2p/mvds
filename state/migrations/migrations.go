package migrations

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

var __1565341329_initial_schema_down_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x72\x09\xf2\x0f\x50\x08\x71\x74\xf2\x71\x55\xc8\x2d\x4b\x29\x8e\x2f\x2e\x49\x2c\x49\x2d\xb6\xe6\x02\x04\x00\x00\xff\xff\xb7\x43\xc1\xc1\x18\x00\x00\x00")

func _1565341329_initial_schema_down_sql() ([]byte, error) {
	return bindata_read(
		__1565341329_initial_schema_down_sql,
		"1565341329_initial_schema.down.sql",
	)
}

var __1565341329_initial_schema_up_sql = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x8f\xc1\x8a\x83\x30\x14\x45\xf7\xf9\x8a\xbb\x54\xf0\x0f\x5c\xe9\x4c\x18\x64\xd2\x58\x42\x0a\x75\x15\xc4\x3c\xac\x0b\x35\x98\x58\xda\xbf\x2f\xb4\x15\xa5\x60\xb7\xf7\x1c\x1e\xef\xfc\x28\x9e\x69\x0e\x9d\xe5\x82\xa3\xbf\x5a\x6f\x7c\xa8\x03\x79\x44\x0c\x00\xc2\xdd\x11\x0a\xa9\xf9\x1f\x57\x90\xa5\x86\x3c\x09\x91\x3c\x91\xa7\xc1\x9a\x66\x9c\x87\xf0\x4d\x20\x37\x36\x97\x1d\xa1\x9d\xc6\xd9\x99\xce\x22\x17\x65\xfe\x9a\x1c\xd1\xb4\x2c\x1f\x76\x4f\xde\xd7\x2d\xed\xd0\xa3\x2a\x0e\x99\xaa\xf0\xcf\x2b\x44\xab\x9a\x2c\x17\x63\x16\xa7\x8c\xbd\x6b\x0b\xf9\xcb\xcf\xe8\xec\xcd\x6c\x7e\x2c\xe5\xb6\x3f\x5a\x49\x9c\xb2\x47\x00\x00\x00\xff\xff\x5e\xe5\x72\x74\x26\x01\x00\x00")

func _1565341329_initial_schema_up_sql() ([]byte, error) {
	return bindata_read(
		__1565341329_initial_schema_up_sql,
		"1565341329_initial_schema.up.sql",
	)
}

var _doc_go = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x84\x8f\xbb\x6e\xc3\x30\x0c\x45\x77\x7f\xc5\x45\x96\x2c\xb5\xb4\x74\xea\xd6\xb1\x7b\x7f\x80\x91\x68\x89\x88\x1e\xae\x48\xe7\xf1\xf7\x85\xd3\x02\xcd\xd6\xf5\x00\xe7\xf0\xd2\x7b\x7c\x66\x51\x2c\x52\x18\xa2\x68\x1c\x58\x95\xc6\x1d\x27\x0e\xb4\x29\xe3\x90\xc4\xf2\x76\x72\xa1\x57\xaf\x46\xb6\xe9\x2c\xd5\x57\x49\x83\x8c\xfd\xe5\xf5\x30\x79\x8f\x40\xed\x68\xc8\xd4\x62\xe1\x47\x4b\xa1\x46\xc3\xa4\x25\x5c\xc5\x32\x08\xeb\xe0\x45\x6e\x0e\xef\x86\xc2\xa4\x06\xcb\x64\x47\x85\x65\x46\x20\xe5\x3d\xb3\xf4\x81\xd4\xe7\x93\xb4\x48\x46\x6e\x47\x1f\xcb\x13\xd9\x17\x06\x2a\x85\x23\x96\xd1\xeb\xc3\x55\xaa\x8c\x28\x83\x83\xf5\x71\x7f\x01\xa9\xb2\xa1\x51\x65\xdd\xfd\x4c\x17\x46\xeb\xbf\xe7\x41\x2d\xfe\xff\x11\xae\x7d\x9c\x15\xa4\xe0\xdb\xca\xc1\x38\xba\x69\x5a\x29\x9c\x29\x31\xf4\xab\x88\xf1\x34\x79\x9f\xfa\x5b\xe2\xc6\xbb\xf5\xbc\x71\x5e\xcf\x09\x3f\x35\xe9\x4d\x31\x77\x38\xe7\xff\x80\x4b\x1d\x6e\xfa\x0e\x00\x00\xff\xff\x9d\x60\x3d\x88\x79\x01\x00\x00")

func doc_go() ([]byte, error) {
	return bindata_read(
		_doc_go,
		"doc.go",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
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
var _bindata = map[string]func() ([]byte, error){
	"1565341329_initial_schema.down.sql": _1565341329_initial_schema_down_sql,
	"1565341329_initial_schema.up.sql": _1565341329_initial_schema_up_sql,
	"doc.go": doc_go,
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
	for name := range node.Children {
		rv = append(rv, name)
	}
	return rv, nil
}

type _bintree_t struct {
	Func func() ([]byte, error)
	Children map[string]*_bintree_t
}
var _bintree = &_bintree_t{nil, map[string]*_bintree_t{
	"1565341329_initial_schema.down.sql": &_bintree_t{_1565341329_initial_schema_down_sql, map[string]*_bintree_t{
	}},
	"1565341329_initial_schema.up.sql": &_bintree_t{_1565341329_initial_schema_up_sql, map[string]*_bintree_t{
	}},
	"doc.go": &_bintree_t{doc_go, map[string]*_bintree_t{
	}},
}}
