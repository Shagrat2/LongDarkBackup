package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	lzf "github.com/zhuyie/golzf"
)

type cacheItem struct {
	DisplayName   string `json:"display_name,omitempty"`
	GameMode      string `json:"game_mode,omitempty"`
	InternalName  string `json:"internal_name,omitempty"`
	Timestamp     string `json:"timestamp,omitempty"`
	XPMode        string `json:"xp_mode,omitempty"`
	Region        string `json:"region,omitempty"`
	Persona       string `json:"persona,omitempty"`
	Condition     string `json:"condition,omitempty"`
	HoursSurvived string `json:"hours_survived,omitempty"`
	WorldExplored string `json:"world_explored,omitempty"`
}

// type Timespan time.Duration

// func (t Timespan) Format(format string) string {
// 	return time.Unix(0, 0).UTC().Add(time.Duration(t)).Format(format)
// }

func loadData(id string) (data cacheItem, img []byte, err error) {

	// Find get filename
	fFileName := ""

	fDir := filepath.Join(BackupFolder, id)

	// Get cached info
	fInfoFile := filepath.Join(fDir, "info.json")
	fJpegFile := filepath.Join(fDir, "info.jpeg")

	_, sErr := os.Stat(fInfoFile)
	if !errors.Is(sErr, os.ErrNotExist) {
		var fInfoFileData []byte
		fInfoFileData, err = os.ReadFile(fInfoFile)
		if err != nil {
			return
		}

		// Get info
		err = json.Unmarshal(fInfoFileData, &data)
		if err != nil {
			return
		}

		// Load image
		img, err = os.ReadFile(fJpegFile)
		return
	}

	err = filepath.WalkDir(fDir, func(path string, d os.DirEntry, err error) error {
		if d == nil || d.IsDir() {
			return nil
		}

		// Skip file
		if !strings.HasPrefix(d.Name(), "sandbox6") {
			return nil
		}

		fFileName = path
		return filepath.SkipAll
	})
	if err != nil {
		return
	}
	if fFileName == "" {
		err = fmt.Errorf("file not found")
		return
	}

	// Read file
	var raw []byte
	raw, err = os.ReadFile(fFileName)
	if err != nil {
		return
	}

	decompressed := make([]byte, len(raw)*10)
	var s int
	s, err = lzf.Decompress(raw, decompressed)
	if err != nil {
		return
	}

	// Save to json
	var jData map[string]interface{}
	err = json.Unmarshal(decompressed[:s], &jData)
	if err != nil {
		return
	}

	//=== Get info
	// m_DisplayName
	data.DisplayName, _ = jData["m_DisplayName"].(string)

	// m_GameMode
	data.GameMode, _ = jData["m_GameMode"].(string)

	// m_InternalName
	data.InternalName, _ = jData["m_InternalName"].(string)

	// m_Timestamp
	data.Timestamp, _ = jData["m_Timestamp"].(string)

	//=== Get ScreenShot
	mDict, ok := jData["m_Dict"].(map[string]interface{})
	if ok {

		//=== Info
		dInfo, ok := mDict["info"].(string)
		if ok {

			var bInfo []byte
			bInfo, err = base64.StdEncoding.DecodeString(dInfo)
			if err != nil {
				return
			}

			// Decompres obj
			dInfo := make([]byte, len(bInfo)*10)
			s, err = lzf.Decompress(bInfo, dInfo)
			if err != nil {
				return
			}

			var iOnfo map[string]interface{}
			err = json.Unmarshal(dInfo[:s], &iOnfo)
			if err != nil {
				return
			}

			// m_XPMode
			data.XPMode, _ = iOnfo["m_XPMode"].(string)

			// m_Region
			data.Region, _ = iOnfo["m_Region"].(string)

			// m_Persona
			data.Persona, _ = iOnfo["m_Persona"].(string)

			// m_Condition
			fCondition, _ := iOnfo["m_Condition"].(float64)
			data.Condition = strconv.FormatFloat(fCondition, 'f', 2, 32)

			// m_HoursSurvived
			fHoursSurvived, _ := iOnfo["m_HoursSurvived"].(float64)
			fTime := time.Unix(0, 0).UTC().Add(time.Duration(fHoursSurvived * float64(time.Hour)))

			if fTime.Day() > 0 {
				data.HoursSurvived = strconv.Itoa(fTime.Day()) + " days "
			}

			data.HoursSurvived += fTime.Format("15:04:05")

			// m_WorldExplored
			fWorldExplored, _ := iOnfo["m_WorldExplored"].(float64)
			data.WorldExplored = strconv.FormatFloat(fWorldExplored, 'f', 2, 32)
		}

		//=== ScreenShot
		dScreen, ok := mDict["screenshot"].(string)
		if ok {
			var bImg []byte
			bImg, err = base64.StdEncoding.DecodeString(dScreen)
			if err != nil {
				return
			}

			// Decompres obj
			dImg := make([]byte, len(bImg)*10)
			s, err = lzf.Decompress(bImg, dImg)
			if err != nil {
				return
			}

			// Unmarshal
			jImg := make(map[string]interface{})
			err = json.Unmarshal(dImg[:s], &jImg)
			if err != nil {
				return
			}

			// m_Encoded
			m_Encoded, ok := jImg["m_Encoded"].(string)
			if ok {
				img, err = base64.StdEncoding.DecodeString(m_Encoded)
				if err != nil {
					return
				}
			}

		}
	}

	// Save info.json
	bInfo, sErr := json.Marshal(data)
	if sErr == nil {
		os.WriteFile(fInfoFile, bInfo, 0644)
	}

	// Save scrennshot
	os.WriteFile(fJpegFile, img, 0644)

	return
}
