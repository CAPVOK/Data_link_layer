package app

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/Network/internal/http/funcs"
)

const (
	GX                  = 19
	PROBABILITY_ONE_BIT = 0.09
)

// @Summary Получение списка периодов
// @Description Возращает список всех активных периодов
// @Tags Период
// @Produce json
// @Param searchName query string false "Название периода" Format(email)
// @Success 200 {object} model.PeriodGetResponse "Список периодов"
// @Failure 500 {object} model.PeriodGetResponse "Ошибка сервера"
// @Router /period [get]
func (app *Application) GetPeriods(c *gin.Context) {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	errArr1Krat := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384}
	text := "Фрукты, много fruit? Длинный text для test ъъъъььья" // Ваш текст
	fmt.Println("Начальный текс:", text)
	bytes := []byte(text)
	bin := funcs.ByteToBin(bytes)
	var boolMatrix [][]bool
	for i := 0; i < len(bin); i += 11 {
		end := i + 11

		if end > len(bin) {
			end = len(bin)
		}

		boolMatrix = append(boolMatrix, bin[i:end])
	}
	last_len := len(boolMatrix[len(boolMatrix)-1])
	decArr := funcs.BinaryArrayToDecimalArray(boolMatrix)
	decArr = encode(decArr)
	// Вызываем функцию с вероятностью 9%
	funcs.CallWithProbability(func() {
		result, i := funcs.FuncToCall(decArr, errArr1Krat)

		// fmt.Println("Ошибка в элементе номер:", i)
		fmt.Printf("Элемент до ошибки:%b \n", decArr[i])
		fmt.Printf("Элемент с ошибкой:%b \n", result)
		decArr[i] = result
	}, PROBABILITY_ONE_BIT)
	decArr = decode(decArr)
	finalBin := funcs.DecimalArrayToBinary(decArr, last_len)
	finalByte := funcs.BinToByte(finalBin)
	finalText := string(finalByte)
	log.Printf("time %s\n", time.Since(start))

	c.JSON(http.StatusOK, finalText)
}

func encode(decArr []int) []int {
	for i, num := range decArr {
		decArr[i] = num << 4
		gx := funcs.DecimalToBinary(GX)
		mod := funcs.DecimalToBinary(decArr[i])
		decArr[i] = decArr[i] + funcs.BinaryToDecimal(funcs.GetRemainder(mod, gx))
	}
	return decArr
}

func decode(decArr []int) []int {
	for i := range decArr {
		gx := funcs.DecimalToBinary(GX)
		mod := funcs.DecimalToBinary(decArr[i])
		// fmt.Println(DecimalToBinary(decArr[i]))
		// fmt.Println(DecimalToBinary(GX))
		// fmt.Println(GetRemainder(mod, gx))
		ex := funcs.Polynom_vector(funcs.GetRemainder(mod, gx))
		decArr[i] = decArr[i] ^ ex
		decArr[i] = funcs.BinaryToDecimal(funcs.DecimalToBinary(decArr[i])[:len(funcs.DecimalToBinary(decArr[i]))-4])
	}
	return decArr
}
