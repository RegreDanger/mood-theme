package data

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func BroadcastTheme(name string) error {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	basePath := filepath.Join(userConfigDir, "Code", "User")
	objectives := []string{filepath.Join(basePath, "settings.json")}
	profilesDir := filepath.Join(basePath, "profiles")
	entries, err := os.ReadDir(profilesDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			path := filepath.Join(profilesDir, entry.Name(), "settings.json")
			objectives = append(objectives, path)
		}
	}

	theme, err := FetchTheme(name)

	if err != nil {
		return err
	}

	for _, path := range objectives {
		if err := injectThemeToFile(path, theme); err != nil {
			continue
		}
	}

	return nil
}

func injectThemeToFile(pathSettings string, theme string) error {

	var config = make(map[string]any)

	jsonConfigBytes, err := os.ReadFile(pathSettings)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonConfigBytes, &config); err != nil {
		return err
	}

	config["workbench.colorTheme"] = theme

	modifiedBytes, err := json.MarshalIndent(config, "", "	")
	if err != nil {
		return err
	}
	err = os.WriteFile(pathSettings, modifiedBytes, 0644)
	if err != nil {
		return err
	}

	return nil

}
