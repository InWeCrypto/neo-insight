package claim

import (
	"fmt"
	"math"
	"sort"

	"github.com/inwecrypto/neogo"
)

type utxoSorter []*neogo.UTXO

func (s utxoSorter) Len() int {
	return len(s)
}
func (s utxoSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s utxoSorter) Less(i, j int) bool {
	return s[i].Block < s[j].Block
}

var generation = []uint{8, 7, 6, 5, 4, 3, 2, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

const decrementInterval = 2000000

const totalNEO = 100000000

func generateGas(id int64) float64 {

	step := int(math.Floor(float64(id) / float64(decrementInterval)))

	if step > len(generation) {
		return 0
	}

	return float64(generation[step])
}

// GetStartBlock get unclaimed utxos start block number
func GetStartBlock(unclaimed []*neogo.UTXO) int64 {

	if len(unclaimed) == 0 {
		panic("unclaimed utxo must be nonzero")
	}

	sort.Sort(utxoSorter(unclaimed))

	return unclaimed[0].Block
}

// GetBlocksFee .
type GetBlocksFee func(start, end int64) (float64, int64, error)

func getUnClaimedGas(start, end int64) float64 {

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

type blockFeeSorter []*neogo.BlockFee

func (s blockFeeSorter) Len() int {
	return len(s)
}
func (s blockFeeSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s blockFeeSorter) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// GetUnClaimedGas .
func GetUnClaimedGas(
	unclaimed []*neogo.UTXO,
	getBlocksFee GetBlocksFee) (unavailable, available float64, err error) {

	for _, utxo := range unclaimed {

		sysfee, end, err := getBlocksFee(utxo.Block, utxo.SpentBlock)

		start := utxo.Block

		gas := sysfee + getUnClaimedGas(start, end)

		val, err := utxo.Value()

		if err != nil {
			return 0, 0, err
		}

		gas = val * gas / totalNEO

		if utxo.SpentBlock != -1 {
			available += gas
		} else {
			unavailable += gas
		}

		utxo.Gas = fmt.Sprintf("%.8f", round(gas, 8))
	}

	return
}

func round(f float64, n int) float64 {
	pow10n := math.Pow10(n)
	return math.Trunc(f*pow10n) / pow10n
}
