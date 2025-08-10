package internal

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

type Downloader interface {
	Download(ctx context.Context, url string, storagePath string) error
	Extract(ctx context.Context, storagePath string, extractPath string) error
}

type HTTPDownloader struct{}

func NewHTTPDownloader() Downloader {
	return &HTTPDownloader{}
}

func (d *HTTPDownloader) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (d *HTTPDownloader) Download(ctx context.Context, url string, storagePath string) error {
	if d.fileExists(storagePath) {
		log.Printf("File already exists: %s", storagePath)
		return nil
	} else {
		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatalf("Failed to create storage directory: %v", err)
		}
	}

	client := grab.NewClient()

	req, err := grab.NewRequest(storagePath, url)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	t := time.NewTicker(2 * time.Second)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			log.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			break Loop
		}
	}

	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Download saved to ./%v \n", resp.Filename)

	return nil
}

func (e *HTTPDownloader) Extract(ctx context.Context, source string, destDir string) error {
	if e.alreadyExtracted(destDir) {
		log.Printf("Files already extracted: %s", destDir)
		return nil
	} else {
		if err := os.MkdirAll(destDir, 0755); err != nil {
			log.Fatalf("Failed to create extract directory: %v", err)
		}
	}

	reader, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("error opening zip file: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		err := e.extractFile(ctx, file, destDir)
		if err != nil {
			return fmt.Errorf("error extracting file %s: %w", file.Name, err)
		}
	}

	log.Printf("Extracted files to %s", destDir)

	return nil
}

func (e *HTTPDownloader) alreadyExtracted(path string) bool {
	files, err := os.ReadDir(path)
	if err != nil {
		return false
	}

	return len(files) > 0
}

func (e *HTTPDownloader) extractFile(ctx context.Context, file *zip.File, destDir string) error {
	filePath := filepath.Join(destDir, file.Name)

	if file.FileInfo().IsDir() {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	sourceFile, err := file.Open()
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}
