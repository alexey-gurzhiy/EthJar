package generate

import (
	log "EthJar/app/log"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)


func NewWallet() string {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "wallet/balance",
		"\n	2. Function": "NewWallet",
	})
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("Failed to generate new keys:", err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	log.Info("Private key generated: ", privateKeyBytes)
	log.Info("Private key generated: ", strings.TrimLeft(hexutil.Encode(privateKeyBytes), "0x"))
	return strings.TrimLeft(hexutil.Encode(privateKeyBytes), "0x")
}
