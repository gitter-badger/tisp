package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPartial(t *testing.T) {
	ifFunc := func(ts ...*Thunk) bool {
		b := App(App(Partial, If, False, True), ts...)
		return bool(b.Eval().(boolType))
	}

	assert.True(t, ifFunc(True))
	assert.True(t, !ifFunc(False))
}
