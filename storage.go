package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

type Store struct {
	structOpts structOpts
}

type pathTransform func(string) PathKey

type structOpts struct {
	pathTransformFunc pathTransform
	Metadata          *Metadata
}

type PathKey struct {
	pathname string
	filename string
}

func NewStore(opts structOpts) *Store {
	return &Store{
		structOpts: opts,
	}
}

func DefaultPathTransformFunc(key string) PathKey {
	return PathKey{
		pathname: key,
		filename: key,
	}
}

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		if to > len(hashStr) {
			to = len(hashStr)
		}
		paths[i] = hashStr[from:to]
	}

	pathKey := strings.Join(paths, "/")

	return PathKey{
		pathname: pathKey,
		filename: hashStr,
	}
}

func (s *Store) ReadStream(key string) (io.Reader, error) {
	// PathKey := s.structOpts.pathTransformFunc(key)

	// filePath := PathKey.pathname

	filePath, ok := s.structOpts.Metadata.Get(key)
	if !ok {
		return nil, os.ErrNotExist
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	buff := new(bytes.Buffer)

	if _, err = io.Copy(buff, file); err != nil {
		return nil, err
	}

	file.Close()

	return buff, nil
}

func (s *Store) WriteStream(key string, w io.Reader) error {
	pathKey := s.structOpts.pathTransformFunc(key)
	//pathKey := s.CASPathTransformFunc(key)

	err := os.MkdirAll(pathKey.pathname, os.ModePerm)
	if err != nil {
		return err
	}

	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, w)
	if err != nil {
		return err
	}

	hash := md5.Sum(buff.Bytes())
	hashStr := hex.EncodeToString(hash[:])
	pathKey.filename = hashStr

	filePath := pathKey.pathname + "/" + pathKey.filename

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()

	n, err := io.Copy(f, buff)
	if err != nil {
		return err
	}

	err = s.structOpts.Metadata.Set(key, filePath)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, filePath)
	return nil
}
