package claim

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func getUnClaimedGas2(start, end int64) float64 {
	generated := float64(0)
	for i := start; i < end; i++ {
		tmp := generateGas(i)

		if tmp == 0 {
			break
		}

		generated += tmp
	}

	return generated
}

func TestGenerateCas(t *testing.T) {
	// val := float64(2)
	generated := getUnClaimedGas2(1992850, 2001227)

	generated2 := getUnClaimedGas(1992850, 2001227)

	require.Equal(t, generated, generated2)

	// gas := round((val * (generated + 2) / totalNEO), 8)

	// fmt.Printf("%v\n", gas)
}
