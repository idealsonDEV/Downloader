package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}
)

var (
	acceptRangeHeader   = "Accept-Ranges"
	contentLengthHeader = "Content-Length"
)

type Part struct {
	url       string
	filename  string
	rangeFrom int64
	rangeTo   int64
	psize     int64
}

type Files struct {
	url       string
	filename  string
	size      int64
	npart     int64
	parts     []Part
	skipTls   bool
	resumable bool
	finish    bool
}

type State struct {
	Url       string `json:"url"`
	Filename  string `json:"filename`
	Size      int64  `json:"size"`
	Npart     int64  `json:"nombre"`
	Resumeale bool
}

func NewDownloader(url string, par int) *Files {
	var resumable = true
	if par == 0 {
		par = 8
	}

	req, err := http.NewRequest("GET", url, nil)
	FatalCheck(err)

	resp, err := client.Do(req)
	FatalCheck(err)

	if resp.Header.Get(acceptRangeHeader) == "" {
		par = 1
		resumable = false
	}

	//get download range
	clen := resp.Header.Get(contentLengthHeader)
	if clen == "" {
		clen = "1" //set 1 because of progress bar not accept 0 length
		par = 1
		resumable = false
	}

	len, err := strconv.ParseInt(clen, 10, 64)
	FatalCheck(err)

	size := int64(len)

	file := filepath.Base(url)
	ret := new(Files)
	ret.url = url
	ret.filename = file
	ret.npart = int64(par)
	ret.size = size
	ret.parts = partCalculate(int64(par), size, url)
	ret.resumable = resumable
	ret.finish = false

	return ret
}

func partCalculate(par int64, len int64, url string) []Part {
	ret := make([]Part, 0)
	offset := len / par
	for index := int64(0); index < par; index++ {
		rangeStart := (index * offset) + 1
		if index == 0 {
			rangeStart = 0
		}
		rangeEnd := (index + 1) * offset
		if index == par-1 {
			rangeEnd = len
		}

		psize := rangeEnd - rangeStart
		file := filepath.Base(url)
		fname := fmt.Sprintf(file+".part%d", index)

		ret = append(ret, Part{url: url, filename: fname, rangeFrom: rangeStart, rangeTo: rangeEnd, psize: psize})
	}
	return ret
}

func (d *Files) Do() {
	var ws sync.WaitGroup
	MkdirIfNotExist("./myget")
	folder := string("./myget/" + d.filename)
	MkdirIfNotExist(folder)
	dj := toJson(State{Url: d.url, Filename: d.filename, Size: d.size, Npart: d.npart, Resumeale: d.resumable})
	ioutil.WriteFile("./myget/"+d.filename+"/state.json", []byte(dj), 0777)
	for _, p := range d.parts {
		ws.Add(1)
		go func(d *Files, part Part, folder string) {
			defer ws.Done()
			pfilepath := string(folder + "/" + part.filename)
			pfile, _ := os.Create(pfilepath)
			defer pfile.Close()
			req, _ := http.NewRequest("GET", part.url, nil)
			req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", part.rangeFrom, part.rangeTo))
			resp, _ := client.Do(req)
			defer resp.Body.Close()
			io.Copy(pfile, resp.Body)
		}(d, p, folder)
	}
	ws.Wait()
}

func (d *Files) Join() {
	folder := string("./finished/")
	MkdirIfNotExist(folder)
	out, _ := os.Create(folder + d.filename)
	offset := d.size / d.npart

	for i := int64(0); i < d.npart; i++ {
		name := fmt.Sprintf("./myget/"+d.filename+"/"+d.filename+".part%d", i)
		body, _ := ioutil.ReadFile(name)
		if i == 0 {
			out.Write(body)
		} else {
			out.WriteAt(body, int64((offset*i)+1))
		}
	}
	os.RemoveAll("./myget/" + d.filename + "/")
	d.finish = true
}
