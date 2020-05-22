package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// input api token
var apiToken = ""

// input jellyfin address
var server = "http://1.1.1.1:8096"

func main() {
	serverAddr := server + "/Sessions"
	req, err := http.NewRequest("Get", serverAddr, nil)
	if err != nil {
		fmt.Printf("error creating http request: %s \n", err)
		os.Exit(1)
	}
	req.Header.Set("X-Emby-Token", apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error connecting to jellyfin server: %s \n", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	var dataSessionsMap []Sessions
	var count int
	err = json.NewDecoder(resp.Body).Decode(&dataSessionsMap)

	for i, s := range dataSessionsMap {
		if s.NowPlaying.Name != "" {
			count++
			fmt.Printf("#%d:\n", i)
			fmt.Printf("user: %v\n", s.UserName)
			fmt.Printf("id: %v\n", s.UserID)
			fmt.Printf("client: %v\n", s.Client)
			fmt.Printf("ip: %v\n", s.IPaddr)
			fmt.Printf("playing: %v\n", s.NowPlaying.Name)
			fmt.Printf("position: %v / %v\n", timeConvert(s.PlaySt.Position).Truncate(time.Second), timeConvert(s.NowPlaying.RunTime).Truncate(time.Second))
			fmt.Printf("method: %v\n", s.PlaySt.PlayMethod)
			fmt.Printf("paused: %v\n", s.PlaySt.IsPaused)
			fmt.Printf("-----------------------\n")
		}
	}
	fmt.Printf("\n%v users online\n", count)
}

type Sessions struct {
	PlaySt     PlayState      `json:"PlayState"`
	IPaddr     string         `json:"RemoteEndPoint"`
	UserID     string         `json:"UserId"`
	UserName   string         `json:"UserName"`
	Client     string         `json:"Client"`
	NowPlaying NowPlayingItem `json:"NowPlayingItem"`
}

type PlayState struct {
	Position   int    `json:"PositionTicks"`
	IsPaused   bool   `json:"IsPaused"`
	IsMuted    bool   `json:"IsMuted"`
	PlayMethod string `json:"PlayMethod"`
}

type NowPlayingItem struct {
	Name    string `json:"Name"`
	RunTime int    `json:"RunTimeTicks"`
}

func timeConvert(ticks int) (length time.Duration) {
	x := ticks / 10
	strX := strconv.Itoa(x)
	strX = strX + "us"
	micro, _ := time.ParseDuration(strX)
	return micro
}
