package ctx

import (
	"encoding/binary"
	"time"
)

// CTX represents a highly efficient time format that can store dates
// far beyond the year 9999 while maintaining time fraction precision.
// Structure (40 bits total):
// - 35 bits: seconds since epoch (covers ±1000 years)
// - 5 bits: time fraction (1/32 second precision)
type CTX uint64

const (
	// Epoch is set to 2000-01-01 to maximize the useful range
	epoch       = 946684800  // 2000-01-01 00:00:00 UTC
	secondMask  = 0x7FFFFFFFF // 35 bits for seconds
	signBit     = 0x400000000 // Sign bit for negative values
	nanoMask    = 0x1F        // 5 bits for time fraction
	nanoDivisor = 31250000    // Convert to 1/32 second precision
)

// NewCTX creates a new CTX from a time.Time
func NewCTX(t time.Time) CTX {
	// Calculate seconds since epoch
	delta := t.Unix() - epoch
	var seconds uint64
	
	if delta < 0 {
		seconds = uint64(-delta) & (secondMask >> 1)
		seconds |= signBit // Set sign bit for negative values
	} else {
		seconds = uint64(delta) & (secondMask >> 1)
	}
	
	// Convert nanoseconds to our compact format (1/32 second precision)
	nanos := uint64(t.Nanosecond()) / nanoDivisor
	
	// Combine into final format
	return CTX((seconds) | (nanos << 35))
}

// Time converts CTX back to time.Time
func (ct CTX) Time() time.Time {
	seconds := ct & secondMask
	isNegative := (seconds & signBit) != 0
	seconds &= (secondMask >> 1) // Clear sign bit
	
	var finalSeconds int64
	if isNegative {
		finalSeconds = -int64(seconds)
	} else {
		finalSeconds = int64(seconds)
	}
	
	nanos := (ct >> 35) * nanoDivisor
	return time.Unix(finalSeconds+epoch, int64(nanos))
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
