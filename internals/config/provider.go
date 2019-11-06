package config

import "github.com/krazybee/internals/config/model"

type Provider interface {
	Initialize() error
	GetConfig() *model.Config
}
