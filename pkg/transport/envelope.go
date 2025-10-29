package transport

import (
	"fmt"

	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	smv1 "github.com/tinywideclouds/go-action-intention-protos/src/action_intention/envelope/v1"
)

// --- Re-exported Protobuf types ---
type SecureEnvelopePb = smv1.SecureEnvelopePb
type SecureEnvelopeListPb = smv1.SecureEnvelopeListPb // ADDED

// --- SecureEnvelope (Single) ---

// SecureEnvelope is the canonical, idiomatic Go struct for a message.
type SecureEnvelope struct {
	MessageID             string
	SenderID              urn.URN
	RecipientID           urn.URN
	GroupID               urn.URN // ADDED (from our agreement)
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
		GroupId:               native.GroupID.String(),        // ADDED
		ConversationId:        native.ConversationID.String(), // ADDED
		EncryptedData:         native.EncryptedData,
		EncryptedSymmetricKey: native.EncryptedSymmetricKey,
		Signature:             native.Signature,
		EncryptedSnippet:      native.EncryptedSnippet, // ADDED
	}
}

// FromProto converts the Protobuf representation into the idiomatic Go struct.
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

	groupID, err := urn.Parse(proto.GroupId) // ADDED
	if err != nil {
		return nil, fmt.Errorf("failed to parse group id: %w", err)
	}

	conversationID, err := urn.Parse(proto.ConversationId) // ADDED
	if err != nil {
		return nil, fmt.Errorf("failed to parse conversation id: %w", err)
	}

	return &SecureEnvelope{
		MessageID:             proto.MessageId,
		SenderID:              senderID,
		RecipientID:           recipientID,
		GroupID:               groupID,        // ADDED
		ConversationID:        conversationID, // ADDED
		EncryptedData:         proto.EncryptedData,
		EncryptedSymmetricKey: proto.EncryptedSymmetricKey,
		Signature:             proto.Signature,
		EncryptedSnippet:      proto.EncryptedSnippet, // ADDED
	}, nil
}

// --- SecureEnvelopeList (List) --- ADDED THIS SECTION ---

// SecureEnvelopeList is the idiomatic Go struct for a list of envelopes.
type SecureEnvelopeList struct {
	Envelopes []*SecureEnvelope
}

// ListToProto converts the idiomatic Go list into its Protobuf representation.
func ListToProto(native *SecureEnvelopeList) *SecureEnvelopeListPb {
	if native == nil {
		return nil
	}
	protoEnvelopes := make([]*SecureEnvelopePb, len(native.Envelopes))
	for i, env := range native.Envelopes {
		protoEnvelopes[i] = ToProto(env)
	}
	return &SecureEnvelopeListPb{
		Envelopes: protoEnvelopes,
	}
}

// ListFromProto converts the Protobuf list into the idiomatic Go struct.
func ListFromProto(proto *SecureEnvelopeListPb) (*SecureEnvelopeList, error) {
	if proto == nil {
		return nil, nil
	}
	nativeEnvelopes := make([]*SecureEnvelope, len(proto.Envelopes))
	for i, pEnv := range proto.Envelopes {
		native, err := FromProto(pEnv)
		if err != nil {
			return nil, fmt.Errorf("failed to parse envelope at index %d: %w", i, err)
		}
		nativeEnvelopes[i] = native
	}
	return &SecureEnvelopeList{
		Envelopes: nativeEnvelopes,
	}, nil
}
