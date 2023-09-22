package internal

import (
	"fmt"
	"whatsApp/core"

	"gorm.io/gorm/logger"
)

type writer struct {
	logger.Writer
}

// NewWriter writer 构造函数
// Author [SliverHorn](https://github.com/SliverHorn)
func NewWriter(w logger.Writer) *writer {
	return &writer{Writer: w}
}

// Printf 格式化打印日志
// Author [SliverHorn](https://github.com/SliverHorn)
func (w *writer) Printf(message string, data ...interface{}) {
	var logZap bool
	switch core.New().Config.System.DbType {
	case "mysql":
		logZap = core.New().Config.Mysql.LogZap
	case "pgsql":
		logZap = core.New().Config.Pgsql.LogZap
	}
	if logZap {
		core.New().ZapLog.Info(fmt.Sprintf(message+"\n", data...))
	} else {
		w.Writer.Printf(message, data...)
	}
}
