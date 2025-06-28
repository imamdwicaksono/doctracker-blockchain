package utils

import (
	"encoding/json"
	"os"
)

func SaveToFile(filename string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, bytes, 0644)
}

func LoadFromFile(filename string, dest any) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, dest)
}
