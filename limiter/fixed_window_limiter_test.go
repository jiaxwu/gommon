package limiter

import (
	"testing"
	"time"
)

func TestNewFixedWindowLimiter(t *testing.T) {
	type args struct {
		limit  int
		window time.Duration
	}
	tests := []struct {
		name string
		args args
		want *FixedWindowLimiter
	}{
		{
			name: "100_second",
			args: args{
				limit:  100,
				window: time.Second,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewFixedWindowLimiter(tt.args.limit, tt.args.window)
			successCount := 0
			for i := 0; i < tt.args.limit*2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tt.args.limit {
				t.Errorf("NewFixedWindowLimiter() = %v, want %v", successCount, tt.args.limit)
			}
			time.Sleep(time.Second)
			successCount = 0
			for i := 0; i < tt.args.limit*2; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tt.args.limit {
				t.Errorf("NewFixedWindowLimiter() = %v, want %v", successCount, tt.args.limit)
			}
		})
	}
}
