package create

import (
	db "EthJar/app/db"
	log "EthJar/app/log"
	generateAddress "EthJar/wallet/address"
	generate "EthJar/wallet/generate"
	"crypto/rand"
	"fmt"
	"math/big"
)

func Create(jar_name, jar_description string) (jar_id int, jar_password string) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "create",
		"\n	2. Function": "Create",
	})
	log.Trace("create started")

	privateKey := generate.NewWallet()
	randomNum, err := rand.Int(rand.Reader, big.NewInt(999999999))
	if err != nil {
		log.Warn("Failed to generate random number:", err)
		return 0, ""
	}
	jar_password = fmt.Sprintf("%x", randomNum)
	log.Info("Password is:", jar_password)
	jar_status := "Created"
	storeJar := db.Jar{
		Jar_password:    jar_password,
		Jar_name:        jar_name,
		Jar_description: jar_description,
		Jar_status:      jar_status,
	}
	jar_id = storeJar.Write()

	wallet, _, err := generateAddress.ConvertedPrivateKeyAndWallet(privateKey)
	if err != nil {
		log.Error(err)
	}
	log.Info("Wallet: ", wallet)
	storeWallet := db.Wallet{
		Wallet:     fmt.Sprintf("%s", wallet),
		Jar_id:     jar_id,
		PrivateKey: privateKey,
	}
	storeWallet.Write()

	return jar_id, jar_password
}
