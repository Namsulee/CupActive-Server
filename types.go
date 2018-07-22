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
}

// UserSettingReq structure contains request information for web server.
type UserSettingReq struct {
	Cmd  		string `json:"cmd"`
	UniqueID 	string `json:"id"`
	Capability	int    `json:"capability"`
}

// Define our message object
type Message struct {
	Cmd 		string `json:"cmd"`
	UniqueID    string `json:"id"`
	Capability  int `json:"capability"`
}

type Device struct {
	WS *websocket.Conn
	UniqueID string
	Capability int
}