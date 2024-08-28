package main

import (
	"fmt"
	"log"
	"time"

	"github.com/getlantern/systray"
	"github.com/kardianos/service"
	"github.com/skratchdot/open-golang/open"

	"git2.jad.ru/LongDarkBackup/icon"
)

const (
	cPeriod  = 5 * time.Second
	cPortNum = "45192"
)

// Server - HTTP server
type Server struct {
	CStop chan bool
}

var s service.Service

func onReady() {
	systray.SetIcon(icon.Data)
	//systray.SetTitle("Awesome App")
	systray.SetTooltip("LD backup")

	go func() {
		mUrl := systray.AddMenuItem("Open", "Open settings")

		mQuit := systray.AddMenuItem("Quit", "Close app")

		for {
			select {
			case <-mUrl.ClickedCh:
				open.Run("http://localhost:" + cPortNum)
			case <-mQuit.ClickedCh:
				s.Stop()
				systray.Quit()
				return
			}
		}

	}()
}

// Start -
func (t *Server) Start(s service.Service) error {

	if service.Interactive() {
		log.Println("Running in terminal.")
	} else {
		log.Println("Running under service manager.")
	}

	go systray.Run(onReady, nil)

	go func() {

		fTimer := time.NewTicker(cPeriod)
		for {
			select {
			case <-t.CStop:
				return
			case <-fTimer.C:
				DoScan()
				DoGarbage()
			}
		}
	}()

	return nil
}

// Stop -
func (t *Server) Stop(s service.Service) error {

	// Stop
	log.Println("Stoping")

	t.CStop <- true

	return nil
}

func main() {

	// Set log format
	log.SetFlags(log.Lmicroseconds)

	// Service
	svcConfig := &service.Config{
		Name:        "HLBackup",
		DisplayName: "LongDark backup",
		Description: "LongDark backup service",
	}

	//== Start service
	var err error
	srv := &Server{}
	s, err = service.New(srv, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		panic(fmt.Errorf("error run server: %v", err))
	}

}
