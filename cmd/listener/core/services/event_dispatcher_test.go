package services

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/core"
	"github.com/raylsnetwork/rayls-privacy-pnh-governance-api/cmd/listener/testutil"
)

// stubEventConsumer returns messages from a queue or an error.
type stubEventConsumer struct {
	messages []core.Message
	idx      int
	err      error
}

func (s *stubEventConsumer) Next(ctx context.Context) (core.Message, error) {
	if s.err != nil {
		err := s.err
		s.err = nil // only return error once, then fall through to context check
		return core.Message{}, err
	}
	if s.idx >= len(s.messages) {
		// Block until context is cancelled
		<-ctx.Done()
		return core.Message{}, ctx.Err()
	}
	msg := s.messages[s.idx]
	s.idx++
	return msg, nil
}

// stubDispatcherHandler records Handle calls and can return a configured error.
type stubDispatcherHandler struct {
	contract    string
	handlerName string
	handleErr   error
	handledLogs []core.ContractLog
}

func (s *stubDispatcherHandler) ContractName() string { return s.contract }
func (s *stubDispatcherHandler) Name() string         { return s.handlerName }

func (s *stubDispatcherHandler) Handle(_ context.Context, log core.ContractLog) error {
	s.handledLogs = append(s.handledLogs, log)
	return s.handleErr
}

// ackTracker records whether Ack was called.
type ackTracker struct {
	called bool
}

func (a *ackTracker) ack(_ context.Context) error {
	a.called = true
	return nil
}

func TestNewEventDispatcher_BuildsHandlerMap(t *testing.T) {
	// Handler map is built from ContractName() of each handler
	h1 := &stubDispatcherHandler{contract: "TokenCore", handlerName: "TokenCoreHandler"}
	h2 := &stubDispatcherHandler{contract: "Teleport", handlerName: "TeleportHandler"}
	consumer := &stubEventConsumer{}

	dispatcher := NewEventDispatcher(consumer, []core.EventHandler{h1, h2}, &testutil.StubLogger{})

	require.NotNil(t, dispatcher)
	assert.Len(t, dispatcher.handlers, 2)
	assert.Equal(t, h1, dispatcher.handlers["TokenCore"])
	assert.Equal(t, h2, dispatcher.handlers["Teleport"])
}

func TestNewEventDispatcher_DuplicateContractNamePanics(t *testing.T) {
	// Registering two handlers for the same contract panics
	h1 := &stubDispatcherHandler{contract: "TokenCore", handlerName: "Handler1"}
	h2 := &stubDispatcherHandler{contract: "TokenCore", handlerName: "Handler2"}
	consumer := &stubEventConsumer{}

	assert.PanicsWithValue(t, `duplicate handler for contract "TokenCore"`, func() {
		NewEventDispatcher(consumer, []core.EventHandler{h1, h2}, &testutil.StubLogger{})
	})
}

func TestEventDispatcher_Run_ContextCancelledReturnsNil(t *testing.T) {
	// Run returns nil when the context is cancelled
	consumer := &stubEventConsumer{}
	dispatcher := NewEventDispatcher(consumer, nil, &testutil.StubLogger{})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := dispatcher.Run(ctx)

	assert.NoError(t, err)
}

func TestEventDispatcher_Run_DispatchesToCorrectHandler(t *testing.T) {
	// Message is dispatched to the handler matching its contract name
	ack := &ackTracker{}
	consumer := &stubEventConsumer{
		messages: []core.Message{
			{
				Log: core.ContractLog{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
				Ack: ack.ack,
			},
		},
	}
	h1 := &stubDispatcherHandler{contract: "TokenCore", handlerName: "TokenCoreHandler"}
	h2 := &stubDispatcherHandler{contract: "Teleport", handlerName: "TeleportHandler"}

	ctx, cancel := context.WithCancel(context.Background())
	dispatcher := NewEventDispatcher(consumer, []core.EventHandler{h1, h2}, &testutil.StubLogger{})

	go func() {
		// Cancel after the consumer runs out of messages
		for consumer.idx < len(consumer.messages) {
			runtime.Gosched()
		}
		cancel()
	}()

	err := dispatcher.Run(ctx)

	assert.NoError(t, err)
	require.Len(t, h1.handledLogs, 1)
	assert.Equal(t, "TokenCore", h1.handledLogs[0].ContractName)
	assert.Empty(t, h2.handledLogs)
	assert.True(t, ack.called)
}

func TestEventDispatcher_Run_UnknownContractSkipsAndAcks(t *testing.T) {
	// Messages for unknown contracts are acked without handler invocation
	ack := &ackTracker{}
	consumer := &stubEventConsumer{
		messages: []core.Message{
			{
				Log: core.ContractLog{ContractName: "Unknown", EventName: "SomeEvent", BlockNumber: 100},
				Ack: ack.ack,
			},
		},
	}
	handler := &stubDispatcherHandler{contract: "TokenCore", handlerName: "TokenCoreHandler"}

	ctx, cancel := context.WithCancel(context.Background())
	dispatcher := NewEventDispatcher(consumer, []core.EventHandler{handler}, &testutil.StubLogger{})

	go func() {
		for consumer.idx < len(consumer.messages) {
			runtime.Gosched()
		}
		cancel()
	}()

	err := dispatcher.Run(ctx)

	assert.NoError(t, err)
	assert.Empty(t, handler.handledLogs)
	assert.True(t, ack.called)
}

func TestEventDispatcher_Run_HandlerErrorDoesNotAck(t *testing.T) {
	// When a handler returns an error, the message is not acked
	ack := &ackTracker{}
	consumer := &stubEventConsumer{
		messages: []core.Message{
			{
				Log: core.ContractLog{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
				Ack: ack.ack,
			},
		},
	}
	handler := &stubDispatcherHandler{
		contract:    "TokenCore",
		handlerName: "TokenCoreHandler",
		handleErr:   errors.New("processing failed"),
	}

	ctx, cancel := context.WithCancel(context.Background())
	dispatcher := NewEventDispatcher(consumer, []core.EventHandler{handler}, &testutil.StubLogger{})

	go func() {
		for consumer.idx < len(consumer.messages) {
			runtime.Gosched()
		}
		cancel()
	}()

	err := dispatcher.Run(ctx)

	assert.NoError(t, err)
	require.Len(t, handler.handledLogs, 1)
	assert.False(t, ack.called)
}

func TestEventDispatcher_Run_HandlerSuccessAcks(t *testing.T) {
	// When a handler succeeds, the message is acked
	ack := &ackTracker{}
	consumer := &stubEventConsumer{
		messages: []core.Message{
			{
				Log: core.ContractLog{ContractName: "Teleport", EventName: "MessageSent", BlockNumber: 200},
				Ack: ack.ack,
			},
		},
	}
	handler := &stubDispatcherHandler{contract: "Teleport", handlerName: "TeleportHandler"}

	ctx, cancel := context.WithCancel(context.Background())
	dispatcher := NewEventDispatcher(consumer, []core.EventHandler{handler}, &testutil.StubLogger{})

	go func() {
		for consumer.idx < len(consumer.messages) {
			runtime.Gosched()
		}
		cancel()
	}()

	err := dispatcher.Run(ctx)

	assert.NoError(t, err)
	require.Len(t, handler.handledLogs, 1)
	assert.True(t, ack.called)
}

func TestEventDispatcher_Run_ConsumerErrorContinues(t *testing.T) {
	// A non-context consumer error is logged and the loop continues
	ack := &ackTracker{}
	consumer := &stubEventConsumer{
		err: errors.New("temporary failure"),
		messages: []core.Message{
			{
				Log: core.ContractLog{ContractName: "TokenCore", EventName: "TokenRegistered", BlockNumber: 100},
				Ack: ack.ack,
			},
		},
	}
	handler := &stubDispatcherHandler{contract: "TokenCore", handlerName: "TokenCoreHandler"}

	ctx, cancel := context.WithCancel(context.Background())
	dispatcher := NewEventDispatcher(consumer, []core.EventHandler{handler}, &testutil.StubLogger{})

	go func() {
		for consumer.idx < len(consumer.messages) {
			runtime.Gosched()
		}
		cancel()
	}()

	err := dispatcher.Run(ctx)

	assert.NoError(t, err)
	require.Len(t, handler.handledLogs, 1)
	assert.True(t, ack.called)
}
