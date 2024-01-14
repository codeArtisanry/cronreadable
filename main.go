package main

import (
	"encoding/csv"
	"flag"
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
	var cronExpressions []string
	var fromFile bool

	flag.BoolVar(&fromFile, "file", false, "Read CRON expressions from file")
	flag.Parse()

	if fromFile {
		// If -file flag is provided, read CRON expressions from file
		fileName := "cron_expressions.txt"
		file, err := os.Open(fileName)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		// Read the content from the file
		buffer := make([]byte, 1024)
		n, err := file.Read(buffer)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		cronExpressions = strings.Fields(string(buffer[:n]))
	} else {
		// If -file flag is not provided, use command line arguments as CRON expressions
		cronExpressions = flag.Args()
	}

	if len(cronExpressions) == 0 {
		fmt.Println("Please provide CRON expressions using command line arguments or use -file flag to read from file.")
		return
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

		// Append each row separately
		for _, time := range nextScheduledTimes {
			row = append(row, "", time)
		}

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
