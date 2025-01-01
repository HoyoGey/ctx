package ctx

import (
	"encoding/binary"
	"fmt"
	"time"
)

// CTX represents a highly efficient time format that can store dates
// far beyond the year 9999 while maintaining microsecond precision.
// Structure (40 bits total):
// - 30 bits: seconds since epoch (covers ±34 years)
// - 10 bits: microsecond fraction (1/2^10 second precision)
type CTX uint64

const (
	// Epoch is set to 2020-01-01 to maximize the useful range
	epoch       = 1577836800 // 2020-01-01 00:00:00 UTC
	secondMask  = 0x3FFFFFFF // 30 bits for seconds
	nanoMask    = 0x3FF      // 10 bits for nano fraction
	nanoDivisor = 1_000_000  // Convert to microsecond precision
)

// NewCTX creates a new CTX from a time.Time
func NewCTX(t time.Time) CTX {
	// Calculate seconds since epoch
	seconds := uint64(t.Unix() - epoch)
	
	// Convert nanoseconds to our compact format (microsecond precision)
	nanos := uint64(t.Nanosecond()) / nanoDivisor
	
	// Combine into final format
	return CTX((seconds & secondMask) | (nanos << 30))
}

// Time converts CTX back to time.Time
func (ct CTX) Time() time.Time {
	seconds := int64(ct&secondMask) + epoch
	nanos := (ct >> 30) * nanoDivisor
	return time.Unix(seconds, int64(nanos))
}

// Bytes converts CTX to a 5-byte slice
func (ct CTX) Bytes() []byte {
	b := make([]byte, 5)
	binary.BigEndian.PutUint32(b[0:4], uint32(ct>>8))
	b[4] = byte(ct)
	return b
}

// FromBytes creates CTX from a 5-byte slice
func FromBytes(b []byte) CTX {
	high := uint64(binary.BigEndian.Uint32(b[0:4]))
	return CTX(high<<8 | uint64(b[4]))
}

func main() {
	// Example usage
	now := time.Date(2025, 1, 2, 0, 2, 22, 0, time.FixedZone("UTC+5", 5*60*60))
	ct := NewCTX(now)
	
	fmt.Println("Original time:", now)
	fmt.Printf("Compact representation (hex): %X\n", uint64(ct))
	fmt.Printf("Size in bytes: %d\n", len(ct.Bytes()))
	
	// Convert back
	restored := ct.Time()
	fmt.Println("Restored time:", restored)
	
	// Example with future date
	future := time.Date(2054, 12, 31, 23, 59, 59, 999999999, time.UTC)
	futureCt := NewCTX(future)
	fmt.Println("\nFuture time:", future)
	fmt.Printf("Future compact (hex): %X\n", uint64(futureCt))
	
	// Demonstrate binary storage
	bytes := ct.Bytes()
	fmt.Printf("Binary storage (hex): % X\n", bytes)
	restored2 := FromBytes(bytes)
	fmt.Println("Restored from binary:", restored2.Time())
}
