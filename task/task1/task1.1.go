package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

//任务 1：区块链读写 任务目标
//使用 Sepolia 测试网络实现基础的区块链交互，包括查询区块和发送交易。
//具体任务
//环境搭建
//安装必要的开发工具，如 Go 语言环境、 go-ethereum 库。
//注册 Infura 账户，获取 Sepolia 测试网络的 API Key。
//查询区块
//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
//输出查询结果到控制台。
//发送交易
//准备一个 Sepolia 测试网络的以太坊账户，并获取其私钥。
//编写 Go 代码，使用 ethclient 连接到 Sepolia 测试网络。
//构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
//对交易进行签名，并将签名后的交易发送到网络。
//输出交易的哈希值。

func task1(client *ethclient.Client) {
	block, err := client.BlockByNumber(context.Background(), big.NewInt(9743987))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Time)
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Hash().Hex())
	fmt.Println(len(block.Transactions()))

	count, err := client.TransactionCount(context.Background(), block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count, "-----------------------")

}

func task2(client *ethclient.Client) {
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(1e9)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	fmt.Println("gasPrice", gasPrice)

	toAddress := common.HexToAddress("0x68e8441ebDac4bE9d43aC5975feC6E45E360563c")

	gasLimit := uint64(21000)
	data := []byte{}

	// 创建 LegacyTx
	legacyTx := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    value,
		Data:     data,
	}
	tx := types.NewTx(legacyTx)
	fmt.Printf("交易哈希(未签名): %s\n", tx.Hash().Hex())
	fmt.Printf("交易类型: %d\n\n", tx.Type())

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tx sent: ", signedTx.Hash().Hex())
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client, err := ethclient.Dial(os.Getenv("SEPOLIA_URL"))
	if err != nil {
		log.Fatal(err)
	}

	//task1(client)
	task2(client)
}
