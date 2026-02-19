# Структура сохранения The Long Dark

## Формат файла

Файл `sandbox<N>` — LZF-сжатый JSON.

```
sandbox файл → LZF decompress → JSON (верхний уровень)
```

### Верхний уровень

| Ключ | Тип | Описание |
|------|-----|----------|
| `m_DisplayName` | string | Имя сохранения (задаётся игроком) |
| `m_GameMode` | string | `"SANDBOX"`, `"STORY"` и т.д. |
| `m_InternalName` | string | Внутреннее имя файла (`sandbox35`) |
| `m_Timestamp` | string | Время сохранения (`2026-02-19 18:49:22`) |
| `m_GameId` | number | Числовой ID сохранения |
| `m_Changelist` | number | Версия билда игры |
| `m_Version` | number | Версия формата сохранения |
| `m_InstalledOptionalContent` | array | DLC (`["DLC01"]`) |
| `m_Dict` | object | Все данные сохранения (33 ключа) |

## m_Dict

Каждое значение в `m_Dict` закодировано:

```
значение → base64 decode → LZF decompress → JSON
```

### Ключи m_Dict

| Ключ | Описание |
|------|----------|
| `info` | Краткая сводка (region, condition, hours) |
| `global` | Состояние игрока, болезни, инвентарь, навыки (91 ключ) |
| `boot` | Текущая сцена (`m_SceneName`) |
| `screenshot` | Скриншот (`m_Encoded` — JPEG в base64) |
| `sandbox_global` | Бункеры (`m_BunkerSetup`) |
| `<SceneName>` | Локации (~28 штук) — контейнеры, предметы, AI, костры |

---

## info

Путь: `m_Dict["info"]`

| Ключ | Тип | Описание | Пример |
|------|-----|----------|--------|
| `m_XPMode` | string | Режим сложности | `"Misery"`, `"Voyageur"` |
| `m_Region` | string | Текущий регион | `"LakeRegion"` |
| `m_Persona` | string | Персонаж | `"Female"`, `"Male"` |
| `m_Condition` | number | HP в процентах | `51.81` |
| `m_HoursSurvived` | number | Игровые часы | `743.06` |
| `m_WorldExplored` | number | Доля исследованного мира (0-1) | `0.1786` |
| `m_LocationOverride` | string | Переопределение локации | `""` |
| `m_ActiveMissionId` | string | Активная миссия | `""` |
| `m_ActiveMissionLocId` | string | Локация миссии | `""` |

---

## global — Базовые параметры

Путь: `m_Dict["global"]` → десериализация → JSON с 91 ключом.

Многие ключи содержат вложенный JSON в строке (Serialized). Для их чтения нужен дополнительный `json.Unmarshal`.

### Здоровье (HP)

Ключ: `m_Condition_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CurrentHPProxy` | number | Текущее HP (0-100) |
| `m_ConditionLevelForPreviousVoiceOver` | string | Уровень: `"Healthy"`, `"SlightlyInjured"`, `"Injured"`, `"NearDeath"` |
| `m_NumSecondsSinceLastVoiceOver` | number | Секунды с последней озвучки |
| `m_CanPlayNearDeathMusic` | bool | Флаг музыки |

### Голод

Ключ: `m_Hunger_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_NumHoursStarving` | number | Часы голодания (0 = сыт) |
| `m_StarvingInLog` | bool | Голодает ли в журнале |
| `m_HungerLevelForPreviousVoiceOver` | string | `"NotHungry"`, `"Hungry"`, `"Starving"` |

### Жажда

Ключ: `m_Thirst_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CurrentThirstProxy` | number | Уровень жажды (0-100, где 0 = не хочет пить) |

### Усталость

Ключ: `m_Fatigue_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CurrentFatigueProxy` | number | Уровень усталости (0-100) |

### Замерзание

Ключ: `m_Freezing_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CurrentFreezingProxy` | number | Уровень замерзания (0-100) |

### Выносливость (спринт)

Ключ: `m_PlayerMovementSerialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_SprintStamina` | number | Стамина спринта (0-100) |

---

## global — Болезни и травмы (Afflictions)

Все affliction-ключи находятся в `m_Dict["global"]` и содержат строку с вложенным JSON.

### Кровотечение

Ключ: `m_BloodLossSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Locations` | array[int] | Зоны тела с кровотечением |
| `m_DurationHoursList` | array[number] | Длительность каждого кровотечения |
| `m_ElapsedHoursList` | array[number] | Прошедшее время лечения |
| `m_CausesLocIDs` | array[string] | Причины (`"GAMEPLAY_Bear"`, `"GAMEPLAY_Wolf"`) |

Активно, если `m_Locations` не пуст.

### Сломанные рёбра

Ключ: `m_BrokenRibSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Locations` | array[int] | Зоны с переломами |
| `m_NumHoursRestForCureList` | array[number] | Часы отдыха для лечения |
| `m_ElapsedRestList` | array[number] | Прошедшее время отдыха |
| `m_BandagesApplied` | array[bool] | Наложены ли бинты |
| `m_PainKillersTaken` | array[bool] | Приняты ли обезболивающие |
| `m_CausesLocIDs` | array[string] | Причины |

Активно, если `m_Locations` не пуст.

### Ожоги

Ключ: `m_BurnsSerialized` (от огня), `m_BurnsElectricSerialized` (электрические)

Пустые (`{}`) когда нет ожогов. При наличии содержат аналогичные поля с зонами и прогрессом.

### Обморожение

Ключ: `m_FrostbiteSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_LocationsCurrentFrostbiteDamage` | array[number] | Урон обморожения по 12 зонам тела (0 = нет) |
| `m_LocationsWithActiveFrostbite` | array[int] | Зоны с активным обморожением |
| `m_LocationsWithFrostbiteRisk` | array[int] | Зоны с риском обморожения |

Зоны тела (индексы 0-11): голова, лицо, шея, руки (L/R), грудь, живот, ноги (L/R), ступни (L/R), спина.

### Растяжение ноги

Ключ: `m_SprainedAnkleSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Locations` | array[int] | Зоны с растяжением |
| `m_DurationHoursList` | array[number] | Общая длительность лечения |
| `m_ElapsedHoursList` | array[number] | Прошедшее время |
| `m_ElapsedRestList` | array[number] | Часы отдыха |
| `m_CausesLocIDs` | array[string] | Причины |

Активно, если `m_Locations` не пуст.

### Растяжение руки

Ключ: `m_SprainedWristSerialized`

Структура идентична `m_SprainedAnkleSerialized`.

### Боль от растяжений

Ключ: `m_SprainPainSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Locations` | array[int] | Зоны с болью |
| `m_EndTimes` | array[number] | Время окончания боли (игровые часы) |
| `m_Causes` | array[string] | Причины |

### Пищевое отравление

Ключ: `m_FoodPoisoningSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CauseLocID` | string | Причина (`"GAMEPLAY_PorkandBeans"`) |
| `m_DurationHours` | number | Общая длительность болезни |
| `m_ElapsedHours` | number | Прошедшее время |
| `m_ElapsedRest` | number | Часы отдыха |
| `m_AntibioticsTaken` | bool | Приняты ли антибиотики |

Активно, если `m_CauseLocID` не пуст.

### Дизентерия

Ключ: `m_DysenterySerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CleanWaterConsumedLiters` | number | Литры чистой воды выпито (2000000000 = не болен / флаг) |

При активной дизентерии добавляются поля `m_DurationHours`, `m_ElapsedHours`, `m_AntibioticsTaken`.

### Инфекция (риск)

Ключ: `m_InfectionRiskSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Locations` | array[int] | Зоны с риском |
| `m_DurationHoursList` | array[number] | Длительность риска |
| `m_ElapsedHoursList` | array[number] | Прошедшее время |
| `m_CurrentInfectionChanceList` | array[number] | Шанс заражения (%) |
| `m_AntisepticTakenList` | array[bool] | Применён ли антисептик |
| `m_ConstantAfflictionIndices` | array[int] | Связи с постоянными болезнями |
| `m_CausesLocIDs` | array[string] | Причины |

### Инфекция

Ключ: `m_InfectionSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Locations` | array[int] | Зоны с инфекцией |
| `m_DurationHoursList` | array[number] | Длительность лечения |
| `m_ElapsedHoursList` | array[number] | Прошедшее время |
| `m_ElapsedRestList` | array[number] | Часы отдыха |
| `m_AntibioticsTakenList` | array[bool] | Приняты ли антибиотики |
| `m_CausesLocIDs` | array[string] | Причины |

### Кишечные паразиты

Ключ: `m_IntestinalParasitesSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CurrentInfectionChance` | number | Текущий шанс заражения (%) |
| `m_NumPiecesEatenThisRiskCycle` | number | Съедено кусков мяса хищника в цикле |
| `m_RiskDurationHours` | number | Длительность цикла риска (часов) |
| `m_RiskElapsedHours` | number | Прошедшее время цикла |
| `m_DayToAllowNextDose` | number | День для следующей дозы лекарства |

При активном заражении добавляются: `m_HasParasites`, `m_ParasiteDurationHours`, `m_ParasiteElapsedHours`, `m_DosesTaken`, `m_DosesRequired`.

### Гипотермия

Ключ: `m_HypothermiaSerialized`

Пустой (`{}`) когда нет гипотермии. При активной:

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Active` | bool | Активна ли |
| `m_ElapsedHours` | number | Прошедшее время |
| `m_WarmHoursElapsed` | number | Часы в тепле (лечение) |
| `m_WarmHoursRequired` | number | Необходимо часов в тепле |

### Cabin Fever (лихорадка затворника)

Ключ: `m_CabinFeverSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_IndoorTimeTracked` | array[number] | Трекинг времени в помещении за каждый день (0-1) |
| `m_HourLastFrame` | number | Час последнего обновления |

Риск рассчитывается на основе `m_IndoorTimeTracked` за последние 6 дней (>80% внутри = риск).

### Головная боль

Ключ: `m_HeadacheSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_ListOfParamData` | array | Активные головные боли (пусто = нет) |

### Бессонница

Ключ: `m_InsomniaSerialized`

| Поле | Тип | Описание |
|------|-----|----------|
| `m_ListOfParamData` | array | Активные случаи бессонницы |
| `m_ResistInsomniaBuffMultiplier` | number | Множитель сопротивления (1 = норма) |

### Цинга

Ключ: `m_ScurvySerialized` (объект, не строка!)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Enabled` | bool | Система цинги включена |
| `m_Repaired` | bool | Вылечена (`true` = здоров) |

### Тяжёлая рваная рана

Ключ: `m_SevereLacerationSerialized`

Пустой объект `{}` когда нет. При наличии — аналогично `m_BloodLossSerialized`.

### Удушье

Ключ: `m_SuffocatingSerialized`

Пустой когда нет. При наличии — содержит данные об удушье (от пещерного воздуха).

### Тревога / Страх

Ключи: `m_AnxietySerialized`, `m_FearSerialized`

Пустые когда нет. Появляются в сюжетных режимах.

### Химическое отравление

Ключ: `m_ChemicalPoisoningSerialized` (объект, не строка!)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Afflictions` | array | Активные отравления (пусто = нет) |

### Ослабленное состояние (Diminished)

Ключ: `m_DiminishedStateAfflictionSerialized` (объект)

Пустой объект `{}` когда нет.

### Чит-смерть

Ключ: `m_CheatDeathAfflictionSerialized` (объект)

Пустой объект `{}` когда нет.

---

## global — Питание и витамины

### Витамин C

Ключ: `m_PlayerNutritionSerialized` (объект, не строка!)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_NutrientNames` | array[string] | `["Nutrient_VitaminC"]` |
| `m_Amounts` | array[number] | Текущий запас (индексы соответствуют `m_NutrientNames`) |
| `m_LossPerDay` | array[number] | Потеря в день |

### Калории

Ключ: `m_PlayerGameStatsSerialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CaloriesEaten` | number | Всего калорий съедено |
| `m_CaloriesBurned` | number | Всего калорий сожжено |
| `m_CaloriesExpendedToday` | number | Калорий сожжено сегодня |
| `m_ConditionGained` | number | Всего HP восстановлено |
| `m_ConditionLost` | number | Всего HP потеряно |
| `m_DistanceTravelledDay` | number | Расстояние днём |
| `m_DistanceTravelledNight` | number | Расстояние ночью |
| `m_BodyTempHigh` | number | Максимальная температура тела |
| `m_BodyTempLow` | number | Минимальная температура тела |

---

## global — Навыки

### Skills Manager

Ключ: `m_SkillsManagerSerialized` (строка → JSON)

Каждый навык — строка с вложенным JSON:

| Поле | Тип | Описание |
|------|-----|----------|
| `m_Skill_ArcherySerialized` | string→JSON | Стрельба из лука (`m_Points`) |
| `m_Skill_CarcassHarvestingSerialized` | string→JSON | Разделка туш (`m_Points`, `m_NumHoursToConvertToSkillPoints`) |
| `m_Skill_ClothingRepairSerialized` | string→JSON | Ремонт одежды |
| `m_Skill_CookingSerialized` | string→JSON | Кулинария |
| `m_Skill_FirestartingSerialized` | string→JSON | Разведение огня |
| `m_Skill_GunsmithingSerialized` | string→JSON | Оружейное дело |
| `m_Skill_IceFishingSerialized` | string→JSON | Подлёдная рыбалка |
| `m_Skill_RevolverSerialized` | string→JSON | Револьвер |
| `m_Skill_RifleSerialized` | string→JSON | Винтовка |

### Player Skills

Ключ: `m_PlayerSkillsSerialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CleanSkill` | number | Навык очистки |
| `m_RepairSkill` | number | Навык ремонта |
| `m_SharpenSkill` | number | Навык заточки |

---

## global — Время и погода

### Время суток

Ключ: `m_TimeOfDay_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_TimeProxy` | number | Время суток (0-1, где 0.5 = полдень) |
| `m_UniStormDayNumberProxy` | number | Номер текущего дня |
| `m_UniStormDayCounterProxy` | number | Счётчик дней |
| `m_HoursPlayedNotPausedProxy` | number | Игровые часы без паузы |
| `m_UniStormMoonPhaseIndexProxy` | number | Фаза луны |

### Погода

Ключ: `m_Weather_Serialized` (строка → JSON)

Содержит данные о температуре, осадках, видимости. Структура зависит от версии.

---

## global — Инвентарь

Ключ: `m_Inventory_Serialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_SerializedItems` | array[object] | Массив предметов |
| `m_QuickSelectInstanceIDs` | array[number] | Слоты быстрого выбора |

Каждый предмет:

| Поле | Тип | Описание |
|------|-----|----------|
| `m_PrefabName` | string | Имя предмета (`"GEAR_WaterSupplyPotable"`) |
| `m_SerializedGear` | string→JSON | Состояние предмета (HP, позиция, кол-во и т.д.) |

---

## global — Журнал

Ключ: `m_LogSerialized` (строка → JSON)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_LogDayInfoList` | array[object] | Записи по дням |
| `m_TodayLogDayInfo` | object | Текущий день |
| `m_DayToLogEndOfDayInfo` | number | Следующий день для записи |
| `m_GeneralNotes` | string | Заметки игрока |

Запись дня:

| Поле | Тип | Описание |
|------|-----|----------|
| `m_DayNumber` | number | Номер дня |
| `m_CaloriesBurned` | number | Калории за день |
| `m_ConditionHigh` | number | Макс HP за день |
| `m_ConditionLow` | number | Мин HP за день |
| `m_HoursRested` | number | Часы сна |
| `m_WorldExplored` | number | % мира исследовано |
| `m_Afflictions` | array[string] | Болезни за день (`"WeakConstitution"`, `"FoodPoisoning"`) |
| `m_LocationLocIDs` | array[string] | Посещённые локации |
| `m_RegionLocIDs` | array[string] | Регионы |

---

## global — Misery Mode

Ключ: `m_MiseryManagerSerialized` (объект, не строка!)

| Поле | Тип | Описание |
|------|-----|----------|
| `m_CurrentStage` | number | Текущая стадия (0-5) |
| `m_MiseryStages` | array[object] | Конфигурация стадий |

Каждая стадия:

| Поле | Тип | Описание |
|------|-----|----------|
| `m_ActiveUntilHours` | number | Активна до N часов выживания |
| `m_Afflictions` | array[string] | Постоянные дебаффы (`"WeakConstitution"`, `"WeakJoints"`, `"SourStomach"`) |

---

## Как извлечь данные (Go)

### 1. Чтение sandbox файла

```go
import lzf "github.com/zhuyie/golzf"

raw, _ := os.ReadFile("sandbox35")
buf := make([]byte, len(raw)*10)
n, _ := lzf.Decompress(raw, buf)

var top map[string]interface{}
json.Unmarshal(buf[:n], &top)
```

### 2. Чтение записи из m_Dict

```go
import "encoding/base64"

mDict := top["m_Dict"].(map[string]interface{})
encoded := mDict["global"].(string)

b64decoded, _ := base64.StdEncoding.DecodeString(encoded)
buf := make([]byte, len(b64decoded)*10)
n, _ := lzf.Decompress(b64decoded, buf)

var global map[string]interface{}
json.Unmarshal(buf[:n], &global)
```

### 3. Чтение Serialized-поля

```go
serialized := global["m_FoodPoisoningSerialized"].(string)

var data map[string]interface{}
json.Unmarshal([]byte(serialized), &data)

// data["m_CauseLocID"] → "GAMEPLAY_PorkandBeans"
// data["m_AntibioticsTaken"] → true
// data["m_ElapsedHours"] → 10.18
```

### 4. Проверка активности affliction

Для массивных afflictions (кровотечение, растяжения, инфекция):
```go
var data map[string]interface{}
json.Unmarshal([]byte(serialized), &data)

locations := data["m_Locations"].([]interface{})
isActive := len(locations) > 0
```

Для одиночных (пищевое отравление):
```go
causeLocID := data["m_CauseLocID"].(string)
isActive := causeLocID != ""
```

Для пустых объектов (ожоги, гипотермия):
```go
isActive := serialized != "{}" && serialized != ""
```

### 5. Чтение скриншота

```go
encoded := mDict["screenshot"].(string)
// base64 decode → LZF decompress → JSON
// → поле m_Encoded → base64 decode → JPEG bytes
```

---

## Зоны тела (Body Locations)

Индексы используемые в `m_Locations`, `m_LocationsCurrentFrostbiteDamage`:

| Индекс | Зона |
|--------|------|
| 0 | Голова |
| 1 | Лицо |
| 2 | Шея |
| 3 | Левая рука |
| 4 | Правая рука |
| 5 | Грудь |
| 6 | Живот |
| 7 | Левая нога |
| 8 | Правая нога |
| 9 | Левая ступня |
| 10 | Правая ступня |
| 11 | Спина |

---

## Коды причин (CauseLocIDs)

| Код | Описание |
|-----|----------|
| `GAMEPLAY_Wolf` | Атака волка |
| `GAMEPLAY_Bear` | Атака медведя |
| `GAMEPLAY_Cougar` | Атака пумы |
| `GAMEPLAY_Fall` | Падение |
| `GAMEPLAY_PorkandBeans` | Свинина с бобами (отравление) |
| `GAMEPLAY_CannedSardines` | Консервированные сардины |
| `GAMEPLAY_Noname` | Неизвестная причина |

---

## Локации (сцены)

Каждая локация в `m_Dict` содержит ~43-46 ключей:

| Ключ | Описание |
|------|----------|
| `m_ContainerManagerSerialized` | Контейнеры (сундуки, шкафы) |
| `m_GearManagerSerialized` | Предметы на земле |
| `m_BaseAiManagerSerialized` | Животные (волки, медведи, олени) |
| `m_FireManagerSerialized` | Костры и печки |
| `m_BreakDownObjectsSerialized` | Разбираемые объекты |
| `m_BodyHarvestManagerSerialized` | Туши животных |
| `m_BedSerialized` | Кровати |
| `m_SpawnRegionManagerSerialized` | Зоны спавна |
| `m_Version` | Версия формата сцены |
