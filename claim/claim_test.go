package claim

import (
	"fmt"
	"testing"

	"github.com/inwecrypto/neogo/tx"
)

func TestGenerateCas(t *testing.T) {
	val := float64(2)
	generated := float64(0)
	for i := int64(1992850); i < 2001227; i++ {
		tmp := generateGas(i)

		generated += tmp
	}

	gas := tx.MakeFixed8(val * (2 + generated) / totalNEO)

	fmt.Printf("%v\n", gas.String())
}
