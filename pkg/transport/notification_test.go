package transport_test

import (
	"testing"

	"github.com/illmade-knight/go-secure-messaging/pkg/transport"
	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	"github.com/stretchr/testify/require"
)

func TestNotificationRequestConversions(t *testing.T) {
	recipientURN, err := urn.New("sm", "user", "recipient-456")
	require.NoError(t, err)

	nativeReq := &transport.NotificationRequest{
		RecipientID: recipientURN,
		Tokens: []transport.DeviceToken{
			{Token: "token-1", Platform: "apns"},
			{Token: "token-2", Platform: "fcm"},
		},
		Content: transport.NotificationContent{
			Title: "New Message",
			Body:  "You have a new secure message.",
			Sound: "default",
		},
		DataPayload: map[string]string{
			"message_id": "msg-789",
		},
	}

	t.Run("ToProto and FromProto Symmetry", func(t *testing.T) {
		// 1. Convert native Go struct to Protobuf message
		protoReq := transport.NotificationRequestToProto(nativeReq)

		// Assert that the Protobuf message has the correct values
		require.Equal(t, "urn:sm:user:recipient-456", protoReq.GetRecipientId())
		require.Len(t, protoReq.GetTokens(), 2)
		require.Equal(t, "token-1", protoReq.GetTokens()[0].GetToken())
		require.Equal(t, "apns", protoReq.GetTokens()[0].GetPlatform())
		require.Equal(t, "New Message", protoReq.GetContent().GetTitle())
		require.Equal(t, "msg-789", protoReq.GetDataPayload()["message_id"])

		// 2. Convert the Protobuf message back to the native Go struct
		convertedNative, err := transport.NotificationRequestFromProto(protoReq)
		require.NoError(t, err)

		// Assert that the round-trip conversion results in the original struct
		require.Equal(t, nativeReq, convertedNative)
	})

	t.Run("FromProto with Invalid URN", func(t *testing.T) {
		protoReq := transport.NotificationRequestToProto(nativeReq)
		protoReq.RecipientId = "invalid-urn" // Introduce an error

		_, err := transport.NotificationRequestFromProto(protoReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to parse recipient URN")
	})

	t.Run("Handling Nil Inputs", func(t *testing.T) {
		// ToProto should handle a nil input gracefully
		require.Nil(t, transport.NotificationRequestToProto(nil))

		// FromProto should also handle a nil input gracefully
		converted, err := transport.NotificationRequestFromProto(nil)
		require.NoError(t, err)
		require.Nil(t, converted)
	})
}
