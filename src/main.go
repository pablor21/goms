package main

import (
	"flag"
	"strings"

	"github.com/pablor21/goms/app/cmd"
	"github.com/pablor21/goms/app/config"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/logger"
)

func main() {
	flag.String("config", "config.yml", "path to config file")
	flag.Parse()
	var cfgFile = flag.Lookup("config").Value.String()
	var files = strings.Split(cfgFile, ",")
	config.InitConfig(files)

	database.NewDbManager(config.GetConfig().Database)
	logger.InitLogger(config.GetConfig().Logger)

	err := cmd.Run()
	if err != nil {
		logger.Fatal().Err(err).Msg("error running command")
	}
}
