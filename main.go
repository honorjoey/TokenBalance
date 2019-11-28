package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/toba/ethtest/contracts"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

func main() {
	contractAddress := "0xdac17f958d2ee523a2206206994597c13d831ec7"
	tokenDecimal := 6
	//Net := "http://52.221.201.15:8232"
	Net := "https://mainnet.infura.io"
	filename := "/root/go/src/TokenBalance/address.txt"


	//read address list from file
	addressList, err := ReadFile(filename)
	if err != nil {
		log.Println("read file failed,", err)
	}
	//check balance of address, print it
	for _, v := range *addressList {
		time.Sleep(500*time.Millisecond)
		bal, err := BalanceToken(Net, contractAddress, strings.ToLower(v))
		if err != nil {
			log.Println("get the balance failed,", err)
		}else{
			balance, _ := new(big.Int).SetString(bal, 10)
			decimal := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(tokenDecimal)), big.NewInt(0))
			fmt.Println(strings.ToLower(v), balance.Div(balance, decimal))
		}
	}
}


func ReadFile(filename string) (*[]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("open the file failed!")
		return nil, err
	}
	defer file.Close()
	buf := make([]byte, 4096*1000)
	_, err = file.Read(buf)
	if err != nil {
		log.Println("read the file failed!")
		return nil, err
	}
	var bytes []byte
	for i, v := range buf {
		if v == 0 {
			bytes = buf[:i]
			break
		}
	}
	slice := strings.Split(string(bytes), "\n")

	return &slice, nil
}

func BalanceToken(Net, tokenAddress, address string) (string, error) {
	client, err := ethclient.Dial(Net)
	if err != nil {
		return "", err
	}

	// Golem (GNT) Address
	tokenAddr := common.HexToAddress(tokenAddress)
	instance, err := contracts.NewTobaToken(tokenAddr, client)
	if err != nil {
		return "", err
	}

	addr := common.HexToAddress(address)
	bal, err := instance.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		return "", err
	}
	//bal.Div(bal, big.NewInt(1000000000000000000))
	return bal.String(), nil
}
