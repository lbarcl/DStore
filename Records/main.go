package Records

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/google/uuid"
)

func NewFile(name string, size int64, path string, totalPart int) (*File, error) {
	id := uuid.New()
	f := File{
		ID:            id.String(),
		Name:          name,
		Size:          size,
		Paused:        false,
		LocalPath:     path,
		TotalParts:    totalPart,
		RecordedParts: 0,
	}
	err := saveToFile(filepath.Join("files", f.ID, "file.json"), f)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func saveToFile(filePath string, data interface{}) error {
	// Ensure the directory exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %v", err)
	}

	file.Close()
	return nil
}

func (f *File) Destroy() error {
	filePath := filepath.Join("files", f.ID)
	return os.Remove(filePath)
}

func (f *File) AddPart(id string, messageID string, sequence int, hash string, size int64) error {
	newPart := part{
		ID:        id,
		MessageID: messageID,
		Sequence:  sequence,
		Hash:      hash,
		Size:      size,
	}
	f.RecordedParts += 1
	saveToFile(filepath.Join("files", f.ID, "file.json"), f)
	return saveToFile(filepath.Join("files", f.ID, "parts", newPart.ID+".json"), newPart)
}

func (f *File) GetParts() ([]part, error) {
	partsDir := filepath.Join("files", f.ID, "parts")
	dirEntry, err := os.ReadDir(partsDir)
	if err != nil {
		return nil, err
	}

	var parts []part
	for _, de := range dirEntry {
		if de.IsDir() {
			// Skip directories, assuming only part files are present
			continue
		}

		partFilePath := filepath.Join(partsDir, de.Name())
		partData, err := ioutil.ReadFile(partFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read part file %s: %v", de.Name(), err)
		}

		var part part
		if err := json.Unmarshal(partData, &part); err != nil {
			return nil, fmt.Errorf("failed to unmarshal part JSON: %v", err)
		}

		parts = append(parts, part)
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].Sequence < parts[j].Sequence
	})

	return parts, nil
}

func (f *File) Pause() {
	f.Paused = true
	saveToFile(filepath.Join("files", f.ID), f)
}

func (f *File) Unpause() {
	f.Paused = true
	saveToFile(filepath.Join("files", f.ID), f)
}

func GetAllFiles() ([]File, error) {
	filesDir := "files"
	dirEntry, err := os.ReadDir(filesDir)
	if err != nil {
		return nil, err
	}

	var files []File
	for _, de := range dirEntry {
		if de.IsDir() {
			fileFilePath := filepath.Join(filesDir, de.Name(), "file.json")
			fileData, err := ioutil.ReadFile(fileFilePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %v", de.Name(), err)
			}

			var f File
			if err := json.Unmarshal(fileData, &f); err != nil {
				return nil, fmt.Errorf("failed to unmarshal file JSON: %v", err)
			}

			files = append(files, f)
		}

	}

	return files, nil
}

func GetFile(id string) (*File, error) {
	files, err := GetAllFiles()
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.ID == id {
			return &f, nil
		}
	}

	return nil, fmt.Errorf("we cannot find file with %s", id)
}

type File struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	Paused        bool   `json:"paused"`
	LocalPath     string `json:"localPath"`
	TotalParts    int    `json:"totalParts"`
	RecordedParts int    `json:"recordedParts"`
}

type part struct {
	ID        string `json:"id"`
	MessageID string `json:"message_id"`
	Sequence  int    `json:"sequence"` // This partid is for the sequence of the parts
	Hash      string `json:"hash"`
	Size      int64  `json:"size"`
}
