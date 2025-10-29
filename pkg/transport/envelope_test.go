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
	groupURN, _ := urn.Parse("urn:sm:group:group-1")
	conversationURN, _ := urn.Parse("urn:sm:convo:convo-123")
	snippet := []byte("this is a test snippet")

	nativeEnvelope := &transport.SecureEnvelope{
		MessageID:             "msg-123",
		SenderID:              senderURN,
		RecipientID:           recipientURN,
		GroupID:               groupURN,
		ConversationID:        conversationURN,
		EncryptedData:         []byte("encrypted-data"),
		EncryptedSymmetricKey: []byte("encrypted-key"),
		Signature:             []byte("signature-data"),
		EncryptedSnippet:      snippet,
	}

	t.Run("Symmetry Test", func(t *testing.T) {
		// Act
		protoEnvelope := transport.ToProto(nativeEnvelope)

		// Assert: Check proto conversion
		assert.Equal(t, nativeEnvelope.MessageID, protoEnvelope.MessageId)
		assert.Equal(t, nativeEnvelope.SenderID.String(), protoEnvelope.SenderId)
		assert.Equal(t, nativeEnvelope.RecipientID.String(), protoEnvelope.RecipientId)
		assert.Equal(t, nativeEnvelope.GroupID.String(), protoEnvelope.GroupId)
		assert.Equal(t, nativeEnvelope.ConversationID.String(), protoEnvelope.ConversationId)
		assert.Equal(t, nativeEnvelope.EncryptedSymmetricKey, protoEnvelope.EncryptedSymmetricKey)

		// Act: Round-trip
		roundTrippedEnvelope, err := transport.FromProto(protoEnvelope)
		require.NoError(t, err)

		// Assert: Check round-trip
		assert.Equal(t, nativeEnvelope, roundTrippedEnvelope)
		assert.Equal(t, nativeEnvelope.EncryptedSnippet, roundTrippedEnvelope.EncryptedSnippet)
	})

	t.Run("FromProto Error Handling", func(t *testing.T) {
		// Base valid proto for modification
		baseProto := func() *transport.SecureEnvelopePb {
			return &transport.SecureEnvelopePb{
				SenderId:       senderURN.String(),
				RecipientId:    recipientURN.String(),
				GroupId:        groupURN.String(),
				ConversationId: conversationURN.String(),
			}
		}

		testCases := []struct {
			name          string
			modifier      func(pb *transport.SecureEnvelopePb)
			expectedError string
		}{
			{
				name: "Invalid SenderID",
				modifier: func(pb *transport.SecureEnvelopePb) {
					pb.SenderId = "not-a-valid-urn"
				},
				expectedError: "failed to parse sender id",
			},
			{
				name: "Invalid RecipientID",
				modifier: func(pb *transport.SecureEnvelopePb) {
					pb.RecipientId = "not-a-valid-urn"
				},
				expectedError: "failed to parse recipient id",
			},
			{
				name: "Invalid GroupID",
				modifier: func(pb *transport.SecureEnvelopePb) {
					pb.GroupId = "not-a-valid-urn"
				},
				expectedError: "failed to parse group id",
			},
			{
				name: "Invalid ConversationID",
				modifier: func(pb *transport.SecureEnvelopePb) {
					pb.ConversationId = "not-a-valid-urn"
				},
				expectedError: "failed to parse conversation id",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				proto := baseProto()
				tc.modifier(proto)
				_, err := transport.FromProto(proto)
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
