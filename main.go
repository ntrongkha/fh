package main

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	ErrFound = errors.New("")
)

func checksum(path string, algo string) (string, error) {
	var (
		out string
		h   hash.Hash
	)

	f, err := os.Open(path)
	if err != nil {
		return out, err
	}

	switch algo {
	case "sha256":
		h = sha256.New()
	default:
		h = md5.New()
	}

	if _, err := io.Copy(h, f); err != nil {
		return out, err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func main() {
	var (
		dir  string
		hash string
		algo string
	)
	flag.StringVar(&dir, "dir", ".", "directory to find")
	flag.StringVar(&hash, "hash", "", "file hash")
	flag.StringVar(&algo, "algo", "md5", "hash algorithm")
	flag.Parse()

	if hash == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalln("resolve path failed:", err)
	}

	err = filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				if info.Name() == ".git" {
					return filepath.SkipDir
				}
				log.Println("searching directory:", info.Name())
				// skip git directory by default
				return nil
			}
			hfile, err := checksum(path, algo)
			if err != nil {
				return err
			}
			if hfile == hash {
				return fmt.Errorf("found hash: %s in path: %s %w", hash, path, ErrFound)
			}
			return nil
		})
	if errors.Is(err, ErrFound) {
		log.Println(err)
	} else if err != nil {
		log.Fatalln("walk failed:", err)
	}
	log.Println("finished")
}
