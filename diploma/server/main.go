package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Indication struct {
	SensorID   int     `json:"sensorid"`
	Indication float64 `json:"indication"`
	Date       string  `json:"date"`
}
type Data struct {
	AvgValue float64 `json:"average"`
}

// Define an array to store the person data
var indications []Indication

func main() {
	// Register the API routes
	http.HandleFunc("/", mainPageHandler)
	http.HandleFunc("/indication/add", addIndicationHandler)
	http.HandleFunc("/indications/sensor/", getIndicationBySensorHandler)
	http.HandleFunc("/average/date", averageValueIndicatorsByDAteHandler)
	http.HandleFunc("/indication/all", getAllIndicationsHandler)
	// Start the server on port 8080
	fmt.Println(http.ListenAndServe(":8080", nil))

}

// Handler for retrieving all people
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	// Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	serverStatus := []byte("Welcome!\nServer is working.\nExample of requests:\n/indications/sensor?sensorid - to SELECT all indications by sensor ID\n" +
		"/average/date&sensorid=XX&start_date=XX-XX-XXXX&end_date=XX-XX-XXXX - to get average values by period from definite sensor\n" +
		"/indication/all - to get all indications from all sensors")
	// Write the JSON response

	w.Write(serverStatus)
}

// Handler for retrieving all people
func getIndicationBySensorHandler(w http.ResponseWriter, r *http.Request) {
	// Установка соединения с базой данных PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1 dbname=weathersensor sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Подготовленный оператор для выборки данных за определенные даты

	stmt, err := db.Prepare(`SELECT * FROM indications WHERE "sensorID" = $1 ORDER BY date`)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	// Задаем даты для выборки
	sensorID := r.URL.Query().Get("sensorid")

	// Выполняем выборку данных за определенные даты
	rows, err := stmt.Query(sensorID)
	if err != nil {
		fmt.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	// Итерируемся по результатам выборки и добавляем их в срез
	for rows.Next() {
		var indication Indication
		err := rows.Scan(&indication.SensorID, &indication.Indication, &indication.Date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error scanning row: %v", err)
			return
		}

		indications = append(indications, indication)
	}

	// Преобразуем результаты выборки в JSON
	jsonData, err := json.Marshal(indications)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON ответ
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// Handler for adding a person
func addIndicationHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1 dbname=weathersensor sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Читаем данные JSON из тела запроса
	var indications []Indication
	err = json.NewDecoder(r.Body).Decode(&indications)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding JSON: %v", err)
		return
	}

	// Вставляем данные в базу данных

	stmt, err := db.Prepare(`INSERT INTO indications ("sensorID", temperature, "date") VALUES ($1, $2, $3)`)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	for _, indication := range indications {
		_, err = stmt.Exec(stmt, indication.SensorID, indication.Indication, indication.Date)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Отправляем успешный ответ
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Data successfully inserted into the database.")
}

func getAllIndicationsHandler(w http.ResponseWriter, r *http.Request) {
	// Установка соединения с базой данных PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1 dbname=weathersensor sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Формируем SQL запрос с условиями выборки
	query := "SELECT * FROM indications "
	args := make([]interface{}, 0)

	// Выполняем SQL запрос и получаем результаты выборки
	rows, err := db.Query(query, args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error querying data: %v", err)
		return
	}
	defer rows.Close()

	// Итерируемся по результатам выборки и добавляем их в срез
	for rows.Next() {
		var indication Indication
		err := rows.Scan(&indication.SensorID, &indication.Indication, &indication.Date)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error scanning row: %v", err)
			return
		}

		indications = append(indications, indication)
	}

	// Преобразуем результаты выборки в JSON
	jsonData, err := json.Marshal(indications)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error marshaling JSON: %v", err)
		return
	}

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON ответ
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func averageValueIndicatorsByDAteHandler(w http.ResponseWriter, r *http.Request) {
	// Установка соединения с базой данных PostgreSQL
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=1 dbname=weathersensor sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Проверка соединения с базой данных
	err = db.Ping()
	if err != nil {
		fmt.Println("Error pinging the database:", err)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Подготовленный оператор для выборки данных за определенные даты
	stmt, err := db.Prepare(`SELECT temperature FROM indications WHERE "sensorID" = $1 AND date >= $2 AND date <= $3`)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return
	}
	defer stmt.Close()

	// Задаем даты для выборки
	sensorID := r.URL.Query().Get("sensorid")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Выполняем выборку данных за определенные даты
	rows, err := stmt.Query(sensorID, startDate, endDate)
	if err != nil {
		fmt.Println("Error querying data:", err)
		return
	}
	defer rows.Close()

	var sum float64
	var count int

	// Итерируемся по результатам выборки и вычисляем сумму и количество значений
	for rows.Next() {
		var value float64
		err := rows.Scan(&value)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}

		sum += value
		count++
	}

	// Рассчитываем среднее значение чисел
	var average float64
	if count > 0 {
		average = sum / float64(count)
	}

	// Создаем структуру среднего значения
	data := Data{
		AvgValue: average,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	sendAverage(w, r, jsonData)
}

func sendAverage(w http.ResponseWriter, r *http.Request, jsonData []byte) {

	// Устанавливаем заголовок Content-Type на application/json
	w.Header().Set("Content-Type", "application/json")

	// Отправляем JSON ответ
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

	fmt.Println("JSON sent successfully.")

}
