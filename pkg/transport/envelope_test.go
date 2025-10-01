package transport_test

import (
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/stretchr/testify/require"
)

func TestSecureEnvelopeConversions(t *testing.T) {
	senderURN, err := urn.New("sm", "user", "sender-123")
	require.NoError(t, err)
	recipientURN, err := urn.New("sm", "user", "recipient-456")
	require.NoError(t, err)

	nativeEnvelope := &transport.SecureEnvelope{
		SenderID:              senderURN,
		RecipientID:           recipientURN,
		MessageID:             "msg-789",
		EncryptedData:         []byte("encrypted-data"),
		EncryptedSymmetricKey: []byte("encrypted-key"),
		Signature:             []byte("signature"),
	}

	t.Run("ToProto and FromProto Symmetry", func(t *testing.T) {
		// 1. Convert native Go struct to Protobuf message
		protoEnvelope := transport.ToProto(nativeEnvelope)

		// Assert that the Protobuf message has the correct values
		require.Equal(t, "urn:sm:user:sender-123", protoEnvelope.GetSenderId())
		require.Equal(t, "urn:sm:user:recipient-456", protoEnvelope.GetRecipientId())
		require.Equal(t, "msg-789", protoEnvelope.GetMessageId())
		require.Equal(t, []byte("encrypted-data"), protoEnvelope.GetEncryptedData())
		require.Equal(t, []byte("encrypted-key"), protoEnvelope.GetEncryptedSymmetricKey())
		require.Equal(t, []byte("signature"), protoEnvelope.GetSignature())

		// 2. Convert the Protobuf message back to the native Go struct
		convertedNative, err := transport.FromProto(protoEnvelope)
		require.NoError(t, err)

		// Assert that the round-trip conversion results in the original struct
		require.Equal(t, nativeEnvelope, convertedNative)
	})

	t.Run("FromProto with Invalid URN", func(t *testing.T) {
		protoEnvelope := transport.ToProto(nativeEnvelope)
		protoEnvelope.SenderId = "invalid-urn" // Introduce an error

		_, err := transport.FromProto(protoEnvelope)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse sender URN")

		// Reset for the next test
		protoEnvelope.SenderId = "urn:sm:user:sender-123"
		protoEnvelope.RecipientId = "invalid-urn"

		_, err = transport.FromProto(protoEnvelope)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse recipient URN")
	})

	t.Run("Handling Nil Inputs", func(t *testing.T) {
		// ToProto should handle a nil input gracefully
		require.Nil(t, transport.ToProto(nil))

		// FromProto should also handle a nil input gracefully
		converted, err := transport.FromProto(nil)
		require.NoError(t, err)
		require.Nil(t, converted)
	})
}
