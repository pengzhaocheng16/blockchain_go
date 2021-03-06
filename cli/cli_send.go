package main

import (
	"fmt"
	"log"
	"../blockchain_go"
	"../p2pprotocol"
	"time"
	"encoding/hex"
	"../blockchain_go/state"
	"github.com/ethereum/go-ethereum/common"
)

func (cli *CLI) send(from, to string, amount *float64, nodeID string, mineNow bool) {
	core.MineNow_ = mineNow
	if !core.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !core.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}
	if from == to {
		log.Panic("ERROR: Wallet from equal Wallet to is not valid")
	}
	log.Println("--start cli send ")

	var bc *core.Blockchain
	bc = core.NewBlockchain(nodeID)

	UTXOSet := core.UTXOSet{bc}
	//defer bc.Db.Close()

	wallets, err := core.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wallet := wallets.GetWallet(from)

	block := bc.CurrentBlock()
	statedb, err := bc.StateAt(block.Root())
	if(err !=nil){
		log.Panic(err)
	}
	var addr = wallet.ToCommonAddress()
	address := addr.String()

	var pendingState = state.ManageState(statedb)
	//var nonce = pendingState.GetNonce(addr)
	var nonce,_ = statedb.GetTransactionNonce(address)
	pendingState.SetNonce(addr,nonce)
	statedb.Finalise(true)
	toaddr := common.HexToAddress(to)
	tx := core.NewUTXOTransaction(nonce,&wallet, toaddr, amount,[]byte{}, &UTXOSet,nodeID)

	if mineNow {
		statedb.PutTransaction(tx.ID,tx.Serialize(),address)

		fmt.Println("==>NewCoinbaseTX ")
		//var nonce = pendingState.GetNonce(core.Base58ToCommonAddress([]byte(from)))
		//var nonce = statedb.GetNonce(addr)
		coinbaseFrom :=  common.Address{}
		var nonce,_ = statedb.GetTransactionNonce(coinbaseFrom.String())

		fromaddr := common.HexToAddress(from)
		cbTx := core.NewCoinbaseTX(nonce,fromaddr, "",nodeID)
		txs := []*core.Transaction{cbTx, tx}
		statedb.PutTransaction(cbTx.ID,cbTx.Serialize(),coinbaseFrom.String())

		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
		bc.Db.Close()
	} else {
		bc.Db.Close()
		//remove comfirmed transaction vin txs from persistent tx queue
		// in confirmTx()
		//In case of double spend check fail need to store prev uncomfirmed transaction input tx
		core.PendingIn(nodeID,tx)
		address := wallet.ToCommonAddress().String()

		statedb.PutTransaction(tx.ID,tx.Serialize(),address)
		go func(){
			if(p2pprotocol.CurrentNodeInfo == nil){
				p2pprotocol.StartServer(nodeID,"","","",0)
			}
		}()
		time.Sleep(2*time.Second)
		select{
			case ch := <- p2pprotocol.Manager.BestTd:

				bc1 := core.NewBlockchain(nodeID)
				td,_ := bc1.GetBestHeight()
				bc1.Db.Close()
				if(td.Cmp(ch) == 0){
					log.Println("---td:",ch)
					break
				}
		}
		//p2pprotocol.SendTx(core.BootNodes[0], tx)
		//go func(){
			for _, p := range p2pprotocol.Manager.Peers.Peers {
				p2pprotocol.SendTx(p, p.Rw, tx)
			}
			p2pprotocol.Manager.Mu.Lock()
			p2pprotocol.Manager.TxMempool[hex.EncodeToString(tx.ID)] = tx
		    p2pprotocol.Manager.Mu.Unlock()
		//}()
		//select{}
		for{
			var send,from,fromaddress,to,toaddress,amount string
			var amountnum float64
			fmt.Scanf("%s %s %s %s %s %s %d", &send,&from,&fromaddress,&to,&toaddress,&amount,&amountnum)

			log.Println("--send ",send)
			if(send == ""||send != "send"){
				log.Panic("need send command")
			}
			if(from == ""||from != "-from"){
				log.Panic("need from command")
			}
			if(fromaddress == ""){
				log.Panic("need fromaddress param")
			}
			if(to == ""||to != "-to"){
				log.Panic("need to command")
			}
			if(toaddress == ""){
				log.Panic("need toaddress param")
			}
			if(amount == ""||amount != "-amount"){
				log.Panic("need amount command")
			}
			if(amountnum == 0){
				log.Panic("need amount  param")
			}

			bc = core.NewBlockchain(nodeID)
			UTXOSet := core.UTXOSet{bc}
			log.Println("--send to",toaddress)

			nonce,_ := statedb.GetTransactionNonce(address)
			//var nonce = statedb.GetNonce(addr)
			pendingState.SetNonce(addr,nonce)
			statedb.Finalise(true)
			toaddr := common.HexToAddress(toaddress)
			tx := core.NewUTXOTransaction(nonce+1,&wallet, toaddr, &amountnum, []byte{},&UTXOSet,nodeID)
			for _, p := range p2pprotocol.Manager.Peers.Peers {
				p2pprotocol.SendTx(p, p.Rw, tx)
			}
			p2pprotocol.Manager.Mu.Lock()
			p2pprotocol.Manager.TxMempool[hex.EncodeToString(tx.ID)] = tx
			p2pprotocol.Manager.Mu.Unlock()
			bc.Db.Close()
			//cli.send(fromaddress,toaddress,amountnum,nodeID,false)
		}
	}

	//fmt.Println("Success!")
}
