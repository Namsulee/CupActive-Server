package main

import (
	"flag"
    "os"
	"log"
    "errors"
    "time"
    "math/rand"
	"net/http"
    "encoding/json"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var deviceList []Device

var webuiDir = "www-static"
var contentsDir = "/www-static/"
const contWebInfoFile = "cups.json"

// Static writeCupsInfo writing web-ui json data which includes cups info
func writeCupsInfo() error {

    var err error
    
    if _, err = os.Stat(webuiDir + "/" + contWebInfoFile); os.IsNotExist(err) {
        log.Println("cups.json is not existed, it will be created")
    } else {
        // file delete and create
        log.Println("cups.json would be removed")
        err = os.Remove(webuiDir + "/" + contWebInfoFile)
        if err != nil {
            log.Printf("Error removing")
        }
    }
       
    var writeData webUIInfo

    writeData.Title = "Cup List"
    for _, dev := range deviceList {
        var List cupInfo

        List.Name = dev.UniqueID
        List.IPAddress = dev.IPAddress
        List.Capability = dev.Capability

        writeData.Lists = append(writeData.Lists, List)
    }

    // Before writing
    for i, list := range writeData.Lists {
        log.Printf("[%d] Name             :[%s]", i, list.Name)
        log.Printf("[%d] IP               :[%s]", i, list.IPAddress)
        log.Printf("[%d] Capability       :[%s]", i, list.Capability)
    }

    //file writing into json
    f, err := os.Create(webuiDir + "/" + contWebInfoFile)
    if err != nil {
        log.Printf("contwebinfo file create error!!!")
        return err
    }

    defer f.Close()

    Lists, err := json.Marshal(writeData)
    _, err = f.Write(Lists)
    if err != nil {
        log.Printf("writeContainerListInfo file write error")
        return err
    }
   
    return err
}

func wsRegister(c *websocket.Conn, message []byte) error {
    var dev Device
    var err error
    var find bool

    rcv := ConnectReq{}

    json.Unmarshal([]byte(message), &rcv)
    log.Printf("cmd : [%s]", rcv.Cmd)
    log.Printf("id : [%s]", rcv.UniqueID)
    log.Printf("ip : [%s]", rcv.IPAddress)

    
    if len(rcv.UniqueID) > 0 {
        for i, cup := range deviceList {
            if cup.UniqueID == rcv.UniqueID {
                find = true
                log.Printf("find same index [%d]", i)
                cup.WS = c
                cup.UniqueID = rcv.UniqueID
                cup.IPAddress = rcv.IPAddress

                deviceList[i] = cup
                break
            }
        }

        if find != true {
            dev.WS = c
            dev.UniqueID = rcv.UniqueID
            dev.IPAddress = rcv.IPAddress

            deviceList = append(deviceList, dev)
        }
        // write cup list in json file
        writeCupsInfo()
    } else {
        log.Printf("UniqueID is not existed")
        err = errors.New("UniqueID is not existed")
    }

    return err
}

func random(min, max int) int {
    return rand.Intn(max - min) + min
}

func start(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
	_, msg, err := c.ReadMessage()
	if err != nil {
	    log.Println("read:", err)
	    break
	}
	log.Printf("recv: %s", msg)

	rcv := Command{}
	json.Unmarshal([]byte(msg), &rcv)
	log.Printf("cmd : [%s]", rcv.Cmd)

	switch rcv.Cmd {
	   case "register":
	    log.Printf("command register device")
            err := wsRegister(c, msg)
            if err == nil {
                // Confirm message to register
                send := Command{}
                send.Cmd = "connected"

                if err = c.WriteJSON(send); err != nil {
                    log.Println(err)
                }
            }
        case "usersetting":
            log.Printf("command usersetting")
            data := UserSetting{}
            json.Unmarshal([]byte(msg), &data)
            info := data.Cap

            for i, cup := range info {
                log.Printf("capability[%d] [%d]", i,  cup)
            }

            for i, dev := range deviceList {
                dev.Capability = data.Cap[i]

                send := UserSettingReq{}
                send.Cmd = "usersetting"
                send.UniqueID = dev.UniqueID 
                send.Capability = dev.Capability
                ws := dev.WS
               
                if err = ws.WriteJSON(send); err != nil {
                    log.Println(err)
                }
                
            }
        case "restart":
            log.Printf("command restart")
            send := Command{}
            send.Cmd = "restart"

            for _, dev := range deviceList {
                ws := dev.WS
                if err = ws.WriteJSON(send); err != nil {
                    log.Println(err)
                }
            }
        case "gamesetting":
            log.Printf("command gamesetting")
            data := Message{}
            json.Unmarshal([]byte(msg), &data)

            send := GameSettingReq{}
            send.Cmd = "gamesetting"
            length := len(deviceList)
            log.Printf("count [%d]", length)
            send.Kind = data.Kind 
            log.Printf("kind [%d]", send.Kind)
            if data.GameState == 0 {  // by image button
                // game ready
                log.Printf("ready state")
                send.GameState = 0
                for _, dev := range deviceList {
                    log.Printf("here")
                    ws := dev.WS
                    if err = ws.WriteJSON(send); err != nil {
                        log.Println(err)
                    }
                }
            } else if data.GameState == 1 { // by start button
                log.Printf("start state")
                send.GameState = 1
                for _, dev := range deviceList {
                    ws := dev.WS
                    if err = ws.WriteJSON(send); err != nil {
                        log.Println(err)
                    }
                }
                // 5 second later
                time.Sleep(time.Second * 5)
                send.GameState = 2
                if send.Kind == 1 { // random game
                    randomNum := rand.Intn(length+1) % (length+1)
                    log.Printf("drink[%d]", randomNum)
                    for i, dev := range deviceList {
                        if i == randomNum {
                            send.Drink = 1
                        } else {
                            send.Drink = 0
                        }
                        ws := dev.WS
                        if err = ws.WriteJSON(send); err != nil {
                            log.Println(err)
                        }
                    }
                } else if send.Kind == 2 { //lovehot game
                    // love shot game has to be started if players are 2 and more
                    var loveShotA, loveShotB int
                    if length >= 2 {
                        loveShotA = rand.Intn(length) % (length)
                        loveShotB = rand.Intn(length) % (length-1)
                        if loveShotB >= loveShotA {
                            loveShotB += 1
                        }
                        log.Printf("A[%d], B[%d]", loveShotA, loveShotB)
                    } else {
                        log.Printf("Not enough to run loveshot game user is [%d]", length)
                    }

                    for i, dev := range deviceList {
                        if i == loveShotA || i == loveShotB {
                            send.Drink = 1
                        } else {
                            send.Drink = 0
                        }
                        ws := dev.WS
                        if err = ws.WriteJSON(send); err != nil {
                            log.Println(err)
                        }
                    }
                }
            }
        default:
            log.Printf("Not support command {%s}", rcv.Cmd)
        }
    }
}

// Add http response headers to a response to disable caching
func addNoCacheHeaders(handler http.Handler) http.HandlerFunc {

    return func(writer http.ResponseWriter, request *http.Request) {

        writer.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
        writer.Header().Add("Pragma", "no-cache")
        writer.Header().Add("Expires", "0")

        handler.ServeHTTP(writer, request)
    }
}

func main() {
	flag.Parse()
	log.SetFlags(0)

    router := http.NewServeMux()

    router.HandleFunc("/ws", start)
    router.Handle("/", addNoCacheHeaders(http.FileServer(http.Dir(webuiDir))))

    log.Fatal(http.ListenAndServe(*addr, router))
}
