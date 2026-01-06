package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazydocker/pkg/commands"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func GetNetworkDisplayStrings(network *commands.Network, isSelected bool) []string {
	return []string{getNetworkSelectionMarker(isSelected), network.Network.Driver, network.Name}
}

func getNetworkSelectionMarker(isSelected bool) string {
	if isSelected {
		return utils.ColoredString("[x]", color.FgGreen)
	}
	return "[ ]"
}
