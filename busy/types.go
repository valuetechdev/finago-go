package busy

import (
	"encoding/json"
	"fmt"
	"net/mail"
)

// Email is an e-mail address that tolerates the empty string.
//
// The schema marks these fields as `format: email`, which oapi-codegen maps to
// [openapi_types.Email]. That type runs every value through
// [mail.ParseAddress] on both marshal and unmarshal, and Busy returns "" for
// users without an address — an invited-but-not-yet-registered user, for
// instance. Parsing "" fails, which breaks decoding of the entire response.
//
// Email keeps the validation for every non-empty value and lets "" through. It
// is wired up in overlay.yaml.
//
// Note that a pointer would not solve this: Busy sends "", not null, so the
// empty string still reaches UnmarshalJSON.
type Email string

// UnmarshalJSON implements [json.Unmarshaler]. An empty address is kept as is,
// anything else must parse as an e-mail address.
func (e *Email) UnmarshalJSON(data []byte) error {
	if e == nil {
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	if s == "" {
		*e = ""
		return nil
	}

	m, err := mail.ParseAddress(s)
	if err != nil {
		return fmt.Errorf("busy: invalid email: %w", err)
	}

	*e = Email(m.Address)
	return nil
}

// MarshalJSON implements [json.Marshaler]. An empty address is encoded as "",
// anything else must parse as an e-mail address.
func (e Email) MarshalJSON() ([]byte, error) {
	if e == "" {
		return []byte(`""`), nil
	}

	m, err := mail.ParseAddress(string(e))
	if err != nil {
		return nil, fmt.Errorf("busy: invalid email: %w", err)
	}

	return json.Marshal(m.Address)
}

// IsEmpty reports whether the address is unset.
func (e Email) IsEmpty() bool {
	return e == ""
}

// String implements the Stringer interface.
func (e Email) String() string {
	return string(e)
}
