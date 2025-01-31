package main

import (
	"Cosmic/executor"
	"Cosmic/helper"
	"Cosmic/parser"
	"fmt"
	"path/filepath"
)

func main() {
	selectTargetFile()
	displayCategories()
	main()
}

func displayCategories() {
	helper.Clear()
	helper.ASCII()
	helper.PrintLine("1", "C#", true)
	helper.PrintLine(">", "", false)
	var category string
	fmt.Scanln(&category)

	switch category {
	case "1":
		displayTools("C%23")
	default:
		displayCategories()
	}
}

func displayTools(toolsList string) {
	helper.Clear()
	helper.ASCII()
	tools, _ := parser.FetchTools(toolsList)

	for _, tool := range tools {
		helper.PrintLine(tool.Metadata.Name, fmt.Sprintf("%s: %s", tool.Metadata.Name, tool.Metadata.Description), true)
	}

	helper.PrintLine(">", "", false)
	var tool string
	fmt.Scanln(&tool)

	if _, ok := tools[tool]; !ok {
		displayTools(toolsList)
		return
	}

	displayData(tools[tool])
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
