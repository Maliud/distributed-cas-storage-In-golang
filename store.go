package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i * blocksize, (i * blocksize) + blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey {
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}
	
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	PathName string
	Filename string
}

func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.Filename)
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) Has(key string) bool {
	PathKey := s.PathTransformFunc(key)

	_, err := os.Stat(PathKey.FullPath())
	if err == fs.ErrNotExist {
		return false
	}

	return true
}

func (s *Store) Delete(key string) error {
	PathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("diskden silindi [%s]", PathKey.Filename)
	}()
	
	//if err := os.RemoveAll(PathKey.FullPath()); err != nil {
	//	return err
	//}
	return os.RemoveAll(PathKey.FirstPathName())
	
	
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, nil
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	PathKey := s.PathTransformFunc(key)
	
	return os.Open(PathKey.FullPath())
	
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	


	fullPath := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("diske (%d) bayt yazıldı: %s", n, fullPath)
	return nil
}