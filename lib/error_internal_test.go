package passgen

import (
	"testing"

	"github.com/ilius/is/v2"
)

func TestErrorSpacedError_markLen0(t *testing.T) {
	is := is.New(t)
	err := NewError("", 0, "")
	err.markLen = 0
	is.Equal(err.SpacedError(), `^  error: `)
}
