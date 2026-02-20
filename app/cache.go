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
		if !strings.HasPrefix(d.Name(), "sandbox") {
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

// decodeDictEntry decodes a base64+LZF entry from m_Dict into JSON map.
func decodeDictEntry(encoded string) (map[string]interface{}, error) {
	b, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, len(b)*10)
	n, err := lzf.Decompress(b, buf)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(buf[:n], &m)
	return m, err
}

// deserialize parses a JSON string value from a map into a map.
func deserialize(m map[string]interface{}, key string) map[string]interface{} {
	s, _ := m[key].(string)
	if s == "" || s == "{}" || s == "[]" {
		return nil
	}
	var out map[string]interface{}
	if json.Unmarshal([]byte(s), &out) != nil {
		return nil
	}
	return out
}

// hasLocations checks if a deserialized affliction has non-empty m_Locations array.
func hasLocations(m map[string]interface{}) bool {
	if m == nil {
		return false
	}
	locs, _ := m["m_Locations"].([]interface{})
	return len(locs) > 0
}

type GlobalData struct {
	Hunger      string
	Thirst      string
	Fatigue     string
	Freezing    string
	Afflictions []string
}

func loadGlobalData(id string) (gd GlobalData) {
	fDir := filepath.Join(BackupFolder, id)

	// Find sandbox file
	var sandboxFile string
	filepath.WalkDir(fDir, func(path string, d os.DirEntry, err error) error {
		if d != nil && !d.IsDir() && strings.HasPrefix(d.Name(), "sandbox") {
			sandboxFile = path
			return filepath.SkipAll
		}
		return nil
	})
	if sandboxFile == "" {
		return
	}

	raw, err := os.ReadFile(sandboxFile)
	if err != nil {
		return
	}
	buf := make([]byte, len(raw)*10)
	n, err := lzf.Decompress(raw, buf)
	if err != nil {
		return
	}

	var top map[string]interface{}
	if json.Unmarshal(buf[:n], &top) != nil {
		return
	}
	mDict, _ := top["m_Dict"].(map[string]interface{})
	if mDict == nil {
		return
	}

	globalEncoded, _ := mDict["global"].(string)
	if globalEncoded == "" {
		return
	}
	global, err := decodeDictEntry(globalEncoded)
	if err != nil {
		return
	}

	// ── Basic stats ──
	if m := deserialize(global, "m_Hunger_Serialized"); m != nil {
		v, _ := m["m_NumHoursStarving"].(float64)
		if v > 0 {
			gd.Hunger = fmt.Sprintf("%.1f %s", v, T(MsgStarvingHours))
		} else {
			gd.Hunger = T(MsgFed)
		}
	}

	if m := deserialize(global, "m_Thirst_Serialized"); m != nil {
		v, _ := m["m_CurrentThirstProxy"].(float64)
		gd.Thirst = fmt.Sprintf("%.0f%%", v)
	}

	if m := deserialize(global, "m_Fatigue_Serialized"); m != nil {
		v, _ := m["m_CurrentFatigueProxy"].(float64)
		gd.Fatigue = fmt.Sprintf("%.0f%%", v)
	}

	if m := deserialize(global, "m_Freezing_Serialized"); m != nil {
		v, _ := m["m_CurrentFreezingProxy"].(float64)
		gd.Freezing = fmt.Sprintf("%.0f%%", v)
	}

	// ── Afflictions ──

	// Blood loss
	if hasLocations(deserialize(global, "m_BloodLossSerialized")) {
		gd.Afflictions = append(gd.Afflictions, T(MsgBloodLoss))
	}

	// Broken ribs
	if hasLocations(deserialize(global, "m_BrokenRibSerialized")) {
		gd.Afflictions = append(gd.Afflictions, T(MsgBrokenRibs))
	}

	// Burns
	if m := deserialize(global, "m_BurnsSerialized"); m != nil {
		gd.Afflictions = append(gd.Afflictions, T(MsgBurns))
	}

	// Electric burns
	if m := deserialize(global, "m_BurnsElectricSerialized"); m != nil {
		gd.Afflictions = append(gd.Afflictions, T(MsgElectricBurns))
	}

	// Frostbite
	if m := deserialize(global, "m_FrostbiteSerialized"); m != nil {
		locs, _ := m["m_LocationsWithActiveFrostbite"].([]interface{})
		if len(locs) > 0 {
			gd.Afflictions = append(gd.Afflictions, T(MsgFrostbite))
		}
	}

	// Sprained ankle
	if hasLocations(deserialize(global, "m_SprainedAnkleSerialized")) {
		gd.Afflictions = append(gd.Afflictions, T(MsgSprainedAnkle))
	}

	// Sprained wrist
	if hasLocations(deserialize(global, "m_SprainedWristSerialized")) {
		gd.Afflictions = append(gd.Afflictions, T(MsgSprainedWrist))
	}

	// Sprain pain
	if m := deserialize(global, "m_SprainPainSerialized"); m != nil {
		locs, _ := m["m_Locations"].([]interface{})
		if len(locs) > 0 {
			gd.Afflictions = append(gd.Afflictions, T(MsgSprainPain))
		}
	}

	// Food poisoning
	if m := deserialize(global, "m_FoodPoisoningSerialized"); m != nil {
		cause, _ := m["m_CauseLocID"].(string)
		if cause != "" {
			gd.Afflictions = append(gd.Afflictions, T(MsgFoodPoisoning))
		}
	}

	// Dysentery
	if m := deserialize(global, "m_DysenterySerialized"); m != nil {
		if _, ok := m["m_DurationHours"]; ok {
			gd.Afflictions = append(gd.Afflictions, T(MsgDysentery))
		}
	}

	// Infection risk
	if hasLocations(deserialize(global, "m_InfectionRiskSerialized")) {
		gd.Afflictions = append(gd.Afflictions, T(MsgInfectionRisk))
	}

	// Infection
	if hasLocations(deserialize(global, "m_InfectionSerialized")) {
		gd.Afflictions = append(gd.Afflictions, T(MsgInfection))
	}

	// Hypothermia
	if m := deserialize(global, "m_HypothermiaSerialized"); m != nil {
		gd.Afflictions = append(gd.Afflictions, T(MsgHypothermia))
	}

	// Intestinal parasites
	if m := deserialize(global, "m_IntestinalParasitesSerialized"); m != nil {
		if hasP, _ := m["m_HasParasites"].(bool); hasP {
			gd.Afflictions = append(gd.Afflictions, T(MsgIntestinalPar))
		}
	}

	// Scurvy
	if m, _ := global["m_ScurvySerialized"].(map[string]interface{}); m != nil {
		enabled, _ := m["m_Enabled"].(bool)
		repaired, _ := m["m_Repaired"].(bool)
		if enabled && !repaired {
			gd.Afflictions = append(gd.Afflictions, T(MsgScurvy))
		}
	}

	// Headache
	if m := deserialize(global, "m_HeadacheSerialized"); m != nil {
		list, _ := m["m_ListOfParamData"].([]interface{})
		if len(list) > 0 {
			gd.Afflictions = append(gd.Afflictions, T(MsgHeadache))
		}
	}

	// Insomnia
	if m := deserialize(global, "m_InsomniaSerialized"); m != nil {
		list, _ := m["m_ListOfParamData"].([]interface{})
		if len(list) > 0 {
			gd.Afflictions = append(gd.Afflictions, T(MsgInsomnia))
		}
	}

	// Severe laceration
	if m, _ := global["m_SevereLacerationSerialized"].(map[string]interface{}); len(m) > 0 {
		gd.Afflictions = append(gd.Afflictions, T(MsgSevereLac))
	}

	// Suffocating
	if m := deserialize(global, "m_SuffocatingSerialized"); m != nil {
		gd.Afflictions = append(gd.Afflictions, T(MsgSuffocating))
	}

	// Misery mode afflictions (permanent debuffs)
	if m, _ := global["m_MiseryManagerSerialized"].(map[string]interface{}); m != nil {
		curStage, _ := m["m_CurrentStage"].(float64)
		stages, _ := m["m_MiseryStages"].([]interface{})
		for i, st := range stages {
			if i >= int(curStage) {
				break
			}
			stMap, _ := st.(map[string]interface{})
			if stMap == nil {
				continue
			}
			affs, _ := stMap["m_Afflictions"].([]interface{})
			for _, a := range affs {
				name, _ := a.(string)
				if tr := miseryAfflictionName(name); tr != "" {
					gd.Afflictions = append(gd.Afflictions, tr)
				}
			}
		}
	}

	return
}

var miseryKeys = map[string]string{
	"WeakConstitution": MsgWeakConst,
	"WeakJoints":       MsgWeakJoints,
	"SourStomach":      MsgSourStomach,
	"PoorCirculation":  MsgPoorCirc,
	"UnsettledSleep":   MsgUnsettledSleep,
	"BrokenBody":       MsgBrokenBody,
}

func miseryAfflictionName(key string) string {
	if msgKey, ok := miseryKeys[key]; ok {
		return T(msgKey)
	}
	if key != "" {
		return key
	}
	return ""
}
