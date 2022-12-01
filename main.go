package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"strings"
)

func main() {
	for _, v := range getProjectFilename() {
		genMD5SumFile(v)
	}
}

func genMD5SumFile(filename string) {
	originData, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("open %s err: %s\n", filename, err)
	}

	md5Data, _ := os.ReadFile(filename + ".md5sum.txt")
	h := md5.Sum(originData)

	if bytes.Equal(h[:], md5Data) {
		return
	}

	dst := make([]byte, hex.EncodedLen(len(h)))
	hex.Encode(dst, h[:])

	if err = os.WriteFile(filename+".md5sum.txt", dst, 0666); err != nil {
		log.Fatalf("write file %s.md5sum.txt err: %s\n", filename, err)
	}
}

func getProjectFilename() []string {
	dirEntry, err := os.ReadDir(".")
	if err != nil {
		panic(err)
	}

	filenames := make([]string, 0, len(dirEntry))
	for _, v := range dirEntry {
		if v.IsDir() ||
			strings.HasPrefix(v.Name(), ".") ||
			strings.HasSuffix(v.Name(), ".md5sum.txt") {
			continue
		}

		filenames = append(filenames, v.Name())
	}

	return filenames
}
