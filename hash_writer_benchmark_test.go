package stream_test

import (
	"bytes"
	"strconv"
	"testing"

	"crypto/sha256"

	stream "github.com/da-moon/go-stream"
	miniosha256 "github.com/minio/sha256-simd"
)

// these benchmarks were taken from stdlib's sha256 tests
func BenchmarkHashWriter(b *testing.B) {
	b.ReportAllocs()
	bs := make([]byte, 8192)
	var sizes = []int{8, 1024, 8192}
	for _, size := range sizes {
		b.SetBytes(int64(size))
		b.Run(strconv.Itoa(size), func(b *testing.B) {
			b.ResetTimer()
			b.Run("SHA256", func(b *testing.B) {
				// SHA256 = 32 bytes
				s := 32
				b.Run("stdlib", func(b *testing.B) {
					bench := sha256.New()
					sum := make([]byte, s)
					for i := 0; i < b.N; i++ {
						bench.Reset()
						bench.Write(bs[:size])
						bench.Sum(sum[:0])
					}
				})
				b.Run("hashwriter-stdlib", func(b *testing.B) {
					sum := make([]byte, s)
					buf := bytes.NewBuffer(sum)
					bench := stream.NewWriter(
						buf,
						stream.WithHasher(stream.SHA256, sha256.New()),
					)
					for i := 0; i < b.N; i++ {
						bench.Reset()
						bench.Write(bs[:size])
						bench.Hash(stream.SHA256)
					}
				})
				b.Run("hashwriter-minio", func(b *testing.B) {
					sum := make([]byte, s)
					buf := bytes.NewBuffer(sum)
					bench := stream.NewWriter(
						buf,
						stream.WithHasher(stream.SHA256, miniosha256.New()),
					)
					for i := 0; i < b.N; i++ {
						bench.Reset()
						bench.Write(bs[:size])
						bench.Hash(stream.SHA256)
					}
				})

			})
			// b.Run("MD5", func(b *testing.B) {
			// 	b.Run("stdlib", func(b *testing.B) {
			// b.ResetTimer()
			// 	})
			// 	b.Run("hashwriter", func(b *testing.B) {
			// b.ResetTimer()
			// 	})
			// })
		})
	}

}
