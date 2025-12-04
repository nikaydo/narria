package plugins

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type Plugins struct {
	Plugins  map[uuid.UUID]PluginData
	Frontend *http.ServeMux
}

func (n *Plugins) LoadPlugin(path string) error {
	dirs, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		file, err := os.Open(path + "/" + dir.Name() + "/" + "plugin.json")
		if err != nil {
			return err
		}
		defer file.Close()
		var plugin PluginData
		err = json.NewDecoder(file).Decode(&plugin.Plugnin)
		if err != nil {
			return err
		}
		pId, err := uuid.Parse(plugin.Plugnin.Id)
		if err != nil {
			return err
		}

		if plugin.Plugnin.FrontendEntry != "" {
			if err := plugin.InitFrontend(); err != nil {
				return err
			}
		}

		n.Plugins[pId] = plugin
	}
	return nil
}
