package repository

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_mapSlice(t *testing.T) {
	t.Parallel()

	require.Equal(t,
		[]string{"1", "2", "3"},
		mapSlice([]int{1, 2, 3}, func(i int) string {
			return strconv.Itoa(i)
		}))
}
