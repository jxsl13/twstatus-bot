package db

import (
	"fmt"
	"strings"
)

const (
	SSLModeDisable    SSLMode = "disable"
	SSLModeRequire    SSLMode = "require"
	SSLModeVerifyCA   SSLMode = "verify-ca"
	SSLModeVerifyFull SSLMode = "verify-full"
	SSLModePreferred  SSLMode = "prefer"
)

type SSLMode string

// implement fmt.Stringer
func (s SSLMode) String() string {
	return string(s)
}

// implemets encoding.TextMarshaler and TextUnmarshaler
func (s *SSLMode) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *SSLMode) UnmarshalText(text []byte) error {
	switch SSLMode(strings.ToLower(string(text))) {
	case SSLModeDisable:
		*s = SSLModeDisable
	case SSLModeRequire:
		*s = SSLModeRequire
	case SSLModeVerifyCA:
		*s = SSLModeVerifyCA
	case SSLModeVerifyFull:
		*s = SSLModeVerifyFull
	case SSLModePreferred:
		*s = SSLModePreferred
	default:
		return fmt.Errorf("invalid ssl mode: %s", string(text))
	}
	return nil
}
