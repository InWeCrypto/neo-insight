package insight

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/dynamicgo/config"
	"github.com/inwecrypto/neogo"
)

type utxoModel struct {
	cnf *config.Config
	db  *sql.DB
}

func newUTXOModel(cnf *config.Config, db *sql.DB) *utxoModel {
	return &utxoModel{
		cnf: cnf,
		db:  db,
	}
}

func (model *utxoModel) Unclaimed(address string) ([]*neogo.UTXO, error) {
	queryStr := `
		select "json","createTime","spentTime","spentBlock","blocks" from neo_utxo 
		where address=$1
		and assert='0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b'
		and claimed=FALSE
	`

	return model.getUTXOs(queryStr, address)
}

func (model *utxoModel) Unspent(address string, asset string) ([]*neogo.UTXO, error) {
	queryStr := `
		select "json","createTime","spentTime","spentBlock","blocks" from neo_utxo 
		where address=$1
		and assert=$2 
		and "spentBlock"=-1
	`

	return model.getUTXOs(queryStr, address, asset)
}

func (model *utxoModel) getUTXOs(query string, args ...interface{}) ([]*neogo.UTXO, error) {

	logger.DebugF("getUTXOs query :%s\n\t%v", query, args)

	rows, err := model.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	assets := make([]*neogo.UTXO, 0)

	for rows.Next() {
		var (
			rawjson    string
			createTime string
			spentTime  sql.NullString
			spentBlock int64
			block      int64
		)

		err := rows.Scan(&rawjson, &createTime, &spentTime, &spentBlock, &block)

		if err != nil {
			return nil, err
		}

		var utxo *neogo.UTXO

		err = json.Unmarshal([]byte(rawjson), &utxo)

		if err != nil {
			return nil, err
		}

		utxo.CreateTime = createTime
		utxo.SpentTime = spentTime.String
		utxo.SpentBlock = spentBlock
		utxo.Block = block

		assets = append(assets, utxo)
	}

	return assets, nil
}

type blockFeeModel struct {
	cnf *config.Config
	db  *sql.DB
}

func newBlockFeeModel(cnf *config.Config, db *sql.DB) *blockFeeModel {
	return &blockFeeModel{
		cnf: cnf,
		db:  db,
	}
}

func (model *blockFeeModel) GetBlockFee(id int64) (*neogo.BlockFee, error) {
	return model.getBlockFee(`select id,sysfree,netfee,"createTime" from neo_block where id=$1`, id)
}

func (model *blockFeeModel) GetBestBlockFee() (*neogo.BlockFee, error) {
	return model.getBlockFee(`select id,sysfree,netfee,"createTime" from neo_block order by id desc limit 1`)
}

func (model *blockFeeModel) getBlockFee(query string, args ...interface{}) (*neogo.BlockFee, error) {

	logger.DebugF("getUTXOs query :%s with args %v", query, args)

	rows, err := model.db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		var (
			id         int64
			sysfee     float64
			netfee     float64
			createTime string
		)

		if err := rows.Scan(&id, &sysfee, &netfee, &createTime); err != nil {
			return nil, err
		}

		return &neogo.BlockFee{
			ID:         id,
			SysFee:     sysfee,
			NetFee:     netfee,
			CreateTime: createTime,
		}, nil
	}

	return nil, fmt.Errorf("block fee not found")
}
