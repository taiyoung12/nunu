package main

import (
	"context"
	"flag"
	"fmt"

	"nunu/cmd/server/wire"
	"nunu/pkg/config"
	"nunu/pkg/log"

	"go.uber.org/zap"
)

func main() {
	var envConf = flag.String("conf", "config/local.yml", "config path, eg: -conf ./config/local.yml")
	flag.Parse()
	conf := config.NewConfig(*envConf)

	logger := log.NewLog(conf)

	app, cleanup, err := wire.NewWire(conf, logger)
	defer cleanup()
	if err != nil {
		panic(err)
	}
	logger.Info("server start", zap.String("host", fmt.Sprintf("http://%s:%d", conf.GetString("http.host"), conf.GetInt("http.port"))))
	logger.Info("nunu agent ready - listening for Slack messages")
	if err = app.Run(context.Background()); err != nil {
		panic(err)
	}
}
