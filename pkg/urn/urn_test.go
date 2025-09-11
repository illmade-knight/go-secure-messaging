package urn_test

import (
	"encoding/json"
	"testing"

	"github.com/illmade-knight/routing-service/pkg/urn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		expectedURN   urn.URN
		expectErr     bool
		expectedErrIs error
	}{
		{
			name:  "Valid User URN",
			input: "urn:sm:user:user-123",
			expectedURN: urn.URN{
				Scheme:     "urn",
				Namespace:  "sm",
				EntityType: "user",
				EntityID:   "user-123",
			},
			expectErr: false,
		},
		{
			name:  "Valid Device URN",
			input: "urn:sm:device:uuid-abc-123",
			expectedURN: urn.URN{
				Scheme:     "urn",
				Namespace:  "sm",
				EntityType: "device",
				EntityID:   "uuid-abc-123",
			},
			expectErr: false,
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
				assert.Equal(t, tc.expectedURN, parsedURN)
				// Test the String() method as well
				assert.Equal(t, tc.input, parsedURN.String())
			}
		})
	}
}

func TestJSONMarshaling(t *testing.T) {
	u := urn.URN{
		Scheme:     "urn",
		Namespace:  "sm",
		EntityType: "user",
		EntityID:   "user-123",
	}
	expectedJSON := `"urn:sm:user:user-123"`

	jsonData, err := json.Marshal(u)
	require.NoError(t, err)
	assert.Equal(t, expectedJSON, string(jsonData))
}

func TestJSONUnmarshaling(t *testing.T) {
	testCases := []struct {
		name        string
		jsonInput   string
		expectedURN urn.URN
		expectErr   bool
	}{
		{
			name:      "Unmarshal Full URN",
			jsonInput: `"urn:sm:user:user-123"`,
			expectedURN: urn.URN{
				Scheme:     "urn",
				Namespace:  "sm",
				EntityType: "user",
				EntityID:   "user-123",
			},
			expectErr: false,
		},
		{
			name:      "Unmarshal Legacy UserID (Backward Compatibility)",
			jsonInput: `"legacy-user-456"`,
			expectedURN: urn.URN{
				Scheme:     "urn",
				Namespace:  "sm",
				EntityType: "user",
				EntityID:   "legacy-user-456",
			},
			expectErr: false,
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
				assert.Equal(t, tc.expectedURN, u)
			}
		})
	}
}
