package main

import (
	"fmt"
	"github.com/josephbmanley/OpenSkins/pluginmanager"
	log "github.com/sirupsen/logrus"
	"os"
	"plugin"
)

const plugindirectory = "./plugins"

func main() {

	pluginFiles, err := pluginmanager.GetPlugins(plugindirectory)
	if err != nil {
		log.Warningln(fmt.Sprintf("Failed to read plugins directory: %v", err.Error()))
	}

	loadedPlugins := []*plugin.Plugin{}

	for _, file := range pluginFiles {
		log.Infoln(fmt.Sprintf("Loading plugin: %v", file))
		plugin, err := plugin.Open(file)
		if err != nil {
			log.Fatalln(fmt.Sprintf("Failed to load plugin '%v': %v", file, err.Error()))
			os.Exit(1)
		}

		loadedPlugins = append(loadedPlugins, plugin)

	}

	err = pluginmanager.LoadSkinstores(loadedPlugins)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to load skinstores '%v'", err.Error()))
		os.Exit(1)
	}

}
