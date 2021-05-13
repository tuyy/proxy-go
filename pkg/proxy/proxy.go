package proxy

import (
	"github.com/tuyy/proxy-go/pkg/log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func handleLog(handler http.Handler, destUrl *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		log.Access("METHOD:%s PATH:%s DEST_URL:%s ELAPSED:%s",
			r.Method,
			r.RequestURI,
			destUrl.String(),
			time.Now().Sub(start))
	}
}

func MakeServeMux(defaultDestUrl string, urlMappings []map[string]string) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	urlCache := make(map[string]*url.URL)

	// Note: "/" 패턴은 모든 요청을 처리한다. Default 핸들링을 위해 추가한다.
	urlMappings = append(urlMappings, map[string]string{"path": "/", "dest_url": defaultDestUrl})

	for _, pathAndUrlMap := range urlMappings {
		path := pathAndUrlMap["path"]
		destUrl := pathAndUrlMap["dest_url"]

		parsedUrl, err := getAndAddUrl(destUrl, urlCache)
		if err != nil {
			return nil, err
		}

		mux.HandleFunc(path, handleLog(httputil.NewSingleHostReverseProxy(parsedUrl), parsedUrl))
	}

	return mux, nil
}

func getAndAddUrl(destUrl string, urlCache map[string]*url.URL) (*url.URL,error) {
	parsedUrl, ok := urlCache[destUrl]
	if ok {
		return parsedUrl, nil
	}

	newUrl, err := url.Parse(destUrl)
	if err != nil {
		return nil, err
	}

	urlCache[destUrl] = newUrl
	return newUrl, nil
}