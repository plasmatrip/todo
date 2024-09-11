package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"todo/configs"
	"todo/model"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	referenceDate, err := time.Parse(configs.DateLayout, date)
	if err != nil {
		log.Printf("%s [date=%s]", err.Error(), date)
		return "", fmt.Errorf("переданное значение не может быть преобразовано в дату")
	}

	if len(repeat) == 0 {
		log.Printf("правило повторения задачи пустое")
		return "", fmt.Errorf("правило повторения задачи пустое")
	}

	switch repeat[0] {
	case 'd':
		if len([]rune(repeat)) < configs.CorrectLenForDays {
			log.Printf("неправильный формат правила повторения задачи [repeat=%s]", repeat)
			return "", fmt.Errorf("неправильный формат правила повторения задачи")
		}

		rule := strings.Split(repeat, " ")

		days, err := strconv.Atoi(rule[1])
		if err != nil {
			log.Printf("недопустимый символ в интервале повторения задачи [repeat=%s]", repeat)
			return "", fmt.Errorf("недопустимый символ в интервале повторения задачи")
		}

		if days > configs.MaxDaysInRule {
			log.Printf("превышен максимально допустимый интервал в правиле повторения задачи [repeat=%s]", repeat)
			return "", fmt.Errorf("превышен максимально допустимый интервал в правиле повторения задачи")
		}

		for referenceDate = referenceDate.AddDate(0, 0, days); referenceDate.Before(now); {
			referenceDate = referenceDate.AddDate(0, 0, days)
		}
	case 'y':
		for referenceDate = referenceDate.AddDate(1, 0, 0); referenceDate.Before(now); {
			referenceDate = referenceDate.AddDate(1, 0, 0)
		}
	case 'w':
		if len([]rune(repeat)) < configs.CorrectLenForWeek {
			log.Printf("неправильный формат правила повторения задачи [repeat=%s]", repeat)
			return "", fmt.Errorf("неправильный формат правила повторения задачи")
		}

		rule := strings.Split(repeat, " ")
		weekdays := strings.Split(rule[1], ",")
		targetDays := []time.Time{}

		var add int
		for _, weekday := range weekdays {
			day, err := strconv.Atoi(weekday)
			if err != nil {
				log.Printf("недопустимый символ в интервале повторения задачи [repeat=%s]", repeat)
				return "", fmt.Errorf("недопустимый символ в интервале повторения задачи")
			}

			if day > configs.DayInWeek {
				log.Printf("недопустимое значени в интервале повторения задачи [repeat=%s]", repeat)
				return "", fmt.Errorf("недопустимое значени в интервале повторения задачи")
			}

			if day <= int(referenceDate.Weekday()) {
				add = configs.DayInWeek - int(referenceDate.Weekday()) + day
			} else {
				add = day - int(referenceDate.Weekday())
			}
			targetDays = append(targetDays, referenceDate.AddDate(0, 0, add))
		}
	outW:
		for {
			for i := range targetDays {
				if targetDays[i].After(now) {
					referenceDate = targetDays[i]
					break outW
				}
				targetDays[i] = targetDays[i].AddDate(0, 0, configs.DayInWeek)
			}
		}
	case 'm':
		ruleD := []string{}
		ruleM := []string{}
		days := []int{}
		months := []int{}

		rule := strings.Split(repeat, " ")
		if len(rule) < configs.CorrectLenForMonth {
			log.Printf("неправильный формат правила повторения задачи [repeat=%s]", repeat)
			return "", fmt.Errorf("неправильный формат правила повторения задачи")
		}

		for i, value := range rule {
			switch i {
			case 1:
				ruleD = append(ruleD, strings.Split(value, ",")...)
			case 2:
				ruleM = append(ruleM, strings.Split(value, ",")...)
			}
		}

		for _, d := range ruleD {
			day, err := strconv.Atoi(d)
			if err != nil {
				log.Printf("недопустимый символ в интервале повторения задачи [repeat=%s]", repeat)
				return "", fmt.Errorf("недопустимый символ в интервале повторения задачи")
			}

			if day > configs.MaxDaysInMonth || day == 0 || day < configs.MaxCountdownDays {
				log.Printf("недопустимое значени в интервале повторения задачи [repeat=%s]", repeat)
				return "", fmt.Errorf("недопустимое значени в интервале повторения задачи")
			}

			days = append(days, day)
		}

		for _, m := range ruleM {
			month, err := strconv.Atoi(m)
			if err != nil {
				log.Printf("недопустимый символ в интервале повторения задачи [repeat=%s]", repeat)
				return "", fmt.Errorf("недопустимый символ в интервале повторения задачи")
			}

			if month < configs.MinMonth || month > configs.MaxMonth {
				log.Printf("недопустимое значени в интервале повторения задачи [repeat=%s]", repeat)
				return "", fmt.Errorf("недопустимое значени в интервале повторения задачи")
			}

			months = append(months, month)
		}

		targetDays := []time.Time{}
		var next time.Time
		if len(months) == 0 {
			for _, day := range days {
				if day < 0 {
					next = time.Date(referenceDate.Year(), referenceDate.Month()+1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, day)
					if next.Before(referenceDate) {
						next = time.Date(referenceDate.Year(), referenceDate.Month()+2, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, day)
					}
				} else {
					if day <= referenceDate.Day() {
						next = time.Date(referenceDate.Year(), referenceDate.Month()+1, day, 0, 0, 0, 0, time.UTC)
					} else {
						next = time.Date(referenceDate.Year(), referenceDate.Month(), day, 0, 0, 0, 0, time.UTC)
					}
					if next.Day() != day {
						next = time.Date(referenceDate.Year(), next.Month(), day, 0, 0, 0, 0, time.UTC)
					}
				}
				targetDays = append(targetDays, next)
			}
			sort.Slice(targetDays, func(i, j int) bool { return targetDays[i].Before(targetDays[j]) })
		outM1:
			for {
				for i := range targetDays {
					if targetDays[i].After(now) {
						referenceDate = targetDays[i]
						break outM1
					}
					targetDays[i] = targetDays[i].AddDate(0, 1, 0)
				}
			}
		} else {
			for _, month := range months {
				for _, day := range days {
					if day < 0 {
						next = time.Date(referenceDate.Year(), time.Month(month+1), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, day)
					} else {
						next = time.Date(referenceDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.UTC)
						if next.Day() != day {
							next = time.Date(referenceDate.Year(), next.Month(), day, 0, 0, 0, 0, time.UTC)
						}
					}
					targetDays = append(targetDays, next)
				}
			}
			sort.Slice(targetDays, func(i, j int) bool { return targetDays[i].Before(targetDays[j]) })
		outM2:
			for {
				for i := range targetDays {
					if targetDays[i].After(now) && checkMonth(int(targetDays[i].Month()), months) {
						referenceDate = targetDays[i]
						break outM2
					}
					targetDays[i] = targetDays[i].AddDate(0, 1, 0)
				}
			}
		}

	default:
		log.Printf("неподдерживаемый формат [repeat=%s]", repeat)
		return "", fmt.Errorf("неподдерживаемый формат")
	}

	return referenceDate.Format(configs.DateLayout), nil
}

func checkMonth(curr int, target []int) bool {
	for _, value := range target {
		if curr == value {
			return true
		}
	}
	return false
}

func ErrorResponse(w http.ResponseWriter, message string, err error) {
	error := model.Error{Message: fmt.Errorf("%s [%s]", message, err.Error()).Error()}
	resp, err := json.Marshal(error)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CheckTask(task *model.Task) error {
	if len(task.Id) > 0 {
		if _, err := strconv.Atoi(task.Id); err != nil {
			return errors.New("идентификатор не может быть преобразован в число")
		}
	}

	if len(task.Title) == 0 {
		return errors.New("заголовок задачи не может быть пустым")
	}

	now := time.Now()

	if len(task.Date) == 0 {
		task.Date = now.Format(configs.DateLayout)
	} else {
		date, err := time.Parse(configs.DateLayout, task.Date)
		if err != nil {
			return errors.New("дата представлена в формате, отличном от 20060102")
		}

		if date.Before(now) {
			task.Date = now.Format(configs.DateLayout)
		}
	}

	if len(task.Repeat) > 0 {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}
	return nil
}
