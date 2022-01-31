package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

type TestsFetch struct {
	name          string
	url           string
	server        *httptest.Server
	md5Response   string
	expectedError error
}

func TestGetArgs(t *testing.T) {
	os.Args = append(os.Args, "http://google.com")
	parallel, urls := getArgs()
	if *parallel != 10 {
		t.Errorf("Failed expected parallel 4, got %v", *parallel)
	}

	if !reflect.DeepEqual(flag.Args(), urls) {
		t.Errorf("Failed expected args 1, got %v", len(urls))
	}
}

func TestGetMd5(t *testing.T) {
	parallel := 3
	urls := []string{"http://google.com", "http://adjust.com", "http://twitter.com"}
	process := make(chan bool)
	go getMD5(&parallel, urls, process)
	<-process
}

func TestFetch(t *testing.T) {
	tests := []TestsFetch{
		{
			name: "test md5 hash of the response",
			url:  "/test1",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("this is a sample web page"))

			})),
			md5Response:   fmt.Sprintf("%x", md5.Sum([]byte("this is a sample web page"))),
			expectedError: nil,
		},
		{
			name: "test no response found",
			url:  "/test2",
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			})),
			md5Response:   "",
			expectedError: fmt.Errorf("bad response from server: 404"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			defer test.server.Close()

			resp, err := fetch(test.server.URL + test.url)
			if !reflect.DeepEqual(resp, test.md5Response) {
				t.Errorf("Failed expected response %v, got %v\n", "", resp)
			}
			if err != nil {
				if !strings.Contains(err.Error(), "404") {
					t.Errorf("Failed expected error: %v, got %v\n", test.expectedError, err)
				}
			}
		})
	}
}
