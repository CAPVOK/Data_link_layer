package funcs

import (
	"fmt"
	"math/rand"
)

const (
	PROBABILITY_ONE_BIT = 0.09
)

func FuncToCall(arr, err []int) (int, int) {
	// Выбираем случайный элемент из списка
	randomIndex := rand.Intn(len(arr))
	randomElement := arr[randomIndex]

	var true_err []int
	for _, num := range err {
		if num <= randomElement {
			true_err = append(true_err, num)
		}
	}
	// Выбираем случайное число из списка, удовлетворяющее условиям
	randomNumberIndex := rand.Intn(len(true_err))
	randomErr := true_err[randomNumberIndex]
	result := randomElement ^ randomErr
	return result, randomIndex

}

func CallWithProbability(fn func(), probability float64) {
	// Генерируем случайное число от 0 до 1
	randomNum := rand.Float64()

	// Проверяем, выполнять ли функцию
	if randomNum < probability {
		fn()
	} else {
		fmt.Println("Ошибки не произошло")
	}
}
