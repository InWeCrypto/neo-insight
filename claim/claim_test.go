package claim

import (
	"math"
	"testing"
)

func TestGenerateCas(t *testing.T) {
	generated := float64(0)
	for i := int64(855786); i < 855792; i++ {
		tmp := generateGas(i)

		generated += tmp
	}

	gas := generated / totalNEO

	math.Trunc(gas * math.Pow10(8))

	math.Floor(4.9)

	println(math.Floor(4.9))
}
