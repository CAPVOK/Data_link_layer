package funcs

import (
	"fmt"
	"math/rand"
)

func FuncToCall(arr, err []int) (int, int) {
	// Выбираем случайный элемент из списка
	randomIndex := rand.Intn(len(arr))
	randomElement := arr[randomIndex]

	var true_err []int
	if randomElement == 0 {
		true_err = append(true_err, err[0])
	}

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

func CallWithProbability(fn func(), probability float64) (err bool) {
	// Генерируем случайное число от 0 до 1
	randomNum := rand.Float64()

	// Проверяем, выполнять ли функцию
	if randomNum < probability {
		fn()
		err = true
	} else {
		fmt.Println("Ошибки не произошло")
		err = false
	}
	return err
}
