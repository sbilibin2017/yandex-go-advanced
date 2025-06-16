package runners

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRunContext_CancelFuncCancelsContext(t *testing.T) {
	ctx, cancel := NewRunContext(context.Background())
	assert.NotNil(t, ctx, "context should not be nil")
	assert.NotNil(t, cancel, "cancel func should not be nil")

	// Context should not be canceled immediately
	select {
	case <-ctx.Done():
		t.Fatal("context should not be canceled yet")
	default:
		// expected
	}

	// Now cancel it manually and check if context is cancelled
	cancel()

	select {
	case <-ctx.Done():
		assert.Equal(t, context.Canceled, ctx.Err())
	case <-time.After(time.Second):
		t.Fatal("context was not cancelled in time")
	}
}
