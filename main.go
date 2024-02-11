package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Users struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

var users []Users

func main() {
	data, err := os.ReadFile("db.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &users)

	http.HandleFunc("/", usersHandler)
	http.ListenAndServe("localhost:8080", nil)
}

func usersHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if r.ContentLength > 0 {
			showUserByNumber(w, r)
		} else {
			showUsers(w, r)
		}

	case http.MethodPost:
		addUser(w, r)
	case http.MethodPut:
		updateNumber(w, r)
	case http.MethodDelete:
		deleteUser(w, r)
	case http.MethodHead:
		showUserByNumber(w, r)

	default:
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
	}

}

// Отображение всего списка телефонной книги методом GET
func showUsers(w http.ResponseWriter, r *http.Request) {

	for i, val := range users {
		fmt.Fprintf(w, "Пользователь № %v\n Имя: %v, Телефон: %v \n", i+1, val.Name, val.Phone)
	}
	log.Println(users)

}

/*
Добавление пользователя методом POST, после поулчаем сообщение
об успешном доавблении или сообщение, что такой пользователь уже сущетсвует
*/
func addUser(w http.ResponseWriter, r *http.Request) {
	var user Users
	var isTrue bool

	// здесь преобразуем введенные json данные в переменную user  с типом данных Users
	err := json.NewDecoder(r.Body).Decode(&user)
	// выдаем сообщение об ошибке, если она есть
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// проверяем есть ли указанный пользователь в базе
	for _, val := range users {
		if val.Name == user.Name || val.Phone == user.Phone {
			fmt.Fprint(w, "Такой пользовтель уже в базе")
			isTrue = true
			break
		}
	}

	if !isTrue {
		// добавляем пользовтаеля
		users = append(users, user)
		// саздаем слайс байтов из нового слайса users с обновленным номером
		pushUpdates(w)

		fmt.Fprintf(w, "Пользователь %v с номером %v добавлен", user.Name, user.Phone)
	}

}

// обновление номер по имени методом PUT
func updateNumber(w http.ResponseWriter, r *http.Request) {
	var user Users
	var isTrue bool

	// Получаем данные в формате json
	json.NewDecoder(r.Body).Decode(&user)

	// из полученных данных ищем пользователя с указнным именем, если совпадение есть, устанавливаем новый номер
	for i, val := range users {
		if val.Name == user.Name {
			// сообщение об успешном обновлении номера
			fmt.Fprintf(w, "У пользователя %v изменен номер. Старый номер: %v, новый номер: %v", user.Name, val.Phone, user.Phone)
			users[i].Phone = user.Phone
			isTrue = true

			// саздаем слайс байтов из нового слайса users с обновленным номером
			data, err := json.MarshalIndent(users, "", "\t")
			if err != nil {
				http.Error(w, err.Error(), 1)
			}
			// записываем в бд обновленные данные
			_ = os.WriteFile("db.json", data, 02)

		}
	}
	// Если заданный пользователь не найден, выдаем сообщение об отсутствии пользователя с таким именем
	if !isTrue {
		fmt.Fprint(w, "Такой пользователь не найден")
	}

}

// удаление пользователя из бд методом DELETE
func deleteUser(w http.ResponseWriter, r *http.Request) {
	var user Users
	var isTrue bool

	json.NewDecoder(r.Body).Decode(&user)

	for i, val := range users {
		if val.Name == user.Name || val.Phone == user.Phone {
			isTrue = true
			users = append(users[:i], users[i+1:]...)
			// Сообщение об успешном удалении пользователя
			fmt.Fprint(w, "Пользователь  удален")
			pushUpdates(w)
			break
		}
	}
	if !isTrue {
		fmt.Fprintf(w, "Пользовтель с именем %v не найден", user.Name)
	}
}

// Отображение пользователя по телефону
func showUserByNumber(w http.ResponseWriter, r *http.Request) {
	var phone Users
	var isTrue bool

	json.NewDecoder(r.Body).Decode(&phone)

	for _, val := range users {
		if val.Phone == phone.Phone {
			fmt.Fprintf(w, "Пользователя с номером %v зовут %v", val.Phone, val.Name)
			isTrue = true
			break
		}
	}
	if !isTrue {
		fmt.Fprintf(w, "Пользователя с номером %v не найден", phone)
	}
}

func pushUpdates(w http.ResponseWriter) {
	data, err := json.MarshalIndent(users, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), 1)
	}

	// записываем в бд обновленные данные
	_ = os.WriteFile("db.json", data, 02)
}
