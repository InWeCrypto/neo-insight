package claim

import (
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

// GetBlockFee .
type GetBlockFee func(id int64) (*neogo.BlockFee, error)

func getUnClaimedGas(
	unclaimed *neogo.UTXO,
	bestBlockFee,
	currentBlockFee *neogo.BlockFee) float64 {

	generated := float64(0)

	for i := unclaimed.Block; i < bestBlockFee.ID; i++ {
		tmp := generateGas(i + 1)

		if tmp == 0 {
			break
		}

		generated += tmp
	}

	return bestBlockFee.SysFee - currentBlockFee.SysFee + generated
}

// GetUnClaimedGas .
func GetUnClaimedGas(
	unclaimed []*neogo.UTXO,
	bestBlockFee *neogo.BlockFee,
	getBlockFee GetBlockFee) (unavailable, available float64, err error) {

	for _, utxo := range unclaimed {

		endBlockFee := bestBlockFee

		if utxo.SpentBlock != -1 {
			spentBlock := utxo.SpentBlock

			if spentBlock != 0 {
				spentBlock--
			}

			endBlockFee, err = getBlockFee(spentBlock)

			if err != nil {
				return 0, 0, err
			}
		}

		block := utxo.Block

		if utxo.Block != 0 {
			block--
		}

		currentBlockFee, err := getBlockFee(block)

		if err != nil {
			return 0, 0, err
		}

		gas := getUnClaimedGas(utxo, endBlockFee, currentBlockFee)

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
	}

	return
}
