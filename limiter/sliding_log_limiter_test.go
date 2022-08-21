package limiter

import (
	"testing"
	"time"
)

func TestNewSlidingLogLimiter(t *testing.T) {
	type args struct {
		smallWindow time.Duration
		strategies  []*SlidingLogLimiterStrategy
	}
	tests := []struct {
		name    string
		args    args
		want    *SlidingLogLimiter
		wantErr bool
	}{
		{
			name: "60_5seconds",
			args: args{
				smallWindow: time.Second,
				strategies: []*SlidingLogLimiterStrategy{
					NewSlidingLogLimiterStrategy(10, time.Minute),
					NewSlidingLogLimiterStrategy(100, time.Hour),
				},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewSlidingLogLimiter(tt.args.smallWindow, tt.args.strategies...)
		})
	}
}
