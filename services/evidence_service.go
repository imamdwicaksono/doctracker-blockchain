package services

import (
	"crypto/sha256"
	"doc-tracker/storage"
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

// SaveEvidenceBase64ToS3 menyimpan evidence dari base64 ke Storj (S3)
func SaveEvidenceBase64ToS3(trackerID, checkpointID string, base64Str *string) (EvidenceInfo, error) {
	if base64Str == nil || *base64Str == "" {
		return EvidenceInfo{}, fmt.Errorf("base64 string is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(*base64Str)
	if err != nil {
		return EvidenceInfo{}, fmt.Errorf("invalid base64: %v", err)
	}

	hash := sha256.Sum256(decoded)

	fileName := fmt.Sprintf("%s_%s.jpg", trackerID, checkpointID)
	s3Key := fmt.Sprintf("evidence/%s", fileName)

	// ✅ FIXED: call global instance
	err = storage.S3.UploadS3File(s3Key, decoded)
	if err != nil {
		return EvidenceInfo{}, fmt.Errorf("upload failed: %v", err)
	}

	return EvidenceInfo{
		FileName: fileName,
		Hash:     fmt.Sprintf("%x", hash[:]),
		Path:     s3Key,
	}, nil
}

func LoadEvidenceFromS3(fileName string) ([]byte, error) {
	// ✅ FIXED: call global instance
	data, err := storage.S3.DownloadS3File(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %v", err)
	}

	return data, nil
}

func SaveBase64EvidenceLocal(trackerID string, checkpointID string, base64Str *string) (*multipart.FileHeader, error) {
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

func SaveEvidenceFileLocal(trackerID, checkpointID string, file *string) (EvidenceInfo, error) {

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
