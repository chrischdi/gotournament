package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/chrischdi/gotournament/pkg/data"
	"github.com/chrischdi/gotournament/pkg/frontend"
)

func main() {
	err := data.RestoreFile()
	if err != nil {
		data.Reset()
	}

	port := ":30081"

	fmt.Printf("Starting application at \"http://127.0.0.1%s\"\n", port)
	openBrowser(fmt.Sprintf("http://127.0.0.1%s/setup", port))
	fmt.Printf("You can stop the application by pressing <STRG>+<C>\n")
	frontend.Run(port)
}

func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}
