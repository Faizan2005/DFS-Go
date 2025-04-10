package main

import (
	// "crypto/sha1"
	// "encoding/hex"
	"encoding/json"
	"log"
	"os"
)

type Metadata struct {
	Key      map[string]string // key = original key, value = hashed file path
	Metapath string
}

func NewMetadata(Metapath string) *Metadata {
	m := &Metadata{
		Key:      make(map[string]string),
		Metapath: Metapath,
	}
	m.load()
	return m
}

func (m *Metadata) Set(key, path string) error {
	// hashedPath := hashString(path) // hash the path for security
	m.Key[key] = path

	if err := m.save(); err != nil {
		return err
	}
	log.Printf("Metadata saved: %s → %s", key, path)
	return nil
}

func (m *Metadata) Get(key string) (string, bool) {
	filePath, ok := m.Key[key]
	if ok {
		log.Printf("Metadata loaded: %s → %s", key, filePath)
	}
	return filePath, ok
}

func (m *Metadata) save() error {
	data, err := json.MarshalIndent(m.Key, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(m.Metapath, data, 0644)
	if err != nil {
		return err
	}
	log.Println("Metadata file saved to disk")
	return nil
}

func (m *Metadata) load() {
	data, err := os.ReadFile(m.Metapath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Metadata file not found. Starting fresh.")
		} else {
			log.Printf("Error reading metadata file: %v", err)
		}
		return
	}
	if err := json.Unmarshal(data, &m.Key); err != nil {
		log.Printf("Failed to parse metadata: %v", err)
	}
	log.Println("Metadata loaded successfully.")
}

// func hashString(s string) string {
// 	hash := sha1.Sum([]byte(s))
// 	return hex.EncodeToString(hash[:])
// }
