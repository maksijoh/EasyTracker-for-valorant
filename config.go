package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type appConfig struct {
	AccessToken  string
	DiscordAppID string
	RiotName     string
	RiotTag      string
	RiotRegion   string
}

func configPath() string {
	cwd, err := os.Getwd()
	if err == nil {
		if _, statErr := os.Stat(filepath.Join(cwd, "go.mod")); statErr == nil {
			return filepath.Join(cwd, ".env")
		}
		if _, statErr := os.Stat(filepath.Join(cwd, ".env")); statErr == nil {
			return filepath.Join(cwd, ".env")
		}
	}

	executable, err := os.Executable()
	if err == nil {
		return filepath.Join(filepath.Dir(executable), ".env")
	}
	return ".env"
}

func currentConfig() appConfig {
	return appConfig{
		AccessToken:  strings.TrimSpace(os.Getenv("ACCESS_TOKEN")),
		DiscordAppID: envOr("DISCORD_APP_ID", defaultDiscordAppID),
		RiotName:     strings.TrimSpace(os.Getenv("RIOT_NAME")),
		RiotTag:      strings.TrimSpace(os.Getenv("RIOT_TAG")),
		RiotRegion:   strings.ToLower(envOr("RIOT_REGION", "eu")),
	}
}

func (c appConfig) valid() bool {
	return c.AccessToken != "" && c.DiscordAppID != "" && c.RiotName != "" && c.RiotTag != ""
}

func saveConfig(path string, c appConfig) error {
	values := map[string]string{
		"ACCESS_TOKEN":        c.AccessToken,
		"DISCORD_APP_ID":      c.DiscordAppID,
		"DISCORD_LARGE_IMAGE": defaultPresenceImage,
		"RIOT_NAME":           c.RiotName,
		"RIOT_TAG":            c.RiotTag,
		"RIOT_REGION":         c.RiotRegion,
	}
	for key, value := range values {
		if strings.ContainsAny(value, "\r\n") {
			return fmt.Errorf("%s contains an invalid line break", key)
		}
	}

	content := fmt.Sprintf(
		"ACCESS_TOKEN=%s\nDISCORD_APP_ID=%s\nDISCORD_LARGE_IMAGE=%s\nRIOT_NAME=%s\nRIOT_TAG=%s\nRIOT_REGION=%s\n",
		c.AccessToken, c.DiscordAppID, defaultPresenceImage, c.RiotName, c.RiotTag, c.RiotRegion,
	)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return err
	}
	for key, value := range values {
		if err := os.Setenv(key, value); err != nil {
			return err
		}
	}
	return nil
}
