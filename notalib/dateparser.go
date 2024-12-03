package notalib

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timers = map[string]string{
		"morning":   "08:00:00",
		"afternoon": "12:00:00",
		"evening":   "18:00:00",
		"default":   "06:00:00",
	}
	dayers = []string{"tomorrow", "today"}
)

func ParseDateTime(input string) (*time.Time, string) {
	if input == "" {
		return nil, ""
	}

	today := time.Now()
	date := today

	parts := strings.Fields(strings.ReplaceAll(input, "+", " "))
	dayPart, timePart := "", ""

	if len(parts) > 0 {
		dayPart = parts[0]
	}
	if len(parts) > 1 {
		timePart = parts[1]
	}

	if isRelativeDate(input) {
		relativeDate, err := parseRelativeDate(dayPart, date)
		if err != nil {
			log.Fatal(err)
		}
		date = relativeDate
	} else if contains(dayers, dayPart) {
		if dayPart == "tomorrow" {
			date = today.AddDate(0, 0, 1)
		}
	} else if dayPart == "now" {
		timePart = date.Format("15:04:05")
	} else if contains(keys(timers), dayPart) {
		timePart = dayPart
	} else if isISODate(dayPart) {
		_date, err := time.Parse("2006-01-02", dayPart)
		if err != nil {
			log.Fatal(err)
		}
		date = _date
	} else {
		log.Print("no date found")
		return nil, input 
	}

	dateStr := date.Format("2006-01-02")
	timeStr := timers["default"]

	if timePart != "" && contains(keys(timers), timePart) {
		timeStr = timers[timePart]
	} else if isISOTime(timePart) {
		parts := strings.Split(timePart, ":")
		timeStr = "";
		for i := 0; i < 3; i++ {
			var part = "0";
			if(i < len(parts)){
				part = parts[i];
			}
			timeStr = fmt.Sprintf("%s%02s:", timeStr, part)
		}
		timeStr = timeStr[:len(timeStr)-1]
	}
	final_date, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", dateStr, timeStr))
	if err != nil {
		log.Fatal(err)
	}
	return &final_date, "";
}

func parseRelativeDate(input string, baseDate time.Time) (time.Time, error) {
	re := regexp.MustCompile(`^\+?(\d+)?(day|days|week|weeks)$`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return baseDate, fmt.Errorf("invalid relative date")
	}

	amountStr := matches[1]
	unit := matches[2]

	amount := 1
	if amountStr != "" {
		var err error
		amount, err = strconv.Atoi(amountStr)
		if err != nil {
			return baseDate, err
		}
	}

	switch unit {
	case "day", "days":
		return baseDate.AddDate(0, 0, amount), nil
	case "week", "weeks":
		return baseDate.AddDate(0, 0, amount*7), nil
	default:
		return baseDate, fmt.Errorf("unknown unit")
	}
}

func isISODate(input string) bool {
	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	return re.MatchString(input)
}

func isISOTime(input string) bool {
	re := regexp.MustCompile(`^\d{1}[\d:]{0,}$`)
	return re.MatchString(input)
}

func constructISODate(date, time string) string {
	if time == "" {
		return fmt.Sprintf("%s 00:00:00", date)
	}
	if matched, _ := regexp.MatchString(`^\d{2}:\d{2}$`, time); matched {
		return fmt.Sprintf("%s %s:00", date, time)
	}
	if matched, _ := regexp.MatchString(`^\d{2}:\d{2}:\d{2}$`, time); matched {
		return fmt.Sprintf("%s %s", date, time)
	}
	return ""
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func keys(m map[string]string) []string {
	var result []string
	for k := range m {
		result = append(result, k)
	}
	return result
}

func isRelativeDate(input string) bool {
	re := regexp.MustCompile(`^\+?\d*(day|days|week|weeks)`)
	return re.MatchString(input)
}
