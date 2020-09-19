# go-stream

<p align="center">
  <a href="https://gitpod.io#https://github.com/da-moon/go-stream">
    <img src="https://img.shields.io/badge/open%20in-gitpod-blue?logo=gitpod" alt="Open In GitPod">
  </a>
  <img src="https://img.shields.io/github/languages/code-size/da-moon/go-stream" alt="GitHub code size in bytes">
  <img src="https://img.shields.io/github/commit-activity/w/da-moon/go-stream" alt="GitHub commit activity">
  <img src="https://img.shields.io/github/last-commit/da-moon/go-stream/master" alt="GitHub last commit">
</p>

utility library with different implementations of I/O primitives for various common tasks.This package is experimental. use at your own risk.

- `HashWriter` an `io.Writer` that stores hash of data as it is getting passed through it. It can use multiple `hash.Hash` types to calculate hash of the stream.

# Benchmark

- `HashWriter`  benchmarks with standalone `stdlib` md5 and sha256 hash functions  , HashWriter backed by `stdlib` md5 and sha256 hash functions and minio's implementation of md5 and sha256 hash functions (go version `go1.15` linux/amd64) 

```
Running tool: /home/gitpod/go/bin/go test -benchmem -run=^$ github.com/da-moon/go-stream -bench ^(BenchmarkHashWriter)$ -v
goos: linux
goarch: amd64
pkg: github.com/da-moon/go-stream
BenchmarkHashWriter
BenchmarkHashWriter/8
BenchmarkHashWriter/8/single-hash
BenchmarkHashWriter/8/single-hash/SHA256
BenchmarkHashWriter/8/single-hash/SHA256/stdlib
BenchmarkHashWriter/8/single-hash/SHA256/stdlib-16         	 4466779	       265 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashWriter/8/single-hash/SHA256/hashwriter-stdlib
BenchmarkHashWriter/8/single-hash/SHA256/hashwriter-stdlib-16         	 2564110	       453 ns/op	      48 B/op	       1 allocs/op
BenchmarkHashWriter/8/single-hash/SHA256/hashwriter-minio
BenchmarkHashWriter/8/single-hash/SHA256/hashwriter-minio-16          	 2858588	       433 ns/op	      61 B/op	       1 allocs/op
BenchmarkHashWriter/8/single-hash/MD5
BenchmarkHashWriter/8/single-hash/MD5/stdlib
BenchmarkHashWriter/8/single-hash/MD5/stdlib-16                       	 8055064	       151 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashWriter/8/single-hash/MD5/hashwriter-stdlib
BenchmarkHashWriter/8/single-hash/MD5/hashwriter-stdlib-16            	  483354	      2642 ns/op	     546 B/op	      13 allocs/op
BenchmarkHashWriter/8/single-hash/MD5/hashwriter-minio
BenchmarkHashWriter/8/single-hash/MD5/hashwriter-minio-16             	  309547	      4026 ns/op	     545 B/op	      13 allocs/op
BenchmarkHashWriter/8/multi-hash
BenchmarkHashWriter/8/multi-hash/stdlib
BenchmarkHashWriter/8/multi-hash/stdlib-16                            	 2363575	       525 ns/op	      48 B/op	       2 allocs/op
BenchmarkHashWriter/8/multi-hash/hashwriter
BenchmarkHashWriter/8/multi-hash/hashwriter-16                        	 1000000	      1022 ns/op	      50 B/op	       1 allocs/op
BenchmarkHashWriter/1024
BenchmarkHashWriter/1024/single-hash
BenchmarkHashWriter/1024/single-hash/SHA256
BenchmarkHashWriter/1024/single-hash/SHA256/stdlib
BenchmarkHashWriter/1024/single-hash/SHA256/stdlib-16                 	  363068	      3343 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashWriter/1024/single-hash/SHA256/hashwriter-stdlib
BenchmarkHashWriter/1024/single-hash/SHA256/hashwriter-stdlib-16      	  270160	      7156 ns/op	    2081 B/op	       1 allocs/op
BenchmarkHashWriter/1024/single-hash/SHA256/hashwriter-minio
BenchmarkHashWriter/1024/single-hash/SHA256/hashwriter-minio-16       	  276546	     29231 ns/op	    4036 B/op	       1 allocs/op
BenchmarkHashWriter/1024/single-hash/MD5
BenchmarkHashWriter/1024/single-hash/MD5/stdlib
BenchmarkHashWriter/1024/single-hash/MD5/stdlib-16                    	  670987	      1921 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashWriter/1024/single-hash/MD5/hashwriter-stdlib
BenchmarkHashWriter/1024/single-hash/MD5/hashwriter-stdlib-16         	  222818	      5038 ns/op	    2967 B/op	      13 allocs/op
BenchmarkHashWriter/1024/single-hash/MD5/hashwriter-minio
BenchmarkHashWriter/1024/single-hash/MD5/hashwriter-minio-16          	  119022	      9511 ns/op	    3688 B/op	      20 allocs/op
BenchmarkHashWriter/1024/multi-hash
BenchmarkHashWriter/1024/multi-hash/stdlib
BenchmarkHashWriter/1024/multi-hash/stdlib-16                         	  209500	      5610 ns/op	      48 B/op	       2 allocs/op
BenchmarkHashWriter/1024/multi-hash/hashwriter
BenchmarkHashWriter/1024/multi-hash/hashwriter-16                     	  120184	      9780 ns/op	    3160 B/op	       8 allocs/op
BenchmarkHashWriter/8192
BenchmarkHashWriter/8192/single-hash
BenchmarkHashWriter/8192/single-hash/SHA256
BenchmarkHashWriter/8192/single-hash/SHA256/stdlib
BenchmarkHashWriter/8192/single-hash/SHA256/stdlib-16                 	   48978	     25139 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashWriter/8192/single-hash/SHA256/hashwriter-stdlib
BenchmarkHashWriter/8192/single-hash/SHA256/hashwriter-stdlib-16      	   32734	     31755 ns/op	   16494 B/op	       1 allocs/op
BenchmarkHashWriter/8192/single-hash/SHA256/hashwriter-minio
BenchmarkHashWriter/8192/single-hash/SHA256/hashwriter-minio-16       	   38716	     44113 ns/op	   27871 B/op	       1 allocs/op
BenchmarkHashWriter/8192/single-hash/MD5
BenchmarkHashWriter/8192/single-hash/MD5/stdlib
BenchmarkHashWriter/8192/single-hash/MD5/stdlib-16                    	   82287	     14096 ns/op	       0 B/op	       0 allocs/op
BenchmarkHashWriter/8192/single-hash/MD5/hashwriter-stdlib
BenchmarkHashWriter/8192/single-hash/MD5/hashwriter-stdlib-16         	   47614	     24452 ns/op	   23113 B/op	      13 allocs/op
BenchmarkHashWriter/8192/single-hash/MD5/hashwriter-minio
BenchmarkHashWriter/8192/single-hash/MD5/hashwriter-minio-16          	   35079	     36130 ns/op	   32095 B/op	      20 allocs/op
BenchmarkHashWriter/8192/multi-hash
BenchmarkHashWriter/8192/multi-hash/stdlib
BenchmarkHashWriter/8192/multi-hash/stdlib-16                         	   29565	     39082 ns/op	      48 B/op	       2 allocs/op
BenchmarkHashWriter/8192/multi-hash/hashwriter
BenchmarkHashWriter/8192/multi-hash/hashwriter-16                     	   22287	     49527 ns/op	   25075 B/op	       8 allocs/op
PASS
ok  	github.com/da-moon/go-stream	41.723s
```