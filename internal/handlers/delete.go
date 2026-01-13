package handlers

import "github.com/lfizzikz/snip/models"

func DeleteCommandByName(name string) error {
	cmds, err := GetSavedCommands()
	if err != nil {
		return err
	}

	out := make([]models.Command, 0, len(cmds))
	for _, c := range cmds {
		if c.Title() != name {
			out = append(out, c)
		}
	}
	return SaveCommands(out)
}
