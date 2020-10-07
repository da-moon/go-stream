package stream

import (
	"hash"
	"log"

	md5 "github.com/minio/md5-simd"
	sha256 "github.com/minio/sha256-simd"

	logger "github.com/da-moon/go-logger"
	stacktrace "github.com/palantir/stacktrace"
)

// HashAlgorithm provides a convenient way
// to have consistency when accessing default hash algorithms
type HashAlgorithm string

// NewHashAlgorithm is a utility function used for casting a string to hash algorithm enum
func NewHashAlgorithm(arg string) HashAlgorithm {
	return HashAlgorithm(arg)
}

//  list of hash algorithms with helper methods
var (
	// SHA256 ...
	SHA256 HashAlgorithm = "SHA256"
	// MD5 ...
	MD5 HashAlgorithm = "MD5"
)

// HashWriterOption ...
type HashWriterOption func(*hashWriter) error

// WithLogger ...
func WithLogger(arg *log.Logger) HashWriterOption {
	return func(s *hashWriter) error {
		if arg == nil {
			return stacktrace.NewError("passed logger was nil")
		}
		s.logger = logger.NewWrappedLogger(arg)
		return nil

	}
}

// WithWrappedLogger ...
func WithWrappedLogger(arg *logger.WrappedLogger) HashWriterOption {
	return func(s *hashWriter) error {
		if arg == nil {
			return stacktrace.NewError("passed logger was nil")
		}
		s.logger = arg
		return nil
	}
}

// WithHasher sets underlying hash algorithms
// [WARN] => make sure s.hasher is allocated in memory before calling this function
func WithHasher(algorithm HashAlgorithm, hasher hash.Hash) HashWriterOption {
	return func(s *hashWriter) error {

		if string(algorithm) == "" {
			return stacktrace.NewError("an empty string was passed as hasher algorithm")
		}
		if hasher == nil {
			return stacktrace.NewError("a nil hasher was passed for '%v' algorithm", string(algorithm))
		}
		_, ok := s.hasher[algorithm]
		if ok {
			return stacktrace.NewError("there is an existing hasher for '%v' algorithm", string(algorithm))
		}
		s.hasher[algorithm] = hasher
		return nil
	}
}

// WithSHA256 is a convenience method that
// adds sha256 hashing based on minio's library.
func WithSHA256() HashWriterOption {
	return WithHasher(SHA256, sha256.New())
}

var (
	// used for minio md5 shasher
	md5Server md5.Server
	md5Hasher md5.Hasher
)

// WithMD5 is a sets up md5 hashing with minio's md5 Hasher.
func WithMD5() HashWriterOption {
	md5Server = md5.NewServer()
	md5Hasher = md5Server.NewHash()
	return WithHasher(MD5, md5Hasher)
}

// ShutdownMD5Hasher must be called when we are done with using the md5hasher
func ShutdownMD5Hasher() {
	md5Hasher.Close()
	md5Server.Close()
}
