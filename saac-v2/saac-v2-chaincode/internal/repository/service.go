package repository

import (
	"github.com/hlf-mipt/saac-v2-core/pkg/model"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AssetsRepository interface {
	SetPrivate(ctx contractapi.TransactionContextInterface, mspId string, item *model.Asset) error
	GetPrivate(ctx contractapi.TransactionContextInterface, mspId, id string) (*model.Asset, error)
}

type Repository struct {
	assets *assetsRepository
}

func NewRepository() *Repository {
	return &Repository{
		assets: newAssetsRepository(),
	}
}

func (svc *Repository) Assets() AssetsRepository {
	return svc.assets
}
