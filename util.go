package main

import (
	"encoding/json"
	"net"
	"os"
)

func FatalCheck(err error) {
	if err != nil {
		panic(err)
	}
}

func FilterIPV4(ips []net.IP) []string {
	var ret = make([]string, 0)
	for _, ip := range ips {
		if ip.To4() != nil {
			ret = append(ret, ip.String())
		}
	}
	return ret
}

func MkdirIfNotExist(folder string) error {
	if _, err := os.Stat(folder); err != nil {
		if err = os.MkdirAll(folder, 0777); err != nil {
			return err
		}
	}
	return nil
}

func toJson(p State) string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}
