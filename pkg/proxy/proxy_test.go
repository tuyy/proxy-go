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

var respCntMap map[string]int

func TestMakeServeMuxHappy(t *testing.T) {
	respCntMap = make(map[string]int)
	testServers, testDestUrl := makeTestServerAndUrls(4)
	mux, err := MakeServeMux(testDestUrl[0], []map[string]string{
		{"path": "/test1", "dest_url": testDestUrl[1]},
		{"path": "/test2", "dest_url": testDestUrl[2]},
		{"path": "/test3", "dest_url": testDestUrl[3]},
	})
	assert.Nil(t, err)

	frontendProxy := httptest.NewServer(mux)
	defer frontendProxy.Close()

	paths := []string{"", "test1", "test2", "test3"}
	for _, path := range paths {
		resp, _ := http.Get(fmt.Sprintf("%s/%s", frontendProxy.URL, path))
		b, _ := io.ReadAll(resp.Body)
		assert.Equal(t, "SUCCESS PASS PROXY\n", string(b))
	}

	assert.Equal(t, len(paths), len(respCntMap))
	closeTestServers(testServers)
}

func closeTestServers(testServers []*httptest.Server) {
	for _, server := range testServers {
		server.Close()
	}
}

func makeTestServerAndUrls(testServerCnt int) ([]*httptest.Server, []string) {
	var testDestUrls []string
	var testServers []*httptest.Server
	for i := 0; i < testServerCnt; i++ {
		newServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			respCntMap[r.RemoteAddr] += 1
			fmt.Fprintf(w, "SUCCESS PASS PROXY\n")
		}))

		testServers = append(testServers, newServer)
		testDestUrls = append(testDestUrls, newServer.URL)
	}
	return testServers, testDestUrls
}
