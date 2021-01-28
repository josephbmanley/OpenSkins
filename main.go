package main

import (
	"fmt"
	"github.com/josephbmanley/OpenSkins/pluginmanager"
	"github.com/josephbmanley/OpenSkins/runtime"
	log "github.com/sirupsen/logrus"
	"os"
	"plugin"
)

var appRuntime = "webserver"

const plugindirectory = "./plugins"

func main() {

	pluginFiles, err := pluginmanager.GetPlugins(plugindirectory)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to read plugins directory: %v", err.Error()))
		os.Exit(1)
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

	switch appRuntime {
	case "webserver":
		runtime.StartWebserver(8081)
	default:
		log.Fatalln("Runtime is currently not implemented!")
		os.Exit(1)
	}

}
