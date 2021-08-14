package main

import (
    "fmt"
    "time"
)

func main() {
    var recieved, sent uint64
    MonitorNetworkUsage(&recieved, &sent)
    for i:=0; i<10; i++ {
        time.Sleep(time.Second)
        fmt.Println(recieved, sent)
    }
}
