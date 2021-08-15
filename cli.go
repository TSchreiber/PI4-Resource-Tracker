package main

import (
    "fmt"
    "strconv"
    "strings"
    "time"
)

var cpuUsage [4]int

func main() {
    // Clear the screen
    fmt.Print("\033[H\033[2J")
    MonitorCpuUsage(&cpuUsage)
    for true {
        fmt.Print(getPrintBuffer())
        time.Sleep(1 * time.Second)
    }
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
        GetThermometer(int(Gettemperature())), "\n")
    output = ""
    for i:=0; i<12; i++ {
        output += cpuLines[i] + tempLines[i] + "\n"
    }
    return output
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
