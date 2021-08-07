package main

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
    "bufio"
)

var cpuUsage [4]int

func main() {
    // Clear the screen
    fmt.Print("\033[H\033[2J")
    MonitorCpuUsage()
    for true {
        fmt.Print(getPrintBuffer())
        time.Sleep(1 * time.Second)
    }
}

/**
 * Creates a mpstat process and updates the int array every time the cpu usage
 * changes.
 */
func MonitorCpuUsage() [4]int {
    // cmd := exec.Command("python","mpstat.py","-P","0-3","1")
    cmd := exec.Command("mpstat","-P","0-3","1")
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
    return cpuUsage
}

/**
 * Calculates and compiles all the information that needs to go onto the console
 * into a single string so it can be printed in one go. This limits tearing
 * during redraws.
 */
func getPrintBuffer() string {
    var output string = "\033[H"
    for cpuNum, cpuPermyriad := range cpuUsage {
        output += GetProgressBar("CPU" + strconv.Itoa(cpuNum), 70, cpuPermyriad / 100)
    }
    cpuLines := strings.Split(output, "\n")
    tempLines := strings.Split(
        GetThermometer(int(GetTempurature())), "\n")
    output = ""
    for i:=0; i<12; i++ {
        output += cpuLines[i] + tempLines[i] + "\n"
    }
    return output
}

/**
 * returns the tempurature reading from the hardwear measurement device. For
 * linux systems, this is located in the file,
 *     "/sys/class/thermal/thermal_zone0/temp"
 * For windows system, wmic needs to be used.
 */
func GetTempurature() float64 {
    // file, err := os.Open("sys.class.thermal.thermal_zone0.temp")
    file, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
    if err != nil {
        panic(err)
    }
    bufTemp := make([]byte,5)
    file.Read(bufTemp)
    file.Close()
    milicelcius,_ := strconv.Atoi(string(bufTemp))
    return float64(milicelcius) / 1000.0
}

/**
 * Creates a string representation for a thermometer in the range of 40c to 60c.
 * Values outside the range will appear roughly correctly if they are exactly 2
 * digits, but anthing greater than 99 or less than 10 will display as -1.
 */
func GetThermometer(temp int) string {
    const height int = 10
    const minTemp int = 40
    const maxTemp int = 60
    const strFilled string = "\033[41m    \033[49m"

    if temp >= 100 || temp < 10 {
        temp = -1
    }
    var filledBarCount = height * (temp - minTemp) / (maxTemp - minTemp)
    var output string = ""
    output += "\u250c\u2500\u2500\u2500\u2500\u2510\n"
    for i:=height; i>0; i-- {
        output += "\u2502"
        if i <= filledBarCount {
            output += "\033[41m"
        }
        if i == height / 2 {
            output += strconv.Itoa(temp) + "°C"
        } else {
            output += "    "
        }
        output += "\033[49m\u2502\n"
    }
    output +=  "\u2514\u2500\u2500\u2500\u2500\u2518\n"
    return output
}

func GetProgressBar(title string, width, percent int) string {
    var output string = ""
    fillWidth := width - 2
    numFilledBlocks := fillWidth * percent / 100
    output += fmt.Sprintf("\u250C\u2500 %s %s\u2510\n", title, strings.Repeat("\u2500",fillWidth - len(title) - 3))
    strPercent := strconv.Itoa(percent) + "%"
    centeredPercent := fmt.Sprintf("%s%s%s",
        strings.Repeat(" ", fillWidth / 2),
        strPercent,
        strings.Repeat(" ", fillWidth / 2 - len(strPercent)))
    output += "\u2502"
    output += "\033[48;2;141;150;83m"
    for i, c := range centeredPercent {
        if i == numFilledBlocks {
            output += "\033[49m"
        }
        output += string(c)
    }
    output += "\033[49m\u2502\n"
    output += fmt.Sprintf("\u2514%s\u2518\n", strings.Repeat("\u2500",fillWidth))
    return output
}
