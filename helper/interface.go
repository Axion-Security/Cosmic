package helper

import (
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"os/exec"
	"strconv"
)

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func SetTitle(title string) {
	cmd := exec.Command("cmd", "/c", "title", title)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ASCII() {
	fmt.Println()
	pterm.DefaultCenter.Println(`  *     ██████╗ ██████╗ ███████╗███╗   ███╗██╗ ██████╗    *
 *     ██╔════╝██╔═══██╗██╔════╝████╗ ████║██║██╔════╝      *
    *  ██║     ██║   ██║███████╗██╔████╔██║██║██║      *    ✧
  ✧    ██║     ██║   ██║╚════██║██║╚██╔╝██║██║██║         *
 *     ╚██████╗╚██████╔╝███████║██║ ╚═╝ ██║██║╚██████╗  ✧   *
   *    ╚═════╝ ╚═════╝ ╚══════╝╚═╝     ╚═╝╚═╝ ╚═════╝     *
  ✧    With great power comes great reverse engineering  ✧`)
}

func PrintLine(option, value string, newLine bool) {
	options := map[string]string{
		"!": "\033[31m", // Red
		"?": "\033[36m", // Cyan
		"~": "\033[33m", // Yellow
		">": "\033[35m", // Magenta
	}

	defaultColor := "\033[36m" // Cyan
	greenColor := "\033[32m"   // Green
	grayColor := "\033[90m"    // Dark Gray
	resetColor := "\033[0m"    // Reset color

	color := defaultColor
	if _, err := strconv.Atoi(option); err == nil {
		color = greenColor
	} else if c, ok := options[option]; ok {
		color = c
	}

	fmt.Print(grayColor, "[", resetColor)
	fmt.Print(color, option, resetColor)
	fmt.Print(grayColor, "] ", resetColor)
	fmt.Print("\033[37m", value, resetColor)

	if newLine {
		fmt.Println()
	}
}
