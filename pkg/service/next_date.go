package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// NextDate возвращает ближайшую дату повторения позже now.
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if len(repeat) > 128 {
		return "", fmt.Errorf("repeat rule must not exceed 128 characters")
	}
	start, err := ParseDate(dstart)
	if errutil.CheckError(err) {
		return "", err
	}
	if now.Year() < 1 || now.Year() > 9999 {
		return "", fmt.Errorf("date year must be between 1 and 9999")
	}

	rule := strings.Fields(repeat)
	if len(rule) == 0 {
		return "", fmt.Errorf("repeat rule is empty")
	}

	var next time.Time
	switch rule[0] {
	case "d":
		next, err = nextByDays(start, now, rule)
	case "y":
		next, err = nextByYears(start, now, rule)
	case "w":
		next, err = nextByWeekdays(start, now, rule)
	case "m":
		next, err = nextByMonths(start, now, rule)
	default:
		err = fmt.Errorf("unsupported repeat rule %q", rule[0])
	}
	if errutil.CheckError(err) {
		return "", err
	}
	if next.Year() > 9999 {
		return "", fmt.Errorf("next date is outside the supported range")
	}
	return next.Format(dateLayout), nil
}

// nextByDays проверяет интервал и вычисляет следующую дату по дням.
func nextByDays(start, now time.Time, rule []string) (time.Time, error) {
	if len(rule) != 2 {
		return time.Time{}, fmt.Errorf("day rule must have one interval")
	}
	interval, err := strconv.Atoi(rule[1])
	if errutil.CheckError(err) || interval < 1 || interval > 400 {
		return time.Time{}, fmt.Errorf("day interval must be between 1 and 400")
	}
	return nextByInterval(start, now, 0, interval), nil
}

// nextByYears вычисляет следующую дату, повторяющуюся раз в год.
func nextByYears(start, now time.Time, rule []string) (time.Time, error) {
	if len(rule) != 1 {
		return time.Time{}, fmt.Errorf("year rule takes no arguments")
	}
	return nextByInterval(start, now, 1, 0), nil
}

// nextByInterval пропускает прошедшие интервалы без перебора каждого дня.
func nextByInterval(start, now time.Time, years, days int) time.Time {
	if days > 0 {
		passedDays := (dateOnly(now).Unix() - start.Unix()) / 86400
		intervals := int(passedDays)/days + 1
		if intervals < 1 {
			intervals = 1
		}
		return start.AddDate(0, 0, intervals*days)
	}
	// Первый перенос сохраняет переход 29 февраля на 1 марта.
	date := start.AddDate(years, 0, 0)
	if skipped := now.Year() - date.Year(); skipped > 0 {
		date = date.AddDate(skipped, 0, 0)
	}
	if !afterNow(date, now) {
		date = date.AddDate(1, 0, 0)
	}
	return date
}

// nextByWeekdays ищет ближайший указанный день недели.
func nextByWeekdays(start, now time.Time, rule []string) (time.Time, error) {
	if len(rule) != 2 {
		return time.Time{}, fmt.Errorf("weekday rule must have one weekday list")
	}
	values, err := parseNumbers(rule[1], 1, 7)
	if errutil.CheckError(err) {
		return time.Time{}, err
	}

	var weekdays [8]bool
	for _, day := range values {
		weekdays[day] = true
	}

	date := repeatSearchStart(start, now)
	for {
		date = date.AddDate(0, 0, 1)
		if weekdays[weekdayNumber(date)] {
			return date, nil
		}
	}
}

// nextByMonths ищет ближайшую дату из заданных дней и месяцев.
func nextByMonths(start, now time.Time, rule []string) (time.Time, error) {
	if len(rule) < 2 || len(rule) > 3 {
		return time.Time{}, fmt.Errorf("month rule must have day and optional month lists")
	}
	values, err := parseNumbers(rule[1], -2, 31)
	if errutil.CheckError(err) {
		return time.Time{}, fmt.Errorf("invalid day list in month rule")
	}

	var days [32]bool
	var lastDay, penultimateDay bool
	for _, day := range values {
		switch day {
		case 0:
			return time.Time{}, fmt.Errorf("invalid day list in month rule")
		case -1:
			lastDay = true
		case -2:
			penultimateDay = true
		default:
			days[day] = true
		}
	}

	var months [13]bool
	for month := 1; month <= 12; month++ {
		months[month] = true
	}
	if len(rule) == 3 {
		values, err = parseNumbers(rule[2], 1, 12)
		if errutil.CheckError(err) {
			return time.Time{}, fmt.Errorf("invalid month list in month rule")
		}
		months = [13]bool{}
		for _, month := range values {
			months[month] = true
		}
	}
	if !monthRuleHasDate(days, months, lastDay, penultimateDay) {
		return time.Time{}, fmt.Errorf("month rule has no possible date")
	}

	date := repeatSearchStart(start, now)
	for {
		date = date.AddDate(0, 0, 1)
		if !months[int(date.Month())] {
			continue
		}
		day := date.Day()
		last := daysInMonth(date)
		if days[day] || (lastDay && day == last) || (penultimateDay && day == last-1) {
			return date, nil
		}
	}
}

// repeatSearchStart пропускает прошедшие дни перед поиском календарного повторения.
func repeatSearchStart(start, now time.Time) time.Time {
	if !afterNow(start.AddDate(0, 0, 1), now) {
		return dateOnly(now)
	}
	return start
}

// monthRuleHasDate проверяет, что правило может совпасть с календарной датой.
func monthRuleHasDate(days [32]bool, months [13]bool, lastDay, penultimateDay bool) bool {
	if lastDay || penultimateDay {
		return true
	}
	for month := 1; month <= 12; month++ {
		if !months[month] {
			continue
		}
		last := daysInMonth(time.Date(2000, time.Month(month), 1, 0, 0, 0, 0, time.UTC))
		for day := 1; day <= last; day++ {
			if days[day] {
				return true
			}
		}
	}
	return false
}

// daysInMonth возвращает количество дней в месяце даты.
func daysInMonth(date time.Time) int {
	return time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// parseNumbers разбирает числа через запятую и проверяет их границы.
func parseNumbers(value string, min, max int) ([]int, error) {
	items := strings.Split(value, ",")
	numbers := make([]int, 0, len(items))
	for _, item := range items {
		number, err := strconv.Atoi(item)
		if errutil.CheckError(err) || number < min || number > max {
			return nil, fmt.Errorf("invalid value %q", item)
		}
		numbers = append(numbers, number)
	}
	return numbers, nil
}

// afterNow сравнивает календарные даты без учёта времени суток.
func afterNow(date, now time.Time) bool {
	return dateOnly(date).After(dateOnly(now))
}

// dateOnly возвращает дату без времени суток.
func dateOnly(date time.Time) time.Time {
	year, month, day := date.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// weekdayNumber переводит номер дня недели в диапазон от 1 до 7.
func weekdayNumber(date time.Time) int {
	weekday := int(date.Weekday())
	if weekday == 0 {
		return 7
	}
	return weekday
}
