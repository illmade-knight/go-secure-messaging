package transport

import (
	"fmt"

	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	smv1 "github.com/tinywideclouds/go-action-intention-protos/src/action_intention/envelope/v1"
)

// SecureEnvelopePb Re-export the Protobuf type for external use.
type SecureEnvelopePb = smv1.SecureEnvelopePb

// SecureEnvelope is the canonical, idiomatic Go struct for a message.
// It uses URN types for identifiers.
type SecureEnvelope struct {
	MessageID             string
	SenderID              urn.URN
	RecipientID           urn.URN
	ConversationID        urn.URN // ADDED
	EncryptedData         []byte
	EncryptedSymmetricKey []byte
	Signature             []byte
	EncryptedSnippet      []byte // ADDED
}

// ToProto converts the idiomatic Go struct into its Protobuf representation.
func ToProto(native *SecureEnvelope) *SecureEnvelopePb {
	if native == nil {
		return nil
	}
	return &SecureEnvelopePb{
		MessageId:             native.MessageID,
		SenderId:              native.SenderID.String(),
		RecipientId:           native.RecipientID.String(),
		EncryptedData:         native.EncryptedData,
		EncryptedSymmetricKey: native.EncryptedSymmetricKey,
		Signature:             native.Signature,
		ConversationId:        native.ConversationID.String(), // ADDED
		EncryptedSnippet:      native.EncryptedSnippet,        // ADDED
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
// It validates that the string identifiers are valid URNs.
func FromProto(proto *SecureEnvelopePb) (*SecureEnvelope, error) {
	if proto == nil {
		return nil, nil
	}

	senderID, err := urn.Parse(proto.SenderId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse sender id: %w", err)
	}

	recipientID, err := urn.Parse(proto.RecipientId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse recipient id: %w", err)
	}

	// ADDED: Parse the ConversationID
	convID, err := urn.Parse(proto.ConversationId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conversation id: %w", err)
	}

	return &SecureEnvelope{
		MessageID:             proto.MessageId,
		SenderID:              senderID,
		RecipientID:           recipientID,
		ConversationID:        convID, // ADDED
		EncryptedData:         proto.EncryptedData,
		EncryptedSymmetricKey: proto.EncryptedSymmetricKey,
		Signature:             proto.Signature,
		EncryptedSnippet:      proto.EncryptedSnippet, // ADDED
	}, nil
}
