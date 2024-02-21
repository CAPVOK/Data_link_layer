package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lud0m4n/Network/internal/http/funcs"
	"github.com/lud0m4n/Network/internal/model"
)

const (
	GX                  = 19
	PROBABILITY_ONE_BIT = 0.09
	SEND_PROBABILITY    = 0.99
)

// @Summary Получение списка периодов
// @Description Возращает список всех активных периодов
// @Tags Период
// @Produce json
// @Param searchName query string false "Название периода" Format(email)
// @Success 200 {object} model.PeriodGetResponse "Список периодов"
// @Failure 500 {object} model.PeriodGetResponse "Ошибка сервера"
// @Router /period [get]
func (app *Application) PostDataLink(c *gin.Context) {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	var segment model.Segment
	if err := c.BindJSON(&segment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	errArr1Krat := []int{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384}

	text := segment.Message
	fmt.Println("Начальный текс:", text)
	byte_code := []byte(text)
	bin := funcs.ByteToBin(byte_code)
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
	err := funcs.CallWithProbability(func() {
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
	segment.Message = finalText
	c.JSON(http.StatusOK, gin.H{"segment": segment, "error": err})

	if rand.Float64() < SEND_PROBABILITY {
		// Кодируем JSON данные
		jsonData, err := json.Marshal(segment)
		if err != nil {
			fmt.Println("Ошибка при кодировании JSON:", err)
			return
		}

		// Выполняем POST запрос на указанный URL
		url := "http://localhost:8082/api/get-message"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка при выполнении запроса:", err)
			return
		}
		defer resp.Body.Close()

		// Проверяем статус код ответа
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Ошибка: неправильный статус код", resp.StatusCode)
			return
		}

		fmt.Println("Запрос успешно выполнен")
	} else {
		fmt.Println("Запрос не был отправлен (вероятность 1%)")
	}
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
