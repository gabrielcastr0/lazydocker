package presentation

import (
	"github.com/fatih/color"
	"github.com/jesseduffield/lazydocker/pkg/commands"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func GetImageDisplayStrings(image *commands.Image, isSelected bool) []string {
	return []string{
		getImageSelectionMarker(isSelected),
		image.Name,
		image.Tag,
		utils.FormatDecimalBytes(int(image.Image.Size)),
	}
}

func getImageSelectionMarker(isSelected bool) string {
	if isSelected {
		return utils.ColoredString("[x]", color.FgGreen)
	}
	return "[ ]"
}
