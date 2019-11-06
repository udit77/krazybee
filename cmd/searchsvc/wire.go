//+build wireinject

package searchsvc

import (
	"github.com/google/wire"
	"github.com/krazybee/internals/config"
	"github.com/krazybee/internals/config/jsonconfig"
	"github.com/krazybee/internals/dbaccess"
	"github.com/krazybee/internals/dbaccess/mysql"
)

func initContainer(path string) (*Container, error) {
	wire.Build(containerSet)
	return nil, nil
}

func ResolveConfigProvider(path string) (config.Provider, error) {
	provider, err := jsonconfig.NewJsonConfigProvider(string(path))
	if err != nil {
		return nil, err
	}
	err = provider.Initialize()

	return provider, nil
}

func ResolveDBProvider(config config.Provider) (dbaccess.Provider, error) {
	provider, err := mysql.NewDBProvider(config)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
