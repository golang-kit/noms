package types

import (
	"bytes"
	"io"
	"testing"

	"github.com/attic-labs/noms/Godeps/_workspace/src/github.com/stretchr/testify/assert"
)

func AssertSymEq(assert *assert.Assertions, a, b Value) {
	assert.True(a.Equals(b))
	assert.True(b.Equals(a))
}

func AssertSymNe(assert *assert.Assertions, a, b Value) {
	assert.False(a.Equals(b))
	assert.False(b.Equals(a))
}

func TestBlobLen(t *testing.T) {
	assert := assert.New(t)
	b, err := NewBlob(&bytes.Buffer{})
	assert.NoError(err)
	assert.Equal(uint64(0), b.Len())
	b, err = NewBlob(bytes.NewBuffer([]byte{0x01}))
	assert.NoError(err)
	assert.Equal(uint64(1), b.Len())
}

func TestBlobEquals(t *testing.T) {
	assert := assert.New(t)
	b1, _ := NewBlob(bytes.NewBuffer([]byte{0x01}))
	b11 := b1
	b12, _ := NewBlob(bytes.NewBuffer([]byte{0x01}))
	b2, _ := NewBlob(bytes.NewBuffer([]byte{0x02}))
	b3, _ := NewBlob(bytes.NewBuffer([]byte{0x02, 0x03}))
	AssertSymEq(assert, b1, b11)
	AssertSymEq(assert, b1, b12)
	AssertSymNe(assert, b1, b2)
	AssertSymNe(assert, b2, b3)
	AssertSymNe(assert, b1, Int32(1))
}

type testReader struct {
	readCount int
	buf       *bytes.Buffer
}

func (r *testReader) Read(p []byte) (n int, err error) {
	r.readCount++

	switch r.readCount {
	case 1:
		for i := 0; i < len(p); i++ {
			p[i] = 0x01
		}
		io.Copy(r.buf, bytes.NewReader(p))
		return len(p), nil
	case 2:
		p[0] = 0x02
		r.buf.WriteByte(p[0])
		return 1, io.EOF
	default:
		return 0, io.EOF
	}
}

func TestBlobFromReaderThatReturnsDataAndError(t *testing.T) {
	// See issue #264.
	// This tests the case of building a Blob from a reader who returns both data and an error for the final Read() call.
	assert := assert.New(t)
	tr := &testReader{buf: &bytes.Buffer{}}

	b, err := NewBlob(tr)
	assert.NoError(err)

	actual := &bytes.Buffer{}
	io.Copy(actual, b.Reader())

	assert.True(bytes.Equal(actual.Bytes(), tr.buf.Bytes()))
	assert.Equal(byte(2), actual.Bytes()[len(actual.Bytes())-1])
}
