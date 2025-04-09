package main

import (
	"bytes"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "test.txt"

	expectedPathname := "4b6fc/b2d52/1ef0f/d442a/5301e/7932d/16cc9/f375a"
	expectedFilename := "4b6fcb2d521ef0fd442a5301e7932d16cc9f375a"

	pathKey := CASPathTransformFunc(key)

	if pathKey.pathname != expectedPathname || pathKey.filename != expectedFilename {
		t.Errorf("unexpected pathKey:\ngot: {pathname: %s, filename: %s}\nwant: {pathname: %s, filename: %s}",
			pathKey.pathname, pathKey.filename,
			expectedPathname, expectedFilename)
	} else {
		t.Logf("pathKey correctly transformed: %+v", pathKey)
	}
}

func TestWriteStream(t *testing.T) {
	key := "test.txt"
	data := []byte("whassup ma boy!")
	reader := bytes.NewReader(data)

	store := NewStore(structOpts{
		pathTransformFunc: CASPathTransformFunc,
	})

	err := store.WriteStream(key, reader)
	if err != nil {
		t.Fatalf("WriteStream failed: %v", err)
	}

	t.Logf("WriteStream succeeded for key: %s", key)
}
