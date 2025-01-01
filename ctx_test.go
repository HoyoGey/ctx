package ctx

import (
	"math"
	"testing"
	"time"
)

func TestCTX(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		maxDiff  time.Duration
	}{
		{
			name:    "current_time",
			time:    time.Now(),
			maxDiff: time.Second,
		},
		{
			name:    "future_time_near",
			time:    time.Now().Add(30 * time.Minute),
			maxDiff: 10 * time.Second,
		},
		{
			name:    "future_time_medium",
			time:    time.Now().Add(12 * time.Hour),
			maxDiff: 20 * time.Minute, // Increased tolerance for medium-range times
		},
		{
			name:    "future_time_far",
			time:    time.Now().Add(365 * 24 * time.Hour),
			maxDiff: 24 * time.Hour, // Full day tolerance for year-range times
		},
		{
			name:    "past_time_near",
			time:    time.Now().Add(-30 * time.Minute),
			maxDiff: 10 * time.Second,
		},
		{
			name:    "past_time_medium",
			time:    time.Now().Add(-12 * time.Hour),
			maxDiff: 20 * time.Minute, // Increased tolerance for medium-range times
		},
		{
			name:    "past_time_far",
			time:    time.Now().Add(-365 * 24 * time.Hour),
			maxDiff: 24 * time.Hour, // Full day tolerance for year-range times
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create CTX
			ct := NewCTX(tt.time)
			
			// Convert to bytes and back
			bytes := ct.Bytes()
			if len(bytes) != 4 {
				t.Errorf("Expected 4 bytes, got %d bytes", len(bytes))
			}
			
			// Print binary representation
			t.Logf("Binary: %02X %02X %02X %02X", bytes[0], bytes[1], bytes[2], bytes[3])
			
			// Restore from bytes
			restored := FromBytes(bytes)
			restoredTime := restored.Time()
			
			// Calculate difference
			diff := tt.time.Sub(restoredTime)
			if math.Abs(float64(diff)) > float64(tt.maxDiff) {
				t.Errorf("Time mismatch: want %v, got %v (diff: %v)", 
					tt.time.Format(time.RFC3339Nano), 
					restoredTime.Format(time.RFC3339Nano),
					diff)
			}
		})
	}
}

func TestPrecision(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name     string
		duration time.Duration
		maxDiff  time.Duration
	}{
		{"100µs", 100 * time.Microsecond, time.Second / 4}, // 1/4 second precision
		{"1ms", time.Millisecond, time.Second / 4},
		{"10ms", 10 * time.Millisecond, time.Second / 4},
		{"100ms", 100 * time.Millisecond, time.Second / 4},
		{"1s", time.Second, time.Second / 4},
		{"-100µs", -100 * time.Microsecond, time.Second / 4},
		{"-1ms", -time.Millisecond, time.Second / 4},
		{"-10ms", -10 * time.Millisecond, time.Second / 4},
		{"-100ms", -100 * time.Millisecond, time.Second / 4},
		{"-1s", -time.Second, time.Second / 4},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			future := now.Add(tt.duration)
			ct := NewCTX(future)
			restored := FromBytes(ct.Bytes()).Time()
			
			diff := future.Sub(restored)
			if math.Abs(float64(diff)) > float64(tt.maxDiff) {
				t.Errorf("Precision test failed for %v: want %v, got %v (diff: %v)",
					tt.duration,
					future.Format(time.RFC3339Nano),
					restored.Format(time.RFC3339Nano),
					diff)
			}
		})
	}
}

func BenchmarkCTX(b *testing.B) {
	now := time.Now()
	
	b.Run("NewCTX", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = NewCTX(now)
		}
	})
	
	ct := NewCTX(now)
	bytes := ct.Bytes()
	
	b.Run("FromBytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = FromBytes(bytes)
		}
	})
	
	b.Run("Time", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ct.Time()
		}
	})
}
