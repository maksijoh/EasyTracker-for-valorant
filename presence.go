package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	discord "github.com/maintainer64/rich-go/client"
)

const defaultDiscordAppID = "1538271778651111514"

type mmrResponse struct {
	Status int `json:"status"`
	Data   struct {
		Account struct {
			Name string `json:"name"`
			Tag  string `json:"tag"`
		} `json:"account"`
		Current struct {
			Tier struct {
				Name string `json:"name"`
			} `json:"tier"`
			RR         int `json:"rr"`
			LastChange int `json:"last_change"`
		} `json:"current"`
	} `json:"data"`
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func fetchMMR(name, tag, region string) (mmrResponse, error) {
	var result mmrResponse
	apiURL := "https://api.henrikdev.xyz/valorant/v3/mmr/" +
		url.PathEscape(region) + "/pc/" + url.PathEscape(name) + "/" + url.PathEscape(tag)

	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return result, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", os.Getenv("ACCESS_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("HenrikDev returned HTTP %d", resp.StatusCode)
	}
	return result, nil
}

func updatePresence(name, tag, region string) error {
	mmr, err := fetchMMR(name, tag, region)
	if err != nil {
		return err
	}

	rank := mmr.Data.Current.Tier.Name
	if rank == "" {
		rank = "Unrated"
	}
	change := fmt.Sprintf("%+d RR за последний матч", mmr.Data.Current.LastChange)

	return discord.SetActivity(discord.Activity{
		Details:   fmt.Sprintf("%s — %d RR", rank, mmr.Data.Current.RR),
		State:     change,
		LargeText: fmt.Sprintf("%s#%s • %s", name, tag, strings.ToUpper(region)),
	})
}

func runPresence() {
	appID := envOr("DISCORD_APP_ID", defaultDiscordAppID)
	name := strings.TrimSpace(os.Getenv("RIOT_NAME"))
	tag := strings.TrimSpace(os.Getenv("RIOT_TAG"))
	region := strings.ToLower(envOr("RIOT_REGION", "eu"))
	if name == "" || tag == "" {
		log.Print("RIOT_NAME and RIOT_TAG are missing: copy .env.example to .env, fill them in, then restart; Rich Presence disabled")
		return
	}

	for {
		if err := discord.Login(appID); err != nil {
			log.Printf("Discord is unavailable, retrying in 15s: %v", err)
			time.Sleep(15 * time.Second)
			continue
		}

		log.Printf("Discord Rich Presence connected for %s#%s", name, tag)
		for {
			if err := updatePresence(name, tag, region); err != nil {
				log.Printf("Rich Presence update failed: %v", err)
				discord.Logout()
				break
			}
			time.Sleep(time.Minute)
		}
		time.Sleep(15 * time.Second)
	}
}
