package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"strings"
)

func main() {
	cleanF()

	for _, v := range getProjectFilename() {
		genMD5SumFile(v)
	}
}

func cleanF() {
	bts, err := os.ReadFile("chnroutes2/chnroutes.txt")
	if err != nil {
		log.Fatalln(err)
	}

	arr := bytes.Split(bts, []byte("\n"))
	buf := make([][]byte, 0, len(arr))

	for i := range arr {
		if len(arr[i]) > 0 && arr[i][0] == '#' {
			continue
		}
		buf = append(buf, arr[i])
	}
	
	var derp []byte
	if derp, err = GetDerp(); err != nil {
		log.Fatalln(err)
	}

	if err = os.WriteFile("chnroutes2/chnroutes.txt", append(bytes.Join(buf, []byte("\n")), derp...), 0666); err != nil {
		log.Fatalln(err)
	}
}

func genMD5SumFile(filename string) {
	originData, err := os.ReadFile("chnroutes2/" + filename)
	if err != nil {
		log.Fatalf("open %s err: %s\n", filename, err)
	}

	h := md5.Sum(originData)

	dst := make([]byte, hex.EncodedLen(len(h)))
	hex.Encode(dst, h[:])

	if err = os.WriteFile("chnroutes2/"+filename+".md5sum.txt", dst, 0666); err != nil {
		log.Fatalf("write file %s.md5sum.txt err: %s\n", filename, err)
	}
}

func getProjectFilename() []string {
	dirEntry, err := os.ReadDir("chnroutes2")
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

func GetDerp() ([]byte, error) {
	resp, err := http.Get("https://login.tailscale.com/derpmap/default")
	if err != nil {
		return nil, err
	}

	var res Derp
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	for _, v := range res.Regions {
		for _, node := range v.Nodes {
			buf.WriteString(node.IPv4)
			buf.Write([]byte("/32\n"))
		}
	}
	return buf.Bytes(), nil
}

type Derp struct {
	Regions map[string]Region `json:"regions"`
}

type Region struct {
	RegionID   int    `json:"RegionID"`
	RegionCode string `json:"RegionCode"`
	RegionName string `json:"RegionName"`
	Nodes      []struct {
		Name     string `json:"Name"`
		RegionID int    `json:"RegionID"`
		HostName string `json:"HostName"`
		IPv4     string `json:"IPv4"`
		IPv6     string `json:"IPv6"`
	} `json:"Nodes"`
}
