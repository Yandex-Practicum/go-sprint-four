package ftracker

import (
	"fmt"
	"math"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	distance := distance(action)
	return distance / duration
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string — вид тренировки(Бег, Ходьба, Плавание).
// duration float64 — длительность тренировки в часах.
// weight float64 - вес пользователя в килограммах.
// height float64 - рост пользователя в сантиметрах.
// lengthPool - длина бассейна в метрах.
// countPool - количество проплытых бассейнов.
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	const (
		resultMessage string = `Тип тренировки: %s
        Длительность: %.2f ч.
        Дистанция: %.2f км.
        Скорость: %.2f км/ч.
        Потрачено энергии: %.2f ккал.
        `
	)
	var trainingDistance, trainingSpeed, spentCalories float64
	switch {
	case trainingType == "Бег":
		trainingDistance = distance(action)
		trainingSpeed = meanSpeed(action, duration)
		spentCalories = RunningSpentCalories(action, weight, duration)
	case trainingType == "Ходьба":
		trainingDistance = distance(action)
		trainingSpeed = meanSpeed(action, duration)
		spentCalories = WalkingSpentCalories(action, duration, weight, height)
	case trainingType == "Плавание":
		trainingDistance = distance(action)
		trainingSpeed = swimmingMeanSpeed(lengthPool, countPool, duration)
		spentCalories = SwimmingSpentCalories(lengthPool, countPool, weight, duration)
	default:
		return "неизвестный тип тренировки"
	}
	return fmt.Sprintf(resultMessage, trainingType, duration, trainingDistance, trainingSpeed, spentCalories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18   // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// weight float64 — вес пользователя.
// duration float64 — длительность тренировки в часах.
func RunningSpentCalories(action int, weight, duration float64) float64 {
	return (runningCaloriesMeanSpeedMultiplier * meanSpeed(action, duration) * runningCaloriesMeanSpeedShift) * weight / mInKm * duration * minInH
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	meanSpeedMS := meanSpeed(action, duration) * kmhInMsec
	heightM := height / cmInM
	return (walkingCaloriesWeightMultiplier*weight + (math.Pow(meanSpeedMS, 2)/heightM)*walkingSpeedHeightMultiplier*weight) * duration * minInH
}

// Константы для расчета калорий, расходуемых при плавании.
const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // среднее количество сжигаемых колорий при плавании относительно скорости.
	swimmingCaloriesWeightMultiplier = 2   // множитель веса при плавании.
)

// swimmingMeanSpeed возвращает среднюю скорость при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// SwimmingSpentCalories возвращает количество потраченных калорий при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна в метрах.
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки в часах.
// weight float64 — вес пользователя.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	return (swimmingMeanSpeed(lengthPool, countPool, duration) + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration
}
