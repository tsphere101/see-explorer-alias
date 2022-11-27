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

// Get user home directory
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

// File name : pathlist.json
var jsonFileName = filepath.Join(getHomeDir(), "pathlist.json")

func main() {

	// if json file not exist then create it
	if _, err := os.Stat(jsonFileName); os.IsNotExist(err) {
		jsonFile := make(map[string]map[string]string)
		jsonFile["program"] = make(map[string]string)
		jsonFile["pathlist"] = make(map[string]string)
		WriteToJsonFile(jsonFile)
	}

	// if no args then print help and exit
	if len(os.Args) == 1 {
		fmt.Println("type see -h for help")
		return
	}

	// get args from command line
	args := os.Args[1:]

	// Get the arguments from the command line

	// If the first argument is "-h" or "--help" print the help message
	if args[0] == "-h" || args[0] == "--help" {
		message := `Usage: see [OPTION]... [PATH]...
Open the specified path in the Windows Explorer.

Options:
-h, --help     display this help and exit
-a, --add      add the path to the pathlist
-l, --list     list all the paths in the pathlist
-r, --remove   remove the path from the pathlist
-re, --rename  rename the path in the pathlist
-wt            open the path in windows terminal
-code  	  	   open the path in vscode
--where 	  location of the pathlist file

`
		fmt.Println(message)

		return
	}

	// If the first argument is "--where" then print the path of the json file
	if args[0] == "--where" {
		fmt.Println(jsonFileName)
		return
	}

	// If the first argument is "-l" or "--list" then list all the paths
	if args[0] == "-l" || args[0] == "--list" {
		// Get the pathlist from json file
		pathList := GetPathList()

		// Print header for the list
		fmt.Println("Name\t\tPath")
		fmt.Println("----\t\t----")

		// Print the list
		for name, path := range pathList {
			fmt.Printf("%s\t\t%s \n", name, path)
		}

		return
	}

	// If the first argument is "-a" or "--add", add the path to the json file
	if args[0] == "-a" || args[0] == "--add" {

		// Get the name and path
		name := args[1]
		path := args[2]

		// Store to json file
		AppendPathAliasToJsonFile(name, path)

		// Print success message
		fmt.Println("Path added successfully.")

		return
	}

	// If the first argument is "-r" or "--remove", remove the path from the json file
	if args[0] == "-r" || args[0] == "--remove" {
		// Get the name
		name := args[1]

		// Remove the path from json file
		RemovePath(name)

		// Print success message
		fmt.Println("Path removed successfully.")

		return
	}

	// If the first argument is "-re" or "--rename", rename the path in the json file
	if args[0] == "-re" || args[0] == "--rename" {
		// Get the name and new name
		name := args[1]
		newName := args[2]

		// Rename the path in json file
		RenamePath(name, newName)

		// Print success message
		fmt.Println("Path renamed successfully.")

		return
	}

	// If the first argument is "-wt", open the path in windows terminal
	if args[0] == "-wt" {

		// if there is no path then, run wt in user's home directory
		if len(args) == 1 {
			wt.RunWt(getHomeDir(), GetProgramPath("wt"))
			return
		}

		// Get the path from lookup function
		path := lookup(args[1])

		// If path not found then print error message and exit
		if path == "" {
			fmt.Println("Path not found")
			return
		}

		// Open the path in windows terminal
		wt.RunWt(path, GetProgramPath("wt"))

		return

	}

	// If the first argument is "-code", open the path in vscode
	if args[0] == "-code" {

		// Get the path from lookup function
		path := lookup(args[1])

		// Run vscode with the path
		vscode.RunCode(path, GetProgramPath("code"))

		return
	}

	// Open the path in explorer
	if len(args) == 1 {

		// Get the path from lookup
		path := lookup(args[0])

		// If the path is in PathList, open the path
		if path != "" {
			openPath(path)
		} else {
			// open the path in explorer
			fmt.Printf("opening %v ", args[0])
			openPath(args[0])
		}

	}

}

func RenamePath(name string, newName string) {
	pathList := GetPathList()
	pathList[newName] = pathList[name]
	delete(pathList, name)
	WritePathAliasToJsonFile(pathList)
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

	jsonData, err := json.Marshal(data)
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

// Remove the path from json file
func RemovePath(name string) {
	// Get PathList from json file
	PathList := OpenJsonFile()

	// Remove the path from PathList
	delete(PathList, name)

	// Write to json file
	WriteToJsonFile(PathList)

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
