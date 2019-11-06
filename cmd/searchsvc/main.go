package searchsvc

import (
	"github.com/krazybee/internals/webservice"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Initialize(configPath *string) {
		if *configPath == "" {
			log.Panic("[ServiceContainer] server expects a --config argument with path to the configuration file")
		}
		var err error
		AppContainer, err = ConstructContainer(*configPath)
		if err != nil {
			log.Panicf("[ServiceContainer] failed to construct app container: %s", err)
		}
		serverParam := &webservice.Param{
			ConfProvider:  AppContainer.ConfigProvider,
			DBProvider: AppContainer.DBProvider,
		}

		go func(serverParam *webservice.Param) {
			if err := webservice.Run(serverParam); err != nil {
				log.Panicf("%s", err)
			}
		}(serverParam)

		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		// gracefully shut down the services
		select {
		case <-done:
			// gracefully shutdown
		}
	}
