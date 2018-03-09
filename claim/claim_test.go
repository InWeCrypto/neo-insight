package claim

import (
	"fmt"
	"testing"
)

func getUnClaimedGas2(start, end int64) float64 {
	generated := float64(0)
	for i := start; i < end; i++ {
		tmp := generateGas(i + 1)

		if tmp == 0 {
			break
		}

		generated += tmp
	}

	return generated
}

func getUnClaimedGas3(start, end int64) float64 {
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
	val := float64(2)
	// generated := getUnClaimedGas2(1992850, 2001227)

	// generated2 := getUnClaimedGas3(1992850, 2001227)

	// require.Equal(t, generated, generated2)

	// for i := int64(1992850); i < int64(2001227); i++ {
	// 	println(i, generateGas(i))
	// }

	// println(generateGas(1992850), generateGas(2001227))

	gas := round((val * (getUnClaimedGas(1992850, 2001227) + 2) / totalNEO), 8)

	fmt.Printf("%v\n", gas)
}
