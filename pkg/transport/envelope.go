package transport

import "github.com/illmade-knight/go-secure-messaging/pkg/urn"

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
