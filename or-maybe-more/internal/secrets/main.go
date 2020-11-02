package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/saracen/go7z"
)

var archiveName string

func openSecretsArchive(path string) (files []string, ok bool) {
	h, err := go7z.OpenReader(path)
	if err != nil {
		log.Println(fmt.Errorf("cant open secrets archive: %w\n", err))
		return nil, false
	}
	defer h.Close()
	var pwprompt string = os.Getenv(archiveName)
	if pwprompt == "" {
		_, fn := filepath.Split(path)
		fname := strings.Split(fn, filepath.Ext(fn))[0]
		fmt.Printf("enter password for %s:\n", fname)
		fmt.Scanln(&pwprompt)
	}
	h.Options.SetPassword(pwprompt)

	for {
		hdr, err := h.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			panic(err)
		}

		// If empty stream (no contents) and isn't specifically an empty file...
		// then it's a directory.
		if hdr.IsEmptyStream && !hdr.IsEmptyFile {
			if err := os.MkdirAll(hdr.Name, os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		// Create file
		f, err := os.Create(hdr.Name)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := io.Copy(f, h); err != nil {
			panic(err)
		}

		files = append(files, hdr.Name)
	}

	return files, true
}

func createVarsAndRm(files []string) error {
	for _, f := range files {
		val, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		if err := os.Setenv(strings.ToUpper(strings.Split(f, ".")[0]), strings.Split(string(val), "\n")[0]); err != nil {
			return err
		}

		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
