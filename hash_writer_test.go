package stream_test

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strconv"
	"testing"

	stream "github.com/da-moon/go-stream"
	"github.com/stretchr/testify/assert"
)

type hashTest struct {
	in string
}

// these cases were taken out of stdlib's sha256 package
var smallGolden = []hashTest{
	{""},
	{"a"},
	{"ab"},
	{"abc"},
	{"abcd"},
	{"abcde"},
	{"abcdef"},
	{"abcdefg"},
	{"abcdefgh"},
	{"abcdefghi"},
	{"abcdefghij"},
	{"Discard medicine more than two years old."},
	{"He who has a shady past knows that nice guys finish last."},
	{"I wouldn't marry him with a ten foot pole."},
	{"Free! Free!/A trip/to Mars/for 900/empty jars/Burma Shave"},
	{"The days of the digital watch are numbered.  -Tom Stoppard"},
	{"Nepal premier won't resign."},
	{"For every action there is an equal and opposite government program."},
	{"His money is twice tainted: 'taint yours and 'taint mine."},
	{"There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"size:  a.out:  bad magic"},
	{"The major problem is with sendmail.  -Mark Horton"},
	{"Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"If the enemy is within range, then so are you."},
	{"It's well we cannot hear the screams/That we create in others' dreams."},
	{"You remind me of a TV show, but that's all right: I watch it anyway."},
	{"C is as portable as Stonehedge!!"},
	{"Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
	{"How can you write a big system without C++?  -Paul Glick"},
}

func TestHashWriter(t *testing.T) {

	buf := bytes.NewBuffer(make([]byte, 0))
	hashWriter := stream.NewWriter(buf, stream.WithSHA256())
	t.Run("io.WriterImpl", func(t *testing.T) {
		var _ io.Writer = hashWriter
	})
	stdlibSha256 := sha256.New()
	stdlibMD5 := md5.New()
	for i := 0; i < len(smallGolden); i++ {
		g := smallGolden[i]
		t.Run("FullWrite-"+strconv.Itoa(i), func(t *testing.T) {
			_, err := hashWriter.Write([]byte(g.in))
			assert.NoError(t, err)
			_, err = stdlibSha256.Write([]byte(g.in))
			assert.NoError(t, err)
			_, err = stdlibMD5.Write([]byte(g.in))
			assert.NoError(t, err)
			t.Run("bytes", func(t *testing.T) {
				t.Run("SHA256", func(t *testing.T) {
					actual, err := hashWriter.Hash(stream.SHA256)
					assert.NoError(t, err)
					assert.NotNil(t, actual)
					expected := stdlibSha256.Sum(nil)
					assert.NotNil(t, expected)
					assert.Equal(t, expected, actual)
				})
			})
			t.Run("hex", func(t *testing.T) {
				t.Run("SHA256", func(t *testing.T) {
					actual, err := hashWriter.HexString(stream.SHA256)
					assert.NoError(t, err)
					assert.NotZero(t, len(actual))
					expected := hex.EncodeToString(stdlibSha256.Sum(nil))
					assert.NotZero(t, len(expected))
					assert.Equal(t, expected, actual)
				})
			})
			t.Run("base64", func(t *testing.T) {
				t.Run("SHA256", func(t *testing.T) {
					actual, err := hashWriter.Base64String(stream.SHA256)
					assert.NoError(t, err)
					assert.NotZero(t, len(actual))
					expected := base64.StdEncoding.EncodeToString(stdlibSha256.Sum(nil))
					assert.NotZero(t, len(expected))
					assert.Equal(t, expected, actual)
				})
			})
			hashWriter.Reset()
			stdlibSha256.Reset()
			stdlibMD5.Reset()
		})
		t.Run("PartialWrite-"+strconv.Itoa(i), func(t *testing.T) {
			var err error
			for j := 0; j < 3; j++ {
				if j < 2 {
					_, err = io.WriteString(hashWriter, g.in)
					assert.NoError(t, err)
					_, err = io.WriteString(stdlibSha256, g.in)
					assert.NoError(t, err)
					_, err = io.WriteString(stdlibMD5, g.in)
					assert.NoError(t, err)
				} else {
					_, err = io.WriteString(hashWriter, g.in[0:len(g.in)/2])
					assert.NoError(t, err)
					_, err = io.WriteString(stdlibSha256, g.in[0:len(g.in)/2])
					assert.NoError(t, err)
					_, err = io.WriteString(stdlibMD5, g.in[0:len(g.in)/2])
					assert.NoError(t, err)
				}
				t.Run("bytes", func(t *testing.T) {
					t.Run("SHA256", func(t *testing.T) {
						actual, err := hashWriter.Hash(stream.SHA256)
						assert.NoError(t, err)
						assert.NotNil(t, actual)
						expected := stdlibSha256.Sum(nil)
						assert.NotNil(t, expected)
						assert.Equal(t, expected, actual)
					})
				})
				t.Run("hex", func(t *testing.T) {
					t.Run("SHA256", func(t *testing.T) {
						actual, err := hashWriter.HexString(stream.SHA256)
						assert.NoError(t, err)
						assert.NotZero(t, len(actual))
						expected := hex.EncodeToString(stdlibSha256.Sum(nil))
						assert.NotZero(t, len(expected))
						assert.Equal(t, expected, actual)
					})
				})
				t.Run("base64", func(t *testing.T) {
					t.Run("SHA256", func(t *testing.T) {
						actual, err := hashWriter.Base64String(stream.SHA256)
						assert.NoError(t, err)
						assert.NotZero(t, len(actual))
						expected := base64.StdEncoding.EncodeToString(stdlibSha256.Sum(nil))
						assert.NotZero(t, len(expected))
						assert.Equal(t, expected, actual)
					})
				})

			}
			hashWriter.Reset()
			stdlibSha256.Reset()
			stdlibMD5.Reset()
		})
	}
}
