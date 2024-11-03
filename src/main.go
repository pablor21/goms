package main

import (
	"flag"
	"strings"

	"github.com/pablor21/goms/app/cmd"
	"github.com/pablor21/goms/app/config"
	_ "github.com/pablor21/goms/app/serializer"
	"github.com/pablor21/goms/pkg/database"
	"github.com/pablor21/goms/pkg/logger"
	"github.com/pablor21/goms/pkg/storage"
)

func main() {

	// Initialize config
	flag.String("config", "config.yml", "path to config file")
	flag.Parse()
	var cfgFile = flag.Lookup("config").Value.String()
	var files = strings.Split(cfgFile, ",")
	config.InitConfig(files)

	// Initialize logger
	logger.InitLogger(config.GetConfig().Logger)

	// Initialize database
	database.NewDbManager(config.GetConfig().Database)
	defer database.Close()
	// Initialize storage
	storage.InitStorage(config.GetConfig().Storage)

	// Run command
	err := cmd.Run()
	if err != nil {
		logger.Fatal().Err(err).Msg("error running command")
	}
}
