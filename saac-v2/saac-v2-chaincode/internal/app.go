package internal

import (
	"context"
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal/handler"
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal/repository"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type AppService struct {
	log        *logrus.Entry
	repository *repository.Repository
}

func NewAppService(
	lc fx.Lifecycle,
	log *logrus.Entry,
	r *repository.Repository,
) *AppService {
	srv := &AppService{
		log:        log,
		repository: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if err := srv.Start(ctx); err != nil {
				return err
			}
			return nil
		},
	})

	return srv
}

func (svc *AppService) Start(ctx context.Context) error {
	go func() error {
		assetChaincode, err := contractapi.NewChaincode(handler.NewSmartContract(svc.log, svc.repository))
		if err != nil {
			svc.log.Errorf("error creating saac-v2 chaincode: %v", err)
		}

		if err := assetChaincode.Start(); err != nil {
			svc.log.Errorf("error starting saac-v2 chaincode: %v", err)
		}
		return nil
	}()

	return nil
}
