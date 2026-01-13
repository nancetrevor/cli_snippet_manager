package handlers

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/lfizzikz/snip/models"
)

func SaveCommands(f []models.Command) error {
	commandFile, err := CommandFileSavePath()
	if err != nil {
		return err
	}
	out, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(commandFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(commandFile, out, 0600)
}

func AppendCommands(f []models.Command) error {
	items, err := GetSavedCommands()
	if err != nil {
		if os.IsNotExist(err) {
			items = []models.Command{}
		} else {
			return err
		}
	}

	commandFile, err := CommandFileSavePath()
	if err != nil {
		return err
	}

	for _, x := range f {
		items = append(items, x)
	}

	out, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(commandFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(commandFile, out, 0600)
}

func GetSavedCommands() ([]models.Command, error) {
	c, err := CommandFileSavePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(c)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Command{}, nil
		}
		return nil, err
	}

	var items []models.Command
	if err = json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func CommandFileSavePath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, "commands.json")
	return p, nil
}
