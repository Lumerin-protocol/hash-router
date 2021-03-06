package contractmanager

import (
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var accountAddress = common.HexToAddress("0x860fB39B5B24c9F974C8f484223eAa573b1D16e0")
var accountPrivateKey = "5a3d62629f54c67fd59bb0dc234f95563f816a7c6d43699656f49db8c92c488d"
var gethNodeAddress = "ws://127.0.0.1:7545"

var clonefactoryAddress common.Address // = common.HexToAddress("0xEA3C21BF6aE276B8d084E79D6Ef45d8BfE1ce7B3")

var hashrateContractAddress common.Address //= common.HexToAddress("0x3ED63115D92a95538EB111D32f07Ef80C455e12b")
var poolUrl = "stratum+tcp://stratum.slushpool.com:3333"

func TestHashrateContractCreation(t *testing.T) {
	// hashrate contract params
	price := 0
	limit := 10
	speed := 100
	length := 100

	client, err := SetUpClient(gethNodeAddress, accountAddress)
	if err != nil {
		log.Fatalf("Error::%v", err)
	}

	CreateHashrateContract(client, accountAddress, accountPrivateKey, clonefactoryAddress, price, limit, speed, length, clonefactoryAddress)

	// subcribe to creation events emitted by clonefactory contract
	cfLogs, cfSub, _ := SubscribeToContractEvents(client, clonefactoryAddress)
	// create event signature to parse out creation event
	contractCreatedSig := []byte("contractCreated(address,string)")
	contractCreatedSigHash := crypto.Keccak256Hash(contractCreatedSig)
	for {
		select {
		case err := <-cfSub.Err():
			log.Fatalf("Error::%v", err)
		case cfLog := <-cfLogs:

			if cfLog.Topics[0].Hex() == contractCreatedSigHash.Hex() {
				hashrateContractAddress := common.HexToAddress(cfLog.Topics[1].Hex())
				fmt.Printf("Address of created Hashrate Contract: %v\n\n", hashrateContractAddress.Hex())
			}
		}
	}
}

func TestHashrateContractPurchase(t *testing.T) {

	client, err := SetUpClient(gethNodeAddress, accountAddress)
	if err != nil {
		log.Fatalf("Error::%v", err)
	}

	PurchaseHashrateContract(client, accountAddress, accountPrivateKey, clonefactoryAddress, hashrateContractAddress, accountAddress, poolUrl)

	// subcribe to purchase events emitted by clonefactory contract
	cfLogs, cfSub, _ := SubscribeToContractEvents(client, clonefactoryAddress)
	// create event signature to parse out purchase event
	clonefactoryContractPurchasedSig := []byte("clonefactoryContractPurchased(address)")
	clonefactoryContractPurchasedSigHash := crypto.Keccak256Hash(clonefactoryContractPurchasedSig)
	for {
		select {
		case err := <-cfSub.Err():
			log.Fatalf("Error::%v", err)
		case cfLog := <-cfLogs:

			if cfLog.Topics[0].Hex() == clonefactoryContractPurchasedSigHash.Hex() {
				hashrateContractAddress := common.HexToAddress(cfLog.Topics[1].Hex())
				fmt.Printf("Address of purchased Hashrate Contract: %v\n\n", hashrateContractAddress.Hex())
			}
		}
	}
}

func TestDeployContracts(t *testing.T) {

	client, err := SetUpClient(gethNodeAddress, accountAddress)
	if err != nil {
		log.Fatalf("Error::%v", err)
	}

	cloneFactoryAddress := DeployContracts(client, accountAddress, accountPrivateKey)

	fmt.Printf("Address of CloneFactory contract: %v\n\n", cloneFactoryAddress.Hex())
}
