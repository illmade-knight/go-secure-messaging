package transport

import (
	"fmt"

	envelopev1 "github.com/illmade-knight/go-action-intention-protos/gen/go/proto/sm/v1"
	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
)

// SecureEnvelopePb Re-export the Protobuf type with a convenient alias.
// Now, any project that imports this 'transport' package can use
// 'transport.SecureEnvelopePb' without needing to import the long proto path.
type SecureEnvelopePb = envelopev1.SecureEnvelopePb

type SecureEnvelopeListPb = envelopev1.SecureEnvelopeListPb

type SecureEnvelope struct {
	SenderID    urn.URN `json:"senderId"`
	RecipientID urn.URN `json:"recipientId"`
	MessageID   string  `json:"messageId"`

	// The AES-encrypted SharedPayload.
	EncryptedData []byte `json:"encryptedData"`

	// The RSA-encrypted AES symmetric key.
	EncryptedSymmetricKey []byte `json:"encryptedSymmetricKey"`

	// The signature of the EncryptedData.
	Signature []byte `json:"signature"`
}

func ToProto(nativeEnvelope *SecureEnvelope) *SecureEnvelopePb {
	if nativeEnvelope == nil {
		return nil
	}

	return &envelopev1.SecureEnvelopePb{
		SenderId:              nativeEnvelope.SenderID.String(),
		RecipientId:           nativeEnvelope.RecipientID.String(),
		MessageId:             nativeEnvelope.MessageID,
		EncryptedData:         nativeEnvelope.EncryptedData,
		EncryptedSymmetricKey: nativeEnvelope.EncryptedSymmetricKey,
		Signature:             nativeEnvelope.Signature,
	}
}

func FromProto(protoEnvelope *SecureEnvelopePb) (*SecureEnvelope, error) {
	if protoEnvelope == nil {
		return nil, nil
	}

	senderURN, err := urn.Parse(protoEnvelope.GetSenderId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse sender URN: %w", err)
	}

	recipientURN, err := urn.Parse(protoEnvelope.GetRecipientId())
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient URN: %w", err)
	}

	return &SecureEnvelope{
		SenderID:              senderURN,
		RecipientID:           recipientURN,
		MessageID:             protoEnvelope.GetMessageId(),
		EncryptedData:         protoEnvelope.GetEncryptedData(),
		EncryptedSymmetricKey: protoEnvelope.GetEncryptedSymmetricKey(),
		Signature:             protoEnvelope.GetSignature(),
	}, nil
}
