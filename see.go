package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		WriteToJsonFile(map[string]string{})
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
		AddThePathToJsonFile(name, path)

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

	// If there is an argument, open the path
	if len(args) == 1 {

		// Get the path from lookup
		path := lookup(args[0])

		// If the path is in PathList, open the path
		if path != "" {
			openPath(path)
		} else {
			// If the path is not in PathList, print error message
			fmt.Println("Path not found.")
		}

	}

}

// Open json file and get data from it and return
func OpenJsonFile() map[string]string {
	filename := jsonFileName

	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]string
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

// Write to json file
func WriteToJsonFile(data map[string]string) {
	filename := jsonFileName

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		fmt.Println(err)
	}

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
	return OpenJsonFile()
}

// Add the path to jsonfile
func AddThePathToJsonFile(name, path string) {
	// Get the pathlist from json file
	PathList := OpenJsonFile()

	// Add the path to PathList
	PathList[name] = path

	// Write to json file
	WriteToJsonFile(PathList)
}

// Path lookup
func lookup(path string) string {
	// Get the pathlist from json file
	PathList := OpenJsonFile()

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
