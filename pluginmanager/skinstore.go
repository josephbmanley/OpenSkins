package pluginmanager

import (
	"errors"
	"fmt"
	"github.com/josephbmanley/OpenSkins-Common/datastore"
	log "github.com/sirupsen/logrus"
	"plugin"
)

// LoadedSkinstores is a list of all loaded skinstores
var LoadedSkinstores []datastore.Skinstore = []datastore.Skinstore{}

// LoadSkinstores searches plugins for a skinstore object
// and adds it to the 'skinstores` array
func LoadSkinstores(plugins []*plugin.Plugin) error {

	for _, plugin := range plugins {

		symSkinstore, err := plugin.Lookup("SkinstoreModule")
		if err != nil {
			log.Fatalln(fmt.Sprintf("Failed to load plugin: %v", err.Error()))
		}

		var skinstore datastore.Skinstore
		skinstore, ok := symSkinstore.(datastore.Skinstore)
		if !ok {
			return errors.New("Invalid type for Skinstore in plugin")
		}

		err = skinstore.Initialize()
		if err != nil {
			return err
		}

		LoadedSkinstores = append(LoadedSkinstores, skinstore)
	}
	return nil
}
