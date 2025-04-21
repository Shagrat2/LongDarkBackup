package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/kardianos/service"
	"github.com/skratchdot/open-golang/open"
	"golang.design/x/mainthread"
)

const (
	cPeriod = 5 * time.Second
)

// Server - HTTP server
type Server struct {
	CStop chan bool
}

var (
	cHost          = "localhost:45192"
	FlagStartOnRun = "yes"
	DataFolder     = ""
	BackupFolder   = ""

	s service.Service
)

// Start -
func (t *Server) Start(s service.Service) error {

	if service.Interactive() {
		log.Println("Running in terminal.")
	} else {
		log.Println("Running under service manager.")
	}

	// Set default
	DataFolder = os.Getenv("LDB_DATA_DIR")
	BackupFolder = os.Getenv("LDB_BACKUP_DIR")

	switch runtime.GOOS {
	case "darwin", "linux":
		fHomeDir, _ := os.UserHomeDir()

		if DataFolder == "" {
			DataFolder = filepath.Join(fHomeDir, ".local/share/Hinterland/TheLongDark/Survival")
		}

		if BackupFolder == "" {
			BackupFolder = filepath.Join(fHomeDir, "Documents/LongDarkBackup")
		}

	case "windows":

		if DataFolder == "" {
			fAppData, _ := os.UserCacheDir()
			DataFolder = filepath.Join(fAppData, "Hinterland/TheLongDark/Survival")
		}

		if BackupFolder == "" {
			fHomeDir, _ := os.UserHomeDir()
			BackupFolder = filepath.Join(fHomeDir, "Documents/LongDarkBackup")
		}

	default:
		return fmt.Errorf("unknown OS")
	}

	if DataFolder == "" {
		return fmt.Errorf("DataFolder is empty")
	}
	if _, err := os.Stat(DataFolder); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("DataFolder is not exist")
	}
	if BackupFolder == "" {
		return fmt.Errorf("BackupFolder is empty")
	}
	if _, err := os.Stat(BackupFolder); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("BackupFolder is not exist")
	}

	t.CStop = make(chan bool)
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

	// app
	err := appSrv()
	if err != nil {
		return err
	}

	if FlagStartOnRun != "" {
		open.Run("http://" + cHost)
	}

	return nil
}

// Stop -
func (t *Server) Stop(s service.Service) error {

	// Stop
	log.Println("Stoping")

	t.CStop <- true

	StopSysTray()

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

	mainthread.Init(func() {

		mainthread.Go(func() {
			StartSysTray()
		})

		err = s.Run()
		if err != nil {
			panic(fmt.Errorf("error run server: %v", err))
		}
	})
}
