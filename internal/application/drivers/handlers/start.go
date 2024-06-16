package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const SendDataPeriod = time.Second * 5 // частота отсылки данных

// переменная нужна для создания websocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Обрабатываем запрос на начало работы
func (h *Handler) Start() http.HandlerFunc {

	// создаем функцию которую будем возвращать как http.HandlerFunc
	fn := func(rw http.ResponseWriter, req *http.Request) {

		// создаем соединение websocket
		conn, err := upgrader.Upgrade(rw, req, nil)

		if err != nil {
			fmt.Println("Ошибка подключения:", err)
			return
		}
		defer conn.Close()

		fmt.Println("Start websocket:")

		// читаем данные
		go readData(conn)
		// go func() {
		// 	// Чтение сообщений от клиента
		// 	for {
		// 		// Чтение сообщения
		// 		_, message, err := conn.ReadMessage()
		// 		if err != nil {
		// 			fmt.Println("Ошибка при чтении сообщения:", err)
		// 			break
		// 		}

		// 		// Вывод полученного сообщения
		// 		fmt.Println("Получено сообщение:", string(message))
		// 	}
		// }()

		// посылаем данные водителю
		sendData(conn)

	}

	return http.HandlerFunc(fn)
}

// Посылаем данные водителю с заданной периодичностью
func sendData(conn *websocket.Conn) {

	ticker := time.NewTicker(SendDataPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteMessage(websocket.TextMessage, []byte("Данные, которые вы хотите отправить клиенту"))
			if err != nil {
				fmt.Println("Ошибка отправки данных:", err)
				return
			}
		}
	}
}

func readData(conn *websocket.Conn) {

	// Чтение сообщений от клиента
	for {
		// Чтение сообщения
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Ошибка при чтении сообщения:", err)
			break
		}

		// Вывод полученного сообщения
		fmt.Println("Получено сообщение:", string(message))
	}
}
