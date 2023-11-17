package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/hashicorp/go-version"
)

const (
	githubVersionApi          = "https://api.github.com/repos/estkme-group/lpac/releases/latest"
	githubReleaseDownloadLink = "https://github.com/estkme-group/lpac/releases/download/%s/%s"
	localVersionFile          = "version"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name string `json:"name"`
	}
}

func fetchRelease() (*GitHubRelease, error) {
	httpClient := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := httpClient.Get(githubVersionApi)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch latest release failed, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()
	release := &GitHubRelease{}
	err = json.NewDecoder(resp.Body).Decode(release)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func shouldDownload(lpacDir string, targetVersion string) (bool, error) {
	b, err := os.ReadFile(filepath.Join(lpacDir, localVersionFile))
	if err != nil {
		return true, err
	}
	localVersion, err := version.NewVersion(string(b))
	if err != nil {
		return true, err
	}
	remoteVersion, err := version.NewVersion(targetVersion)
	if err != nil {
		return true, err
	}
	return localVersion.LessThan(remoteVersion), nil
}

func assetName() string {
	// TODO: support by arch
	switch runtime.GOOS {
	case "windows":
		return "lpac-windows-x86_64.zip"
	case "darwin":
		return "lpac-macos-universal.zip"
	default:
		return "lpac-linux-x86_64.tar.gz"
	}
}

func downloadFile(dir string, githubRelease *GitHubRelease) (err error) {
	downloadUrl := fmt.Sprintf(githubReleaseDownloadLink, githubRelease.TagName, assetName())
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download lpac failed, status code: %d", resp.StatusCode)
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	destPath := filepath.Join(dir, assetName())
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	if err = decompress(destPath); err != nil {
		return err
	}

	if err := os.Remove(destPath); err != nil {
		return err
	}
	return writeVersionFile(dir, githubRelease.TagName)
}

func writeVersionFile(dir string, version string) error {
	return os.WriteFile(filepath.Join(dir, localVersionFile), []byte(version), 0644)
}

func unzip(path string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return err
	}

	targetPath := filepath.Dir(path)
	for _, file := range archive.File {
		if file.FileInfo().IsDir() {
			continue
		}

		targetFilePath := filepath.Join(targetPath, file.Name)
		if _, err := os.Stat(targetFilePath); os.IsExist(err) {
			if err := os.Remove(targetFilePath); err != nil {
				return fmt.Errorf("unzip: Remove() failed: %s", err.Error())
			}
		}

		outFile, err := os.Create(filepath.Join(targetPath, file.Name))
		if err != nil {
			return fmt.Errorf("unzip: Create() failed: %s", err.Error())
		}

		inFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("unzip: Open() failed: %s", err.Error())
		}
		if _, err := io.Copy(outFile, inFile); err != nil {
			return fmt.Errorf("tarx: Copy() failed: %s", err.Error())
		}
		outFile.Close()
		inFile.Close()
	}
	return nil
}

func tarx(path string) error {
	tarFile, err := os.Open(path)
	if err != nil {
		return err
	}

	archive, err := gzip.NewReader(tarFile)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(archive)
	targetPath := filepath.Dir(path)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("tarx: Next() failed: %s", err.Error())
		}

		if header.FileInfo().IsDir() {
			continue
		}

		targetFilePath := filepath.Join(targetPath, header.Name)
		if _, err := os.Stat(targetFilePath); os.IsExist(err) {
			if err := os.Remove(targetFilePath); err != nil {
				return fmt.Errorf("tarx: Remove() failed: %s", err.Error())
			}
		}
		outFile, err := os.Create(filepath.Join(targetPath, header.Name))
		if err != nil {
			return fmt.Errorf("tarx: Create() failed: %s", err.Error())
		}
		if _, err := io.Copy(outFile, tarReader); err != nil {
			return fmt.Errorf("tarx: Copy() failed: %s", err.Error())
		}
		outFile.Close()
	}

	return nil
}

func decompress(path string) error {
	if filepath.Ext(path) == ".zip" {
		return unzip(path)
	}
	return tarx(path)
}

func isSupportedArch() bool {
	switch runtime.GOOS {
	case "windows":
		return runtime.GOARCH == "amd64" || runtime.GOARCH == "386"
	case "darwin":
		return true
	default:
		return runtime.GOARCH == "amd64" || runtime.GOARCH == "386"
	}
}

func Download(lpacDir string) error {
	if !isSupportedArch() {
		return fmt.Errorf("not supported arch yet: %s %s", runtime.GOOS, runtime.GOARCH)
	}

	githubRelease, err := fetchRelease()
	if err != nil {
		return err
	}

	should, err := shouldDownload(lpacDir, githubRelease.TagName)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if should {
		slog.Info("download lpac", "version", githubRelease.TagName)
		return downloadFile(lpacDir, githubRelease)
	}
	return nil
}
