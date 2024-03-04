package propagationManager

import (
	"context"

	"github.com/dominant-strategies/go-quai/consensus"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"
)

func makeValidatorFunc() func(context.Context, peer.ID, *pubsub.Message) pubsub.ValidationResult {
	return func(ctx context.Context, id peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
		return pubsub.ValidationAccept
	}
}

type PropagationManager interface {
	SetValidatorFunc(topic string, validatorFunc func(context.Context, peer.ID, *pubsub.Message) pubsub.ValidationResult)
	ValidateMessage(context.Context, peer.ID, *pubsub.Message) pubsub.ValidationResult
}

// PropagationManager is responsible for managing the propagation of messages
// Each topic has its own validator function that uses the backend to validate the message
type propagationFilter struct {
	validatorFunc func(context.Context, peer.ID, *pubsub.Message) pubsub.ValidationResult
	engine        *consensus.Engine
}

// NewPropagationManager creates a new PropagationManager
func NewPropagationManager(engine *consensus.Engine) PropagationManager {
	return &propagationFilter{
		engine: engine,
	}
}

func (p *propagationFilter) SetValidatorFunc(topic string, validatorFunc func(context.Context, peer.ID, *pubsub.Message) pubsub.ValidationResult) {
	p.validatorFunc = validatorFunc
}

func (p *propagationFilter) ValidateMessage(ctx context.Context, id peer.ID, msg *pubsub.Message) pubsub.ValidationResult {
	if p.validatorFunc != nil {
		return p.validatorFunc(ctx, id, msg)
	}
	return pubsub.ValidationAccept
}
