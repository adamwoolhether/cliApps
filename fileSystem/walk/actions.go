package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func filterOut(path, ext string, minSize int64, minAge int, nameMatch string, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if nameMatch != "" {
		if match, _ := filepath.Match(nameMatch, filepath.Base(path)); !match {
			return true
		}
	}

	if minAge != 0 {
		if info.ModTime().After(time.Now().Add(-24 * time.Hour * time.Duration(minAge))) {
			return true
		}
	}

	if ext != "" {
		exts := strings.Split(ext, ",")
		for _, e := range exts {
			if filepath.Ext(path) == e {
				return false
			}
		}
		// File doesn't match any specified extensions.
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, delLoger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLoger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	// Check that destDir is a directory
	info, err := os.Stat(destDir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", destDir)
	}

	// Determine relative directory of file to be archived in relation to its source root path.
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	dest := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(destDir, relDir, dest)

	if err = os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	zw := gzip.NewWriter(out)
	zw.Name = filepath.Base(path)

	if _, err = io.Copy(zw, in); err != nil {
		return nil
	}
	if err = zw.Close(); err != nil {
		return err
	}

	return out.Close()
}
