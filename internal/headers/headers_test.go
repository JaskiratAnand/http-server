package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:        barbar            \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069", host)

	foofoo, _ := headers.Get("FooFoo")
	assert.Equal(t, "barbar", foofoo)

	missingKey, _ := headers.Get("MissingKey")
	assert.Equal(t, "", missingKey)
	assert.Equal(t, 60, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Header Name
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Duplicate Header Fields
	headers = NewHeaders()
	data = []byte("Host: localhost:42069\r\nHost: localhost:8000\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ = headers.Get("Host")
	assert.Equal(t, "localhost:42069,localhost:8000", host)
	assert.False(t, done)
}
