package send

import (
	db "EthJar/app/db"
	log "EthJar/app/log"
	"EthJar/node/connect"
	block "EthJar/transaction/block"
	generateAddress "EthJar/wallet/address"
	ballance "EthJar/wallet/balance"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
)

type MyError struct {
	message string
}

func (e *MyError) Error() string {
	return e.message
}

type TxData struct {
	Nonce               uint64
	GasTipCap           *big.Int
	GasFeeCap           *big.Int
	Gas                 uint64
	To                  common.Address
	Value               *big.Int
	Data                []byte
	ConvertedPrivateKey *ecdsa.PrivateKey
	Conn                *ethclient.Client
	Cntxt               context.Context
}

type WithdrawData struct {
	TxData
	Jar_id   int
	Wallet   string
	Password string
	NodeURL  string
}

type SendData struct {
	PrivateKey string
	Comment    string
	TxData
	Jar_id   int
	Wallet   string
	Password string
	NodeURL  string
}

type TxSender interface {
	Send() (string, error)
	Validate() error
}

func Transfer(t TxSender) (string, error) {
	err := t.Validate()
	if err != nil {
		log.Info(err)
		return "", err
	}
	txHash, err := t.Send()
	if err != nil {
		log.Info(err)
		return "", err
	}
	return txHash, nil
}

func (s *SendData) Validate() error {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "transaction/withdraw",
		"\n	2. Function": "Validate (inp *SendData)",
	})

	conn, cntxt := connect.LiveConnectionToEthereumNode(s.NodeURL)

	fromAddress, convertedPrivateKey, err := generateAddress.ConvertedPrivateKeyAndWallet(s.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Trace("From address and convertedPrivateKey generated")

	nonce, err := conn.PendingNonceAt(cntxt, fromAddress)
	if err != nil {
		log.Fatal("Unable to get nonce")
	}
	log.Info("Nonce: ", nonce)

	wallet := generateAddress.GetWalletFromID(s.Jar_id)
	toAddress := common.HexToAddress(wallet)
	data := []byte(s.Comment)
	log.Trace(data)

	s.Nonce = nonce
	// s.GasTipCap =
	// s.GasFeeCap =
	// s.Gas =
	s.To = toAddress
	// s.Value =
	s.Data = data
	s.ConvertedPrivateKey = convertedPrivateKey
	s.Conn = conn
	s.Cntxt = cntxt
	return nil
}

func (inp *WithdrawData) Validate() error {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "transaction/withdraw",
		"\n	2. Function": "Validate (inp *InputData)",
	})

	jarDb := db.Jar{
		Jar_id: inp.Jar_id,
	}

	jar_password := jarDb.Read().Jar_password
	if jar_password != inp.Password {
		log.Warn("Incorrect password")
		err := &MyError{
			message: "Invalid Password",
		}
		return err
	}
	log.Trace("Password is correct")

	Wallet := db.Wallet{
		Jar_id: inp.Jar_id,
	}
	privateKey := Wallet.Read().PrivateKey
	log.Info(privateKey)

	conn, cntxt := connect.LiveConnectionToEthereumNode(inp.NodeURL)

	fromAddress, convertedPrivateKey, err := generateAddress.ConvertedPrivateKeyAndWallet(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	log.Trace("From address and convertedPrivateKey generated")
	nonce, err := conn.PendingNonceAt(cntxt, fromAddress)
	if err != nil {
		log.Fatal("Unable to get nonce")
	}
	log.Info("Nonce: ", nonce)
	toAddress := common.HexToAddress(inp.Wallet)

	gas, gasAverageTotal, gasTip := block.CalculateGas(inp.NodeURL)
	log.Info("gasLimit: ", gas, "\ngasAverageTotal:", gasAverageTotal, "\ngasTip: ", gasTip)
	gasMax := gasAverageTotal
	gasMax.Mul(gasMax, big.NewInt(2))
	log.Info("GasMax :", gasMax)
	// calculate amount to be sent

	var value *big.Int
	walletBalance := ballance.GetWalletBalance(fromAddress.String(), inp.NodeURL)
	log.Info("walletBalance :", walletBalance)
	if walletBalance.Cmp(gasMax) <= 0 {
		err := &MyError{
			message: "Not enough Eth to cover transaction fee",
		}
		return err
	}
	value.Sub(walletBalance, gasMax)
	log.Info("value :", value)
	inp.Nonce = nonce
	inp.GasTipCap = gasMax
	inp.GasFeeCap = gasTip
	inp.Gas = gas
	inp.To = toAddress
	inp.Value = value
	// inp.Data     =
	inp.ConvertedPrivateKey = convertedPrivateKey
	inp.Conn = conn
	inp.Cntxt = cntxt
	return nil
}

func (t *TxData) Send() (string, error) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "transaction/withdraw",
		"\n	2. Function": "Send (t *TxData)",
	})
	transaction := types.NewTx(&types.DynamicFeeTx{
		Nonce:     t.Nonce,
		GasTipCap: t.GasTipCap,
		GasFeeCap: t.GasFeeCap,
		Gas:       t.Gas,
		To:        &t.To,
		Value:     t.Value,
		Data:      t.Data,
	})
	config, block := params.SepoliaChainConfig, params.SepoliaChainConfig.LondonBlock
	// params.MainnetChainConfig, params.MainnetChainConfig.LondonBlock
	// params.GoerliChainConfig, params.GoerliChainConfig.LondonBlock
	signer := types.MakeSigner(config, block)
	signedTransaction, err := types.SignTx(transaction, signer, t.ConvertedPrivateKey)
	if err != nil {
		log.Fatal("Unable to sign a transaction")
	}
	hash := signedTransaction.Hash().Bytes()
	raw, err := rlp.EncodeToBytes(signedTransaction)
	if err != nil {
		log.Fatal("Unable to cast raw to raw transaction")
	}
	log.Info("Hash: ", hash, "\nRaw: ", raw, "\n")
	// Увімкнути щоб транзакція реально надсилалась
	err = t.Conn.SendTransaction(t.Cntxt, signedTransaction)
	if err != nil {
		log.Warn("Unable to submit transaction: ", err)
	}
	txHash := fmt.Sprintf("Transaction sent: %s", signedTransaction.Hash().Hex())
	fmt.Println(txHash)
	return txHash, err
}

func GasPriceFnAPI(nodeURL string) (gas uint64, gasPrice, gasMax *big.Int) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "transaction/send",
		"\n	2. Function": "GasPrice",
	})
	gas, GasAverageTotal, gasPrice := block.CalculateGas(nodeURL)

	log.Info("\nGas: ", gas, "\nGasPrice: ", gasPrice)
	var miltiplier int64 = 4
	gasMax = GasAverageTotal.Mul(GasAverageTotal, big.NewInt((miltiplier)))
	log.Info("\nGasMax: ", gasMax)
	return
}
