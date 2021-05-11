package proxy

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

func TestReverseProxy(t *testing.T) {
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "this call was relayed by the reverse proxy")
	}))
	defer backendServer.Close()

	rpURL, err := url.Parse(backendServer.URL)
	if err != nil {
		log.Fatal(err)
	}
	frontendProxy := httptest.NewServer(httputil.NewSingleHostReverseProxy(rpURL))
	defer frontendProxy.Close()

	resp, err := http.Get(frontendProxy.URL)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)
}

func TestMakeServeMuxHappy(t *testing.T) {
	pathAndUrlMap := map[string]string{
		"/path1": "http://bye.com",
		"/path2": "http://hello.com",
		"/path3": "http://bye.com",
	}

	_, err := MakeServeMux("http://hello.com", pathAndUrlMap)
	assert.Nil(t, err)
}
