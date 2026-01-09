package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Checkpoints []Checkpoint

func (c Checkpoints) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Checkpoints) Scan(value interface{}) error {
	if value == nil {
		*c = Checkpoints{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Checkpoints")
	}

	return json.Unmarshal(bytes, c)
}
