// REFACTOR: This file is updated to introduce a validating constructor, New(),
// which prevents the creation of invalid URN structs. The Parse function is
// now a wrapper around this constructor, and the struct fields are unexported
// to enforce the use of the constructor.

package urn

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	// Scheme is the required scheme for all URNs in the system.
	Scheme = "urn"
	// Namespace is the required namespace for all URNs in the system.
	Namespace      = "sm" // Secure Messaging
	urnParts       = 4
	urnDelimiter   = ":"
	EntityTypeUser = "user"
)

var (
	// ErrInvalidFormat is returned when a string or components do not conform
	// to the expected URN structure.
	ErrInvalidFormat = errors.New("invalid URN format")
)

// URN represents a parsed, validated Uniform Resource Name.
// Its fields are unexported to ensure that all instances are created via the
// validating New() constructor.
type URN struct {
	scheme     string
	namespace  string
	entityType string
	entityID   string
}

// New is the constructor for a URN. It validates that the provided entity type
// and ID are not empty, ensuring that no invalid URNs can be created.
func New(entityType, entityID string) (URN, error) {
	if entityType == "" {
		return URN{}, fmt.Errorf("%w: entity type cannot be empty", ErrInvalidFormat)
	}
	if entityID == "" {
		return URN{}, fmt.Errorf("%w: entity ID cannot be empty", ErrInvalidFormat)
	}
	return URN{
		scheme:     Scheme,
		namespace:  Namespace,
		entityType: entityType,
		entityID:   entityID,
	}, nil
}

// Parse converts a raw string into a structured URN, validating its format.
func Parse(s string) (URN, error) {
	parts := strings.Split(s, urnDelimiter)
	if len(parts) != urnParts {
		return URN{}, fmt.Errorf("%w: expected %d parts, but got %d", ErrInvalidFormat, urnParts, len(parts))
	}

	if parts[0] != Scheme {
		return URN{}, fmt.Errorf("%w: invalid scheme '%s', expected '%s'", ErrInvalidFormat, parts[0], Scheme)
	}

	if parts[1] != Namespace {
		return URN{}, fmt.Errorf("%w: invalid namespace '%s', expected '%s'", ErrInvalidFormat, parts[1], Namespace)
	}

	// Delegate final validation to the constructor.
	return New(parts[2], parts[3])
}

// String reassembles the URN into its canonical string representation.
func (u URN) String() string {
	return strings.Join([]string{u.scheme, u.namespace, u.entityType, u.entityID}, urnDelimiter)
}

// IsZero returns true if the URN has not been initialized.
func (u URN) IsZero() bool {
	return u.scheme == "" && u.namespace == "" && u.entityType == "" && u.entityID == ""
}

// MarshalJSON implements the json.Marshaler interface.
func (u URN) MarshalJSON() ([]byte, error) {
	if u.IsZero() {
		// Prevent serialization of an uninitialized URN.
		return []byte("null"), nil
	}
	return json.Marshal(u.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (u *URN) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("URN should be a string, but got %s: %w", string(data), err)
	}

	if strings.HasPrefix(s, Scheme+urnDelimiter) {
		parsedURN, parseErr := Parse(s)
		if parseErr != nil {
			return parseErr
		}
		*u = parsedURN
		return nil
	}

	if s != "" {
		// For legacy IDs, we now use the validating constructor.
		legacyURN, err := New(EntityTypeUser, s)
		if err != nil {
			return err
		}
		*u = legacyURN
		return nil
	}

	return ErrInvalidFormat
}
