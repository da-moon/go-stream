package stream_test

import (
	"bytes"
	"strconv"
	"testing"

	"crypto/md5"
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
			b.Run("single-hash", func(b *testing.B) {
				b.Run("SHA256", func(b *testing.B) {
					// SHA256 = 32 bytes
					s := 32
					b.Run("stdlib", func(b *testing.B) {
						b.ResetTimer()
						bench := sha256.New()
						sum := make([]byte, s)
						for i := 0; i < b.N; i++ {
							bench.Reset()
							bench.Write(bs[:size])
							bench.Sum(sum[:0])
						}
					})
					b.Run("hashwriter-stdlib", func(b *testing.B) {
						b.ResetTimer()
						sum := make([]byte, s)
						buf := bytes.NewBuffer(sum)
						bench, _ := stream.NewHashWriter(
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
						b.ResetTimer()
						sum := make([]byte, s)
						buf := bytes.NewBuffer(sum)
						bench, _ := stream.NewHashWriter(
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
				b.Run("MD5", func(b *testing.B) {
					// md5 chucksum size
					s := 16
					b.Run("stdlib", func(b *testing.B) {
						b.ResetTimer()
						b.ResetTimer()
						bench := md5.New()
						sum := make([]byte, s)
						for i := 0; i < b.N; i++ {
							bench.Reset()
							bench.Write(bs[:size])
							bench.Sum(sum[:0])
						}
					})
					b.Run("hashwriter-stdlib", func(b *testing.B) {
						b.ResetTimer()
						sum := make([]byte, s)
						buf := bytes.NewBuffer(sum)
						bench, _ := stream.NewHashWriter(
							buf,
							stream.WithHasher(stream.MD5, md5.New()),
						)
						for i := 0; i < b.N; i++ {
							bench.Reset()
							bench.Write(bs[:size])
							bench.Hash(stream.SHA256)
						}
					})
					b.Run("hashwriter-minio", func(b *testing.B) {
						b.ResetTimer()
						sum := make([]byte, s)
						buf := bytes.NewBuffer(sum)
						bench, _ := stream.NewHashWriter(
							buf,
							stream.WithMD5(),
						)
						for i := 0; i < b.N; i++ {
							bench.Reset()
							bench.Write(bs[:size])
							bench.Hash(stream.SHA256)
						}
					})
				})
			})
			b.Run("multi-hash", func(b *testing.B) {
				s := 8
				b.Run("stdlib", func(b *testing.B) {
					b.ResetTimer()
					shabench := sha256.New()
					md5bench := md5.New()
					sum := make([]byte, s)
					for i := 0; i < b.N; i++ {
						shabench.Reset()
						md5bench.Reset()
						shabench.Write(bs[:size])
						md5bench.Write(bs[:size])
						shabench.Sum(sum[:0])
						md5bench.Sum(sum[:0])
					}
				})
				b.Run("hashwriter", func(b *testing.B) {
					b.ResetTimer()
					sum := make([]byte, s)
					buf := bytes.NewBuffer(sum)
					bench, _ := stream.NewHashWriter(
						buf,
						stream.WithMD5(),
						stream.WithSHA256(),
					)
					for i := 0; i < b.N; i++ {
						bench.Reset()
						bench.Write(bs[:size])
						bench.Hash(stream.SHA256)
					}
				})

			})

		})
	}

}
