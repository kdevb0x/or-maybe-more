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

var archivePass string

type secrets struct {
	files []string
	vars  []string
}

func (s *secrets) openSecretsArchive(path string) error {
	f, err := openSecretsArchive(path)
	if err != nil {
		return err
	}
	s.files = f
	return nil
}

// creates env vars using s.files, deleting the file afterwards if rmFiles is
// set to true.
func (s *secrets) createVars(rmFiles bool) error {
	for _, f := range s.files {
		val, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		if err := os.Setenv(strings.ToUpper(strings.Split(f, ".")[0]), strings.Split(string(val), "\n")[0]); err != nil {
			return err
		}

		if rmFiles {
			if err := os.Remove(f); err != nil {
				return fmt.Errorf("unable to remove %s: %w\n", f, err)
			}
		}
	}
	return nil
}

func openSecretsArchive(path string) (files []string, err error) {
	h, err := go7z.OpenReader(path)
	if err != nil {
		log.Println(fmt.Errorf("cant open secrets archive: %w\n", err))
		return nil, err
	}
	defer h.Close()
	var pwprompt string = os.Getenv(archivePass)

	// loop until the user gives us something
	for pwprompt == "" {
		_, fn := filepath.Split(path)
		fmt.Printf("enter password for %s:\n", fn)
		fmt.Scanln(&pwprompt)
	}
	h.Options.SetPassword(pwprompt)

	for {
		hdr, err := h.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return nil, err
		}

		// If empty stream (no contents) and isn't specifically an empty file...
		// then it's a directory.
		if hdr.IsEmptyStream && !hdr.IsEmptyFile {
			if err := os.MkdirAll(hdr.Name, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		// Create file
		f, err := os.Create(hdr.Name)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if _, err := io.Copy(f, h); err != nil {
			return nil, err
		}

		files = append(files, f.Name())
	}
	return files, nil
}

func main() {
	archivePass = "OOM_SECRETS_AR_PASS"
}
