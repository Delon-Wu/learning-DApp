package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"learning-dapp/contracts"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func deploy(client *ethclient.Client) {
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	fmt.Println("gasPrice", gasPrice)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	authSigner, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	authSigner.From = fromAddress
	authSigner.Nonce = big.NewInt(int64(nonce))
	authSigner.Value = big.NewInt(0)
	authSigner.GasLimit = uint64(3e5)
	authSigner.GasPrice = gasPrice
	address, tx, instance, err := contracts.DeployCouter(authSigner, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Deployed: ", address)
	fmt.Println("TxHash: ", tx.Hash().Hex())
	fmt.Println("Nonce: ", tx.Nonce())
	fmt.Println("GasPrice: ", gasPrice)

	_ = instance
}

// waitForTransaction 等待交易被确认
func waitForTransaction(client *ethclient.Client, txHash common.Hash) error {
	fmt.Printf("等待交易确认: %s\n", txHash.Hex())

	ctx := context.Background()
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			if receipt.Status == 1 {
				fmt.Printf("交易已确认，区块号: %d\n", receipt.BlockNumber.Uint64())
				return nil
			} else {
				return fmt.Errorf("交易失败")
			}
		}

		// 等待一段时间后重试
		time.Sleep(2 * time.Second)
	}
}

func callContract(client *ethclient.Client) {
	const contractAddress = "0xd787D2928c09444267CbCe3b84Eeb23798a2d5BB"
	counterContract, err := contracts.NewCouter(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111))
	if err != nil {
		log.Fatal(err)
	}

	// 第一次增加
	tx, err := counterContract.Increase(opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TxHash1: ", tx.Hash().Hex())

	// 等待第一次交易确认
	if err := waitForTransaction(client, tx.Hash()); err != nil {
		log.Fatal(err)
	}

	// 第二次增加
	tx, err = counterContract.Increase(opt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TxHash2: ", tx.Hash().Hex())

	// 等待第二次交易确认
	if err := waitForTransaction(client, tx.Hash()); err != nil {
		log.Fatal(err)
	}

	// 现在读取值
	callOption := &bind.CallOpts{Context: context.Background()}
	value, err := counterContract.GetNumber(callOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Counter value: ", value)
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

	//deploy(client)

	callContract(client)
}
