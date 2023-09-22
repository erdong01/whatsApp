package core

import (
	"sync"
	"whatsApp/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Info struct {
	Port    int
	Env     string
	Name    string
	version string
}

type Core struct {
	Db     *gorm.DB
	Config config.Server
	Info   *Info
	Viper  *viper.Viper
	ZapLog *zap.Logger
}

var (
	c    *Core
	once sync.Once
)

func New() *Core {
	once.Do(func() {
		c = &Core{}
	})
	return c

}

func Copy() *Core {
	var core = new(Core)
	*core = *New()
	return core
}

func Set(newCore *Core) {
	c = newCore
}
