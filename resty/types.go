package resty

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MultilineString covers the case where the fields is either string or []string.
type MultilineString string

func (m *MultilineString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*m = MultilineString(str)
		return nil
	}

	var strSlice []string
	if err := json.Unmarshal(data, &strSlice); err == nil {
		*m = MultilineString(strings.Join(strSlice, "\n"))
		return nil
	}

	return fmt.Errorf("failed to unmarshal to MultilineString")
}

func (m MultilineString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(m))
}

func (m MultilineString) String() string {
	return string(m)
}
