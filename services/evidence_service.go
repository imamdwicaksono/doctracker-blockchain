package services

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

type EvidenceInfo struct {
	FileName string `json:"file_name"`
	Hash     string `json:"hash"`
	Path     string `json:"path"`
}

func SaveBase64Evidence(trackerID string, checkpointID string, base64Str *string) (*multipart.FileHeader, error) {
	// Simpan base64Str ke file sementara
	tempFile, err := os.CreateTemp("storage/evidence", fmt.Sprintf("%s_%s_*.txt", trackerID, checkpointID))
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	if _, err := tempFile.WriteString(*base64Str); err != nil {
		return nil, err
	}

	return &multipart.FileHeader{
		Filename: tempFile.Name(),
		Size:     int64(len(*base64Str)),
	}, nil
}

func SaveEvidenceFile(trackerID, checkpointID string, file *string) (EvidenceInfo, error) {

	// Simpan file
	savePath := fmt.Sprintf("storage/evidence/%s_%s.jpg", trackerID, checkpointID)
	dst, err := os.Create(savePath)
	if err != nil {
		return EvidenceInfo{}, err
	}
	defer dst.Close()

	hash := sha256.New()
	writer := io.MultiWriter(dst, hash)

	decoded, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, strings.NewReader(*file)))
	if err != nil {
		return EvidenceInfo{}, err
	}
	if _, err := writer.Write(decoded); err != nil {
		return EvidenceInfo{}, err
	}

	return EvidenceInfo{
		FileName: fmt.Sprintf("%s_%s", trackerID, checkpointID),
		Hash:     fmt.Sprintf("%x", hash.Sum(nil)),
		Path:     savePath,
	}, nil
}
