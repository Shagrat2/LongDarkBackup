package main

import "strings"

type Lang string

const (
	LangRU Lang = "ru"
	LangEN Lang = "en"
)

var currentLang = LangRU

func SetLang(l Lang) { currentLang = l }
func GetLang() Lang  { return currentLang }

// Message keys
const (
	MsgAppTitle       = "app_title"
	MsgNoSaves        = "no_saves"
	MsgBack           = "back"
	MsgSurvival       = "survival"
	MsgCondition      = "condition"
	MsgHunger         = "hunger"
	MsgThirst         = "thirst"
	MsgFatigue        = "fatigue"
	MsgFreezing       = "freezing"
	MsgWorldExplored  = "world_explored"
	MsgCharacter      = "character"
	MsgBackupLabel    = "backup_label"
	MsgSaveLabel      = "save_label"
	MsgRestore        = "restore"
	MsgRestoreConfirm = "restore_confirm"
	MsgBackupFrom     = "backup_from"
	MsgOverwriteWarn  = "overwrite_warn"
	MsgCancel         = "cancel"
	MsgFed            = "fed"
	MsgStarvingHours  = "starving_hours"
	MsgBloodLoss      = "blood_loss"
	MsgBrokenRibs     = "broken_ribs"
	MsgBurns          = "burns"
	MsgElectricBurns  = "electric_burns"
	MsgFrostbite      = "frostbite"
	MsgSprainedAnkle  = "sprained_ankle"
	MsgSprainedWrist  = "sprained_wrist"
	MsgSprainPain     = "sprain_pain"
	MsgFoodPoisoning  = "food_poisoning"
	MsgDysentery      = "dysentery"
	MsgInfectionRisk  = "infection_risk"
	MsgInfection      = "infection"
	MsgHypothermia    = "hypothermia"
	MsgIntestinalPar  = "intestinal_parasites"
	MsgScurvy         = "scurvy"
	MsgHeadache       = "headache"
	MsgInsomnia       = "insomnia"
	MsgSevereLac      = "severe_laceration"
	MsgSuffocating    = "suffocating"
	MsgWeakConst      = "weak_constitution"
	MsgWeakJoints     = "weak_joints"
	MsgSourStomach    = "sour_stomach"
	MsgPoorCirc       = "poor_circulation"
	MsgUnsettledSleep = "unsettled_sleep"
	MsgBrokenBody     = "broken_body"
	MsgSettings       = "settings"
	MsgBackupDir      = "backup_dir"
	MsgBackupDirHint  = "backup_dir_hint"
	MsgSave           = "save"
	MsgSettingsSaved  = "settings_saved"
	MsgPathRequired   = "path_required"
	MsgPathRelative   = "path_relative"
	MsgPathRoot       = "path_root"
	MsgDirCreateErr   = "dir_create_error"
	MsgDirNotWritable = "dir_not_writable"
	MsgConfigReadErr  = "config_read_error"
	MsgConfigSaveErr  = "config_save_error"
)

var translations = map[Lang]map[string]string{
	LangRU: {
		MsgAppTitle:       "LONGDARK BACKUP",
		MsgNoSaves:        "Нет сохранений",
		MsgBack:           "НАЗАД",
		MsgSurvival:       "Выживание",
		MsgCondition:      "Состояние",
		MsgHunger:         "Голод",
		MsgThirst:         "Жажда",
		MsgFatigue:        "Усталость",
		MsgFreezing:       "Замерзание",
		MsgWorldExplored:  "Мир исследован",
		MsgCharacter:      "Персонаж",
		MsgBackupLabel:    "Бэкап",
		MsgSaveLabel:      "Сохранение",
		MsgRestore:        "ВОССТАНОВИТЬ",
		MsgRestoreConfirm: "ВОССТАНОВИТЬ СОХРАНЕНИЕ?",
		MsgBackupFrom:     "Бэкап от",
		MsgOverwriteWarn:  "Текущее сохранение будет перезаписано.",
		MsgCancel:         "ОТМЕНА",
		MsgFed:            "Сыт",
		MsgStarvingHours:  "ч голодания",
		MsgBloodLoss:      "Кровотечение",
		MsgBrokenRibs:     "Сломанные рёбра",
		MsgBurns:          "Ожог",
		MsgElectricBurns:  "Электрический ожог",
		MsgFrostbite:      "Обморожение",
		MsgSprainedAnkle:  "Растяжение ноги",
		MsgSprainedWrist:  "Растяжение руки",
		MsgSprainPain:     "Боль от растяжения",
		MsgFoodPoisoning:  "Пищевое отравление",
		MsgDysentery:      "Дизентерия",
		MsgInfectionRisk:  "Риск инфекции",
		MsgInfection:      "Инфекция",
		MsgHypothermia:    "Гипотермия",
		MsgIntestinalPar:  "Кишечные паразиты",
		MsgScurvy:         "Цинга",
		MsgHeadache:       "Головная боль",
		MsgInsomnia:       "Бессонница",
		MsgSevereLac:      "Рваная рана",
		MsgSuffocating:    "Удушье",
		MsgWeakConst:      "Слабое телосложение",
		MsgWeakJoints:     "Слабые суставы",
		MsgSourStomach:    "Слабый желудок",
		MsgPoorCirc:       "Плохое кровообращение",
		MsgUnsettledSleep: "Беспокойный сон",
		MsgBrokenBody:     "Измождённое тело",
		MsgSettings:       "НАСТРОЙКИ",
		MsgBackupDir:      "Папка для бэкапов",
		MsgBackupDirHint:  "Абсолютный путь к папке, где хранятся резервные копии. Например: /Volumes/Disk/LongDarkBackup",
		MsgSave:           "СОХРАНИТЬ",
		MsgSettingsSaved:  "Настройки сохранены",
		MsgPathRequired:   "Укажите путь к папке",
		MsgPathRelative:   "Путь должен быть абсолютным",
		MsgPathRoot:       "Нельзя использовать корень диска",
		MsgDirCreateErr:   "Не удалось создать папку",
		MsgDirNotWritable: "Нет прав на запись в папку",
		MsgConfigReadErr:  "Не удалось прочитать конфигурацию",
		MsgConfigSaveErr:  "Не удалось сохранить настройки",
	},
	LangEN: {
		MsgAppTitle:       "LONGDARK BACKUP",
		MsgNoSaves:        "No saves found",
		MsgBack:           "BACK",
		MsgSurvival:       "Survival",
		MsgCondition:      "Condition",
		MsgHunger:         "Hunger",
		MsgThirst:         "Thirst",
		MsgFatigue:        "Fatigue",
		MsgFreezing:       "Freezing",
		MsgWorldExplored:  "World explored",
		MsgCharacter:      "Character",
		MsgBackupLabel:    "Backup",
		MsgSaveLabel:      "Save",
		MsgRestore:        "RESTORE",
		MsgRestoreConfirm: "RESTORE SAVE?",
		MsgBackupFrom:     "Backup from",
		MsgOverwriteWarn:  "Current save will be overwritten.",
		MsgCancel:         "CANCEL",
		MsgFed:            "Fed",
		MsgStarvingHours:  "h starving",
		MsgBloodLoss:      "Blood loss",
		MsgBrokenRibs:     "Broken ribs",
		MsgBurns:          "Burns",
		MsgElectricBurns:  "Electric burns",
		MsgFrostbite:      "Frostbite",
		MsgSprainedAnkle:  "Sprained ankle",
		MsgSprainedWrist:  "Sprained wrist",
		MsgSprainPain:     "Sprain pain",
		MsgFoodPoisoning:  "Food poisoning",
		MsgDysentery:      "Dysentery",
		MsgInfectionRisk:  "Infection risk",
		MsgInfection:      "Infection",
		MsgHypothermia:    "Hypothermia",
		MsgIntestinalPar:  "Intestinal parasites",
		MsgScurvy:         "Scurvy",
		MsgHeadache:       "Headache",
		MsgInsomnia:       "Insomnia",
		MsgSevereLac:      "Severe laceration",
		MsgSuffocating:    "Suffocating",
		MsgWeakConst:      "Weak constitution",
		MsgWeakJoints:     "Weak joints",
		MsgSourStomach:    "Sour stomach",
		MsgPoorCirc:       "Poor circulation",
		MsgUnsettledSleep: "Unsettled sleep",
		MsgBrokenBody:     "Broken body",
		MsgSettings:       "SETTINGS",
		MsgBackupDir:      "Backup folder",
		MsgBackupDirHint:  "Absolute path to the folder where backups are stored. For example: /Volumes/Disk/LongDarkBackup",
		MsgSave:           "SAVE",
		MsgSettingsSaved:  "Settings saved",
		MsgPathRequired:   "Folder path is required",
		MsgPathRelative:   "Path must be absolute",
		MsgPathRoot:       "Disk root cannot be used",
		MsgDirCreateErr:   "Cannot create folder",
		MsgDirNotWritable: "Folder is not writable",
		MsgConfigReadErr:  "Cannot read settings",
		MsgConfigSaveErr:  "Cannot save settings",
	},
}

func T(key string) string {
	if m, ok := translations[currentLang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	if m, ok := translations[LangEN]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	return key
}

func DetectLang(cookie, acceptLang string) Lang {
	if cookie == "en" {
		return LangEN
	}
	if cookie == "ru" {
		return LangRU
	}
	if strings.Contains(acceptLang, "ru") {
		return LangRU
	}
	return LangEN
}
