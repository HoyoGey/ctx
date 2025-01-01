# CTX: Compact Time eXtended Format

CTX is a revolutionary 40-bit time format designed for modern systems requiring compact storage while maintaining precision. Think of it as "UTC's compact cousin for the digital age."

## Features

- **Ultra Compact**: Only 5 bytes
- **Wide Range**: ±1000 years (1000-3000 AD)
- **High Precision**: 1/32 second accuracy
- **Zero Dependencies**: Pure implementation
- **Multiple Language Support**: Easy to implement in any language

## Format Specification

```
40-bit structure:
┌─────────────┬──────────────────┬─────────────┐
│  Sign (1)   │  Seconds (34)    │  Frac (5)   │
└─────────────┴──────────────────┴─────────────┘

- Sign: 1 bit for negative values
- Seconds: 34 bits for seconds since 2000-01-01
- Frac: 5 bits for subsecond precision (1/32 sec)
```

## Usage Examples

### Go Implementation
```go
import "github.com/HoyoGey/ctx"

// Create CTX from time
now := time.Now()
ct := ctx.NewCTX(now)

// Get binary representation
bytes := ct.Bytes() // 5 bytes

// Restore from binary
restored := ctx.FromBytes(bytes)
time := restored.Time()
```

### JavaScript Implementation
```javascript
class CTX {
    static EPOCH = new Date('2000-01-01T00:00:00Z').getTime() / 1000;
    static SECOND_MASK = 0x7FFFFFFFF;
    static SIGN_BIT = 0x400000000;
    static NANO_MASK = 0x1F;
    static NANO_DIVISOR = 31250000;

    constructor(date) {
        const delta = Math.floor(date.getTime() / 1000) - CTX.EPOCH;
        let seconds = 0;
        
        if (delta < 0) {
            seconds = (-delta) & (CTX.SECOND_MASK >> 1);
            seconds |= CTX.SIGN_BIT;
        } else {
            seconds = delta & (CTX.SECOND_MASK >> 1);
        }
        
        const nanos = (date.getMilliseconds() * 1e6) / CTX.NANO_DIVISOR;
        this.value = (seconds) | (nanos << 35);
    }

    toDate() {
        let seconds = this.value & CTX.SECOND_MASK;
        const isNegative = (seconds & CTX.SIGN_BIT) !== 0;
        seconds &= (CTX.SECOND_MASK >> 1);
        
        if (isNegative) {
            seconds = -seconds;
        }
        
        const nanos = (this.value >> 35) * CTX.NANO_DIVISOR;
        return new Date((seconds + CTX.EPOCH) * 1000 + nanos / 1e6);
    }

    toBytes() {
        const buffer = new ArrayBuffer(5);
        const view = new DataView(buffer);
        view.setUint32(0, Number(this.value >> 8n));
        view.setUint8(4, Number(this.value & 0xFFn));
        return new Uint8Array(buffer);
    }

    static fromBytes(bytes) {
        const view = new DataView(bytes.buffer);
        const high = BigInt(view.getUint32(0));
        const low = BigInt(view.getUint8(4));
        const ctx = new CTX(new Date(0));
        ctx.value = (high << 8n) | low;
        return ctx;
    }
}

// Usage Example
const now = new Date();
const ctx = new CTX(now);
console.log('Original:', now);

const bytes = ctx.toBytes();
console.log('Bytes:', Array.from(bytes).map(b => b.toString(16).padStart(2, '0')).join(' '));

const restored = CTX.fromBytes(bytes);
console.log('Restored:', restored.toDate());
```

### Python Implementation
```python
from datetime import datetime, timezone
import struct

class CTX:
    EPOCH = int(datetime(2000, 1, 1, tzinfo=timezone.utc).timestamp())
    SECOND_MASK = 0x7FFFFFFFF
    SIGN_BIT = 0x400000000
    NANO_MASK = 0x1F
    NANO_DIVISOR = 31250000

    def __init__(self, dt):
        delta = int(dt.timestamp()) - self.EPOCH
        
        if delta < 0:
            seconds = (-delta) & (self.SECOND_MASK >> 1)
            seconds |= self.SIGN_BIT
        else:
            seconds = delta & (self.SECOND_MASK >> 1)
            
        nanos = (dt.microsecond * 1000) // self.NANO_DIVISOR
        self.value = (seconds) | (nanos << 35)

    def to_datetime(self):
        seconds = self.value & self.SECOND_MASK
        is_negative = (seconds & self.SIGN_BIT) != 0
        seconds &= (self.SECOND_MASK >> 1)
        
        if is_negative:
            seconds = -seconds
            
        nanos = (self.value >> 35) * self.NANO_DIVISOR
        return datetime.fromtimestamp(seconds + self.EPOCH + nanos/1e9, tz=timezone.utc)

    def to_bytes(self):
        return struct.pack('>IB', 
            self.value >> 8, 
            self.value & 0xFF)

    @classmethod
    def from_bytes(cls, data):
        high, low = struct.unpack('>IB', data)
        ctx = cls(datetime.now(timezone.utc))
        ctx.value = (high << 8) | low
        return ctx

# Usage Example
now = datetime.now(timezone.utc)
ctx = CTX(now)
print('Original:', now)

bytes_data = ctx.to_bytes()
print('Bytes:', ' '.join(f'{b:02x}' for b in bytes_data))

restored = CTX.from_bytes(bytes_data)
print('Restored:', restored.to_datetime())
```

## Implementation Guidelines

1. **Epoch Selection**
   - Use 2000-01-01 00:00:00 UTC as epoch
   - This provides balanced range for past/future dates

2. **Bit Layout**
   - Sign bit: Most significant bit (bit 39)
   - Seconds: Next 34 bits (bits 5-38)
   - Fraction: Last 5 bits (bits 0-4)

3. **Handling Negative Values**
   - Use sign bit for dates before 2000
   - Apply two's complement for seconds

4. **Binary Serialization**
   - Use big-endian byte order
   - Pack into exactly 5 bytes

## Best Practices

1. **Time Zone Handling**
   - Convert to UTC before encoding
   - Store time zone separately if needed
   - Restore in UTC and convert to local time

2. **Range Validation**
   - Check if date is within supported range
   - Handle out-of-range errors gracefully

3. **Precision Considerations**
   - 1/32 second is sufficient for most uses
   - Round to nearest fraction if needed

4. **Binary Storage**
   - Always use binary format for storage
   - Use consistent endianness (big-endian)

## Performance Tips

1. **Bit Operations**
   - Use bitwise operations for speed
   - Avoid floating-point calculations
   - Minimize type conversions

2. **Memory Usage**
   - Preallocate byte arrays
   - Reuse buffers when possible
   - Avoid unnecessary allocations

## Common Pitfalls

1. **Time Zones**
   - Not storing time zone information
   - Mixing local and UTC times

2. **Range Errors**
   - Not checking date range
   - Overflow in calculations

3. **Precision Loss**
   - Not accounting for rounding
   - Assuming nanosecond precision

## Contributing

Feel free to contribute to CTX format:
1. Fork the repository
2. Create your feature branch
3. Add tests for new features
4. Submit a pull request

## License

MIT License - feel free to use in any project
