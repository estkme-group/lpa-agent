package main

import (
	"encoding/json"
	"flag"
	"github.com/esimclub/lpa-agent/frontend"
	"github.com/esimclub/lpa-agent/lpac"
	"net/http"
	"os"
)

var config Configuration

func init() {
	config.Listen = ":9527"
	var configFile string
	flag.StringVar(&configFile, "c", "./config.json", "Configuration filepath")
	flag.Parse()
	file, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	if err = json.Unmarshal(file, &config); err != nil {
		panic(err)
	}
}

func main() {
	lpaApiHandler := NewAPIHTTPHandler(&lpac.CommandLine{
		Program: config.Program,
		EnvMap:  config.EnvMap,
	})
	mux := http.NewServeMux()
	mux.Handle("/api/lpa/", http.StripPrefix("/api/lpa", lpaApiHandler))
	mux.Handle("/", http.FileServer(http.FS(frontend.DistFS())))
	if err := http.ListenAndServe(config.Listen, mux); err != nil {
		panic(err)
	}
}
