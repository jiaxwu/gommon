package limiter

import (
	"testing"
	"time"
)

func TestNewTokenBucketLimiter(t *testing.T) {
	type args struct {
		capacity int
		rate     int
	}
	tests := []struct {
		name string
		args args
		want *TokenBucketLimiter
	}{
		{
			name: "60",
			args: args{
				capacity: 60,
				rate:     10,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewTokenBucketLimiter(tt.args.capacity, tt.args.rate)
			time.Sleep(time.Second)
			successCount := 0
			for i := 0; i < tt.args.rate; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tt.args.rate {
				t.Errorf("NewTokenBucketLimiter() got = %v, want %v", successCount, tt.args.rate)
				return
			}

			successCount = 0
			for i := 0; i < tt.args.capacity; i++ {
				if l.TryAcquire() {
					successCount++
				}
				time.Sleep(time.Second / 10)
			}
			if successCount != tt.args.capacity-tt.args.rate {
				t.Errorf("NewTokenBucketLimiter() got = %v, want %v", successCount, tt.args.capacity-tt.args.rate)
				return
			}
		})
	}
}
