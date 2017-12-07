package claim

import (
	"math"
	"testing"
)

func TestGenerateCas(t *testing.T) {
	generated := float64(0)
	for i := int64(855983); i < 855989; i++ {
		tmp := generateGas(i)

		generated += tmp
	}

	gas := generated / totalNEO

	math.Trunc(gas * math.Pow10(8))

	data := 0.00000048

	println(round(data, 8))
}
