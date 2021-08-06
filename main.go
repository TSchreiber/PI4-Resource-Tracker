package main

import (
    "fmt"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
)

func main() {
    // Clear the screen
    fmt.Print("\033[H\033[2J")
    for true {
        fmt.Print(getPrintBuffer())
        time.Sleep(1 * time.Second)
    }
}

func getPrintBuffer() string {
    var output string = "\033[H"
    output += getTempuratureString();
    stdout, _ := exec.Command("mpstat", "-P", "0-3").Output()
    // stdout, _ := exec.Command("python", "mpstat.py").Output()
    for _, line := range strings.Split(string(stdout), "\n")[3:7] {
        fields := strings.Fields(line)
        cpuNum, idlePercent := fields[2], fields[12]
        fIdlePercent, _ := strconv.ParseFloat(idlePercent, 64)
        cpuPercent := int(100 - fIdlePercent)
        output += fmt.Sprint(SprintProgressBar("CPU" + cpuNum, 70, cpuPercent))
    }
    return output
}

func getTempuratureString() string {
    //file, err := os.Open("sys.class.thermal.thermal_zone0.temp")
    file, err := os.Open("/sys/class/thermal/thermal_zone0/temp")
    if err != nil {
        panic(err)
    }
    bufTemp := make([]byte,5)
    file.Read(bufTemp)
    file.Close()
    milicelcius,_ := strconv.Atoi(string(bufTemp))
    return fmt.Sprintf("%.1f°C\n", float64(milicelcius) / 1000.0)
}

func SprintProgressBar(title string, width, percent int) string {
    var output string = ""
    fillWidth := width - 2
    numFilledBlocks := fillWidth * percent / 100
    output += fmt.Sprintf("\u250C\u2500 %s %s\u2510\n", title, strings.Repeat("\u2500",fillWidth - len(title) - 3))
    strPercent := strconv.Itoa(percent) + "%"
    centeredPercent := fmt.Sprintf("%s%s%s",
        strings.Repeat(" ", fillWidth / 2),
        strPercent,
        strings.Repeat(" ", fillWidth / 2 - len(strPercent)))
    output += fmt.Sprint("\u2502")
    output += fmt.Sprint("\033[48;2;141;150;83m")
    for i, c := range centeredPercent {
        if i == numFilledBlocks {
            output += fmt.Sprint("\033[49m")
        }
        output += fmt.Sprint(string(c))
    }
    output += fmt.Sprintln("\033[49m\u2502")
    output += fmt.Sprintf("\u2514%s\u2518\n", strings.Repeat("\u2500",fillWidth))
    return output
}
