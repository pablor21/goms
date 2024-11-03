package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type Metadata map[string]interface{} // @name Metadata

func (m Metadata) GormDataType() string {
	return "json"
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (j *Metadata) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := make(map[string]interface{})
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Value return json value, implement driver.Valuer interface
func (j Metadata) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j Metadata) ToMap() map[string]interface{} {
	return j
}

func (j Metadata) GetString(key string) string {
	if val, ok := j[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (j Metadata) GetStringOrDefault(key string, def string) string {
	if val, ok := j[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return def
}

func (j Metadata) Set(key string, value interface{}) {
	j[key] = value
}
