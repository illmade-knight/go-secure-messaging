package transport

type SecureEnvelope struct {
	SenderID    string `json:"sender_id"`
	RecipientID string `json:"recipient_id"`

	// The AES-encrypted SharedPayload.
	EncryptedData []byte `json:"encrypted_data"`

	// The RSA-encrypted AES symmetric key.
	EncryptedSymmetricKey []byte `json:"encrypted_symmetric_key"`

	// The signature of the EncryptedData.
	Signature []byte `json:"signature"`
}
