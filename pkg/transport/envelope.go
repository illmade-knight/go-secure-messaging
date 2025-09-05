package transport

// SecureEnvelope is the data structure passed over the network.
// It is understood by both the client and the routing-service.
type SecureEnvelope struct {
	SenderID         string `json:"sender_id"`
	RecipientID      string `json:"recipient_id"`
	EncryptedPayload []byte `json:"encrypted_payload"`
	Signature        []byte `json:"signature"`
}
