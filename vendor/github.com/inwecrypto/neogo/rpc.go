package neogo

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"

	"github.com/dynamicgo/slf4go"
	"github.com/inwecrypto/jsonrpc"
)

// Client neo jsonrpc 2.0 client
type Client struct {
	jsonrpcclient *jsonrpc.RPCClient
	slf4go.Logger
}

// NewClient create new neo client
func NewClient(url string) *Client {
	return &Client{
		jsonrpcclient: jsonrpc.NewRPCClient(url),
		Logger:        slf4go.Get("neo-rpc-client"),
	}
}

func (client *Client) call(method string, result interface{}, args ...interface{}) error {

	var buff bytes.Buffer

	buff.WriteString(fmt.Sprintf("jsonrpc call: %s\n", method))
	buff.WriteString(fmt.Sprintf("\tresult: %v\n", reflect.TypeOf(result)))

	for i, arg := range args {
		buff.WriteString(fmt.Sprintf("\targ(%d): %v\n", i, arg))
	}

	client.Debug(buff.String())

	response, err := client.jsonrpcclient.Call(method, args...)

	if err != nil {
		return err
	}

	if response.Error != nil {
		return fmt.Errorf("rpc error : %d %s %v", response.Error.Code, response.Error.Message, response.Error.Data)
	}

	return response.GetObject(result)
}

// GetAccountState get account state using jsonrpc :http://docs.neo.org/zh-cn/node/api/getaccountstate.html
func (client *Client) GetAccountState(account string) (state *AccountSate, err error) {

	err = client.call("getaccountstate", &state, account)

	return
}

// GetAssetState get asset state using jsonrpc :http://docs.neo.org/zh-cn/node/api/getassetstate.html
func (client *Client) GetAssetState(asset string) (state *AssetState, err error) {

	err = client.call("getassetstate", &state, asset)

	return
}

// GetConnectionCount get node's connection count http://docs.neo.org/zh-cn/node/api/getconnectioncount.html
func (client *Client) GetConnectionCount() (count int, err error) {

	err = client.call("getbalance", &count)

	return
}

// GetBestBlockHash get neo last block hash http://docs.neo.org/zh-cn/node/api/getbestblockhash.html
func (client *Client) GetBestBlockHash() (hash string, err error) {

	err = client.call("getbestblockhash", &hash)

	return
}

// GetTxOut get tx vout http://docs.neo.org/zh-cn/node/api/gettxout.html
func (client *Client) GetTxOut(txid string, n uint) (vout *Vout, err error) {
	err = client.call("gettxout", &vout, txid, n)

	return
}

// GetBlock get neo block info http://docs.neo.org/zh-cn/node/api/getblock.html
func (client *Client) GetBlock(hash string) (block *Block, err error) {
	err = client.call("getblock", &block, hash, 1)

	return
}

// GetBlockByIndex get neo block info http://docs.neo.org/zh-cn/node/api/getblock2.html
func (client *Client) GetBlockByIndex(index int64) (block *Block, err error) {
	err = client.call("getblock", &block, index, 1)

	return
}

// GetBlockCount get neo count info http://docs.neo.org/zh-cn/node/api/getblockcount.html
func (client *Client) GetBlockCount() (count int64, err error) {
	err = client.call("getblockcount", &count)

	return
}

// GetRawTransaction get transaction with txid http://docs.neo.org/zh-cn/node/api/getrawtransaction.html
func (client *Client) GetRawTransaction(txid string) (trans *Transaction, err error) {
	err = client.call("getrawtransaction", &trans, txid, 1)

	return
}

// GetPeers  .
func (client *Client) GetPeers() (data interface{}, err error) {
	err = client.call("getpeers", &data)

	return
}

// SendRawTransaction send raw transaction with jsonrpc api:http://docs.neo.org/zh-cn/node/api/sendrawtransaction.html
func (client *Client) SendRawTransaction(data []byte) (status bool, err error) {
	err = client.call("sendrawtransaction", &status, hex.EncodeToString(data))

	return
}

// GetBalance extend rpc method get address's utxos
func (client *Client) GetBalance(address string, asset string) (utxos []*UTXO, err error) {
	err = client.call("balance", &utxos, address, asset)

	return
}

// GetClaim get unclaimed utxos
func (client *Client) GetClaim(address string) (unclaimed *Unclaimed, err error) {
	err = client.call("claim", &unclaimed, address)

	return
}
