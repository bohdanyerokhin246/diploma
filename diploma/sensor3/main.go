package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

type Indication struct {
	SensorID    string `json:"sensor_id"`
	Temperature string `json:"temperature"`
	Date        string `json:"date"`
}

func main() {
	// Открываем CSV файл для чтения
	file, err := os.Open("data1.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Создаем читатель CSV файла
	reader := csv.NewReader(file)

	// Создаем срез для хранения данных
	var indications []Indication

	// Читаем данные из CSV файла
	for {
		// Считываем строку из CSV файла
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV row:", err)
			return
		}

		// Создаем структуру Person из данных строки
		indication := Indication{
			SensorID:    row[0],
			Date:        row[1],
			Temperature: row[2],
		}

		// Добавляем Person в срез
		indications = append(indications, indication)
	}

	//// Преобразуем срез в JSON
	//jsonData, err := json.Marshal(indications)
	//if err != nil {
	//	fmt.Println("Error marshaling JSON:", err)
	//	return
	//}

	//// Отправляем JSON на сервер
	//url := "http://127.0.0.1:8080/indication/add"
	//req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(jsonData)))
	//if err != nil {
	//	fmt.Println("Error creating HTTP request:", err)
	//	return
	//}
	//req.Header.Set("Content-Type", "application/json")
	//
	//client := http.DefaultClient
	//resp, err := client.Do(req)
	//if err != nil {
	//	fmt.Println("Error sending HTTP request:", err)
	//	return
	//}
	//defer resp.Body.Close()
	//
	//// Проверяем код статуса ответа сервера
	//if resp.StatusCode != http.StatusOK {
	//	fmt.Println("Server returned non-OK status:", resp.Status)
	//	return
	//}
	for i := 0; i < len(indications); i++ {
		fmt.Println(fmt.Sprintf("Indication with date %s is sent to server", indications[i].Date))
	}
	fmt.Scanln()
	fmt.Scanln()
}
