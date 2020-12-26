package main

import (
	"fmt"
	"github.com/josephbmanley/OpenSkins-Common/datastore"
	"github.com/josephbmanley/OpenSkins/pluginmanager"
	log "github.com/sirupsen/logrus"
	"os"
	pas "plugin"
)

const plugindirectory = "./plugins"

func main() {

	pluginFiles, err := pluginmanager.GetPlugins(plugindirectory)
	if err != nil {
		log.Warningln(fmt.Sprintf("Failed to read plugins directory: %v", err.Error()))
	}

	for _, file := range pluginFiles {
		log.Infoln(fmt.Sprintf("Loading plugin: %v", file))
		plugin, err := pas.Open(file)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Failed to load plugin '%v': %v", file, err.Error()))
			os.Exit(1)
		}

		symSkinstore, err := plugin.Lookup("SkinstoreModule")
		if err != nil {
			log.Fatalln(fmt.Sprintf("Failed to load plugin '%v': %v", file, err.Error()))
			os.Exit(1)
		}

		var skinstore datastore.Skinstore
		skinstore, ok := symSkinstore.(datastore.Skinstore)
		if !ok {
			log.Fatalln(fmt.Sprintf("Invalid type for Skinstore in plugin '%v'", file))
			os.Exit(1)
		}

		err = skinstore.Initialize()
		if err != nil {
			log.Fatalln(fmt.Sprintf("Failed to intialize Skinstore in plugin '%v'", file))
			os.Exit(1)
		}

	}

}
