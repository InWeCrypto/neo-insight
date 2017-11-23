package claim

import "testing"

func TestGenerateCas(t *testing.T) {
	generated := float64(0)
	for i := int64(802873); i < 805360; i++ {
		tmp := generateGas(i)

		generated += tmp
	}

	gas := generated / totalNEO

	println(gas)
}
