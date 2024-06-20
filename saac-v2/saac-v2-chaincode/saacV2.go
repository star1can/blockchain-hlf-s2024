package main

import (
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal"
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal/logger"
	"github.com/hlf-mipt/saac-v2-chaincode-go/internal/repository"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(logger.NewLogger),
		fx.Provide(internal.NewAppService),
		fx.Provide(repository.NewRepository),

		fx.Invoke(func(*internal.AppService) {}),
	).Run()
}
