package plugins

type PluginData struct {
	Plugnin PluginInfo `json:"plugnin"`
	Tokens  *Tokens
}

type PluginInfo struct {
	Id            string      `json:"id"`
	Description   Description `json:"description"`
	FrontendEntry string      `json:"frontend_entry"`
}

type Description struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

func (n *PluginInfo) GetId() string {
	return n.Id
}

func (n *PluginInfo) GetDescription() Description {
	return n.Description
}

func (n *PluginData) InitFrontend() error {
	if err := n.Tokens.GenSecret(); err != nil {
		return err
	}

	if err := n.Tokens.GetToken(); err != nil {
		return err
	}
	return nil
}
