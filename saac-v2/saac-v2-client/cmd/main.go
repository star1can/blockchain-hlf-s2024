package main

import (
	"github.com/hlf-mipt/saac-v2-client/internal"
	"github.com/hlf-mipt/saac-v2-client/internal/config"
	"github.com/hlf-mipt/saac-v2-client/internal/httphandler"
	"github.com/hlf-mipt/saac-v2-client/internal/logger"
	"github.com/hlf-mipt/saac-v2-client/internal/service/hlf"
	"go.uber.org/fx"
)

func main() {

	fx.New(
		fx.Provide(logger.NewLogger),
		fx.Provide(httphandler.NewHttpHandler),
		fx.Provide(config.NewWebApiConfig),
		fx.Provide(hlf.NewHLFService),
		fx.Provide(internal.NewAppService),

		fx.Invoke(func(*hlf.HLFService) {}),
		fx.Invoke(func(*internal.AppService) {}),
	).Run()

}
