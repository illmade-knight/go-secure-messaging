package transport

import (
	"fmt"

	"github.com/illmade-knight/go-secure-messaging/pkg/urn"
	smv1 "github.com/tinywideclouds/go-action-intention-protos/src/action_intention/envelope/v1"
)

// --- Re-exported Protobuf types ---
type EncryptedDigestPb = smv1.EncryptedDigestPb
type EncryptedDigestItemPb = smv1.EncryptedDigestItemPb

// --- EncryptedDigest (List) ---

// EncryptedDigest is the idiomatic Go struct for a list of digest items.
type EncryptedDigest struct {
	Items []*EncryptedDigestItem
}

// EncryptedDigestItem is the idiomatic Go struct for a single digest item.
type EncryptedDigestItem struct {
	ConversationID        urn.URN
	EncryptedSnippet      []byte
	EncryptedSymmetricKey []byte
}

// DigestToProto converts the idiomatic Go digest struct into its Protobuf representation.
func DigestToProto(native *EncryptedDigest) *EncryptedDigestPb {
	if native == nil {
		return nil
	}

	protoItems := make([]*EncryptedDigestItemPb, len(native.Items))
	for i, item := range native.Items {
		if item == nil {
			continue // Skip nil items in the slice
		}
		protoItems[i] = &EncryptedDigestItemPb{
			ConversationId:        item.ConversationID.String(),
			EncryptedSnippet:      item.EncryptedSnippet,
			EncryptedSymmetricKey: item.EncryptedSymmetricKey,
		}
	}

	return &EncryptedDigestPb{
		Items: protoItems,
	}
}

// DigestFromProto converts the Protobuf digest representation into the idiomatic Go struct.
// It validates that the string identifiers are valid URNs.
func DigestFromProto(proto *EncryptedDigestPb) (*EncryptedDigest, error) {
	if proto == nil {
		return nil, nil
	}

	nativeItems := make([]*EncryptedDigestItem, len(proto.Items))
	for i, item := range proto.Items {
		if item == nil {
			continue // Skip nil items in the slice
		}

		conversationID, err := urn.Parse(item.ConversationId)
		if err != nil {
			return nil, fmt.Errorf("item %d: failed to parse conversation id: %w", i, err)
		}

		nativeItems[i] = &EncryptedDigestItem{
			ConversationID:        conversationID,
			EncryptedSnippet:      item.EncryptedSnippet,
			EncryptedSymmetricKey: item.EncryptedSymmetricKey,
		}
	}

	return &EncryptedDigest{
		Items: nativeItems,
	}, nil
}
