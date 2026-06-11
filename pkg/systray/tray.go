package systray

import (
	"log"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
)

var (
	serverURL = "http://localhost:3333"
	quitChan  = make(chan struct{})
)

func OnReady() {
	systray.SetTitle("PDF Bookmarker")

	libItem := systray.AddMenuItem("Library", "Open library in browser")
	//	lastItem := systray.AddMenuItem("Resume Last Book", "Open most recent book")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Stop server and exit")

	go func() {
		for {
			select {
			case <-libItem.ClickedCh:
				openBrowser(serverURL + "/")
			case <-quitItem.ClickedCh:
				systray.Quit()
				close(quitChan)
				return

			}
		}
	}()
}

func OnExit() {}
func GetQuitChan() <-chan struct{} {
	return quitChan

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

	}
	if err != nil {
		log.Println("Failed to open browser:", err)
	}
}
func Start() {
	go systray.Run(OnReady, OnExit)
}
