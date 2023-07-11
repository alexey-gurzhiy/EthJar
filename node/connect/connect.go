package connect

import (
	"context"

	log "EthJar/app/log"

	"github.com/ethereum/go-ethereum/ethclient"
)

func LiveConnectionToEthereumNode(nodeURL string) (conn *ethclient.Client, cntxt context.Context) {

	log := log.WithFields(log.Fields{
		"\n	1. Module":   "node/connect",
		"\n	2. Function": "LiveConnectionToEthereumNode",
	})
	conn, err := ethclient.Dial(nodeURL)

	if err != nil {
		log.Fatal("Failed to connect to Eth Node:", err)
	} else {
		log.Trace("connection to ", nodeURL, " was established")
	}
	cntxt = context.Background()
	return conn, cntxt

}
