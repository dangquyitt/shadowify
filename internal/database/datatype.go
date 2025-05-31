package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type JSONType[T any] struct {
	Data T
}

// Value return json value, implement driver.Valuer interface
func (j JSONType[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

// Scan scan value into JSONType[T], implements sql.Scanner interface
func (j *JSONType[T]) Scan(value interface{}) error {
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}
	return json.Unmarshal(bytes, &j.Data)
}

// MarshalJSON to output non base64 encoded []byte
func (j JSONType[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}

// UnmarshalJSON to deserialize []byte
func (j *JSONType[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.Data)
}
