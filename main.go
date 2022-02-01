package main

import (
	"flag"

	"github.com/platform9/nodelet2go/pkg/cluster"
)

func initConfig() (string, error) {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "/root/nodeletConfig.yml", "path to bootstrap file")
	flag.Parse()
	return cfgFile, nil
}

func main() {

	cfgFile, _ := initConfig()
	cluster.CreateCluster(cfgFile)
	return
}
