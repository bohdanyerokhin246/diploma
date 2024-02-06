package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

type Person struct {
	SensorID   string `json:"sensor_id"`
	Indication string `json:"indication"`
	Date       string `json:"date"`
}

func main() {
	// Создаем срез объектов Person
	//winterMax, winterMin := 0, -20
	//springMax, springMin := 20, 0
	//summerMax, summerMin := 35, 20
	//autumnMax, autumnMin := 20, 0
	people := make([]Person, 400)
	k := 0
	for i := 0; i < 13; i++ {

		for j := 0; j < 30; j++ {

			if i >= 0 && i < 3 {
				people[k] = Person{
					SensorID:   fmt.Sprintf("%d", 1),
					Indication: fmt.Sprintf("%d", rand.Intn(5)),
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else if i > 2 && i < 6 {
				people[k] = Person{
					SensorID:   fmt.Sprintf("%d", 1),
					Indication: fmt.Sprintf("%d", rand.Intn(20)),
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else if i > 5 && i < 9 {
				people[k] = Person{
					SensorID:   fmt.Sprintf("%d", 1),
					Indication: fmt.Sprintf("%d", rand.Intn(30)),
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else if i > 8 && i < 12 {
				people[k] = Person{
					SensorID:   fmt.Sprintf("%d", 1),
					Indication: fmt.Sprintf("%d", rand.Intn(15)),
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
				fmt.Println(i, j)
			} else {
				people[k] = Person{
					SensorID:   fmt.Sprintf("%d", 1),
					Indication: fmt.Sprintf("%d", rand.Intn(5)),
					Date:       fmt.Sprintf("%d-%d-2023", i, j),
				}
			}
			k++
		}
	}

	// Отправляем каждый JSON объект на сервер
	for _, person := range people {
		// Конвертируем Person объект в JSON
		jsonBytes, err := json.Marshal(person)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Создаем HTTP запрос с JSON данными
		req, err := http.NewRequest("POST", "http://127.0.0.1:8080/indication/add", bytes.NewBuffer(jsonBytes))
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
