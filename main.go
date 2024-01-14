package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func cronToHumanReadable(cronExpression string) (string, error) {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)

	schedule, err := parser.Parse(cronExpression)
	if err != nil {
		return "", err
	}

	// Get the next 5 scheduled times
	var scheduleTimes []string
	now := schedule.Next(time.Now())
	for i := 0; i < 5; i++ {
		nextTime := schedule.Next(now)
		scheduleTimes = append(scheduleTimes, nextTime.Format("2006-01-02 15:04:05"))
		now = nextTime
	}

	// Extract human-readable details
	humanReadable := fmt.Sprintf("Expression: %s\n", cronExpression)
	humanReadable += fmt.Sprintf("Next 5 scheduled times:\n%s\n", strings.Join(scheduleTimes, "\n"))

	return humanReadable, nil
}

func writeToCSV(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cronExpressions := []string{
		"0 * * * *",        // Every hour
		"0 0 * * *",        // Every day at midnight
		"0 0 1,15 * *",     // Every 1st and 15th of the month
		"0 8-17 * * 1-5",   // Every day at 8 AM to 5 PM from Monday through Friday
		"*/15 * * * *",     // Every 15 minutes
		"0 12 * JAN,MAR *", // Every day at noon in January and March
		"0 0 1-7 * MON",    // First 7 days of the month if it's Monday
	}

	var csvData [][]string
	for _, cronExpression := range cronExpressions {
		humanReadable, err := cronToHumanReadable(cronExpression)
		if err != nil {
			fmt.Printf("Error parsing CRON expression '%s': %v\n", cronExpression, err)
			continue
		}

		fmt.Println(humanReadable)

		// Create CSV data
		row := []string{"CRON Expression", "Next Scheduled Time 1", "Next Scheduled Time 2", "Next Scheduled Time 3", "Next Scheduled Time 4", "Next Scheduled Time 5"}
		row = append(row, cronExpression)

		schedule, err := cronToHumanReadable(cronExpression)
		if err != nil {
			fmt.Printf("Error parsing CRON expression '%s': %v\n", cronExpression, err)
			continue
		}

		nextScheduledTimes := strings.Split(schedule, "\n")[2:]
		row = append(row, nextScheduledTimes...)

		csvData = append(csvData, row)
	}

	// Write to CSV file
	err := writeToCSV("cron_schedule.csv", csvData)
	if err != nil {
		fmt.Println("Error writing to CSV:", err)
		return
	}

	fmt.Println("CSV file created successfully.")
}
