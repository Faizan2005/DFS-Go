package main

import (
	"bytes"
	"fmt"
	"io"

	//"os"
	"testing"
)

var Key = "test.txt"
var Data = []byte("whassup ma boy!")

func TestPathTransformFunc(t *testing.T) {
	expectedPathname := "4b6fc/b2d52/1ef0f/d442a/5301e/7932d/16cc9/f375a"
	expectedFilename := "4b6fcb2d521ef0fd442a5301e7932d16cc9f375a"

	pathKey := CASPathTransformFunc(Key)

	if pathKey.pathname != expectedPathname || pathKey.filename != expectedFilename {
		t.Errorf("unexpected pathKey:\ngot: {pathname: %s, filename: %s}\nwant: {pathname: %s, filename: %s}",
			pathKey.pathname, pathKey.filename,
			expectedPathname, expectedFilename)
	} else {
		t.Logf("pathKey correctly transformed: %+v", pathKey)
	}
}

func TestWriteStream(t *testing.T) {
	reader := bytes.NewReader(Data)

	metaPath := "test_metadata.json"
	//defer os.Remove(metaPath) // clean after test

	store := NewStore(StructOpts{
		PathTransformFunc: CASPathTransformFunc,
		Metadata:          NewMetadata(metaPath),
	})

	err := store.WriteStream(Key, reader)
	if err != nil {
		t.Fatalf("WriteStream failed: %v", err)
	}

	t.Logf("WriteStream succeeded for key: %s", Key)
}

func TestReadStream(t *testing.T) {
	metaPath := "test_metadata.json"
	//defer os.Remove(metaPath) // clean metadata file

	store := NewStore(StructOpts{
		PathTransformFunc: CASPathTransformFunc,
		Metadata:          NewMetadata(metaPath),
	})

	for i := 0; i < 50; i++ {

		key := fmt.Sprintf("file%d", i)
		// First write to ensure file exists
		reader := bytes.NewReader(Data)
		err := store.WriteStream(key, reader)
		if err != nil {
			t.Fatalf("WriteStream failed: %v", err)
		}

		// Now test read
		readStream, err := store.ReadStream(key)
		if err != nil {
			t.Fatalf("ReadStream failed: %v", err)
		}

		readData, err := io.ReadAll(readStream)
		if err != nil {
			t.Fatalf("Reading data failed: %v", err)
		}

		if !bytes.Equal(readData, Data) {
			t.Errorf("ReadStream data mismatch:\ngot: %s\nwant: %s", string(readData), string(Data))
		} else {
			t.Logf("ReadStream succeeded and data matched.")
		}
	}

	err := store.TearDown()
	if err != nil {
		t.Error(err)
	}

	// err = store.Remove(key)
	// if err != nil {
	// 	t.Errorf("Removing data failed: %v", err)
	// }
}
