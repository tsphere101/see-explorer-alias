package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"see/vscode"
	"see/wt"
)

// File name : pathlist.json
var jsonFileName = filepath.Join(getHomeDir(), "see.json")

// Get user home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

// If json file not exist then create it
func checkJsonFile() {
	// Get the json file name
	filename := jsonFileName

	// Check if json file exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Create the json file
		os.Create(filename)

		// Create the json data
		jsonData := map[string]map[string]string{
			"program":  {},
			"pathlist": {},
		}

		// Write the json data to json file
		WriteToJsonFile(jsonData)
	}

}

func main() {

	// Check if json file exist
	checkJsonFile()

	parseArgs()

}

// Add the program path to jsonfile
func AppendProgramPathToJsonFile(name, path string) {
	// Get the pathlist from json file
	ProgramPathList := GetProgramPathList()

	// if the path list is empty, create a new one
	if ProgramPathList == nil {
		ProgramPathList = make(map[string]string)
	}

	// Add the path to PathList
	ProgramPathList[name] = path

	// Write to json file
	WriteProgramPathToJsonFile(ProgramPathList)
}

func WriteProgramPathToJsonFile(data map[string]string) {
	jsonFile := OpenJsonFile()

	// Write the data to json file
	jsonFile["program"] = data

	// Call WriteToJsonFile function to write the data to json file
	WriteToJsonFile(jsonFile)
}

// Open json file and get data from it and return
func OpenJsonFile() map[string]map[string]string {
	filename := jsonFileName

	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	var result map[string]map[string]string
	json.NewDecoder(jsonFile).Decode(&result)

	return result
}

// Write to json file
func WriteToJsonFile(data map[string]map[string]string) {
	filename := jsonFileName

	// turn the data to json format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	// Write to json file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func GetProgramPath(name string) string {
	pathList := GetProgramPathList()
	return pathList[name]
}

func GetProgramPathList() map[string]string {
	return OpenJsonFile()["program"]
}

func WritePathAliasToJsonFile(data map[string]string) {
	jsonFile := OpenJsonFile()

	// Write the data to json file
	jsonFile["pathlist"] = data

	// Call WriteToJsonFile function to write the data to json file
	WriteToJsonFile(jsonFile)
}

func RenamePath(name string, newName string) {
	pathList := GetPathList()
	pathList[newName] = pathList[name]
	delete(pathList, name)
	WritePathAliasToJsonFile(pathList)
}

// Remove the path from json file
func RemovePath(name string) {
	// Get PathList from json file
	PathList := GetPathList()

	// Remove the path from PathList
	delete(PathList, name)

	// Call WritePathAliasToJsonFile function to write the data to json file
	WritePathAliasToJsonFile(PathList)

}

// Get the pathlist from json file
func GetPathList() map[string]string {
	return OpenJsonFile()["pathlist"]
}

// Add the path to jsonfile
func AppendPathAliasToJsonFile(name, path string) {
	// Get the pathlist from json file
	PathList := GetPathList()

	// if the path list is empty, create a new one
	if PathList == nil {
		PathList = make(map[string]string)
	}

	// Add the path to PathList
	PathList[name] = path

	// Write to json file
	WritePathAliasToJsonFile(PathList)
}

// Path lookup
func lookup(path string) string {
	// Get the pathlist from json file
	PathList := GetPathList()

	// Check if path in PathList
	if PathList[path] != "" {
		return PathList[path]
	}

	return ""
}

func openPath(path string) {
	// Open the path
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", path).Start()
	case "windows":
		exec.Command("explorer", path).Start()
	case "darwin":
		exec.Command("open", path).Start()
	default:
		fmt.Println("Unsupported platform")
	}

}

var version = "see v0.1.0"

func parseArgs() {

	// Debug mode
	debug := false

	if debug {
		fmt.Println("Debug mode")
	}

	// If there is no argument, open explorer
	if len(os.Args) == 1 {
		openPath("")
		return
	}

	args := os.Args[1:]

	// If the first argument is "-h" or "--help", print help message and exit
	if args[0] == "-h" || args[0] == "--help" {

		if debug {
			fmt.Println("Print help message")
		}

		fmt.Println(helpMessage())
		return
	}

	// If the first argument is "-v" or "--version", print version and exit
	if args[0] == "-v" || args[0] == "--version" {
		fmt.Println(version)
		return
	}

	// If the first argument is "-a" or "--add", check if the path is provided
	if args[0] == "-a" || args[0] == "--add" {
		if len(args) == 1 {
			fmt.Println("No path provided")
			return
		}

		// If the path is provided, check if the name is provided
		if len(args) == 2 {
			fmt.Println("No name provided")
			return
		}

		// If the name is provided, add the path to json file
		if len(args) == 3 {
			AppendPathAliasToJsonFile(args[2], args[1])

			// Print success message
			fmt.Println("Path added successfully")
			return
		}
	}

	// If the first argument is "-r" or "--remove", check if the name is provided
	if args[0] == "-r" || args[0] == "--remove" {
		if len(args) == 1 {
			fmt.Println("No name provided")
			return
		}

		// If the name is provided,
		if len(args) == 2 {
			// check if the name is in the pathlist
			if GetPathList()[args[1]] == "" {
				fmt.Printf("Path with name \"%s\" not found\n", args[1])
				return
			}

			// If the name is in the pathlist, remove the path from json file
			RemovePath(args[1])

			// Print success message
			fmt.Println("Path removed successfully")

			return
		}
	}

	// If the first argument is "-re" or "--rename", check if the name is provided
	if args[0] == "-re" || args[0] == "--rename" {
		if len(args) == 1 {
			fmt.Println("No name provided")
			return
		}

		// If the name is provided, check if the new name is provided
		if len(args) == 2 {
			fmt.Println("No new name provided")
			return
		}

		// If the new name is provided, check if the name is in the pathlist
		if len(args) == 3 {
			if GetPathList()[args[1]] == "" {
				fmt.Printf("Name \"%s\" not found\n", args[1])
				return
			}

			// If the name is in the pathlist, rename the path
			RenamePath(args[1], args[2])

			// Print success message
			fmt.Println("Path renamed successfully")
			return
		}
	}

	// If the first argument is "-l" or "--list", print the pathlist
	if args[0] == "-l" || args[0] == "--list" {

		if debug {
			fmt.Println("List mode")
		}

		// Get the pathlist from json file
		PathList := GetPathList()

		// Print the header of the table
		fmt.Println("Name\tPath")

		// Print the pathlist
		for name, path := range PathList {
			fmt.Println(name + "\t" + path)
		}

		return
	}

	// If the first argument is "-wt"
	if args[0] == "-wt" {
		// if there is no path then, run wt in user's home directory
		if len(args) == 1 {
			wt.RunWt(getHomeDir(), GetProgramPath("wt"))
			return
		}

		// Get the path from lookup function
		path := lookup(args[1])

		// Open the path in windows terminal
		wt.RunWt(path, GetProgramPath("wt"))

		return
	}

	// If the first argument is "-code", open the path in vscode
	if args[0] == "-code" {
		// if there is no path then, run vscode in user's home directory
		if len(args) == 1 {
			openPath(getHomeDir())
			return
		}

		// Get the path from lookup function
		path := lookup(args[1])

		// Open the path in vscode
		vscode.RunCode(path, GetProgramPath("code"))

		return
	}

	// If there is one argument
	if len(args) == 2 {
		// Get the path from lookup function
		path := lookup(args[0])
		if path == "" {
			path = args[0]
		}

		// Open the path
		fmt.Println("opening", path)
		openPath(path)
		return
	}

	// If there are more than one argument, print help message
	if len(args) > 1 {
		fmt.Println(helpMessage())
		return
	}

}

func helpMessage() string {
	message := `
Usage: path [options] [path]

Options:
  -h, --help      Show this help message and exit
  -v, --version   Show version and exit
  -l, --list      List all path
  -a, --add       Add a path to pathlist
  -r, --remove    Remove a path from pathlist
  -re, --rename   Rename a path in pathlist
  -wt             Open the path in windows terminal
  -code           Open the path in vscode
`
	return message
}
