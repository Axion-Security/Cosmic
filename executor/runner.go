package executor

import (
	"Cosmic/helper"
	"Cosmic/parser"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func checkCompatibility(tool parser.Application) bool {
	var OS = tool.Compatibility.OS
	var Arch = tool.Compatibility.Architectures

	// windows = Windows
	// darwin = macOS
	// linux = Linux
	var UserOS = runtime.GOOS

	// amd64 = 64-bit
	// 386 = 32-bit
	var UserArch = os.Getenv("PROCESSOR_ARCHITECTURE")

	for _, o := range OS {
		if strings.ToLower(o) == strings.ToLower(UserOS) {
			for _, arch := range Arch {
				if strings.ToLower(arch) == strings.ToLower(UserArch) {
					return true
				}
			}
		}
	}

	helper.PrintLine("!", "Incompatible OS or architecture", true)
	helper.PrintLine("!", "Supported OS: "+fmt.Sprintf("%v", tool.Compatibility.OS), true)
	helper.PrintLine("!", "Supported Architectures: "+fmt.Sprintf("%v", tool.Compatibility.Architectures), true)
	helper.PrintLine("!", "Your OS: "+UserOS, true)
	helper.PrintLine("!", "Your Architecture: "+UserArch, true)

	return false
}

func ExecuteFile(filePath string, args []string, tool parser.Application) error {
	if !checkCompatibility(tool) {
		return fmt.Errorf("compatibility check failed")
	}

	if !CheckRequirements(tool.Execution.Requirements) {
		return fmt.Errorf("requirements check failed")
	}

	var cmd *exec.Cmd
	if tool.Execution.RunAsAdmin {
		powershellCommand := fmt.Sprintf("Start-Process -FilePath '%s' -ArgumentList '%s' -Verb RunAs", filePath, strings.Join(ReplaceArgs(args), "', '"))
		cmd = exec.Command("powershell", "-Command", powershellCommand)
	} else {
		cmd = exec.Command(filePath, ReplaceArgs(args)...)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
