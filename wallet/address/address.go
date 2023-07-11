package address

import (
	db "EthJar/app/db"
	log "EthJar/app/log"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func ConvertedPrivateKeyAndWallet(privateKey string) (wallet common.Address, convertedPrivateKey *ecdsa.PrivateKey, err error) {
	log := log.WithFields(log.Fields{
		"\n1. Module":   "wallet/address",
		"\n2. Function": "ConvertedPrivateKeyAndWallet\n",
	})

	convertedPrivateKey, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err, "PrivateKey: ", privateKey)
	}
	log.Trace("PrivateKey address converted to ECDSA")

	publicKey := convertedPrivateKey.Public()
	log.Trace("Publickey generated", publicKey)

	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Unable to cast publicKey to ECDSA")
	}
	log.Info("publicKeyECDSA: ", publicKeyECDSA)

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Info(fromAddress)
	return fromAddress, convertedPrivateKey, err
}

func GetWalletFromID(jar_id int) (wallet string) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "wallet/address",
		"\n	2. Function": "GetWalletFromId",
	})

	getWallet := db.Wallet{
		Jar_id: jar_id,
	}
	wallet = getWallet.Read().Wallet
	log.Info(wallet)
	return wallet
}
