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

type ResponseMessage struct {
	Message string `json:"message"`
}

// @Summary Передача данных на канальном уровне
// @Description Кодирует данные, вносит ошибку, декодирует с исправлением, теряет с заданной вероятностью и отправляет в ответ
// @Tags DataLink
// @Accept json
// @Produce json
// @Param segment body model.Segment true "Пользовательский объект в формате JSON"
// @Success 200 {object} ResponseMessage "Успешно"
// @Failure 400 {object} ResponseMessage "Некорректный запрос"
// @Failure 500 {object} ResponseMessage "Внутренняя ошибка сервера"
// @Router /api/datalink [post]
func (app *Application) PostDataLink(c *gin.Context) {
	start := time.Now()
	rand.Seed(time.Now().UnixNano())
	var segment model.Segment
	if err := c.BindJSON(&segment); err != nil {
		c.JSON(http.StatusBadRequest, ResponseMessage{Message: "Неправильный формат"})
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
		result, i := funcs.FuncToCall(decArr, errArr1Krat, last_len)
		// fmt.Println("Ошибка в элементе номер:", i)
		fmt.Printf("Элемент до ошибки:%b \n", decArr[i])
		fmt.Printf("Элемент с ошибкой:%b \n", result)
		decArr[i] = result
	}, model.PROBABILITY_ONE_BIT)
	decArr = decode(decArr)
	finalBin := funcs.DecimalArrayToBinary(decArr, last_len)
	finalByte := funcs.BinToByte(finalBin)
	finalText := string(finalByte)
	log.Printf("time %s\n", time.Since(start))
	segment.Message = finalText
	fmt.Println(err)
	fmt.Println(finalText)
	// c.JSON(http.StatusOK, gin.H{"segment": segment, "error": err})
	var segmentSend model.SegmentSend
	segmentSend.AmountOfSegments = segment.AmountOfSegments
	segmentSend.Error = err
	segmentSend.Message = segment.Message
	segmentSend.SegmentNum = segment.SegmentNum
	segmentSend.Sender = segment.Sender
	segmentSend.Timestamp = segment.Timestamp
	if rand.Float64() < model.SEND_PROBABILITY {
		// Кодируем JSON данные
		jsonData, err := json.Marshal(segmentSend)
		if err != nil {
			fmt.Println("Ошибка при кодировании JSON:", err)
			c.JSON(500, ResponseMessage{Message: "Внутренняя ошибка сервера"})
			return
		}

		// Выполняем POST запрос на указанный URL
		url := "http://localhost:8082/api/get-message"
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println("Ошибка при выполнении запроса:", err)
			c.JSON(400, ResponseMessage{Message: "Неправильный формат"})
			return
		}
		defer resp.Body.Close()

		// Проверяем статус код ответа
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Ошибка: неправильный статус код", resp.StatusCode)
			c.JSON(400, ResponseMessage{Message: "Неправильный формат"})
			return
		}

		fmt.Println("Запрос успешно выполнен")
		c.JSON(http.StatusOK, ResponseMessage{Message: "Успешно"})
	} else {
		fmt.Println("Запрос не был отправлен (вероятность 1%)")
	}
}

func encode(decArr []int) []int {
	for i, num := range decArr {
		decArr[i] = num << 4
		gx := funcs.DecimalToBinary(model.GX)
		mod := funcs.DecimalToBinary(decArr[i])
		if decArr[i] != 0 {
			decArr[i] = funcs.BinaryToDecimal(funcs.DecimalToBinary(decArr[i])[:(len(funcs.DecimalToBinary(decArr[i]))-4)] + funcs.GetRemainder(mod, gx))
		}
	}
	return decArr
}

func decode(decArr []int) []int {
	for i := range decArr {
		gx := funcs.DecimalToBinary(model.GX)
		mod := funcs.DecimalToBinary(decArr[i])
		// fmt.Println(decimalToBinary(decArr[i]))
		// fmt.Println(decimalToBinary(GX))
		// fmt.Println(GetRemainder(mod, gx))
		ex := funcs.Polynom_vector(funcs.GetRemainder(mod, gx))
		decArr[i] = decArr[i] ^ ex
		if decArr[i] != 0 {
			decArr[i] = funcs.BinaryToDecimal(funcs.DecimalToBinary(decArr[i])[:len(funcs.DecimalToBinary(decArr[i]))-4])
		}

	}
	return decArr
}
