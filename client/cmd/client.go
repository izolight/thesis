package main

import (
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	log "github.com/sirupsen/logrus"
	"os"
)

func setUpEnv() {
	err := os.Setenv("FYNE_SCALE", "1.0")
	if err != nil {
		log.Info(fmt.Sprintf("Unable to set environment variable: %s", err))
	}

}

func setUpMainWindow() {

	a := app.New()

	w := a.NewWindow("Offline Signature Verification")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			a.Quit()
		})))

	w.ShowAndRun()
}

func main() {
	setUpEnv()
	setUpMainWindow()
}
