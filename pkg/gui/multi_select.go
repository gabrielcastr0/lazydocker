package gui

import (
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazydocker/pkg/commands"
)

// hasAnySelections returns true if any items are selected across all panels
func (gui *Gui) hasAnySelections() bool {
	return len(gui.State.MultiSelect.Containers) > 0 ||
		len(gui.State.MultiSelect.Images) > 0 ||
		len(gui.State.MultiSelect.Volumes) > 0 ||
		len(gui.State.MultiSelect.Networks) > 0
}

// clearAllSelections clears all multi-select state
func (gui *Gui) clearAllSelections() {
	gui.State.MultiSelect.Containers = make(map[string]bool)
	gui.State.MultiSelect.Images = make(map[string]bool)
	gui.State.MultiSelect.Volumes = make(map[string]bool)
	gui.State.MultiSelect.Networks = make(map[string]bool)
}

// rerenderAllPanels re-renders all panels to reflect selection state changes
func (gui *Gui) rerenderAllPanels() error {
	if err := gui.Panels.Containers.RerenderList(); err != nil {
		return err
	}
	if err := gui.Panels.Images.RerenderList(); err != nil {
		return err
	}
	if err := gui.Panels.Volumes.RerenderList(); err != nil {
		return err
	}
	if err := gui.Panels.Networks.RerenderList(); err != nil {
		return err
	}
	return nil
}

// getSelectionCount returns the total number of selected items across all panels
func (gui *Gui) getSelectionCount() int {
	return len(gui.State.MultiSelect.Containers) +
		len(gui.State.MultiSelect.Images) +
		len(gui.State.MultiSelect.Volumes) +
		len(gui.State.MultiSelect.Networks)
}

// Select all handlers for each panel type

func (gui *Gui) handleSelectAllContainers() error {
	for _, ctr := range gui.Panels.Containers.List.GetItems() {
		gui.State.MultiSelect.Containers[ctr.ID] = true
	}
	return gui.Panels.Containers.RerenderList()
}

func (gui *Gui) handleSelectAllImages() error {
	for _, img := range gui.Panels.Images.List.GetItems() {
		gui.State.MultiSelect.Images[img.ID] = true
	}
	return gui.Panels.Images.RerenderList()
}

func (gui *Gui) handleSelectAllVolumes() error {
	for _, vol := range gui.Panels.Volumes.List.GetItems() {
		gui.State.MultiSelect.Volumes[vol.Name] = true
	}
	return gui.Panels.Volumes.RerenderList()
}

func (gui *Gui) handleSelectAllNetworks() error {
	for _, net := range gui.Panels.Networks.List.GetItems() {
		gui.State.MultiSelect.Networks[net.Name] = true
	}
	return gui.Panels.Networks.RerenderList()
}

// Deselect all handlers for each panel type

func (gui *Gui) handleDeselectAllContainers() error {
	gui.State.MultiSelect.Containers = make(map[string]bool)
	return gui.Panels.Containers.RerenderList()
}

func (gui *Gui) handleDeselectAllImages() error {
	gui.State.MultiSelect.Images = make(map[string]bool)
	return gui.Panels.Images.RerenderList()
}

func (gui *Gui) handleDeselectAllVolumes() error {
	gui.State.MultiSelect.Volumes = make(map[string]bool)
	return gui.Panels.Volumes.RerenderList()
}

func (gui *Gui) handleDeselectAllNetworks() error {
	gui.State.MultiSelect.Networks = make(map[string]bool)
	return gui.Panels.Networks.RerenderList()
}

// Toggle selection handlers for each panel type

func (gui *Gui) handleToggleContainerSelection() error {
	ctr, err := gui.Panels.Containers.GetSelectedItem()
	if err != nil {
		return nil
	}

	if gui.State.MultiSelect.Containers[ctr.ID] {
		delete(gui.State.MultiSelect.Containers, ctr.ID)
	} else {
		gui.State.MultiSelect.Containers[ctr.ID] = true
	}

	// Move to next item after toggling for easier batch selection
	gui.Panels.Containers.SelectNextLine()

	return gui.Panels.Containers.RerenderList()
}

func (gui *Gui) handleToggleImageSelection() error {
	img, err := gui.Panels.Images.GetSelectedItem()
	if err != nil {
		return nil
	}

	if gui.State.MultiSelect.Images[img.ID] {
		delete(gui.State.MultiSelect.Images, img.ID)
	} else {
		gui.State.MultiSelect.Images[img.ID] = true
	}

	gui.Panels.Images.SelectNextLine()

	return gui.Panels.Images.RerenderList()
}

func (gui *Gui) handleToggleVolumeSelection() error {
	vol, err := gui.Panels.Volumes.GetSelectedItem()
	if err != nil {
		return nil
	}

	if gui.State.MultiSelect.Volumes[vol.Name] {
		delete(gui.State.MultiSelect.Volumes, vol.Name)
	} else {
		gui.State.MultiSelect.Volumes[vol.Name] = true
	}

	gui.Panels.Volumes.SelectNextLine()

	return gui.Panels.Volumes.RerenderList()
}

func (gui *Gui) handleToggleNetworkSelection() error {
	net, err := gui.Panels.Networks.GetSelectedItem()
	if err != nil {
		return nil
	}

	if gui.State.MultiSelect.Networks[net.Name] {
		delete(gui.State.MultiSelect.Networks, net.Name)
	} else {
		gui.State.MultiSelect.Networks[net.Name] = true
	}

	gui.Panels.Networks.SelectNextLine()

	return gui.Panels.Networks.RerenderList()
}

// handleDeleteSelected handles the global delete selected command
func (gui *Gui) handleDeleteSelected(g *gocui.Gui, v *gocui.View) error {
	totalSelected := gui.getSelectionCount()

	if totalSelected == 0 {
		return nil // Nothing selected
	}

	confirmMsg := gui.buildDeleteConfirmationMessage()

	return gui.createConfirmationPanel(
		gui.Tr.Confirm,
		confirmMsg,
		func(g *gocui.Gui, v *gocui.View) error {
			return gui.deleteSelectedItems()
		},
		nil,
	)
}

// buildDeleteConfirmationMessage builds a message showing what will be deleted
func (gui *Gui) buildDeleteConfirmationMessage() string {
	var parts []string

	if n := len(gui.State.MultiSelect.Containers); n > 0 {
		parts = append(parts, fmt.Sprintf("%d container(s)", n))
	}
	if n := len(gui.State.MultiSelect.Images); n > 0 {
		parts = append(parts, fmt.Sprintf("%d image(s)", n))
	}
	if n := len(gui.State.MultiSelect.Volumes); n > 0 {
		parts = append(parts, fmt.Sprintf("%d volume(s)", n))
	}
	if n := len(gui.State.MultiSelect.Networks); n > 0 {
		parts = append(parts, fmt.Sprintf("%d network(s)", n))
	}

	return fmt.Sprintf("%s\n\n%s", gui.Tr.ConfirmDeleteSelected, strings.Join(parts, ", "))
}

// deleteSelectedItems deletes all selected items across all panels
func (gui *Gui) deleteSelectedItems() error {
	return gui.WithWaitingStatus(gui.Tr.RemovingStatus, func() error {
		var errors []string

		// 1. Delete containers first (they may reference images/volumes)
		for containerID := range gui.State.MultiSelect.Containers {
			ctr := gui.findContainerByID(containerID)
			if ctr != nil {
				if err := ctr.Remove(container.RemoveOptions{Force: true}); err != nil {
					errors = append(errors, fmt.Sprintf("Container %s: %v", ctr.Name, err))
				}
			}
		}

		// 2. Delete images
		for imageID := range gui.State.MultiSelect.Images {
			img := gui.findImageByID(imageID)
			if img != nil {
				if err := img.Remove(image.RemoveOptions{Force: true, PruneChildren: true}); err != nil {
					errors = append(errors, fmt.Sprintf("Image %s: %v", img.Name, err))
				}
			}
		}

		// 3. Delete volumes
		for volumeName := range gui.State.MultiSelect.Volumes {
			vol := gui.findVolumeByName(volumeName)
			if vol != nil {
				if err := vol.Remove(true); err != nil { // force=true
					errors = append(errors, fmt.Sprintf("Volume %s: %v", vol.Name, err))
				}
			}
		}

		// 4. Delete networks
		for networkName := range gui.State.MultiSelect.Networks {
			net := gui.findNetworkByName(networkName)
			if net != nil {
				if err := net.Remove(); err != nil {
					errors = append(errors, fmt.Sprintf("Network %s: %v", net.Name, err))
				}
			}
		}

		// Clear selections after deletion
		gui.clearAllSelections()

		// Handle errors
		if len(errors) > 0 {
			return gui.createErrorPanel(strings.Join(errors, "\n"))
		}

		return nil
	})
}

// Helper functions to find items by ID/name

func (gui *Gui) findContainerByID(id string) *commands.Container {
	for _, c := range gui.Panels.Containers.List.GetAllItems() {
		if c.ID == id {
			return c
		}
	}
	return nil
}

func (gui *Gui) findImageByID(id string) *commands.Image {
	for _, img := range gui.Panels.Images.List.GetAllItems() {
		if img.ID == id {
			return img
		}
	}
	return nil
}

func (gui *Gui) findVolumeByName(name string) *commands.Volume {
	for _, vol := range gui.Panels.Volumes.List.GetAllItems() {
		if vol.Name == name {
			return vol
		}
	}
	return nil
}

func (gui *Gui) findNetworkByName(name string) *commands.Network {
	for _, net := range gui.Panels.Networks.List.GetAllItems() {
		if net.Name == name {
			return net
		}
	}
	return nil
}
