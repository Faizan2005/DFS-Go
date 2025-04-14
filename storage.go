package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRoot = "DFSNetworkRoot"

type Store struct {
	structOpts StructOpts
}

type pathTransform func(string) PathKey

type StructOpts struct {
	PathTransformFunc pathTransform
	Metadata          *Metadata
	Root              string
}

type PathKey struct {
	pathname string
	filename string
}

func NewStore(opts StructOpts) *Store {
	store := &Store{
		structOpts: opts,
	}

	if store.structOpts.PathTransformFunc == nil {
		store.structOpts.PathTransformFunc = DefaultPathTransformFunc
	}

	if store.structOpts.Root == "" {
		store.structOpts.Root = defaultRoot
	}

	return store
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
	pathKey := s.structOpts.PathTransformFunc(key)
	//pathKey := s.CASPathTransformFunc(key)

	err := os.MkdirAll(s.structOpts.Root+"/"+pathKey.pathname, os.ModePerm)
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

	filePath := s.structOpts.Root + "/" + pathKey.pathname + "/" + pathKey.filename

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

func (s *Store) Remove(key string) error {
	filePath, ok := s.structOpts.Metadata.Get(key)
	if !ok {
		return os.ErrNotExist
	}

	fmt.Println(filePath)
	paths := strings.Split(filePath, "/")

	err := os.RemoveAll(paths[1])
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) TearDown() error {
	return os.RemoveAll(s.structOpts.Root)
}
