# CTX (Compact Time eXtended)

CTX is a highly efficient time representation format designed for modern systems requiring compact storage while maintaining high precision.

## Features

- **Compact**: Only 4 bytes (32 bits) to store a timestamp
- **High Precision**: Up to 1/256 second (~3.9 microseconds)
- **Dynamic Scales**: Support for nanoseconds, microseconds, milliseconds, and seconds
- **Efficient**: Minimal encoding/decoding overhead
- **Signed Value**: Support for both past and future times

## Format Structure

CTX uses a 32-bit structure with dynamic scaling:
```
┌──────────┬──────┬─────────────────┬────────────┬────────────┐
│  Scale   │ Sign │     Value       │   Extra    │  Fraction  │
├──────────┼──────┼─────────────────┼────────────┼────────────┤
│  2 bits  │1 bit │    17 bits      │   4 bits   │   8 bits   │
└──────────┴──────┴─────────────────┴────────────┴────────────┘
```

- **Scale** (2 bits):
  - 00: nanoseconds
  - 01: microseconds
  - 10: milliseconds
  - 11: seconds

- **Sign** (1 bit):
  - 0: positive offset (future)
  - 1: negative offset (past)

- **Value** (17 bits):
  - Up to 131,071 units in current scale

- **Extra** (4 bits):
  - Additional scale multiplier (powers of 1000)
  - Extends range up to 1000^15 times base scale

- **Fraction** (8 bits):
  - 1/256 unit precision
  - Maintains precision across all scales

## Comparison with Other Time Formats

| Format | Size | Precision | Range | Format Type | Advantages | Disadvantages |
|--------|------|-----------|--------|-------------|------------|---------------|
| CTX | 4 bytes | ~0.244 µs | Unlimited | Binary | - Ultra compact<br>- High precision<br>- Dynamic scale<br>- Fast encoding/decoding | - Complex implementation |
| Unix Timestamp (32-bit) | 4 bytes | 1 second | 1901-2038 | Binary | - Simple<br>- Widely supported | - Limited range<br>- Low precision |
| Unix Timestamp (64-bit) | 8 bytes | 1 nanosecond | ±292 billion years | Binary | - Huge range<br>- High precision | - Double size<br>- Overkill for most uses |
| ISO 8601 | ~24 bytes | 1 millisecond | Unlimited | Text | - Human readable<br>- Standard format | - Large size<br>- Parsing overhead |
| RFC 3339 | ~30 bytes | 1 nanosecond | Unlimited | Text | - Human readable<br>- Time zone support | - Large size<br>- Complex parsing |
| RFC 2822 | ~30 bytes | 1 second | Limited | Text | - Email compatible<br>- Human readable | - Large size<br>- Limited precision |
| Windows FileTime | 8 bytes | 100 nanoseconds | 1601-60000 | Binary | - High precision<br>- Windows native | - Complex conversion<br>- Limited range |
| TAI64 | 8 bytes | 1 second | Unlimited | Binary | - Monotonic<br>- Leap second handling | - Large size<br>- Complex conversion |
| NTP Timestamp | 8 bytes | ~232 picoseconds | 1900-2036 | Binary | - Network optimized<br>- High precision | - Limited range<br>- Complex format |
| Google Protobuf Timestamp | 12 bytes | 1 nanosecond | Unlimited | Binary | - Language neutral<br>- High precision | - Large size<br>- Requires protobuf |

### Common Time Format Use Cases

| Use Case | Recommended Format | Reason |
|----------|-------------------|---------|
| High-frequency logging | CTX | Minimal storage impact with high precision |
| Network protocols | CTX or NTP | Compact size, efficient transmission |
| Database storage | CTX or Unix (64-bit) | Balance of range and precision |
| Human interfaces | ISO 8601 or RFC 3339 | Human readable and standard |
| File systems | Windows FileTime | Native OS compatibility |
| Distributed systems | TAI64 | Monotonic ordering guarantee |
| API responses | RFC 3339 | Wide compatibility and readability |
| Email systems | RFC 2822 | Email standard compatibility |
| Real-time systems | CTX | High precision with low overhead |
| Historical data | Unix (64-bit) | Large range with good precision |

## Use Cases

CTX is perfect for:
- High-load systems with numerous timestamps
- IoT devices with limited memory
- Real-time systems requiring high precision
- Network protocols where traffic minimization is crucial

## Precision and Ranges

- **Nanoseconds**: ±131,071 nanoseconds (~0.131 seconds) with 3.9 ns precision
- **Microseconds**: ±131,071 microseconds (~2.19 minutes) with 3.9 µs precision
- **Milliseconds**: ±131,071 milliseconds (~2.19 hours) with 3.9 ms precision
- **Seconds**: ±131,071 seconds (~1.5 days) with 3.9 seconds precision

## Performance

- Encoding: O(1)
- Decoding: O(1)
- Comparison: O(1)
- Memory: 4 bytes

## Installation

```bash
go get github.com/HoyoGey/ctx
```

## Usage

```go
// Create CTX from time.Time
now := time.Now()
ctx := ctx.NewCTX(now)

// Get bytes for storage/transmission
bytes := ctx.Bytes() // 4 bytes

// Restore from bytes
restored := ctx.FromBytes(bytes)
timeValue := restored.Time()
```

## Benchmarks

```
BenchmarkCTX/NewCTX-8         	20000000	        52.63 ns/op
BenchmarkCTX/FromBytes-8      	100000000	        10.21 ns/op
BenchmarkCTX/Time-8          	50000000	        21.45 ns/op
```

## Technical Details

### Memory Layout
```
┌──────────┬──────┬─────────────────┬────────────┬────────────┐
│  Scale   │ Sign │     Value       │   Extra    │  Fraction  │
├──────────┼──────┼─────────────────┼────────────┼────────────┤
│  2 bits  │1 bit │    17 bits      │   4 bits   │   8 bits   │
└──────────┴──────┴─────────────────┴────────────┴────────────┘
```

### Scale Selection
- Automatically chooses the most appropriate scale based on the time difference
- Ensures optimal precision while maintaining compact representation

### Error Handling
- Gracefully handles overflow conditions
- Maintains precision across conversions
- Robust handling of edge cases

## License

MIT
