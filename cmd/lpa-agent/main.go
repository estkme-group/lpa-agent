package main

import (
	"github.com/esimclub/lpa-agent/lpac"
	"net/http"
)

func main() {
	lpaApiHandler := NewAPIHTTPHandler(&lpac.CommandLine{
		Program: "/home/septs/tools/lpac/build/output/lpac",
		EnvMap: map[string]string{
			"APDU_INTERFACE": "/home/septs/tools/lpac/build/output/libapduinterface_pcsc.so",
			"HTTP_INTERFACE": "/home/septs/tools/lpac/build/output/libhttpinterface_curl.so",
		},
	})
	mux := http.NewServeMux()
	mux.Handle("/api/lpa/", http.StripPrefix("/api/lpa", lpaApiHandler))
	mux.Handle("/", http.FileServer(http.Dir("/home/septs/Projects/lpa-agent/frontend/dist")))
	if err := http.ListenAndServe(":10240", mux); err != nil {
		panic(err)
	}
}
