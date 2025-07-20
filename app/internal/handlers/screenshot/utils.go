package screenshot

import (
	"context"
	"time"
)

func SleepContext(ctx context.Context, duration int) error {
	if duration == 0 {
		return nil
	}
	timer := time.NewTimer(time.Duration(duration) * time.Second)
	select {
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
