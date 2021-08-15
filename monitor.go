package main

import (
    "os"
    "os/exec"
    "strconv"
    "strings"
    "bufio"
    "runtime"
    "text/scanner"
    "time"
)

/**
 * Creates a mpstat process and updates the int array every time the cpu usage
 * changes.
 */
func MonitorCpuUsage(cpuUsage *[4]int) {
    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        cmd = exec.Command("python","mpstat.py","-P","0-3","1")
    } else {
        cmd = exec.Command("mpstat","-P","0-3","1")
    }
    r, _ := cmd.StdoutPipe()
    cmd.Stderr = cmd.Stdout
    scanner := bufio.NewScanner(r)
    go func() {
        for scanner.Scan() {
            line := scanner.Text()
            fields := strings.Fields(line)
            if len(fields) == 13 {
                cpuNum, err := strconv.Atoi(fields[2])
                if err == nil {
                    fIdleCpu, err := strconv.ParseFloat(fields[12], 64)
                    if err == nil {
                        cpuUsage[cpuNum] = int((100 - fIdleCpu) * 100)
                    }
                }
            }
        }
    }()
    cmd.Start()
}

/**
 * returns the temperature reading from the hardwear measurement device. For
 * linux systems, this is located in the file,
 *     "/sys/class/thermal/thermal_zone0/temp"
 * For windows system, wmic needs to be used.
 */
func GetTemperature() float64 {
    var file *os.File
    var err error
    if runtime.GOOS == "windows" {
        file, err = os.Open("sys.class.thermal.thermal_zone0.temp")
    } else {
        file, err = os.Open("/sys/class/thermal/thermal_zone0/temp")
    }
    if err != nil {
        panic(err)
    }
    bufTemp := make([]byte,5)
    file.Read(bufTemp)
    file.Close()
    milicelcius,_ := strconv.Atoi(string(bufTemp))
    return float64(milicelcius) / 1000.0
}

func MonitorTemperature(temp *int) {
    go func() {
        for true {
            *temp = int(GetTemperature())
            time.Sleep(1 * time.Second)
        }
    }()
}

/**
 * starts a goroutine that populates the provided variable with network usage
 * data every second.
 */
func MonitorNetworkUsage(recieved, sent *uint64) {
    go func() {
        var totalRecieved, totalSent uint64
        for true {
            var r, s, x uint64
            var scan scanner.Scanner
            var file *os.File
            if (runtime.GOOS == "windows") {
                file, _ = os.Open("proc.net.dev")
            } else {
                file, _ = os.Open("/proc/net/dev")
            }
            scan.Init(file)
            skip := func(count int) {
                for i:=0; i<count; i++ {
                    scan.Scan();
                }
            }
            skip(27)
            for i:=0; i<3; i++ {
                scan.Scan()
                x,_= strconv.ParseUint(scan.TokenText(), 10, 64)
                r += uint64(x)
                skip(7)
                scan.Scan()
                x,_= strconv.ParseUint(scan.TokenText(), 10, 64)
                s += uint64(x)
                skip(9)
            }
            defer file.Close()
            *recieved = r - totalRecieved
            *sent = s - totalSent
            totalRecieved = r
            totalSent = s
            time.Sleep(time.Second * 1)
        }
    }()
}
