package httpclient

import (
	"bytes"
	"net/http"

	"log"
)

// HTTPGet do http post by url and param
func HTTPGet(url, token string) (*http.Response, error) {
	// glog.Info(url)
	// glog.Info(bytes.NewBuffer(json))
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.do err: %v", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Fatalf("statusCode != 200: %v", resp)
		return nil, err
	}
	return resp, nil
}

// HTTPPost do http post by url and param
func HTTPPost(url string, json []byte) (*http.Response, error) {
	// glog.Info(url)
	// glog.Info(bytes.NewBuffer(json))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.do err: %v", err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Fatalf("statusCode != 200: %v", resp)
		return nil, err
	}
	return resp, nil
}
