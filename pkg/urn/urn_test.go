package urn_test

import (
	"encoding/json"
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewURN validates the behavior of the new constructor function.
func TestNewURN(t *testing.T) {
	t.Run("Valid URN", func(t *testing.T) {
		u, err := urn.New("user", "user-123", urn.SecureMessaging)
		require.NoError(t, err)
		assert.Equal(t, "urn:sm:user:user-123", u.String())
	})

	t.Run("Empty Entity Type", func(t *testing.T) {
		_, err := urn.New("", "user-123", urn.SecureMessaging)
		require.Error(t, err)
		assert.ErrorIs(t, err, urn.ErrInvalidFormat)
	})

	t.Run("Empty Entity ID", func(t *testing.T) {
		_, err := urn.New("user", "", "")
		require.Error(t, err)
		assert.ErrorIs(t, err, urn.ErrInvalidFormat)
	})
}

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedURN   string // We check the string representation
		expectErr     bool
		expectedErrIs error
	}{
		{
			name:        "Valid User URN",
			input:       "urn:sm:user:user-123",
			expectedURN: "urn:sm:user:user-123",
			expectErr:   false,
		},
		{
			name:        "Valid Device URN",
			input:       "urn:sm:device:uuid-abc-123",
			expectedURN: "urn:sm:device:uuid-abc-123",
			expectErr:   false,
		},
		{
			name:          "Invalid Scheme",
			input:         "foo:sm:user:user-123",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Invalid Namespace",
			input:         "urn:foo:user:user-123",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Too Few Parts",
			input:         "urn:sm:user",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Too Many Parts",
			input:         "urn:sm:user:user-123:extra",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Empty Entity Type",
			input:         "urn:sm::user-123",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Empty Entity ID",
			input:         "urn:sm:user:",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
		{
			name:          "Empty Input",
			input:         "",
			expectErr:     true,
			expectedErrIs: urn.ErrInvalidFormat,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parsedURN, err := urn.Parse(tc.input)
			if tc.expectErr {
				require.Error(t, err)
				if tc.expectedErrIs != nil {
					assert.ErrorIs(t, err, tc.expectedErrIs)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedURN, parsedURN.String())
			}
		})
	}
}

func TestJSONMarshaling(t *testing.T) {
	u, err := urn.New("user", "user-123", urn.SecureMessaging)
	require.NoError(t, err)
	expectedJSON := `"urn:sm:user:user-123"`

	jsonData, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Equal(t, expectedJSON, string(jsonData))

	// Test marshaling a zero-value URN
	var zeroURN urn.URN
	zeroJSON, err := json.Marshal(zeroURN)
	require.NoError(t, err)
	assert.Equal(t, "null", string(zeroJSON))
}

func TestJSONUnmarshaling(t *testing.T) {
	testCases := []struct {
		name        string
		jsonInput   string
		expectedURN string
		expectErr   bool
	}{
		{
			name:        "Unmarshal Full URN",
			jsonInput:   `"urn:sm:user:user-123"`,
			expectedURN: "urn:sm:user:user-123",
			expectErr:   false,
		},
		{
			name:        "Unmarshal Legacy UserID (Backward Compatibility)",
			jsonInput:   `"legacy-user-456"`,
			expectedURN: "urn:sm:user:legacy-user-456",
			expectErr:   false,
		},
		{
			name:      "Unmarshal Invalid URN",
			jsonInput: `"urn:sm:user"`, // Too short
			expectErr: true,
		},
		{
			name:      "Unmarshal Empty String",
			jsonInput: `""`,
			expectErr: true,
		},
		{
			name:      "Unmarshal Non-String Type",
			jsonInput: `123`,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var u urn.URN
			err := json.Unmarshal([]byte(tc.jsonInput), &u)

			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedURN, u.String())
			}
		})
	}
}
