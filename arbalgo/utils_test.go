package arbalgo

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrangeCycle(t *testing.T) {
	cases := []struct {
		in  []int
		out []int
	}{
		{
			in:  []int{},
			out: []int{},
		},
		{
			in:  []int{1},
			out: []int{1},
		},
		{
			in:  []int{4, 5, 6, 1, 3},
			out: []int{1, 3, 4, 5, 6},
		},
	}
	for _, c := range cases {
		assert.True(t, reflect.DeepEqual(c.out, ArrangeCycle(c.in)))
	}
}
