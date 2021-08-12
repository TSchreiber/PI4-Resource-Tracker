package main

import (
    "os"
    "os/exec"
    "strconv"
    "strings"
    "bufio"
    "runtime"
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
 * returns the tempurature reading from the hardwear measurement device. For
 * linux systems, this is located in the file,
 *     "/sys/class/thermal/thermal_zone0/temp"
 * For windows system, wmic needs to be used.
 */
func GetTempurature() float64 {
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
