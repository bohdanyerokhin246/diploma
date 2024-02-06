package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type Indication struct {
	SensorID   int    `json:"sensorid"`
	Indication int    `json:"indication"`
	Date       string `json:"date"`
}

func main() {
	indications := generateJSON()
	sendJSON(indications)
}
func generateJSON() []Indication {

	indications := make([]Indication, 400)
	k := 0
	for i := 0; i < 13; i++ {
		for j := 0; j < 30; j++ {
			if i >= 0 && i < 3 {
				temperature := rand.Intn(20) + 5
				indications[k] = Indication{
					SensorID:   2,
					Indication: -temperature,
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else if i > 2 && i < 6 {
				indications[k] = Indication{
					SensorID:   2,
					Indication: rand.Intn(12) + 5,
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else if i > 5 && i < 9 {
				indications[k] = Indication{
					SensorID:   2,
					Indication: rand.Intn(27) + 5,
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else if i > 8 && i < 12 {
				indications[k] = Indication{
					SensorID:   2,
					Indication: rand.Intn(15) + 5,
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else {
				temperature := rand.Intn(20) + 5
				indications[k] = Indication{
					SensorID:   2,
					Indication: -temperature,
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
			}
			k++
		}
	}
	return indications
}

func sendJSON(indications []Indication) {
	// Отправляем каждый JSON объект на сервер
	for _, indication := range indications {
		// Конвертируем Person объект в JSON
		jsonBytes, err := json.Marshal(indication)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Создаем HTTP запрос с JSON данными
		req, err := http.NewRequest("POST", "http://10.20.77.3:8080/indication/add", bytes.NewBuffer(jsonBytes))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Устанавливаем заголовок Content-Type на application/json
		req.Header.Set("Content-Type", "application/json")

		// Создаем HTTP клиент и отправляем запрос
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()

		// Проверяем код статуса ответа
		if resp.StatusCode != http.StatusOK {
			fmt.Println("Unexpected response:", resp.Status)
			return
		}

		// Читаем тело ответа
		var responseJSON map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&responseJSON)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Обрабатываем JSON ответ
		fmt.Println("Response:", responseJSON)
	}
}
