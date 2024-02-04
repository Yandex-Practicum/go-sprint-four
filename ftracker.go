package ftracker

import (
	"fmt"
	"math"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага в метрах.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// distance возвращает дистанцию (км), которую преодолел пользователь за время тренировки. (км)
//
// Параметры:
//
// action int — количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки. (км/ч)
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки (часы).
func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	distance := distance(action)
	return distance / duration
}

// WalkingSpentCalories возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string — вид тренировки(Бег, Ходьба, Плавание).
// duration float64 — длительность тренировки (часы).
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	speed := 0.0
	calories := 0.0

	distance := distance(action)

	switch {
	case trainingType == "Бег":
		speed = meanSpeed(action, duration)
		calories = RunningSpentCalories(action, weight, duration)
	case trainingType == "Ходьба":
		speed = meanSpeed(action, duration)
		calories = WalkingSpentCalories(action, duration, weight, height)
	case trainingType == "Плавание":
		speed = swimMeanSpeed(lengthPool, countPool, duration)
		calories = SwimmingSpentCalories(lengthPool, countPool, duration, weight)
	default:
		return "неизвестный тип тренировки"
	}
	return fmt.Sprintf("Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		trainingType,
		duration,
		distance,
		speed,
		calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	speedRunMult  = 18   // множитель средней скорости.
	speedRunShift = 1.79 // среднее количество сжигаемых калорий при беге.
)

// WalkingSpentCalories возвращает количе\ство потраченных колорий при беге.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// weight float64 — вес пользователя.
// duration float64 — длительность тренировки (часы).
func RunningSpentCalories(action int, weight, duration float64) float64 {
	meanSpeed := meanSpeed(action, duration)
	return speedRunMult * meanSpeed * speedRunShift * weight / mInKm * duration * minInH
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	weightWalkMult = 0.035 // множитель массы тела.
	heightWalkMult = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// action int — количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 — длительность тренировки (часы).
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	speed := meanSpeed(action, duration) * kmhInMsec // m/s
	speedSqr := math.Pow(speed, 2)
	height = height / cmInM
	return (weightWalkMult + (speedSqr/height)*heightWalkMult) * weight * duration * minInH
}

// swimMeanSpeed возвращает среднюю скорость при плавании (км/ч).
//
// Параметры:
//
// lengthPool int — длина бассейна (метры).
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки (часы).
func swimMeanSpeed(lengthPool, countPool int, duration float64) float64 {

	if duration == 0 {
		return 0
	}
	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// Константы для расчета калорий, расходуемых при плавании.
const (
	speedSwimShift = 1.1 // среднее количество сжигаемых колорий при плавании относительно скорости.
	weightSwimMult = 2   // множитель веса при плавании.
)

// SwimmingSpentCalories возвращает количество потраченных калорий при плавании.
//
// Параметры:
//
// lengthPool int — длина бассейна (метры).
// countPool int — сколько раз пользователь переплыл бассейн.
// duration float64 — длительность тренировки (часы).
// weight float64 — вес пользователя.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	meanSpeed := swimMeanSpeed(lengthPool, countPool, duration)
	return (meanSpeed + speedSwimShift) * weightSwimMult * weight * duration
}
