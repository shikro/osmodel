package events_test

import (
	"context"
	"testing"
	"time"

	"github.com/shikro/osmodel/events"
)

func TestPingEvent(t *testing.T) {
	t.Run("Should send event", func(t *testing.T) {
		ctx := context.Background()
		pingEvent := events.NewPingEvent(10 * time.Millisecond)
		pingEvent.Start(ctx)
		time.Sleep(100 * time.Millisecond)
		select {
		case <-pingEvent.Event():
		default:
			t.Fail()
		}
	})
}
