package claim

import (
	"testing"
)

func TestGenerateCas(t *testing.T) {
	generated := float64(0)
	for i := int64(858279); i < 858335; i++ {
		tmp := generateGas(i)

		generated += tmp
	}

	gas := generated / totalNEO

	println(gas)
}
