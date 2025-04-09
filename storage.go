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

	filePath := pathKey.pathname + "/" + hashStr

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buff)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, filePath)
	return nil
}
