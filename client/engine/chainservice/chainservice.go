// Package chainservice is a chain service responsible for submitting blockchain transactions and relaying blockchain events.
package chainservice // import "github.com/statechannels/go-nitro/client/chainservice"

import (
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type Event interface {
	GetChannelId() types.Destination
}

// DepositedEvent is an internal representation of the deposited blockchain event
type DepositedEvent struct {
	ChannelId          types.Destination
	Holdings           types.Funds // indexed by asset
	AdjudicationStatus protocols.AdjudicationStatus
	BlockNum           uint64
}

func (de DepositedEvent) GetChannelId() types.Destination {
	return de.ChannelId
}

// todo implement other event types
// AllocationUpdated
// Concluded
// ChallengeRegistered
// ChallengeCleared

type ChainEventHandler interface {
	UpdateWithChainEvent(event Event) (protocols.Objective, error)
}

type ChainService interface {
	Out() <-chan Event
	In() chan<- protocols.ChainTransaction
}
