package twitchmodels

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func (u *UserAccessToken) Scan(value interface{}) error {
	if value == nil {
		*u = UserAccessToken{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return json.Unmarshal(bytes, u)
}

func (u *UserAccessToken) Value() (driver.Value, error) {
	if u == nil {
		return nil, nil
	}
	return json.Marshal(u)
}
