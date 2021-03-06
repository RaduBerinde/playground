// Results on a Haswell (i5-4670 3.4Ghz)
// go 1.6.2:
//   BenchmarkHash_FHV32_10-4    	100000000	        14.0 ns/op
//   BenchmarkHash_FHV32_1000-4  	 1000000	      1158 ns/op
//   BenchmarkHash_FHV32_10000-4 	  100000	     11622 ns/op
//   BenchmarkHash_FHV32a_10-4   	100000000	        13.2 ns/op
//   BenchmarkHash_FHV32a_1000-4 	 1000000	      1163 ns/op
//   BenchmarkHash_FHV32a_10000-4	  200000	     11629 ns/op
//   BenchmarkHash_FHV64_10-4    	100000000	        14.9 ns/op
//   BenchmarkHash_FHV64_1000-4  	 1000000	      1165 ns/op
//   BenchmarkHash_FHV64_10000-4 	  100000	     11640 ns/op
//   BenchmarkHash_FHV64a_10-4   	100000000	        14.9 ns/op
//   BenchmarkHash_FHV64a_1000-4 	 1000000	      1167 ns/op
//   BenchmarkHash_FHV64a_10000-4	  100000	     11668 ns/op
//   BenchmarkHash_CRC32C_10-4   	100000000	        16.1 ns/op
//   BenchmarkHash_CRC32C_1000-4 	20000000	       113 ns/op
//   BenchmarkHash_CRC32C_10000-4	 1000000	      1118 ns/op
//   BenchmarkHash_XX32_10-4     	100000000	        14.1 ns/op
//   BenchmarkHash_XX32_1000-4   	 3000000	       408 ns/op
//   BenchmarkHash_XX32_10000-4  	  300000	      3998 ns/op
//   BenchmarkHash_XX64_10-4     	100000000	        14.4 ns/op
//   BenchmarkHash_XX64_1000-4   	 5000000	       242 ns/op
//   BenchmarkHash_XX64_10000-4  	 1000000	      2297 ns/op
//
// go 1.7:
//   BenchmarkHash_FHV32_10-4       	100000000	        12.3 ns/op
//   BenchmarkHash_FHV32_1000-4     	 1000000	      1160 ns/op
//   BenchmarkHash_FHV32_10000-4    	  200000	     11614 ns/op
//   BenchmarkHash_FHV32a_10-4      	100000000	        12.1 ns/op
//   BenchmarkHash_FHV32a_1000-4    	 1000000	      1182 ns/op
//   BenchmarkHash_FHV32a_10000-4   	  200000	     11640 ns/op
//   BenchmarkHash_FHV64_10-4       	100000000	        13.6 ns/op
//   BenchmarkHash_FHV64_1000-4     	 1000000	      1163 ns/op
//   BenchmarkHash_FHV64_10000-4    	  200000	     11621 ns/op
//   BenchmarkHash_FHV64a_10-4      	100000000	        14.0 ns/op
//   BenchmarkHash_FHV64a_1000-4    	 1000000	      1185 ns/op
//   BenchmarkHash_FHV64a_10000-4   	  100000	     11654 ns/op
//   BenchmarkHash_CRC32C_10-4      	100000000	        15.9 ns/op
//   BenchmarkHash_CRC32C_1000-4    	20000000	       113 ns/op
//   BenchmarkHash_CRC32C_10000-4   	 1000000	      1095 ns/op
//   BenchmarkHash_XX32_10-4        	100000000	        11.1 ns/op
//   BenchmarkHash_XX32_1000-4      	 5000000	       301 ns/op
//   BenchmarkHash_XX32_10000-4     	  500000	      2945 ns/op
//   BenchmarkHash_XX64_10-4        	100000000	        12.1 ns/op
//   BenchmarkHash_XX64_1000-4      	20000000	       108 ns/op
//   BenchmarkHash_XX64_10000-4     	 2000000	       988 ns/op

package hash

import (
	"hash/crc32"
	"hash/fnv"
	"testing"

	xxhash "github.com/OneOfOne/xxhash/native"
)

const (
	FHV32 = iota
	FHV32a
	FHV64
	FHV64a
	CRC32C
	XX32
	XX64
)

var hashFNV32 = fnv.New32()
var hashFNV32a = fnv.New32a()
var hashFNV64 = fnv.New64()
var hashFNV64a = fnv.New64a()
var hashCRC32C = crc32.New(crc32.MakeTable(crc32.Castagnoli))
var hashXX32 = xxhash.New32()
var hashXX64 = xxhash.New64()

func hash(method int, buf []byte) {
	switch method {
	case FHV32:
		hashFNV32.Reset()
		hashFNV32.Write(buf)
	case FHV32a:
		hashFNV32a.Reset()
		hashFNV32a.Write(buf)
	case FHV64:
		hashFNV64.Reset()
		hashFNV64.Write(buf)
	case FHV64a:
		hashFNV64a.Reset()
		hashFNV64a.Write(buf)
	case CRC32C:
		hashCRC32C.Reset()
		hashCRC32C.Write(buf)
	case XX32:
		hashXX32.Reset()
		hashXX32.Write(buf)
	case XX64:
		hashXX64.Reset()
		hashXX64.Write(buf)
	default:
		panic(method)
	}
}

func benchmarkHash(b *testing.B, method, len int) {
	buf := make([]byte, len)
	for i := range buf {
		buf[i] = byte(i)
	}
	for i := 0; i < b.N; i++ {
		hash(method, buf)
	}
}

func BenchmarkHash_FHV32_10(b *testing.B) {
	benchmarkHash(b, FHV32, 10)
}

func BenchmarkHash_FHV32_1000(b *testing.B) {
	benchmarkHash(b, FHV32, 1000)
}

func BenchmarkHash_FHV32_10000(b *testing.B) {
	benchmarkHash(b, FHV32, 10000)
}

func BenchmarkHash_FHV32a_10(b *testing.B) {
	benchmarkHash(b, FHV32a, 10)
}

func BenchmarkHash_FHV32a_1000(b *testing.B) {
	benchmarkHash(b, FHV32a, 1000)
}

func BenchmarkHash_FHV32a_10000(b *testing.B) {
	benchmarkHash(b, FHV32a, 10000)
}

func BenchmarkHash_FHV64_10(b *testing.B) {
	benchmarkHash(b, FHV64, 10)
}

func BenchmarkHash_FHV64_1000(b *testing.B) {
	benchmarkHash(b, FHV64, 1000)
}

func BenchmarkHash_FHV64_10000(b *testing.B) {
	benchmarkHash(b, FHV64, 10000)
}

func BenchmarkHash_FHV64a_10(b *testing.B) {
	benchmarkHash(b, FHV64a, 10)
}

func BenchmarkHash_FHV64a_1000(b *testing.B) {
	benchmarkHash(b, FHV64a, 1000)
}

func BenchmarkHash_FHV64a_10000(b *testing.B) {
	benchmarkHash(b, FHV64a, 10000)
}

func BenchmarkHash_CRC32C_10(b *testing.B) {
	benchmarkHash(b, CRC32C, 10)
}

func BenchmarkHash_CRC32C_1000(b *testing.B) {
	benchmarkHash(b, CRC32C, 1000)
}

func BenchmarkHash_CRC32C_10000(b *testing.B) {
	benchmarkHash(b, CRC32C, 10000)
}

func BenchmarkHash_XX32_10(b *testing.B) {
	benchmarkHash(b, XX32, 10)
}

func BenchmarkHash_XX32_1000(b *testing.B) {
	benchmarkHash(b, XX32, 1000)
}

func BenchmarkHash_XX32_10000(b *testing.B) {
	benchmarkHash(b, XX32, 10000)
}

func BenchmarkHash_XX64_10(b *testing.B) {
	benchmarkHash(b, XX64, 10)
}

func BenchmarkHash_XX64_1000(b *testing.B) {
	benchmarkHash(b, XX64, 1000)
}

func BenchmarkHash_XX64_10000(b *testing.B) {
	benchmarkHash(b, XX64, 10000)
}
