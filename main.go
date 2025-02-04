package main

import (
	"Cosmic/executor"
	"Cosmic/helper"
	"Cosmic/parser"
	"fmt"
	"github.com/pterm/pterm"
	"path/filepath"
	"strconv"
)

var Version = "1.3"

func main() {
	helper.SetTitle(fmt.Sprintf("Cosmic | %s", Version))
	selectTargetFile()
	displayCategories()
	main()
}

func displayCategories() {
	helper.Clear()
	helper.ASCII()
	helper.PrintLine("!", "Select a category", true)
	helper.PrintLine("1", "C#", true)
	helper.PrintLine("2", "File Information", true)
	helper.PrintLine("X", "Clear Downloads", true)
	helper.PrintLine(">", "", false)
	var category string
	fmt.Scanln(&category)

	switch category {
	case "1":
		displayTools("C%23")
	case "2":
		displayTools("Info")
	case "X":
		ClearDownloadFolder()
	default:
		displayCategories()
	}
}

func ClearDownloadFolder() {
	helper.Clear()
	helper.ASCII()
	helper.PrintLine(">", "Clearing Cosmic Download Folder...", true)
	executor.ClearDownloadFolder()
	helper.PrintLine("~", "Press enter to go back...", false)
	fmt.Scanln()
	displayCategories()
}

func displayTools(toolsList string) {
	helper.Clear()
	helper.ASCII()
	tools, _ := parser.FetchTools(toolsList)

	tableData := pterm.TableData{
		{"Option", "Name", "Description"},
	}

	toolMap := make(map[int]parser.Application)
	option := 1

	for _, tool := range tools {
		tableData = append(tableData, []string{strconv.Itoa(option), tool.Metadata.Name, tool.Metadata.Description})
		toolMap[option] = tool
		option++
	}

	renderedTable, err := pterm.DefaultTable.WithHasHeader().WithData(tableData).Srender()
	if err != nil {
		fmt.Println("Error rendering table:", err)
		return
	}

	pterm.DefaultCenter.WithCenterEachLineSeparately().Println(renderedTable)

	helper.PrintLine(">", "", false)
	var toolNumber int
	fmt.Scanln(&toolNumber)

	if selectedTool, ok := toolMap[toolNumber]; ok {
		displayData(selectedTool)
	} else {
		displayTools(toolsList)
	}
}

func selectTargetFile() {
	helper.Clear()
	helper.ASCII()
	helper.PrintLine("~", "Enter the target file path...", true)
	helper.PrintLine(">", "", false)
	fmt.Scanln(&executor.TargetFile)
}

func displayData(tool parser.Application) {
	helper.Clear()
	helper.ASCII()
	helper.PrintLine("Name", tool.Metadata.Name, true)
	helper.PrintLine("Description", tool.Metadata.Description, true)
	helper.PrintLine("Author", tool.Metadata.Author, true)
	helper.PrintLine("Tags", fmt.Sprintf("%v", tool.Metadata.Tags), true)
	helper.PrintLine("Supported OS", fmt.Sprintf("%v", tool.Compatibility.OS), true)
	helper.PrintLine("Supported Architectures", fmt.Sprintf("%v", tool.Compatibility.Architectures), true)
	helper.PrintLine("Requires Admin", fmt.Sprintf("%v", tool.Execution.RunAsAdmin), true)
	helper.PrintLine("~", "Press enter to continue...", false)
	fmt.Scanln()

	handleTool(tool)
	helper.PrintLine("~", "Press enter to go back...", false)
	fmt.Scanln()
	main()
}

func handleTool(tool parser.Application) {
	helper.Clear()
	helper.ASCII()
	path, isFolder := executor.DownloadFile(tool.Download.IsCompressed, tool.Download.URL)
	helper.PrintLine(">", "Executing...", true)
	if isFolder {
		filePath := filepath.Join(path, tool.Execution.FileLocation)
		err := executor.ExecuteFile(filePath, tool.Execution.Arguments, tool)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
	err := executor.ExecuteFile(path, tool.Execution.Arguments, tool)
	if err != nil {
		fmt.Println(err)
	}
}
