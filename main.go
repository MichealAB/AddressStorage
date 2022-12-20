package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	host          = os.Getenv("host")
	database      = os.Getenv("database")
	user          = os.Getenv("user")
	password      = os.Getenv("password")
	tgBotKey      = os.Getenv("TG_BOT_KEY")
	sku           string
	title         string
	place         string
	count         string
	id_sku        string
	output        string
	MessageChatId int
	id            int
	chat_id       int
	Text          string
)

type JsonFile struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateId int `json:"update_id"`
		Message  struct {
			MessageId int `json:"message_id"`
			From      struct {
				Id           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				Id        int    `json:"id"`
				FirstName string `json:"first_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date int    `json:"date"`
			Text string `json:"text"`
		} `json:"message"`
	} `json:"result"`
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	// todo: Просмотр ячейки и кол-во по скю. ДА. Поиск по "Скю 123456"
	// todo: Просмотр скю и кол-ва по ячейке. ДА. Поиск по "Ячейка 123/456"
	// todo: Добавление строк в stocks (скю, ячейка, кол-во). ДА. Поиск по "Добавить 123456 123/456 100"
	// todo: Удаление скю из ячейки. ДА. Поиск по "Удалить 123456 123/456"
	// todo: Изменение кол-ва скю в ячейке. ДА. Поиск по "Изменить 123456 123/456 100"
	// todo: Изменение ячейки у скю. ДА. Поиск по "Вывести 123456 123/456"
	// todo: Очистка всей ячейки. ДА. Поиск по "Удалить 123/456"
	TGBot()
	TGChat()
}

func TGBot() {
	readFile := ReadChatID()
	var telegramChatID []string
	_ = json.Unmarshal(readFile, &telegramChatID)
	for i := range telegramChatID {
		text := "База данных активна"
		url := "https://api.telegram.org/bot" + tgBotKey + "/sendMessage?chat_id=" + telegramChatID[i] + "&text=" + text
		method := "GET"
		payload := strings.NewReader(`{` + "	" + `"chat_id": 5658090622,` + "" + `"text": "Кнопки для товаров",` + "" + `"reply_markup": {` + "" + `"keyboard": [` + "" + `[` + "" + `{"text": "Просмотр скю"},` + "" + `{"text": "Просмотр ячейки"},` + "" + `{"text": "Добавить товар в ячейку"}` + "" + `],` + "" + `[` + "" + `{"text": "Удалить скю из ячейки"},` + "" + `
{"text": "Изменить кол-во в ячейке"},` + "" + `{"text": "Очистить ячейку"}` + "" + `]` + "" + `]` + "" + `}` + "" + `}`)
		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)
		if err != nil {
			fmt.Println(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))
	}
}

func TGChat() {
	MessageOffset := ReadUpdatesID()
	for true {
		url2 := "https://api.telegram.org/bot" + tgBotKey + "/getUpdates?offset=" + MessageOffset
		fmt.Println(url2)
		method := "GET"
		client := &http.Client{}
		req, err := http.NewRequest(method, url2, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(body))

		var A JsonFile
		_ = json.Unmarshal(body, &A)
		ResultStruct := A.Result
		for i := range ResultStruct {
			MessageUpdateId := A.Result[i].UpdateId
			if strconv.Itoa(MessageUpdateId) == MessageOffset {
				continue
			}

			//readFile := ReadMessage()
			//var telegramMessageList []JsonFile
			//var telegramMessage JsonFile
			//_ = json.Unmarshal(readFile, &telegramMessageList)
			//_ = json.Unmarshal(body, &telegramMessage)
			//telegramMessageList = append(telegramMessageList, telegramMessage)
			//telegramMessageListJson, _ := json.Marshal(telegramMessageList)
			//err = os.WriteFile("./message.json", telegramMessageListJson, 0666)
			//if err != nil {
			//	log.Fatalf("Unable to open file:", err)
			//}
			//fmt.Println(err)

			MessageOffset = strconv.Itoa(MessageUpdateId)
			MessageChatId = A.Result[i].Message.Chat.Id
			MessageText := A.Result[i].Message.Text
			WriteUpdatesID([]byte(MessageOffset))
			var output1 string

			//var B []JsonFile
			//readFile1 := ReadMessage()
			//err = json.Unmarshal(readFile1, &B)
			//if err != nil {
			//	fmt.Println("Не удалось преобразовать массив байт в структуру В", err)
			//	return
			//}
			//fmt.Println("Преобразовали массив в структуру")
			//for i, v := range B {
			//	fmt.Println(i)
			//	for _, v2 := range v.Result {
			//		fmt.Println(v2.Message.Chat.Id, v2.Message.Text)
			//	}
			//}

			//q := A.Result[i-1].Message.Text
			InsertTextInSql(A.Result[i].Message.Chat.Id, A.Result[i].Message.Text)
			FindText(A.Result[i].Message.Chat.Id)

			vvod := Find(MessageText)
			if vvod == "" {
				switch FindText(A.Result[i].Message.Chat.Id) {
				case "Просмотр скю":
					output1 = FindAddressAndCountOnSKU(MessageText)
				case "Просмотр ячейки":
					output1 = FindSKUonAddress(MessageText)
				case "Добавить товар в ячейку":
					words := strings.Split(MessageText, " ")
					for _, word := range words {
						fmt.Println(word)
					}
					output1 = CreateNewRowStocks(words[0], words[1], words[2])
				case "Удалить скю из ячейки":
					words := strings.Split(MessageText, " ")
					for _, word := range words {
						fmt.Println(word)
					}
					output1 = DeleteSKUFromPlace(words[0], words[1])
				case "Изменить кол-во в ячейке":
					words := strings.Split(MessageText, " ")
					for _, word := range words {
						fmt.Println(word)
					}
					output1 = ChangeCount(words[0], words[1], words[2])
				case "Очистить ячейку":
					output1 = ClearPlace(MessageText)
				default:
					output1 = "Неверный запрос"
				}
				url1 := "https://api.telegram.org/bot" + tgBotKey + "/sendMessage?chat_id=" + strconv.Itoa(MessageChatId) + "&text=" + url.QueryEscape(output1)
				fmt.Println(url1)
				method := "GET"
				client := &http.Client{}
				req, err := http.NewRequest(method, url1, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				res, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					return
				}
				defer res.Body.Close()
				body1, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(body1)

			} else if vvod == "Добро пожаловать" {
				readFile := ReadChatID()
				var telegramChatList []string
				_ = json.Unmarshal(readFile, &telegramChatList)
				Contains(telegramChatList, MessageChatId)
			} else {
			}

		}
		time.Sleep(1 * time.Second)
	}
}

func ReadUpdatesID() string {
	readUpdatesId, err := os.ReadFile("./updatesid.txt") //читаем текстовый файл с айдишниками
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return string(readUpdatesId)
}

func WriteUpdatesID(JsonUpdates []byte) bool {
	err := os.WriteFile("./updatesid.txt", JsonUpdates, 0666)
	if err != nil {
		return false
	}
	return true
}

func ReadChatID() []byte {
	readFile, err := os.ReadFile("./chatid.json") //читаем текстовый файл с айдишниками
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return readFile
}
func ReadMessage() []byte {
	readFile, err := os.ReadFile("./message.json") //читаем текстовый файл с айдишниками
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	return readFile
}

func Contains(telegramChatList []string, MessageChatId int) string {
	for _, n := range telegramChatList {
		if strconv.Itoa(MessageChatId) == n {
			fmt.Println(MessageChatId, "ID существует в базе")
			return ""
		}
	}
	telegramChatList = append(telegramChatList, strconv.Itoa(MessageChatId))
	fmt.Println(MessageChatId, "добавили новый ID")
	telegramChatListJson, _ := json.Marshal(telegramChatList)
	err := os.WriteFile("./chatid.json", telegramChatListJson, 0666)
	if err != nil {
		log.Fatalf("Unable to open file:", err)
	}
	fmt.Println(err)
	return ""
}

func FindAddressAndCountOnSKU(q string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Query("SELECT sku, title, place, count, id_sku from goods inner join stocks on sku = id_sku WHERE sku = ?;", q)
	checkError(err)
	defer rows.Close()
	fmt.Println("Reading data:")
	var output1 string
	for rows.Next() {
		err := rows.Scan(&sku, &title, &place, &count, &id_sku)
		checkError(err)
		output1 += "Ячейка:" + place + " Кол-во: " + count + "\n"
		fmt.Printf("Data row = (%d, %s, %v, %v, %v, %v)\n", sku, title, place, count, id_sku)
	}
	err = rows.Err()
	checkError(err)
	fmt.Println("Done.")
	if output1 == "" {
		return "Инфо не найдено/ошибка при вводе"
	} else {
		return "Инфо об скю: " + sku + " <" + title + ">" + "\n" + output1
	}
}
func FindSKUonAddress(q string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Query("SELECT sku, title, place, count, id_sku from goods inner join stocks on sku = id_sku WHERE place = ?;", q)
	checkError(err)
	defer rows.Close()
	fmt.Println("Reading data:")
	var output1 string
	for rows.Next() {
		err := rows.Scan(&sku, &title, &place, &count, &id_sku)
		checkError(err)
		output1 += "Арт.: " + sku + " <" + title + ">" + " Кол-во: " + count + "\n"
		fmt.Printf("Data row = (%d, %s, %v, %v, %v, %v)\n", sku, title, place, count, id_sku)
	}
	err = rows.Err()
	checkError(err)
	fmt.Println("Done.")
	if output1 == "" {
		return "Инфо не найдено/ошибка при вводе"
	} else {
		return "Инфо о ячейке: " + place + "\n" + output1
	}
}

func CreateNewRowStocks(sku, place, count string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Exec("INSERT INTO stocks (id_sku, place, count) VALUES (?, ?, ?);", sku, place, count)
	checkError(err)
	rowCount, err := rows.RowsAffected()
	fmt.Printf("Updated %d row(s) of data.\n", rowCount)
	fmt.Println("Done.")
	return "В базу внесена новая запись." + "\n" + "Арт.: " + sku + " Place " + place + " Count " + count
}

func DeleteSKUFromPlace(sku, place string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Exec("DELETE FROM stocks WHERE id_sku = ? AND place = ?", sku, place)
	checkError(err)
	rowCount, err := rows.RowsAffected()
	fmt.Printf("Deleted %d row(s) of data.\n", rowCount)
	fmt.Println("Done.")
	return "Скю " + sku + " удалена из ячейки " + place
}

func ChangeCount(sku, place, count string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Exec("UPDATE stocks SET count = ? WHERE id_sku = ? AND place = ?", count, sku, place)
	checkError(err)
	rowCount, err := rows.RowsAffected()
	fmt.Printf("Deleted %d row(s) of data.\n", rowCount)
	fmt.Println("Done.")
	return "Кол-во " + sku + " в ячейке " + place + "\nизменено на " + count
}

func ClearPlace(place string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Exec("DELETE FROM stocks WHERE place = ?", place)
	checkError(err)
	rowCount, err := rows.RowsAffected()
	fmt.Printf("Deleted %d row(s) of data.\n", rowCount)
	fmt.Println("Done.")
	return "Ячейка " + place + " очищена"
}

func Find(q string) string {
	var TextMassage string
	switch q {
	case "Просмотр скю":
		TextMassage = "Введите скю"
	case "Просмотр ячейки":
		TextMassage = "Введите номер ячейки"
	case "Добавить товар в ячейку":
		TextMassage = "Введите скю, ячейку, кол-во через пробел"
	case "Удалить скю из ячейки":
		TextMassage = "Введите скю и ячейку через пробел"
	case "Изменить кол-во в ячейке":
		TextMassage = "Введите скю, ячейку и итоговое кол-во через пробел"
	case "Очистить ячейку":
		TextMassage = "Введите номер ячейки"
	case "/start":
		TextMassage = "Добро пожаловать"
	default:
		TextMassage = ""
	}
	url := "https://api.telegram.org/bot" + tgBotKey + "/sendMessage?chat_id=" + strconv.Itoa(MessageChatId) + "&text=" + TextMassage
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	fmt.Println(string(body))
	return TextMassage
}

func InsertTextInSql(ChatId int, Text string) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Exec("INSERT INTO ChatMessage (chat_id, Text) VALUES (?, ?);", ChatId, Text)
	checkError(err)
	rowCount, err := rows.RowsAffected()
	fmt.Printf("Updated %d row(s) of data.\n", rowCount)
	fmt.Println("Done.")
	return "Новое сообщение записано\n" + "Чат ИД: " + string(ChatId) + "\nТекст: " + Text
}

func FindText(MessageChatId int) string {
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?allowNativePasswords=true", user, password, host, database)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()
	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database.")
	rows, err := db.Query("SELECT * from ChatMessage where chat_id=? order by id desc limit 1,1;", MessageChatId)
	checkError(err)
	defer rows.Close()
	fmt.Println("Reading data:")
	for rows.Next() {
		err := rows.Scan(&id, &chat_id, &Text)
		checkError(err)
		fmt.Printf("Data row =", Text)
	}
	err = rows.Err()
	checkError(err)
	fmt.Println("Done.")
	return Text
}
