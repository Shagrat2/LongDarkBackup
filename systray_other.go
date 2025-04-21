//go:build !wasm
// +build !wasm

package main

import (
	"log"

	"git2.jad.ru/LongDarkBackup/icon"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
)

var chQuit chan bool

func onReady() {
	systray.SetIcon(icon.Data)
	//- np nide: systray.SetTitle("LD backup")
	systray.SetTooltip("LD backup")

	chQuit = make(chan bool)
	go func() {
		mUrl := systray.AddMenuItem("Open", "Open settings")
		mQuit := systray.AddMenuItem("Quit", "Close app")

		for {
			select {
			case <-mUrl.ClickedCh:
				open.Run("http://" + cHost)

			case <-mQuit.ClickedCh:
				s.Stop()
				systray.Quit()
				return

			case <-chQuit:
				log.Println("@2")
				systray.Quit()

				return
			}
		}

	}()
}

func StartSysTray() {
	systray.Run(onReady, nil)
}

func StopSysTray() {
	log.Println("@1")
	chQuit <- true
}
