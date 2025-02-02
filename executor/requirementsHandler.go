package executor

import (
	"os/exec"
	"strconv"
	"strings"
)

func CheckRequirements(requirements []string) bool {
	for _, req := range requirements {
		if strings.HasPrefix(req, "DOTNET(") && strings.HasSuffix(req, ")") {
			version := strings.TrimPrefix(req, "DOTNET(")
			version = strings.TrimSuffix(version, ")")
			if !CheckDotNet(version) {
				return false
			}
		} else if strings.HasPrefix(req, "PYTHON(") && strings.HasSuffix(req, ")") {
			version := strings.TrimPrefix(req, "PYTHON(")
			version = strings.TrimSuffix(version, ")")
			if !CheckPython(version) {
				return false
			}
		}
	}
	return true
}

func CheckDotNet(requiredVersion string) bool {
	cmd := exec.Command("dotnet", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	installedVersion := strings.TrimSpace(string(output))
	return isVersionCompatible(installedVersion, requiredVersion)
}

func CheckPython(requiredVersion string) bool {
	var output []byte
	var err error

	cmd := exec.Command("python3", "--version")
	output, err = cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command("python", "--version")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return false
		}
	}

	versionLine := strings.TrimSpace(string(output))
	parts := strings.Split(versionLine, " ")
	if len(parts) < 2 {
		return false
	}
	installedVersion := parts[1]
	return isVersionCompatible(installedVersion, requiredVersion)
}

func isVersionCompatible(installed, required string) bool {
	installedParts := strings.Split(installed, ".")
	requiredParts := strings.Split(required, ".")

	for i := 0; i < len(requiredParts); i++ {
		if i >= len(installedParts) {
			return false
		}

		installedPart, err1 := strconv.Atoi(installedParts[i])
		requiredPart, err2 := strconv.Atoi(requiredParts[i])

		if err1 != nil || err2 != nil {
			if installedParts[i] < requiredParts[i] {
				return false
			} else if installedParts[i] > requiredParts[i] {
				return true
			}
			continue
		}

		if installedPart < requiredPart {
			return false
		} else if installedPart > requiredPart {
			return true
		}
	}

	return true
}
