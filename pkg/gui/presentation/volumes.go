package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazydocker/pkg/commands"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func GetVolumeDisplayStrings(volume *commands.Volume, isSelected bool) []string {
	return []string{getVolumeSelectionMarker(isSelected), volume.Volume.Driver, volume.Name}
}

func getVolumeSelectionMarker(isSelected bool) string {
	if isSelected {
		return utils.ColoredString("[x]", color.FgGreen)
	}
	return "[ ]"
}
