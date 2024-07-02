package main

import (
	"fmt"
	g "github.com/AllenDang/giu"
)

const (
	MMURL = "ws://localhost:8080/v1/mm/server/connect?instanceID="
)

func getGameModeIndex(mode string) int {
	switch mode {
	case "BR":
		return 0
	case "TDM":
		return 3
	default:
		return 0
	}
}

var (
	Modes             = []string{"BR", "TDM"}
	totalPlayers      int32
	selectedMode      string
	selectedModeIndex int32
	previewModeValue  int32
	GameServerPort    int32
)

type AddServerRequestData struct {
	Type       string `json:"type"`
	RoomID     string `json:"roomID"`
	InstanceID string `json:"instanceID"`
	MatchID    int64  `json:"matchID"`
	Version    string `json:"version"`
	TeamSize   int    `json:"teamSize"`
	ServerType string `json:"serverType"`
	GameMode   int    `json:"gameMode"`
	MapName    string `json:"mapName"`
}

func addServer() {
	fmt.Println("port: ", GameServerPort)
	//spin the gameServer here
	server := NewServer(int(GameServerPort))
	server.addServer()
	server.connectToMM()
}
func loop() {
	g.SingleWindowWithMenuBar().Layout(
		g.Row(
			g.InputInt(&totalPlayers).Size(50),
			g.Button("Add Server").OnClick(addServer),
			g.Combo("Mode", Modes[previewModeValue], Modes, &selectedModeIndex).Size(100).OnChange(func() {
				previewModeValue = selectedModeIndex
			}),
			g.InputInt(&GameServerPort),
		),
	)
}

func main() {
	w := g.NewMasterWindow("Overview", 1000, 800, 0)
	w.Run(loop)
}
