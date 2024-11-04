package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gen2brain/go-unarr"
)

var fileName = "MinGW"
var filePath = "C:\\cc"

func getUrl() string {
	url := "https://api.github.com/repos/niXman/mingw-builds-binaries/releases/latest"
	res, err := http.Get(url)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
	}

	assets := response["assets"].([]interface{})
	json.Unmarshal(body, &assets)

	for _, i := range assets {
		asset := i.(map[string]interface{})
		if strings.Contains(asset["name"].(string), "win32-seh-msvcrt") {
			url = asset["browser_download_url"].(string)
			break
		}
	}
	return url
}

func downloadFile() {
	url := getUrl()

	file, err := os.Create("./" + fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func extractFile() {
	file, err := unarr.NewArchive(fileName)
	if err != nil {
		fmt.Println(err)
	}

	file.Extract(filePath)
	defer file.Close()
}

func addPath() {
	path := os.Getenv("Path")

	_, err := os.Stat(filePath + "\\mingw64\\bin")
	if os.IsNotExist(err) {
		os.Mkdir(filepath.Join(filePath), 0755)
	}

	if !strings.Contains(path, filePath) {
		exec.Command("powershell", `[Environment]::SetEnvironmentVariable("Path",[Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User) + ";`+filePath+"\\mingw64\\bin"+`", [EnvironmentVariableTarget]::User)`).Run()
	}
}

func renameMake() {
	dest, err := os.Open(filePath + "\\mingw64\\bin\\mingw32-make.exe")
	if err != nil {
		fmt.Println(err)
	}
	defer dest.Close()

	src, err := os.Create(filePath + "\\mingw64\\bin\\make.exe")
	if err != nil {
		fmt.Println(err)
	}
	defer src.Close()

	_, err = io.Copy(src, dest)
	if err != nil {
		fmt.Printf("Error due to: %s\n", err)
	}
}
func setup() {
	addPath()
	downloadFile()
	extractFile()
	renameMake()
}

func main() {
	setup()
}
