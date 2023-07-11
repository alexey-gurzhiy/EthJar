package main

import (
	log "EthJar/app/log"
	create "EthJar/create"
	block "EthJar/transaction/block"
	send "EthJar/transaction/send"
	generateAddress "EthJar/wallet/address"
	balance "EthJar/wallet/balance"
	"math/big"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	transactionID  string = "0x3cf25c6a0a4333d33276b1eee17ef7495fc8a88cdf86123c890af03d0f867ca3"
	nodeURL        string = "https://mainnet.infura.io/v3/bfdf85f781ae493aaf2773a6b495c36a"
	nodeURLGoerli  string = "https://goerli.infura.io/v3/bfdf85f781ae493aaf2773a6b495c36a"
	nodeURLSepolia string = "https://sepolia.infura.io/v3/bfdf85f781ae493aaf2773a6b495c36a"
)

func main() {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "main",
		"\n	2. Function": "main",
	})
	log.Trace("App started")
	router := gin.Default()
	router.GET("/create", getAlbums)
	router.POST("/create", createJar)
	router.GET("/withdraw", getAlbums)
	router.POST("/withdraw", withdrawJar)
	router.GET("/send", sendGetGas)
	router.POST("/send", sendToJar)
	router.GET("/transactions", getAlbums)
	router.POST("/transactions", getTransactions)
	log.Trace("Api started")

	// Azure App Service sets the port as an Environment Variable
	// This can be random, so needs to be loaded at startup
	port := os.Getenv("HTTP_PLATFORM_PORT")

	// default back to 8080 for local dev
	if port == "" {
		port = "8080"
	}

	err := router.Run("127.0.0.1:" + port)
	if err != nil {
		log.Panic(err)
	}

}

func createJar(c *gin.Context) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "web/api",
		"\n	2. Function": "createJar",
	})
	type addJar struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	var newJar addJar

	if err := c.BindJSON(&newJar); err != nil {
		return
	}
	log.Info("newJar was binded: ", newJar)

	jar_name := newJar.Name
	jar_description := newJar.Description
	log.Info("Name: ", jar_name, "\nDescription: ", jar_description)
	jar_id, jar_password := create.Create(jar_name, jar_description)

	type jarIDAndPassword struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}

	newjarIDAndPassword := jarIDAndPassword{
		ID:       jar_id,
		Password: jar_password,
	}
	c.IndentedJSON(http.StatusCreated, newjarIDAndPassword)
}

func withdrawJar(c *gin.Context) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "web/api",
		"\n	2. Function": "withdrawJar",
	})
	type withdraJar struct {
		JarID    int    `json:"jar_id"`
		Password string `json:"password"`
		Wallet   string `json:"wallet"`
	}

	var newWithdraw withdraJar

	if err := c.BindJSON(&newWithdraw); err != nil {
		return
	}
	log.Info("newJar was binded: ", newWithdraw)

	nodeURL := nodeURLSepolia
	withdrawData := send.WithdrawData{
		Jar_id:   newWithdraw.JarID,
		Wallet:   newWithdraw.Wallet,
		Password: newWithdraw.Password,
		NodeURL:  nodeURL,
	}

	log.Info("\nID: ", withdrawData.Jar_id, "\nPassword: ", withdrawData.Password, "\nWallet: ", withdrawData.Wallet)

	txHash, err := send.Transfer(&withdrawData)
	if err != nil {
		c.IndentedJSON(http.StatusCreated, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, txHash)
}

func sendToJar(c *gin.Context) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "web/api",
		"\n	2. Function": "sendToJar",
	})
	type sendToJar struct {
		JarId      int      `json:"jar_id"`
		SendOnOwn  bool     `json:"sendOnOwn"`
		PrivateKey string   `json:"privateKey"`
		Amount     *big.Int `json:"amount"`
		Comment    string   `json:"comment"`
		Gas        uint64   `json:"gas"`
		GasPrice   *big.Int `json:"gasPrice"`
		GasMax     *big.Int `json:"gasMax"`
	}

	var newSendToJar sendToJar

	if err := c.BindJSON(&newSendToJar); err != nil {
		return
	}
	log.Info("newJar was binded: ", newSendToJar)

	jar_id := newSendToJar.JarId
	if newSendToJar.SendOnOwn {
		log.Trace("Selected to send on their on")
		wallet := generateAddress.GetWalletFromID(jar_id)
		c.IndentedJSON(http.StatusCreated, wallet)
		return
	}
	privateKey := newSendToJar.PrivateKey
	amount := newSendToJar.Amount
	comment := newSendToJar.Comment
	gas := newSendToJar.Gas
	gasPrice := newSendToJar.GasPrice
	gasMax := newSendToJar.GasMax

	log.Infof(`
jar_id: %d
privateKey: %s
amount: %d
comment: %s
gas: %d
gasPrice: %d
gasMax: %d`, jar_id, privateKey, amount, comment, gas, gasPrice, gasMax)

	nodeURL := nodeURLSepolia
	sendData := send.SendData{
		Jar_id:     jar_id,
		PrivateKey: privateKey,
		NodeURL:    nodeURL,
		Comment:    comment,
		TxData: send.TxData{
			Value:     amount,
			GasTipCap: gasMax,
			GasFeeCap: gasPrice,
			Gas:       gas,
		},
	}
	txHash, err := send.Transfer(&sendData)

	// txHash, err := send.SendApi(nodeUrl, jar_id, privateKey, amount, comment, gas, gasPrice, gasMax)

	if err != nil {
		c.IndentedJSON(http.StatusCreated, err.Error())
		return
	}
	c.IndentedJSON(http.StatusCreated, txHash)
}

func sendGetGas(c *gin.Context) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "web/api",
		"\n	2. Function": "sendToJar",
	})
	nodeURL := nodeURLSepolia
	// gasTip = gasPrice			gasMax = intMaxCapGasFee			gasLimit = gas
	gasLimit, gasPrice, gasMax := send.GasPriceFnAPI(nodeURL)
	log.Trace("gas price sent")
	c.JSON(http.StatusOK, gin.H{
		"gasLimit": gasLimit, "gasPrice": gasPrice, "gasMax": gasMax})
}

func getTransactions(c *gin.Context) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "web/api",
		"\n	2. Function": "getTransactions",
	})

	var jar_id int

	if err := c.BindJSON(&jar_id); err != nil {
		return
	}
	log.Info("jar_id was binded: ", jar_id)

	wallet := generateAddress.GetWalletFromID(jar_id)
	// walletToTest:
	// wallet = "0x19588531AD56920058c2d6D175538636cdCc2F98"
	nodeUrl := nodeURLSepolia
	balance := balance.GetWalletBalance(wallet, nodeUrl)

	txList, err := block.EtherscanListTransactionsAPI(wallet)

	if err != nil {
		c.IndentedJSON(http.StatusCreated, err.Error())
		return
	}
	log.Info(balance)

	type reply struct {
		Balance *big.Int
		TxList  []block.TxStruct
	}
	postReply := reply{Balance: balance, TxList: txList}
	log.Debug(postReply)
	c.IndentedJSON(http.StatusCreated, postReply)
}

func getAlbums(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
