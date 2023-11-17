package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/esimclub/lpa-agent/frontend"
	"github.com/esimclub/lpa-agent/lpac"
)

var (
	listen     string
	lpacDir    string
	download   bool
	reader     string
	readerName string
)

func init() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&listen, "listen", ":9527", "listen address")
	flag.StringVar(&lpacDir, "lpac-dir", filepath.Join(workDir, "lpac-cli"), "lpac directory")
	flag.BoolVar(&download, "download", true, "download lpac")
	flag.StringVar(&reader, "reader", "pcsc", "pscs or at")
	flag.StringVar(&readerName, "readerName", "", "reader name")
	flag.Parse()
}

func dylibWithExtension(filename string) string {
	switch runtime.GOOS {
	case "windows":
		return filename + ".dll"
	case "darwin":
		return filename + ".dylib"
	default:
		return filename + ".so"
	}
}

func envs() map[string]string {
	envs := map[string]string{
		"HTTP_INTERFACE": filepath.Join(lpacDir, dylibWithExtension("libhttpinterface_curl")),
		"APDU_INTERFACE": filepath.Join(lpacDir, dylibWithExtension(fmt.Sprintf("libapduinterface_%s", reader))),
	}

	if reader == "at" {
		envs["AT_DEVICE"] = readerName
	}
	return envs
}

func lpacProgram() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(lpacDir, "lpac.exe")
	}

	if err := os.Chmod(filepath.Join(lpacDir, "lpac"), 0755); err != nil {
		panic(err)
	}
	return filepath.Join(lpacDir, "lpac")
}

func main() {
	if download {
		if err := Download(lpacDir); err != nil {
			panic(err)
		}
	}

	lpaApiHandler := NewAPIHTTPHandler(&lpac.CommandLine{
		Program: lpacProgram(),
		EnvMap:  envs(),
	})
	mux := http.NewServeMux()
	mux.Handle("/api/lpa/", http.StripPrefix("/api/lpa", lpaApiHandler))
	mux.Handle("/", http.FileServer(http.FS(frontend.DistFS())))

	slog.Info("lpac-agent has started", "listen", listen, "lpacDir", lpacDir, "download", download, "reader", reader, "readerName", readerName)
	if strings.Contains(listen, ".") {
		slog.Info(fmt.Sprintf("please visit http://%s", listen))
	} else {
		slog.Info(fmt.Sprintf("please visit http://127.0.0.1:%s", strings.TrimLeft(listen, ":")))
	}

	if err := http.ListenAndServe(listen, mux); err != nil {
		panic(err)
	}
}
