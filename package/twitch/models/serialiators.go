package twitchmodels

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func (u *UserAccessTokens) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, u)
}

func (u *UserAccessTokens) Value() (driver.Value, error) {
	return json.Marshal(u)
}
