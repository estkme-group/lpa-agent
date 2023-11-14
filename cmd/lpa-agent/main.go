package main

import (
	"github.com/esimclub/lpa-agent/lpac"
	"net/http"
)

func main() {
	apiHandler := NewAPIHTTPHandler(&lpac.CommandLine{
		Program: "/home/septs/tools/lpac/build/output/lpac",
		EnvMap: map[string]string{
			"APDU_INTERFACE": "./libapduinterface_pcsc.so",
			"HTTP_INTERFACE": "./libhttpinterface_curl.so",
		},
	})
	mux := http.NewServeMux()
	mux.Handle("/api/", apiHandler)
	mux.Handle("/", http.FileServer(http.Dir("./public")))
	if err := http.ListenAndServe(":10240", mux); err != nil {
		panic(err)
	}
}
