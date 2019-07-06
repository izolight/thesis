package client

import (
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/sirupsen/logrus"
	"os"
)

func SetUpEnv() {
	err := os.Setenv("FYNE_SCALE", "1.0")
	if err != nil {
		logrus.Info(fmt.Sprintf("Unable to set environment variable: %s", err))
	}

}

func SetUpMainWindow() {

	a := app.New()

	w := a.NewWindow("Offline Signature Verification")
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			a.Quit()
		})))

	w.ShowAndRun()
}
