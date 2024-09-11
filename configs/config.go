package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const MaxDaysInRule int = 400
const DayInWeek int = 7
const CorrectLenForDays int = 3
const CorrectLenForWeek int = 3
const CorrectLenForMonth int = 2
const MaxDaysInMonth int = 31
const MaxCountdownDays int = -2
const MinMonth int = 1
const MaxMonth int = 12
const dirPerm = 0755

var DBFile string
var DBDir string
var WebDir string
var Port string
var DateLayout string = "20060102"
var SearchLayout string = "02.01.2006"

var logFile *os.File

func LoadEnv() {
	var exists bool

	if err := godotenv.Load("./.env"); err != nil {
		log.Println("не найден .env файл")
	}

	WebDir, exists = os.LookupEnv("WEB_DIR")
	if !exists {
		log.Fatal("не найдена переменная окружения WEB_DIR")
	}

	Port, exists = os.LookupEnv("TODO_PORT")
	if !exists {
		log.Fatal("не найдена переменная окружения TODO_PORT")
	}

	DBFile, exists = os.LookupEnv("TODO_DBFILE")
	if !exists {
		log.Fatal("не найдена переменная окружения TODO_DBFILE")
	}

	DBDir, exists = os.LookupEnv("TODO_DB_DIR")
	if !exists {
		log.Fatal("не найдена переменная окружения TODO_DB_DIR")
	}

	_, exists = os.LookupEnv("TODO_PASSWORD")
	if !exists {
		log.Fatal("не найдена переменная окружения TODO_PASSWORD")
	}

	_, exists = os.LookupEnv("APP_LOG_DIR")
	if !exists {
		log.Fatal("не найдена переменная окружения APP_LOG_DIR")
	}

	_, exists = os.LookupEnv("APP_LOG_FILE")
	if !exists {
		log.Fatal("не найдена переменная окружения APP_LOG_FILE")
	}
}

func StartLog() {
	logDir, _ := os.LookupEnv("APP_LOG_DIR")
	logFile, _ := os.LookupEnv("APP_LOG_FILE")

	if _, err := os.Stat(logDir); err != nil {
		if err := os.Mkdir(logDir, dirPerm); err != nil {
			log.Println("не удалось создать каталог для log-файла", err)
		}
	}
	l, err := os.OpenFile(logDir+logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("не удалось открыть файл ", err)
	}
	log.SetOutput(l)
	log.SetFlags(log.Ldate | log.Ltime)
	log.Println("логирование начато")
}

func StopLog() {
	log.Println("логирование окончено")
	logFile.Close()
}
