package stream

import (
	"encoding/base64"
	"encoding/hex"
	"hash"
	"io"
	"log"
	"os"
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
	writer io.Writer
	hasher map[HashAlgorithm]hash.Hash
	logger *logger.WrappedLogger `json:"-"`
}

// NewHashWriter ...
func NewHashWriter(writer io.Writer, opts ...HashWriterOption) (HashWriter, error) {
	var err error
	result := &hashWriter{
		hasher: make(map[HashAlgorithm]hash.Hash),
		writer: writer,
	}
	for _, opt := range opts {
		opt(result)
	}
	if result.logger == nil {
		l := log.New(logger.NewLevelFilter(
			logger.WithWriter(os.Stderr),
		), "", log.LstdFlags)
		result.logger = logger.NewWrappedLogger(l)
	}
	if len(result.hasher) == 0 {
		err = stacktrace.NewError("no underlying hash functions were provided")
		return nil, err
	}
	if result.writer == nil {
		err = stacktrace.NewError("underlying io.Writer is nil")
		return nil, err
	}
	return result, nil
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
		for k, v := range w.hasher {
			_, err = v.Write(p[:n])
			if err != nil {
				err = stacktrace.Propagate(err, "could not calculate hash of the written data for '%s' algorithm", string(k))
				return -1, err
			}
		}
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

// WriteBackground : background concurrent processing of hash jobs
// [TODO] => optimize allocs ... maybe it can be more performant than
// synchronous Write
func (w *hashWriter) writeBackground(p []byte) (int, error) {
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

	}
	return n, nil
}
