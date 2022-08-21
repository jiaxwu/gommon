package limiter

import (
	"testing"
	"time"
)

func TestNewLeakyBucketLimiter(t *testing.T) {
	type args struct {
		peakLevel       int
		currentVelocity int
	}
	tests := []struct {
		name    string
		args    args
		want    *LeakyBucketLimiter
		wantErr bool
	}{
		{
			name: "60",
			args: args{
				peakLevel:       60,
				currentVelocity: 10,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLeakyBucketLimiter(tt.args.peakLevel, tt.args.currentVelocity)
			successCount := 0
			for i := 0; i < tt.args.peakLevel; i++ {
				if l.TryAcquire() {
					successCount++
				}
			}
			if successCount != tt.args.peakLevel {
				t.Errorf("NewLeakyBucketLimiter() got = %v, want %v", successCount, tt.args.peakLevel)
				return
			}

			successCount = 0
			for i := 0; i < tt.args.peakLevel; i++ {
				if l.TryAcquire() {
					successCount++
				}
				time.Sleep(time.Second / 10)
			}
			if successCount != tt.args.peakLevel-tt.args.currentVelocity {
				t.Errorf("NewLeakyBucketLimiter() got = %v, want %v", successCount, tt.args.peakLevel-tt.args.currentVelocity)
				return
			}
		})
	}
}
