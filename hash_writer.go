package stream

import (
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"sync"

	logger "github.com/da-moon/go-logger"
	stacktrace "github.com/palantir/stacktrace"
)

// HashWriter implements an io.writer that hashes
// the payload as operating on data
// while any hash.Hasher can be used as underlying hashing algorithm,
// at this point utility functions for md5 and sha256n are written
type HashWriter interface {
	io.Writer
	Hash(h HashAlgorithm) ([]byte, error)
	HexString(h HashAlgorithm) (string, error)
	Base64String(h HashAlgorithm) (string, error)
	Reset()
}

type hashWriter struct {
	writer     io.Writer
	md5Hash    hash.Hash
	sha256Hash hash.Hash
	hasher     map[HashAlgorithm]hash.Hash
	logger     *logger.WrappedLogger `json:"-"`
}

// NewWriter ...
func NewWriter(writer io.Writer, opts ...HashWriterOption) HashWriter {
	result := &hashWriter{
		hasher: make(map[HashAlgorithm]hash.Hash),
		writer: writer,
	}
	for _, opt := range opts {
		opt(result)
	}
	// [TODO] => add safety checks
	return result
}

// Write writes to the underlying stream and hashes the data as it is writing it
func (w *hashWriter) Write(p []byte) (int, error) {
	var (
		n   int
		err error
	)
	n, err = w.writer.Write(p)
	if err != nil {
		err = stacktrace.Propagate(err, "could not write to underlying stream")
		return -1, err
	}
	if n > 0 {
		// background concurrent processing of hash jobs
		wgErr := make(chan error)
		wgDone := make(chan bool)
		wg := sync.WaitGroup{}
		for k, v := range w.hasher {
			wg.Add(1)
			go func(data []byte, algorithm HashAlgorithm, h hash.Hash) {
				_, err = h.Write(p[:n])
				if err != nil {
					err = stacktrace.Propagate(err, "could not calculate hash of the written data for '%s' algorithm", string(algorithm))
					wgErr <- err
				}
				wg.Done()
			}(p[:n], k, v)
			// _, err = v.Write(p[:n])
			// if err != nil {
			// 	err = stacktrace.Propagate(err, "could not calculate hash of the written data for '%s' algorithm", string(k))
			// 	return -1, err
			// }
		}
		// waits to make sure waitgroup it done
		go func() {
			wg.Wait()
			close(wgDone)
		}()
		// blocks until a response
		select {
		case <-wgDone:
			break
		case err := <-wgErr:
			close(wgErr)
			return -1., err
		}

		// _, err = w.md5Hash.Write(p[:n])
		// if err != nil {
		// 	err = stacktrace.Propagate(err, "could not calculate md5 hash of the written data")
		// 	return -1, err
		// }
		// _, err = w.sha256Hash.Write(p[:n])
		// if err != nil {
		// 	err = stacktrace.Propagate(err, "could not calculate sha256 hash of the written data")
		// 	return -1, err
		// }
	}
	return n, nil
}

// Hash ...
func (w *hashWriter) Hash(h HashAlgorithm) ([]byte, error) {
	hasher, ok := w.hasher[h]
	if !ok {
		err := stacktrace.NewError("underlying ")
		return nil, err
	}
	return hasher.Sum(nil), nil
}

// Reset function resets underlying hashers
// [NOTICE] => this is not thread safe
func (w *hashWriter) Reset() {
	for _, v := range w.hasher {
		v.Reset()
	}
}

// HexString ...
func (w *hashWriter) HexString(h HashAlgorithm) (string, error) {
	res, err := w.Hash(h)
	if err != nil {
		err = stacktrace.Propagate(err, "could not return hex string '%s' hash of the data", string(h))
		return "", err
	}
	return hex.EncodeToString(res), nil
}

// Base64String ...
func (w *hashWriter) Base64String(h HashAlgorithm) (string, error) {
	res, err := w.Hash(h)
	if err != nil {
		err = stacktrace.Propagate(err, "could not return base64 string '%s' hash of the data", string(h))
		return "", err
	}
	return base64.StdEncoding.EncodeToString(res), nil
}
