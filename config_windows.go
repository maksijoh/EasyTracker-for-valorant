//go:build windows

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	henrikKeysURL   = "https://api.henrikdev.xyz/dashboard/api-keys"
	discordAppsURL  = "https://discord.com/developers/applications"
	discordDownload = "https://discord.com/download"
	riotAccountURL  = "https://account.riotgames.com/"
)

var (
	colorBackground = walk.RGB(15, 25, 35)
	colorSurface    = walk.RGB(31, 39, 49)
	colorInput      = walk.RGB(39, 49, 61)
	colorAccent     = walk.RGB(255, 70, 85)
	colorText       = walk.RGB(236, 232, 225)
	colorMuted      = walk.RGB(170, 178, 189)
	colorSuccess    = walk.RGB(69, 214, 143)
)

func solid(color walk.Color) SolidColorBrush {
	return SolidColorBrush{Color: color}
}

func openBrowser(target string) {
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", target).Start()
}

func resourceCard(title, description, buttonText, target string) Composite {
	return Composite{
		Background: solid(colorSurface),
		Layout:     VBox{Margins: Margins{Left: 12, Top: 10, Right: 12, Bottom: 10}, Spacing: 5},
		Children: []Widget{
			Label{Text: title, TextColor: colorText, Font: Font{Family: "Segoe UI", PointSize: 10, Bold: true}},
			Label{Text: description, TextColor: colorMuted},
			PushButton{Text: buttonText, Background: solid(colorAccent), Font: Font{Family: "Segoe UI", PointSize: 9, Bold: true}, OnClicked: func() { openBrowser(target) }},
		},
	}
}

func showConfigWindow(path string, existing appConfig) (bool, error) {
	var window *walk.MainWindow
	var tokenEdit, appIDEdit, nameEdit, tagEdit *walk.LineEdit
	var regionBox *walk.ComboBox
	var statusLabel *walk.Label
	var testButton *walk.PushButton
	regions := []string{"eu", "na", "latam", "br", "ap", "kr"}
	regionIndex := 0
	for index, region := range regions {
		if region == existing.RiotRegion {
			regionIndex = index
			break
		}
	}

	readForm := func() appConfig {
		region := "eu"
		if index := regionBox.CurrentIndex(); index >= 0 && index < len(regions) {
			region = regions[index]
		}
		return appConfig{
			AccessToken:  strings.TrimSpace(tokenEdit.Text()),
			DiscordAppID: strings.TrimSpace(appIDEdit.Text()),
			RiotName:     strings.TrimSpace(nameEdit.Text()),
			RiotTag:      strings.TrimSpace(strings.TrimPrefix(tagEdit.Text(), "#")),
			RiotRegion:   region,
		}
	}

	saved := false
	_, err := (MainWindow{
		AssignTo:   &window,
		Title:      fmt.Sprintf("EasyTracker %s — настройка", version),
		MinSize:    Size{Width: 680, Height: 640},
		Size:       Size{Width: 760, Height: 710},
		Background: solid(colorBackground),
		Font:       Font{Family: "Segoe UI", PointSize: 9},
		Layout:     VBox{Margins: Margins{Left: 22, Top: 20, Right: 22, Bottom: 20}, Spacing: 12},
		Children: []Widget{
			Composite{
				Background: solid(colorAccent),
				Layout:     VBox{Margins: Margins{Left: 18, Top: 14, Right: 18, Bottom: 14}, Spacing: 3},
				Children: []Widget{
					Label{Text: "EASYTRACKER", TextColor: colorText, Font: Font{Family: "Segoe UI", PointSize: 18, Bold: true}},
					Label{Text: "VALORANT  •  DISCORD RICH PRESENCE", TextColor: colorText, Font: Font{Family: "Segoe UI", PointSize: 9, Bold: true}},
				},
			},
			Label{Text: "Подключение аккаунта", TextColor: colorText, Font: Font{Family: "Segoe UI", PointSize: 13, Bold: true}},
			Label{Text: "Заполни данные один раз — они останутся только на этом компьютере.", TextColor: colorMuted},
			Composite{
				Background: solid(colorSurface),
				Layout:     Grid{Columns: 2, Spacing: 9, Margins: Margins{Left: 14, Top: 14, Right: 14, Bottom: 14}},
				Children: []Widget{
					Label{Text: "HenrikDev API key", TextColor: colorText, Font: Font{Bold: true}},
					LineEdit{AssignTo: &tokenEdit, Text: existing.AccessToken, PasswordMode: true, Background: solid(colorInput), TextColor: colorText, ToolTipText: "Ключ из HenrikDev Dashboard"},
					Label{Text: "Discord Application ID", TextColor: colorText, Font: Font{Bold: true}},
					LineEdit{AssignTo: &appIDEdit, Text: existing.DiscordAppID, Background: solid(colorInput), TextColor: colorText, ToolTipText: "Уже заполнено — обычно менять не нужно"},
					Label{Text: "Riot Name (до #)", TextColor: colorText, Font: Font{Bold: true}},
					LineEdit{AssignTo: &nameEdit, Text: existing.RiotName, Background: solid(colorInput), TextColor: colorText, ToolTipText: "Например: PlayerName"},
					Label{Text: "Riot Tag (после #)", TextColor: colorText, Font: Font{Bold: true}},
					LineEdit{AssignTo: &tagEdit, Text: existing.RiotTag, Background: solid(colorInput), TextColor: colorText, ToolTipText: "Например: EUW"},
					Label{Text: "Регион", TextColor: colorText, Font: Font{Bold: true}},
					ComboBox{AssignTo: &regionBox, Model: regions, CurrentIndex: regionIndex, Background: solid(colorInput)},
				},
			},
			Label{Text: "Где взять необходимые данные", TextColor: colorText, Font: Font{Family: "Segoe UI", PointSize: 13, Bold: true}},
			Composite{
				Background: solid(colorBackground),
				Layout:     Grid{Columns: 2, Spacing: 10},
				Children: []Widget{
					resourceCard("1. HenrikDev API key", "Создай бесплатный ключ для Valorant API.", "ПОЛУЧИТЬ API KEY", henrikKeysURL),
					resourceCard("2. Riot ID", "Посмотри имя и тег аккаунта Riot.", "ОТКРЫТЬ RIOT ACCOUNT", riotAccountURL),
					resourceCard("3. Discord Application", "Application ID находится в General Information.", "ОТКРЫТЬ DEVELOPER PORTAL", discordAppsURL),
					resourceCard("4. Discord Desktop", "Rich Presence работает только с запущенным Discord.", "СКАЧАТЬ DISCORD", discordDownload),
				},
			},
			Label{AssignTo: &statusLabel, Text: "Можно сохранить сразу или сначала проверить API и Riot ID.", TextColor: colorMuted},
			Composite{
				Background: solid(colorBackground),
				Layout:     HBox{Spacing: 10},
				Children: []Widget{
					PushButton{
						AssignTo: &testButton, Text: "ПРОВЕРИТЬ ДАННЫЕ", Background: solid(colorSurface), Font: Font{Family: "Segoe UI", PointSize: 9, Bold: true},
						OnClicked: func() {
							config := readForm()
							if !config.valid() {
								statusLabel.SetText("Заполни все обязательные поля перед проверкой.")
								statusLabel.SetTextColor(colorAccent)
								return
							}
							testButton.SetEnabled(false)
							statusLabel.SetText("Проверяем HenrikDev API и Riot ID…")
							statusLabel.SetTextColor(colorMuted)
							go func() {
								oldToken := os.Getenv("ACCESS_TOKEN")
								_ = os.Setenv("ACCESS_TOKEN", config.AccessToken)
								mmr, checkErr := fetchMMR(config.RiotName, config.RiotTag, config.RiotRegion)
								_ = os.Setenv("ACCESS_TOKEN", oldToken)
								window.Synchronize(func() {
									testButton.SetEnabled(true)
									if checkErr != nil {
										statusLabel.SetText("Ошибка проверки: " + checkErr.Error())
										statusLabel.SetTextColor(colorAccent)
										return
									}
									rank := mmr.Data.Current.Tier.Name
									if rank == "" {
										rank = "Unrated"
									}
									statusLabel.SetText(fmt.Sprintf("Готово: %s#%s найден • %s • %d RR", mmr.Data.Account.Name, mmr.Data.Account.Tag, rank, mmr.Data.Current.RR))
									statusLabel.SetTextColor(colorSuccess)
								})
							}()
						},
					},
					HSpacer{},
					PushButton{
						Text: "СОХРАНИТЬ И ЗАПУСТИТЬ", Background: solid(colorAccent), Font: Font{Family: "Segoe UI", PointSize: 10, Bold: true},
						OnClicked: func() {
							config := readForm()
							if !config.valid() {
								statusLabel.SetText("Заполни API key, Riot Name, Riot Tag и Discord Application ID.")
								statusLabel.SetTextColor(colorAccent)
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
