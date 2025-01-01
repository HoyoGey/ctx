package ctx

import (
	"math"
	"time"
)

type CTX uint32

const (
	scaleMask  = 0xC0000000 // 2 bits for scale
	signMask   = 0x20000000 // 1 bit for sign
	valueMask  = 0x1FFFF000 // 17 bits for value
	extraMask  = 0x00000F00 // 4 bits for extra scale
	fracMask   = 0x000000FF // 8 bits for fraction

	scaleShift  = 30
	signShift   = 29
	valueShift  = 12
	extraShift  = 8
	fracShift   = 0

	fracBits     = 8
	fracMultiple = 1 << fracBits // 256 for 8 bits

	// Scale values
	scaleNano  = 0 // nanoseconds
	scaleMicro = 1 // microseconds
	scaleMilli = 2 // milliseconds
	scaleSecond = 3 // seconds
)

var scaleFactors = []float64{
	1e-9,  // nanoseconds
	1e-6,  // microseconds
	1e-3,  // milliseconds
	1,     // seconds
}

func NewCTX(t time.Time) CTX {
	// Calculate difference from Unix epoch
	diff := t.UnixNano()
	
	// Find the most appropriate scale
	var scale, extra uint32
	absDiff := math.Abs(float64(diff))
	
	if absDiff < 1e9 { // < 1 second
		scale = scaleNano
	} else if absDiff < 1e12 { // < 1000 seconds
		scale = scaleMicro
	} else if absDiff < 1e15 { // < 1M seconds
		scale = scaleMilli
	} else {
		scale = scaleSecond
	}
	
	// Calculate extra scale (powers of 1000)
	for absDiff >= float64(math.MaxInt32) {
		absDiff /= 1000
		extra++
		if extra >= 15 { // 15 is max value for 4 bits
			break
		}
	}

	// Convert to selected scale
	scaleFactor := scaleFactors[scale] * math.Pow(1000, float64(extra))
	value := float64(diff) * scaleFactor

	// Split into integer and fractional parts
	intPart := uint32(math.Abs(float64(int64(value))))
	fracPart := uint32((math.Abs(value) - float64(intPart)) * fracMultiple)

	// Combine all parts
	var result uint32
	result |= scale << scaleShift
	if diff < 0 {
		result |= 1 << signShift
	}
	result |= (intPart & 0x1FFFF) << valueShift
	result |= (extra & 0xF) << extraShift
	result |= fracPart & 0xFF

	return CTX(result)
}

func (c CTX) Time() time.Time {
	// Extract components
	scale := (uint32(c) & scaleMask) >> scaleShift
	isNegative := (uint32(c) & signMask) != 0
	value := (uint32(c) & valueMask) >> valueShift
	extra := (uint32(c) & extraMask) >> extraShift
	frac := float64(uint32(c)&fracMask) / fracMultiple

	// Calculate total value
	scaleFactor := scaleFactors[scale] * math.Pow(1000, float64(extra))
	totalValue := (float64(value) + frac) / scaleFactor

	if isNegative {
		totalValue = -totalValue
	}

	// Convert to time
	return time.Unix(0, int64(totalValue))
}

func (c CTX) Bytes() []byte {
	return []byte{
		byte(uint32(c) >> 24),
		byte(uint32(c) >> 16),
		byte(uint32(c) >> 8),
		byte(uint32(c)),
	}
}

func FromBytes(b []byte) CTX {
	if len(b) != 4 {
		return 0
	}
	return CTX(uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3]))
}
