package balance

import (
	"EthJar/node/connect"


	"math/big"

	"github.com/ethereum/go-ethereum/common"
	log "EthJar/app/log"
)




func GetWalletBalance(wallet, nodeURL string) *big.Int {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "wallet/balance",
		"\n	2. Function": "GetWalletBalance",
	})
	conn, cntxt := connect.LiveConnectionToEthereumNode(nodeURL)
	log.Trace("client and context were created")

	account := common.HexToAddress(wallet)
	balance, err := conn.BalanceAt(cntxt, account, nil)

	if err != nil {
		log.Fatal("Unable to get balance:", err)
	}
	log.Trace("Balance: ", balance)

	return balance
}
