
package main

import (
	"fmt"
	"../rpc"
	"../internal/swcapi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common"
	"log"
)

func main() {

	client, err := rpc.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("rpc.Dial err", err)
		return
	}

	/*var  blockNumber uint64
	var  swc_coinbase string
	var  clientversion string
	var  block map[string]interface{}

	err = client.Call(&blockNumber, "eth_blockNumber")
	err = client.Call(&swc_coinbase, "eth_coinbase")
	err = client.Call(&clientversion, "web3_clientVersion")
	var bigint = new(big.Int).SetInt64(0)
	var bigNumber = (*hexutil.Big)(bigint)
	err = client.Call(&block, "eth_getBlockByNumber",bigNumber,false)

	var  balance = int64(0)

	var bigint1 = new(big.Int).SetInt64(0)
	var bigNumber1 = (*hexutil.Big)(bigint1)
	err = client.Call(&balance, "eth_getBalance","1Q1oECL9rvC642THhNB6QZMqU55fDieXDK",bigNumber1)
*/
	var  txhash common.Hash
	var  txhash1 common.Hash
	sendTx := new(swcapi.SendTxArgs)
	sendTx1 := new(swcapi.SendTxArgs)
	var from = common.HexToAddress("0x6bC1B2C682c66B046903131e64B2FA4e58ae4ec8")
	//var from = core.Base58ToCommonAddress([]byte("1G7EmF7Umd96FLMKh3PhqZCi3bfMzqC4tH"))
	log.Println("==>from：", from.String())

	//var to = core.Base58ToCommonAddress([]byte("1Q1oECL9rvC642THhNB6QZMqU55fDieXDK"))

	var to = common.HexToAddress("0xAEe4dF5C595900708c0676445bFfd67B0E4C7B66")
	var from1 = common.HexToAddress("0x6bC1B2C682c66B046903131e64B2FA4e58ae4ec8")
	//var to1 = core.Base58ToCommonAddress([]byte("1G7EmF7Umd96FLMKh3PhqZCi3bfMzqC4tH"))
	var to1 = common.HexToAddress("0xe82412b3ea6c345142db5d67060800d7d74269b2")


	//var nonce = hexutil.Uint64(1)
	//var bigi = new(big.Int).SetInt64(5)
	//var value = (*hexutil.Big)(bigi)
	var value =  float64(0.00000001)
	var value1 = float64(40)
	var data = hexutil.Bytes{}
	sendTx.From = from
	sendTx.To = &to
	sendTx.Nonce = nil
	sendTx.Value = value
	sendTx.Data = &data

	sendTx1.From = from1
	sendTx1.To = &to1
	sendTx1.Nonce = nil
	sendTx1.Value = value1
	sendTx1.Data = &data

	err = client.Call(&txhash, "personal_sendTransaction",sendTx,"")
	err = client.Call(&txhash1, "personal_sendTransaction",sendTx1,"")
/*
	fmt.Println("eth_blockNumber ", blockNumber)
	fmt.Println("eth_coinbase ", swc_coinbase)
	fmt.Println("web3_clientversion ", clientversion)
	fmt.Println("eth_getBlockByNumber ", block)
	fmt.Println("eth_getBalance ",balance )*/
	fmt.Println("personal_sendTransaction ",txhash )
	fmt.Println("personal_sendTransaction1 ",txhash1 )

	/*var account[]string
	err = client.Call(&account, "eth_accounts")
	var result string
	//var result hexutil.Big
	err = client.Call(&result, "eth_getBalance", account[0], "latest")
	//err = ec.c.CallContext(ctx, &result, "eth_getBalance", account, "latest")

	if err != nil {
		fmt.Println("client.Call err", err)
		return
	}
*/
	//fmt.Printf("account[0]: %s\nbalance[0]: %s\n", account[0], result)
	//fmt.Printf("accounts: %s\n", account[0])

	/*
	curl --header "Content-Type:application/json"  --data {\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\":[\"0x407\", \"latest\"],\"id\":1} http://localhost:8545
	 curl -l -H "Content-Type:application/json" -H "Accept:application/json" -X POST -d {\"jsonrpc\":\"2.0\",\"method\":\"swc_coinbase\",\"id\":1} http://localhost:8545

	*/
}
