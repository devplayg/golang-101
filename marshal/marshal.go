package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var resourceDir = "data"

var extMap = map[string]bool{
	".js":    true,
	".png":   true,
	".css":   true,
	".eot":   true,
	".svg":   true, // image/svg+xml
	".ttf":   true,
	".woff":  true,
	".woff2": true, // font/woff2
}

func main() {

	assets := getAssets()
	if len(assets) < 1 {
		return
	}

	b, m := encodeAssets(assets)

	compressed, err := compressData(b)
	if err != nil {
		panic(err)
	}

	path, err := writeData(compressed, m)
	if err != nil {
		panic(err)
	}
	fmt.Println(path)

	//encodedStr, m := encodeAssets(list)

	//f.WriteString("/*\n")
	////fmt.Println("Keys---------------------------------")
	//
	keys := make([]string, 0)
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("	%s\n", k)
		//f.WriteString("\t" + k + "\n")
	}
	//f.WriteString("*/\n")
	//f.WriteString(encodedStr)
	//
	//fmt.Printf("output: %s\n", f.Name())
}

func writeData(b []byte, m map[string][]byte) (string, error) {
	f, err := ioutil.TempFile(".", "result")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("/*\n")
	keys := make([]string, 0)
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, k := range keys {
		//fmt.Printf("	%s\n", k)
		f.WriteString("\t" + k + "\n")
	}
	f.WriteString("*/\n")

	f.WriteString(base64.StdEncoding.EncodeToString(b))

	return f.Name(), nil
}

func compressData(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil

}

func encodeAssets(list []string) ([]byte, map[string][]byte) {
	m, err := encode(list)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	return b, m

}

func getAssets() []string {

	list := make([]string, 0, 10)

	err := filepath.Walk(resourceDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return nil
		}
		if f.IsDir() {
			return nil
		}
		if _, exists := extMap[filepath.Ext(path)]; !exists {
			return nil
		}

		list = append(list, filepath.ToSlash(path))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return list

}

func encode(files []string) (map[string][]byte, error) {
	m := make(map[string][]byte)
	for _, path := range files {
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
		m[strings.TrimPrefix(path, resourceDir)] = b
		fmt.Printf("[%s] is encoded successfully. len=%d\n", filepath.Base(path), len(b))
	}
	return m, nil
}
