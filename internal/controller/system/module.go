package system

import (
	"github.com/formancehq/go-libs/logging"
	ledgercontroller "github.com/formancehq/ledger/internal/controller/ledger"
	"go.uber.org/fx"
)

type ModuleConfiguration struct {
	NSCacheConfiguration ledgercontroller.CacheConfiguration
}

func NewFXModule(configuration ModuleConfiguration) fx.Option {
	return fx.Options(
		fx.Provide(func(controller *DefaultController) Controller {
			return controller
		}),
		fx.Provide(func(
			store Store,
			listener ledgercontroller.Listener,
			logger logging.Logger,
		) *DefaultController {
			options := make([]Option, 0)
			if configuration.NSCacheConfiguration.MaxCount != 0 {
				options = append(options, WithCompiler(ledgercontroller.NewCachedCompiler(
					ledgercontroller.NewDefaultCompiler(),
					configuration.NSCacheConfiguration,
				)))
			}
			return NewDefaultController(store, listener, options...)
		}),
	)
}