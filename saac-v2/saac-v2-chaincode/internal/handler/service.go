package handler

import (
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal/repository"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/sirupsen/logrus"
)

type SmartContract struct {
	contractapi.Contract
	log        *logrus.Entry
	repository *repository.Repository
}

func NewSmartContract(log *logrus.Entry, r *repository.Repository) *SmartContract {
	return &SmartContract{
		log:        log,
		repository: r,
	}
}
