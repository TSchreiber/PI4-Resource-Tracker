package main

import (
	"fmt"
	"log"
	"net/http"
    "time"
	"encoding/json"
	"os"
	socketio "github.com/googollee/go-socket.io"
)

const PORT string = "8080"
var server string

func main() {
	config := ParseConfig()

    var cpuUsage [4]int
    MonitorCpuUsage(&cpuUsage)

    sockets := make(map[string]socketio.Conn)

    go func() {
        for true {
            for _,s := range sockets {
                s.Emit("data", fmt.Sprintf("%c%c%c%c%c",
                    int(cpuUsage[0] / 100),
                    int(cpuUsage[1] / 100),
                    int(cpuUsage[2] / 100),
                    int(cpuUsage[3] / 100),
                    int(GetTempurature())))
            }
            time.Sleep(1 * time.Second)
        }
    }()

	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
        sockets[s.ID()] = s
		return nil
	})

    server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
        delete(sockets, s.ID())
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./static")))
	log.Println("Serving at "+config.Server.GetURL()+"...")
	log.Fatal(http.ListenAndServe(config.Server.GetURL(), nil))
}

type Config struct {
	Server ServerInfo
}

type ServerInfo struct {
	Host, Port string
}
func (server ServerInfo) GetURL() string {
	return server.Host + ":" + server.Port
}

func ParseConfig() Config {
	b,_ := os.ReadFile("config.json")
	con := Config{}
	json.Unmarshal([]byte(b), &con)
	return con
}
