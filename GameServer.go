package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	Port       int
	httpServer *http.Server
	InstanceId string
	mux        *http.ServeMux
}

func NewServer(port int) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", handler)

	return &Server{
		Port: port,
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		InstanceId: strconv.Itoa(rand.Intn(10000)),
		mux:        mux,
	}
}

func (s *Server) addServer() {
	go func() {
		log.Printf("Starting server on port %d\n", s.Port)
		if err := s.httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *Server) connectToMM() {
	conn, _, err := websocket.DefaultDialer.Dial(MMURL+s.InstanceId, nil)
	if err != nil {
		fmt.Println("Error connecting to MM:", err)
	}
	roomId := rand.Intn(10000)
	request := AddServerRequestData{
		Type:       "AddServer",
		RoomID:     strconv.Itoa(roomId),
		InstanceID: s.InstanceId,
		MatchID:    1000,
		Version:    "0.32.xx #yyyy",
		TeamSize:   4,
		ServerType: "ni",
		GameMode:   getGameModeIndex(Modes[selectedModeIndex]),
		MapName:    "VR",
	}
	err = conn.WriteJSON(request)
	if err != nil {
		log.Println("Error sending add server request: ", err)
		return
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
