package ctx

import (
	"testing"
	"time"
)

func TestCTX(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
	}{
		{
			name: "current time",
			time: time.Date(2025, 1, 2, 0, 16, 27, 0, time.FixedZone("UTC+5", 5*60*60)),
		},
		{
			name: "future time 2099",
			time: time.Date(2099, 12, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "very future time",
			time: time.Date(2054, 12, 31, 23, 59, 59, 999999999, time.UTC),
		},
		{
			name: "past time",
			time: time.Date(1986, 1, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to CTX
			ctx := NewCTX(tt.time)

			// Convert back to time.Time
			got := ctx.Time()

			// Compare with microsecond precision
			if got.Unix() != tt.time.Unix() {
				t.Errorf("Time = %v, want %v", got, tt.time)
			}

			// Test binary serialization
			bytes := ctx.Bytes()
			if len(bytes) != 5 {
				t.Errorf("Bytes length = %v, want 5", len(bytes))
			}

			// Test deserialization
			restored := FromBytes(bytes)
			if restored != ctx {
				t.Errorf("FromBytes = %v, want %v", restored, ctx)
			}

			// Print the binary representation for the 2099 test case
			if tt.name == "future time 2099" {
				t.Logf("2099-12-01 binary: % X", bytes)
				t.Logf("2099-12-01 restored: %v", restored.Time())
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

	ctx := NewCTX(now)
	
	b.Run("Time", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ctx.Time()
		}
	})

	b.Run("Bytes", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = ctx.Bytes()
		}
	})
}
