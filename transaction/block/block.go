package block

import (
	"context"
	"io"
	"strings"

	log "EthJar/app/log"
	"EthJar/node/connect"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"encoding/hex"
	"encoding/json"
	"net/http"
)

type MyError struct {
	message string
}

func (e *MyError) Error() string {
	return e.message
}

type TxStruct struct {
	From   string `json:"from"`
	TxHash string `json:"txHash"`
	Value  string `json:"value"`
	Data   string `json:"data"`
	Err    string `json:"error"`
}

func blockFunction(blockNumber *big.Int, nodeURL string) (conn *ethclient.Client, cntxt context.Context, header *types.Header, block *types.Block) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "transaction/block",
		"\n	2. Function": "blockFunction",
	})
	conn, cntxt = connect.LiveConnectionToEthereumNode(nodeURL)
	log.Trace("client and context were created")

	header, err := conn.HeaderByNumber(cntxt, nil) 
	if err != nil {
		log.Fatal("Unable to get header")
	} else {
		log.Trace("Block header was created")
	}

	block, err = conn.BlockByNumber(cntxt, blockNumber)
	if err != nil {
		log.Fatal("Unable to get block by number,", blockNumber, err)
	} else {
		log.Trace("Block section was created")
	}

	return conn, cntxt, header, block
}

func CalculateGas(nodeURL string) (gasAverage uint64, gasAverageTotal, gasPriceAverage *big.Int) {
	log := log.WithFields(log.Fields{
		"\n	1.	Module":   "transaction/block",
		"\n 2.	Function": "BlockTransactions",
	})

	_, _, _, block := blockFunction(nil, nodeURL)
	log.Trace("Connection, context and block function imported")
	log.Info("Our selected block is:", block.Number().Uint64())
	txAmount := uint64(len(block.Transactions()))
	log.Info("Amount of transactions in selected block:", txAmount)
	var totalGas uint64 = 0
	totalGasPrice := big.NewInt(0)

	for _, tx := range block.Transactions() {
		log.Debugf(`
		Value: %s
		Gas: %d
		Cost: %s
		GasPrice: %s
		`, tx.Value().String(), tx.Gas(), tx.Cost(), tx.GasPrice())
		gas, gasPrice := tx.Gas(), tx.GasPrice()
		totalGas += + gas
		totalGasPrice.Add(totalGasPrice, gasPrice)
	}
	gasAverage = totalGas / txAmount
	log.Info("Gas average: ", gasAverage)
	gasPriceAverage = totalGasPrice.Div(totalGasPrice, big.NewInt(int64(txAmount)))
	log.Info("Gas Price average: ", gasPriceAverage)
	gasAverageTotal = gasPriceAverage.Mul(gasPriceAverage, big.NewInt(int64(gasAverage)))
	log.Info("Total recommended price for Gas: ", gasAverageTotal)
	return gasAverage, gasAverageTotal, gasPriceAverage
}

// -------------- Etherscan API ---------------//
type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	GasPrice string `json:"gasPrice"`
	GasLimit string `json:"gasLimit"`
	Nonce    string `json:"nonce"`
	Data     string `json:"input"`
}

func EtherscanListTransactionsAPI(wallet string) ([]TxStruct, error) {
	log := log.WithFields(log.Fields{
		"\n	1. Module":   "transaction/block",
		"\n	2. Function": "EtherscanListTransactionsApi",
	})

	apiKey := "2ZZTPE87N8D9U8GEKNWPF3Y7V3M58G9621"

	url := fmt.Sprintf("https://api.etherscan.io/api?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s", wallet, apiKey)
	log.Trace("Got result from url")

	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()
	log.Trace("Got response from url")

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Trace("Read the data and logged to body")

	var response struct {
		Status  string        `json:"status"`
		Message string        `json:"message"`
		Result  []Transaction `json:"result"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Trace("parsed JSON")

	if response.Status != "1" {
		log.Warnf("API error: %s", response.Message)
		err = &MyError{
			message: fmt.Sprintf("API error: %s", response.Message),
		}
		return nil, err
	}

	var txList = []TxStruct{}

	for _, tx := range response.Result {
		if tx.To == strings.ToLower(wallet) {
			byteData, err := hex.DecodeString(tx.Data[2:])
			if err != nil {
				log.Error("Error decoding hex string:", err)
				err = &MyError{
					message: fmt.Sprint("Error decoding hex string: ", err),
				}
				return nil, err
			}

			data := string(byteData)
			txAdd := []TxStruct{
				{From: tx.From, TxHash: tx.Hash, Value: tx.Value, Data: data, Err: fmt.Sprint(err)},
			}

			txList = append(txList, txAdd...)
		}
	}
	return txList, nil
}
