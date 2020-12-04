package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/awnumar/memguard"
	"github.com/saracen/go7z"
)

var archivePass = "OOM_SECRETS_ARCHIVE_P"

var archivePathKey = "OOM_SECRETS_ARCHIVE_PATH"

type VarLoader interface {
	LoadVars(o Opener)
	UnloadVars()
}

type Opener interface {
	Open(path string) (io.Reader, error)
}

type secrets struct {
	files   []string
	vars    map[string]string
	cleanup func() error
}

/*
func (s *secrets) openSecretsArchive(path string) error {
	f, err := openSecretsArchive(path)
	if err != nil {
		return err
	}
	s.files = f
	return nil
}
*/

// clean calls s.cleanup if it is present, otherwise it iteratively unsets all
// of the environment variables previously set by s.
func (s *secrets) clean() error {
	if s.cleanup != nil {
		return s.cleanup()
	}

	//else
	for k := range s.vars {
		if err := os.Unsetenv(k); err != nil {
			return err
		}
	}

	return nil
}

// creates env vars using s.files, deleting the file afterwards if rmFiles is
// set to true.
func (s *secrets) createVars() error {
	for _, f := range s.files {
		val, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		// get env var name from filename without ".pw" extension, and
		// its value from the first line of the file contents.
		if err := os.Setenv(strings.ToUpper(strings.Split(f, ".")[0]), strings.Split(string(val), "\n")[0]); err != nil {
			return err
		}

		if err := os.Remove(f); err != nil {
			return fmt.Errorf("unable to remove %s: %w\n", f, err)
		}
	}
	return nil
}

// like createVars, but for a single file.
func (s *secrets) createVar() error {

}

func openSecretsArchive(path string) (*secrets, error) {
	var s = new(secrets)
	s.vars = make(map[string]string)

	h, err := go7z.OpenReader(path)
	if err != nil {
		log.Println(fmt.Errorf("cant open secrets archive: %w\n", err))
		return nil, err
	}
	defer h.Close()
	var pwprompt = os.Getenv(archivePass)

	// if no archivePass is set,
	if pwprompt == "" {
		// loop until the user gives us something
		for pwprompt == "" {
			_, fn := filepath.Split(path)
			fmt.Printf("$%s is unset; please enter password for %s:\n", "OOM_SECRETS_ARCHIVE_PASS", fn)
			fmt.Scanln(&pwprompt)
		}
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

		// create temp dir
		// tmp, err := ioutil.TempDir("", "./.secrets")
		buf, err := memguard.NewBufferFromEntireReader(h)
		if err != nil {
			fmt.Printf("unable to create temporary enclave for extracting secret files: %s\n", err.Error())
			return nil, err
		}
		// Create file
		f, err := os.Create(filepath.Join(tmp, hdr.Name))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if _, err := io.Copy(f, h); err != nil {
			return nil, err
		}

		s.files = append(s.files, f.Name())
	}
	return s, nil
}

func main() {
	archivePass = "OOM_SECRETS_ARCHIVE_PASS"
}
