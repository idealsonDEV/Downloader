package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

func Resume(task string, url string) *Files {
	raw, _ := ioutil.ReadFile("./myget/" + task + "/state.json")
	var d State
	json.Unmarshal(raw, &d)
	if url == "" {
		url = d.Url
	}
	if d.Resumeale == false {
		os.RemoveAll("./myget/" + task + "/")
		return NewDownloader(url, 1)
	} else {
		tpart := make([]Part, 0)
		offset := d.Size / d.Npart
		for i := int64(0); i < d.Npart; i++ {
			name := fmt.Sprintf("./myget/"+task+"/"+task+".part%d", i)
			file, _ := os.Open(name)
			defer file.Close()
			stat, _ := file.Stat()
			bytes := int64(stat.Size())
			rangeStart := (i * offset) + 1 + bytes
			if i == 0 {
				rangeStart = 0 + bytes
			}
			rangeEnd := (i + 1) * offset
			if i == d.Npart-1 {
				rangeEnd = d.Size
			}
			psize := rangeEnd - rangeStart + bytes
			tpart = append(tpart, Part{url: url, filename: name, rangeFrom: rangeStart, rangeTo: rangeEnd, psize: psize})
		}
		ret := &Files{}
		ret.finish = false
		ret.filename = task
		ret.url = url
		ret.resumable = d.Resumeale
		ret.size = d.Size
		ret.npart = d.Npart
		ret.parts = tpart

		return ret
	}
}

func (d *Files) DoResume() {
	var ws sync.WaitGroup
	MkdirIfNotExist("./myget")
	folder := string("./myget/" + d.filename)
	MkdirIfNotExist(folder)
	for _, p := range d.parts {
		ws.Add(1)
		if p.rangeFrom < p.rangeTo {
			go func(d *Files, part Part, folder string) {
				defer ws.Done()
				pfilepath := string(part.filename)
				pfile, _ := os.OpenFile(pfilepath, os.O_WRONLY|os.O_APPEND, 0777)
				defer pfile.Close()
				req, _ := http.NewRequest("GET", part.url, nil)
				req.Header.Add("Range", fmt.Sprintf("bytes=%d-%d", part.rangeFrom, part.rangeTo))
				resp, _ := client.Do(req)
				defer resp.Body.Close()
				io.Copy(pfile, resp.Body)
			}(d, p, folder)
		} else {
			ws.Done()
		}
	}
	ws.Wait()
}
