package main

import (
	"flag"
	"html/template"
	"log"
    "errors"
	"net/http"
    "encoding/json"
	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", ":8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options
var deviceList = make(map[string]Device)

var webuiDir = "www-static"
var contentsDir = "/www-static/"
const contWebInfoFile = "cups.json"

func wsRegister(c *websocket.Conn, message []byte) error {
    var dev Device
    var err error

    rcv := ConnectReq{}

    json.Unmarshal([]byte(message), &rcv)
    log.Printf("cmd : [%s]", rcv.Cmd)
    log.Printf("id : [%s]", rcv.UniqueID)

    if len(rcv.UniqueID) > 0 {
        dev.WS = c
        dev.UniqueID = rcv.UniqueID

        deviceList[rcv.UniqueID] = dev
    } else {
        log.Printf("UniqueID is not existed")
        err = errors.New("UniqueID is not existed")
    }

    return err
}
/*
// Static writeContainerListInfo writing web-ui json data which includes container info
func writeContainerListInfo() error {

    var err error
    if _, err = os.Stat(contentsDir + "/" + contWebInfoFile); os.IsNotExist(err) {
        log.Println("Error contwebinfo.json is not existed")
    } else {
        info, err := containermgt.GetAppInfo()

        if err != nil {
            log.Printf("[%s] Get Container List Info error [%s]", logPrefix, err)
        } else {

            var writeData webUIInfo

            writeData.Title = containerLists

            for _, container := range info.ContainerLists {

                var List containerListsInfo

                if container.Name != containermgt.DockzenAgentName {
                    List.Name = container.Name
                    List.IPAddress = container.IPAddress

                    // web-ui can be shown only one port/hostport
                    for _, port := range container.Ports {
                        List.Port = fmt.Sprint(port.ContainerPort)
                        List.HostPort = fmt.Sprint(port.HostPort)
                    }

                    writeData.Lists = append(writeData.Lists, List)

                }
            }

            // Before writing
            for i, list := range writeData.Lists {
                log.Printf("[%s] C[%d] Name    :[%s]", logPrefix, i, list.Name)
                log.Printf("[%s] C[%d] IP      :[%s]", logPrefix, i, list.IPAddress)
            }

            //file writing into json
            f, err := os.Create(contentsDir + "/" + contWebInfoFile)
            if err != nil {
                log.Printf("[%s] contwebinfo file create error!!!", logPrefix)
                return err
            }

            defer f.Close()

            Lists, err := json.Marshal(writeData)

            _, err = f.Write(Lists)
            if err != nil {
                log.Printf("[%s] writeContainerListInfo file write error", logPrefix)
                return err
            }

        }
    }

    return err
}*/

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
            data := Message{}
            json.Unmarshal([]byte(msg), &data)

            log.Printf("id [%s]", data.UniqueID)
            log.Printf("cap [%d]", data.Capability)

            if dev, ok := deviceList[data.UniqueID]; ok {
                dev.UniqueID = data.UniqueID
                dev.Capability = data.Capability

                send := UserSettingReq{}
                send.Cmd = "usersetting"
                send.UniqueID = dev.UniqueID 
                send.Capability = dev.Capability

                ws := dev.WS
               
                if err = ws.WriteJSON(send); err != nil {
                    log.Println(err)
                }
            } else {
                log.Printf("device is not registered with the UniqueID [%s] ", data.UniqueID)
            }
        default:
            log.Printf("Not support command {%s}", rcv.Cmd)
        }
    }
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/ws")
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

    //http.HandleFunc("/ws", start)
    //http.HandleFunc("/", home)
    //log.Fatal(http.ListenAndServe(*addr, nil))

    router.HandleFunc("/ws", start)
    router.Handle("/", addNoCacheHeaders(http.FileServer(http.Dir(webuiDir))))

    log.Fatal(http.ListenAndServe(*addr, router))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var name = document.getElementById("name");
    var cap = document.getElementById("cap");

    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.innerHTML = message;
        output.appendChild(d);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + name.value);
        print("SEND: " + cap.value);
        
        var obj = new Object();
        obj.cmd = "usersetting";
        obj.id  = name.value;
        obj.capability = Number(cap.value);
        var jsonString = JSON.stringify(obj);

        ws.send(jsonString);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="name" type="text" value="namsu">
<p><input id="cap" type="number" value=10>
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))