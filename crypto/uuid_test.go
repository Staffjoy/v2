package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const format = "^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$"

func TestParse(t *testing.T) {
	assert := assert.New(t)
	_, err := Parse([]byte{1, 2, 3, 4, 5})
	assert.NotNil(err)
	base, err := NewUUID()
	assert.NoError(err)

	u, err := Parse(base[:])
	assert.NoError(err, "Expected to parse UUID sequence without problems")

	assert.Equal(u.String(), base.String(), "Expected parsed UUID to be the same as base, %s != %s", u.String(), base.String())
}

func TestParseString(t *testing.T) {
	assert := assert.New(t)
	_, err := ParseHex("foo")
	assert.NotNil(err, "Expected error due to invalid UUID string")

	base, err := NewUUID()
	assert.NoError(err)

	u, err := ParseHex(base.String())
	assert.NoError(err, "Expected to parse UUID sequence without problems")
	assert.Equal(u.String(), base.String(), "Expected parsed UUID to be the same as base, %s != %s", u.String(), base.String())

}

func BenchmarkParseHex(b *testing.B) {
	s := "f3593cff-ee92-40df-4086-87825b523f13"
	for i := 0; i < b.N; i++ {
		_, err := ParseHex(s)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.ReportAllocs()
}
