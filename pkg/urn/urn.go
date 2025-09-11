// Package urn provides a type-safe implementation for handling Uniform
// Resource Names (URNs) within the messaging ecosystem. It ensures that all
// entity identifiers are validated and consistently structured.
package urn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	urnScheme      = "urn"
	urnNamespace   = "sm" // Secure Messaging
	urnParts       = 4
	urnDelimiter   = ":"
	EntityTypeUser = "user"
)

var (
	// ErrInvalidFormat is returned when a string does not conform to the
	// expected URN structure.
	ErrInvalidFormat = errors.New("invalid URN format")
)

// URN represents a parsed, validated Uniform Resource Name.
// It is the standard identifier for all addressable entities in the system.
type URN struct {
	Scheme     string
	Namespace  string
	EntityType string
	EntityID   string
}

// Parse converts a raw string into a structured URN, validating its format.
// It expects the format "urn:sm:<type>:<id>".
func Parse(s string) (URN, error) {
	parts := strings.Split(s, urnDelimiter)
	if len(parts) != urnParts {
		return URN{}, fmt.Errorf("%w: expected %d parts, but got %d", ErrInvalidFormat, urnParts, len(parts))
	}

	if parts[0] != urnScheme {
		return URN{}, fmt.Errorf("%w: invalid scheme '%s', expected '%s'", ErrInvalidFormat, parts[0], urnScheme)
	}

	if parts[1] != urnNamespace {
		return URN{}, fmt.Errorf("%w: invalid namespace '%s', expected '%s'", ErrInvalidFormat, parts[1], urnNamespace)
	}

	if parts[2] == "" {
		return URN{}, fmt.Errorf("%w: entity type cannot be empty", ErrInvalidFormat)
	}

	if parts[3] == "" {
		return URN{}, fmt.Errorf("%w: entity ID cannot be empty", ErrInvalidFormat)
	}

	return URN{
		Scheme:     parts[0],
		Namespace:  parts[1],
		EntityType: parts[2],
		EntityID:   parts[3],
	}, nil
}

// String reassembles the URN into its canonical string representation.
func (u URN) String() string {
	return strings.Join([]string{u.Scheme, u.Namespace, u.EntityType, u.EntityID}, urnDelimiter)
}

// IsZero returns true if the URN has not been initialized.
func (u URN) IsZero() bool {
	return u.Scheme == "" && u.Namespace == "" && u.EntityType == "" && u.EntityID == ""
}

// MarshalJSON implements the json.Marshaler interface.
func (u URN) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This is the core of our backward-compatibility strategy. It can unmarshal
// both a full URN string ("urn:sm:user:123") and a legacy userID string ("123"),
// defaulting the legacy ID to a user entity.
func (u *URN) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("URN should be a string, but got %s: %w", string(data), err)
	}

	// REFACTOR: If the string looks like a URN, it MUST be valid.
	if strings.HasPrefix(s, urnScheme+urnDelimiter) {
		parsedURN, parseErr := Parse(s)
		if parseErr != nil {
			return parseErr // It tried to be a URN and failed.
		}
		*u = parsedURN
		return nil
	}

	// If it's not empty and doesn't look like a URN, treat it as a legacy ID.
	if s != "" {
		*u = URN{
			Scheme:     urnScheme,
			Namespace:  urnNamespace,
			EntityType: EntityTypeUser,
			EntityID:   s,
		}
		return nil
	}

	// The string was empty, which is invalid.
	return ErrInvalidFormat
}
