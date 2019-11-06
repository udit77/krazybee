package searchsvc

import (
	"github.com/google/wire"
	"github.com/krazybee/internals/config"
	"github.com/krazybee/internals/dbaccess"
	"log"
)

type Container struct {
	ConfigProvider config.Provider
	DBProvider     dbaccess.Provider
}

var AppContainer *Container

var containerSet = wire.NewSet(
	ResolveConfigProvider,
	ResolveDBProvider,
	wire.Struct(new(Container), "*"),
)

func ConstructContainer(configPath string) (*Container, error) {
	log.Println("[SearchSvc] constructing container with config path",configPath)
	return initContainer(configPath)
}
