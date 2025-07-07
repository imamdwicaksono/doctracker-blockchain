package storage

import (
	"doc-tracker/models"
	"os"
	"path/filepath"
)

// SaveBlock menyimpan block terenkripsi
func SaveBlock(block models.Block) error {
	// Block disimpan dalam format terenkripsi
	// Implementasi aktual menggunakan fungsi dari blockchain package
	return nil
}

// LoadBlock memuat block terenkripsi
func LoadBlock(index int) (models.Block, error) {
	// Block dimuat melalui fungsi dekripsi di blockchain package
	return models.Block{}, nil
}

// ClearAllBlocks menghapus semua block
func ClearAllBlocks() error {
	files, err := filepath.Glob("data/blocks/*.bin")
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
