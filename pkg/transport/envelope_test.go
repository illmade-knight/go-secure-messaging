package transport_test

import (
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvelopeConversions(t *testing.T) {
	senderURN, _ := urn.Parse("urn:sm:user:user-alice")
	recipientURN, _ := urn.Parse("urn:sm:user:user-bob")
	convURN, _ := urn.Parse("urn:sm:convo:alice-bob-123")
	snippet := []byte("this is a test snippet")

	nativeEnvelope := &transport.SecureEnvelope{
		MessageID:             "msg-123",
		SenderID:              senderURN,
		RecipientID:           recipientURN,
		ConversationID:        convURN,
		EncryptedData:         []byte("encrypted-data"),
		EncryptedSymmetricKey: []byte("encrypted-key"),
		Signature:             []byte("signature-data"),
		EncryptedSnippet:      snippet,
	}

	t.Run("Symmetry Test", func(t *testing.T) {
		// Act
		protoEnvelope := transport.ToProto(nativeEnvelope)
		roundTrippedEnvelope, err := transport.FromProto(protoEnvelope)
		require.NoError(t, err)

		// Assert
		assert.Equal(t, nativeEnvelope, roundTrippedEnvelope)
		// Explicitly check new byte slice
		assert.Equal(t, nativeEnvelope.EncryptedSnippet, roundTrippedEnvelope.EncryptedSnippet)
	})

	t.Run("FromProto Error Handling", func(t *testing.T) {
		testCases := []struct {
			name          string
			proto         *transport.SecureEnvelopePb
			expectedError string
		}{
			{
				name: "Invalid SenderID",
				proto: &transport.SecureEnvelopePb{
					SenderId:    "not-a-valid-urn",
					RecipientId: recipientURN.String(),
				},
				expectedError: "failed to parse sender id",
			},
			{
				name: "Invalid RecipientID",
				proto: &transport.SecureEnvelopePb{
					SenderId:    senderURN.String(),
					RecipientId: "not-a-valid-urn",
				},
				expectedError: "failed to parse recipient id",
			},
			{
				name: "Invalid ConversationID",
				proto: &transport.SecureEnvelopePb{
					SenderId:       senderURN.String(),
					RecipientId:    recipientURN.String(),
					ConversationId: "not-a-valid-urn",
				},
				expectedError: "failed to parse conversation id",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := transport.FromProto(tc.proto)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			})
		}
	})

	t.Run("Nil Handling", func(t *testing.T) {
		assert.Nil(t, transport.ToProto(nil))
		native, err := transport.FromProto(nil)
		assert.NoError(t, err)
		assert.Nil(t, native)
	})
}
