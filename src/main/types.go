package main

import (
	"github.com/gorilla/websocket"
)


// Command structure contains command information.
type Command struct {
    Cmd string `json:"cmd"`
}

// ConnectReq structure contains request information for web server.
type ConnectReq struct {
	Cmd  		string `json:"cmd"`
	UniqueID 	string `json:"id"`
	IPAddress	string `json:"ipaddress"`
}

// UserSettingReq structure contains request information for web server.
type UserSettingReq struct {
	Cmd  		string `json:"cmd"`
	UniqueID 	string `json:"id"`
	Capability	int    `json:"capability"`
}

// GameSettingReq structure contains request information for web server.
type GameSettingReq struct {
	Cmd  		string `json:"cmd"`
	UniqueID 	string `json:"id"`
	Kind		int    `json:"kind"`
	GameState   int    `json:"state"`
	Drink		int    `json:"drink"`
}
// Define our message object
type Message struct {
	Cmd 		string `json:"cmd"`
	UniqueID    string `json:"id"`
	Capability  int    `json:"capability"`
	Kind		int    `json:"kind"`
	GameState   int    `json:"state"`
	Drink		int    `json:"drink"`
}

type UserSetting struct {
	Cmd 		string `json:"cmd"`
	Cap         []int  `json:"cap"`
}

type Device struct {
	WS *websocket.Conn
	UniqueID string
	IPAddress string
	Capability int
}

type cupInfo struct {
	Name 		string `json:"Name"`
	IPAddress	string `json:"IPAddress"`
	Capability  int    `json:"Capability"`
	HostPort    string `json:"HostPort"`
}

type webUIInfo struct {
	Title	string 		`json:"title"`
	Lists 	[]cupInfo   `json:"lists`
}

