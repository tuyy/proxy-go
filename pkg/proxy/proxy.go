package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func MakeServeMux(defaultDestUrl string, urlMappings []map[string]string) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	urlCache := make(map[string]*url.URL)

	// Note: "/" pattern은 모든 요청을 처리한다. default handling을 위해 추가한다.
	urlMappings = append(urlMappings, map[string]string{"/": defaultDestUrl})

	var err error

	for _, pathAndUrlMap := range urlMappings {
		path := pathAndUrlMap["path"]
		destUrl := pathAndUrlMap["dest_url"]

		parsedUrl, ok := urlCache[destUrl]
		if !ok {
			parsedUrl, err = url.Parse(destUrl)
			if err != nil {
				return nil, err
			}
			urlCache[destUrl] = parsedUrl
		}

		mux.Handle(path, httputil.NewSingleHostReverseProxy(parsedUrl))
	}

	return mux, nil
}
