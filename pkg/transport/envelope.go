package transport

type SecureEnvelope struct {
	SenderID    string `json:"senderId"`
	RecipientID string `json:"recipientId"`

	// The AES-encrypted SharedPayload.
	EncryptedData []byte `json:"encryptedData"`

	// The RSA-encrypted AES symmetric key.
	EncryptedSymmetricKey []byte `json:"encryptedSymmetricKey"`

	// The signature of the EncryptedData.
	Signature []byte `json:"signature"`
}
