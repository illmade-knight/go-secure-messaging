package transport_test

import (
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDigestConversions(t *testing.T) {
	convURN1, _ := urn.Parse("urn:sm:convo:convo-1")
	convURN2, _ := urn.Parse("urn:sm:convo:convo-2")

	nativeDigest := &transport.EncryptedDigest{
		Items: []*transport.EncryptedDigestItem{
			{
				ConversationID:        convURN1,
				EncryptedSnippet:      []byte("snippet-1"),
				EncryptedSymmetricKey: []byte("key-1"),
			},
			{
				ConversationID:        convURN2,
				EncryptedSnippet:      []byte("snippet-2"),
				EncryptedSymmetricKey: []byte("key-2"),
			},
		},
	}

	t.Run("Symmetry Test", func(t *testing.T) {
		// Act
		protoDigest := transport.DigestToProto(nativeDigest)
		roundTrippedDigest, err := transport.DigestFromProto(protoDigest)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, nativeDigest, roundTrippedDigest)
	})

	t.Run("FromProto Error Handling", func(t *testing.T) {
		testCases := []struct {
			name          string
			proto         *transport.EncryptedDigestPb
			expectedError string
		}{
			{
				name: "Invalid ConversationID in item",
				proto: &transport.EncryptedDigestPb{
					Items: []*transport.EncryptedDigestItemPb{
						{ConversationId: convURN1.String()},
						{ConversationId: "not-a-valid-urn"},
					},
				},
				expectedError: "failed to parse conversation id",
			},
			{
				name: "Empty ConversationID in item",
				proto: &transport.EncryptedDigestPb{
					Items: []*transport.EncryptedDigestItemPb{
						{ConversationId: ""},
					},
				},
				expectedError: "failed to parse conversation id",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := transport.DigestFromProto(tc.proto)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})

	t.Run("Nil Handling", func(t *testing.T) {
		// Test ToProto with nil input
		assert.Nil(t, transport.DigestToProto(nil))

		// Test FromProto with nil input
		native, err := transport.DigestFromProto(nil)
		assert.NoError(t, err)
		assert.Nil(t, native)

		// Test with nil items in slices
		nativeWithNil := &transport.EncryptedDigest{
			Items: []*transport.EncryptedDigestItem{
				nil,
				{
					ConversationID:        convURN1,
					EncryptedSnippet:      []byte("snippet-1"),
					EncryptedSymmetricKey: []byte("key-1"),
				},
				nil,
			},
		}
		protoWithNil := transport.DigestToProto(nativeWithNil)
		assert.Len(t, protoWithNil.Items, 3)
		assert.Nil(t, protoWithNil.Items[0])
		assert.NotNil(t, protoWithNil.Items[1])
		assert.Nil(t, protoWithNil.Items[2])

		roundTripped, err := transport.DigestFromProto(protoWithNil)
		require.NoError(t, err)
		assert.Equal(t, nativeWithNil, roundTripped)
	})
}
