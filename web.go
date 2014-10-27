package goutils

import (
	"io/ioutil"
	"net/http"
)

func DownloadHttpFile(url string, target string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(target, bytes, 0444); err != nil {
		return nil, err
	}

	return bytes, nil
}
