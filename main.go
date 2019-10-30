package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

var messages [3]string
var action byte

// Модель данных для хранения телефонной книги
type Record struct {
	Name   string `json:"Name"`
	Number uint64 `json:"Phone"`
}

type AddressBook struct {
	Records map[uint64]Record `json:"AddressBook"`
}

// Функция печатает список записей в телефонной книге
func printList(list *AddressBook) {
	for key, value := range (*list).Records {
		log.Println("№", key, "- Имя:", value.Name, "/ Номер телефона:", value.Number)
	}
}

// Функция записывает данные в файл
func saveToFile(tempData *AddressBook) {
	p := log.Println // the alias for log.Println in order to simplify the code
	dataToWrite, err := json.Marshal(*tempData)
	if err == nil {
		if err = ioutil.WriteFile("addressbook", dataToWrite, 0777); err == nil {
			p("Запись данных выполнена")
		} else {
			p("Ошибка записи данных")
			p(err)
		}
	} else {
		p("Ошибка обращения данных в JSON")
		p(err)
	}
}

// Функция вычитывает данные из файла
func list() (tempData AddressBook, err error) {
	content, err := ioutil.ReadFile("addressbook")
	if err == nil {
		err = json.Unmarshal(content, &tempData)
	}
	return
}

// Функция удаляем данные в телефонной книге
func deleteRecord() {
	p := log.Println // the alias for log.Println in order to simplify the code

	p("Enter an index key of the user you want to delete, enter 0 to cancel:")

	var key uint64
	for {
		fmt.Scan(&key)
		if key == 0 {
			p("Deletion is cancelled...")
			return
		}

		// получаем данные из файла
		tempData, err := list()
		if err != nil {
			p("Ошибка:", err)
			return
		}

		// проверяем есть ли введённый ключ в базе
		recordToDelete, ok := tempData.Records[key]
		if ok {
			p("Ключ будет удалён:", recordToDelete.Name, "Телефонный номер:", recordToDelete.Number)
			delete(tempData.Records, key) // удаляем из карты запись
			saveToFile(&tempData)         // и записываем обновлённые данные в файл
			return
		}
		p("Введённый ключ отсутствует в базе телефонной книги, повторите ввод")
	}
}

// Функция добавляет запись в телефонную книгу
func add() {
	p := log.Println // the alias for log.Println in order to simplify the code

	// вводим имя
	p("Please, enter a name:")
	var name string
	for {
		fmt.Scan(&name)
		if name != "" {
			p("The data input is done")
			break
		} else {
			p("The data input has been wrong. Please, repeat:")
		}
	}

	// вводим номер
	p("Please, enter a number (10 digits):")
	var number uint64
	for {
		fmt.Scan(&number)
		if number < 9999999999 && number > 1000000000 {
			p("The data input is done")
			break
		} else {
			p("The data input has been wrong. Please, repeat:")
		}
	}

	// получаем данные из файла
	tempData, err := list()
	if err != nil {
		p("Ошибка:", err)
		return
	}

	// Определяем свободный ключ
	var key uint64 = 1 // создаём переменную для ключа
	for {
		_, ok := tempData.Records[key]
		if !ok { // выходим из цикла если ключа нет, key сохраняется
			break
		}
		key++
	}

	// Записываем в файл базу
	tempData.Records[key] = Record{Name: name, Number: number}

	p("Текущее состояние телефонной книги:")
	printList(&tempData)
	saveToFile(&tempData)
}

func main() {

	p := log.Println // the alias for log.Print in order to simplify the code

	p("An address book example")

	// список доступных действий
	messages = [3]string{
		"Добавить пользователя в справочник",
		"Просмотр справочника",
		"Удалить пользователя из справчника",
	}

	for {
		p("--------------------------------")
		p("Список доступных действий:")
		for key, value := range messages {
			p(key+1, "-", value)
		}
		p("Определите ваше действие:")
		fmt.Scan(&action)
		switch action {
		case 1:
			add()
		case 2:
			{
				p("--------------------------------")
				tempData, err := list()
				if err == nil {
					printList(&tempData)
				} else {
					p("Ошибка:", err)
				}
			}
		case 3:
			deleteRecord()
		default:
			p("Неверный ввод. Повторите пожалуйста")
		}
	}
}
