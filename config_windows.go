//go:build windows

package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	henrikKeysURL  = "https://api.henrikdev.xyz/dashboard/api-keys"
	discordAppsURL = "https://discord.com/developers/applications"
)

func openBrowser(target string) {
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
}

func showConfigWindow(path string, existing appConfig) (bool, error) {
	var window *walk.MainWindow
	var tokenEdit, appIDEdit, nameEdit, tagEdit *walk.LineEdit
	var regionBox *walk.ComboBox
	regions := []string{"eu", "na", "latam", "br", "ap", "kr"}
	regionIndex := 0
	for index, region := range regions {
		if region == existing.RiotRegion {
			regionIndex = index
			break
		}
	}

	saved := false
	_, err := (MainWindow{
		AssignTo: &window,
		Title:    fmt.Sprintf("EasyTracker %s — настройка", version),
		MinSize:  Size{Width: 560, Height: 470},
		Size:     Size{Width: 620, Height: 520},
		Layout:   VBox{Margins: Margins{Left: 22, Top: 20, Right: 22, Bottom: 20}, Spacing: 10},
		Children: []Widget{
			Label{Text: "Настройка Valorant Discord Rich Presence", Font: Font{Family: "Segoe UI", PointSize: 15, Bold: true}},
			Label{Text: "Данные сохраняются только рядом с EasyTracker.exe в локальном файле .env."},
			VSpacer{Size: 4},
			Composite{
				Layout: Grid{Columns: 2, Spacing: 8},
				Children: []Widget{
					Label{Text: "HenrikDev API key:"},
					LineEdit{AssignTo: &tokenEdit, Text: existing.AccessToken, PasswordMode: true},
					Label{Text: "Discord Application ID:"},
					LineEdit{AssignTo: &appIDEdit, Text: existing.DiscordAppID},
					Label{Text: "Riot Name (до #):"},
					LineEdit{AssignTo: &nameEdit, Text: existing.RiotName},
					Label{Text: "Riot Tag (после #):"},
					LineEdit{AssignTo: &tagEdit, Text: existing.RiotTag},
					Label{Text: "Регион:"},
					ComboBox{AssignTo: &regionBox, Model: regions, CurrentIndex: regionIndex},
				},
			},
			GroupBox{
				Title:  "Где взять данные",
				Layout: VBox{Spacing: 7},
				Children: []Widget{
					Label{Text: "API key: открой HenrikDev Dashboard и создай ключ."},
					PushButton{Text: "Открыть HenrikDev API Keys", OnClicked: func() { openBrowser(henrikKeysURL) }},
					Label{Text: "Application ID: Discord Developer Portal → ваше приложение → General Information."},
					PushButton{Text: "Открыть Discord Developer Portal", OnClicked: func() { openBrowser(discordAppsURL) }},
				},
			},
			VSpacer{},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						Text: "Сохранить и запустить",
						OnClicked: func() {
							region := "eu"
							if index := regionBox.CurrentIndex(); index >= 0 && index < len(regions) {
								region = regions[index]
							}
							config := appConfig{
								AccessToken:  strings.TrimSpace(tokenEdit.Text()),
								DiscordAppID: strings.TrimSpace(appIDEdit.Text()),
								RiotName:     strings.TrimSpace(nameEdit.Text()),
								RiotTag:      strings.TrimSpace(strings.TrimPrefix(tagEdit.Text(), "#")),
								RiotRegion:   region,
							}
							if !config.valid() {
								walk.MsgBox(window, "Не все поля заполнены", "Заполни API key, Discord Application ID, Riot Name и Riot Tag.", walk.MsgBoxIconWarning)
								return
							}
							if err := saveConfig(path, config); err != nil {
								walk.MsgBox(window, "Ошибка сохранения", err.Error(), walk.MsgBoxIconError)
								return
							}
							saved = true
							walk.MsgBox(window, "EasyTracker запущен", "Настройки сохранены. Оставь EasyTracker запущенным вместе с Discord.", walk.MsgBoxIconInformation)
							window.Close()
						},
					},
				},
			},
		},
	}).Run()
	if err != nil {
		return false, fmt.Errorf("open configuration window: %w", err)
	}
	return saved, nil
}

func showFatalError(message string) {
	walk.MsgBox(nil, "EasyTracker — ошибка", message, walk.MsgBoxIconError)
}
