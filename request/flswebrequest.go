package request

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	base       = "https://fundingletter.substack.com/p/the-funding-letter-" //572-january-26
	startMonth = 0                                                          // 2019: 4, 2020: 1                                                         // 2019: 4, 2020: 1
	// start count = 108 for apr16 2019, 296 for jan2 2020
	startDate = 0                            // 2019: 16, 2020: 2
	directory = "webpages/fundingletter/%d/" // TODO: config year
)

func BatchFetchWebpages(year int, startCount int) *html.Node {
	if year == 2019 {
		startMonth = 4
		startDate = 16
	} else if year == 2020 {
		startMonth = 1
		startDate = 2
	} else {
		fmt.Println("Year not valid to collect info from")
		os.Exit(1)
	}
	directory = fmt.Sprintf(directory, year)
	fmt.Println(directory)

	generatePaths(year, startCount)

	return nil
}

func generatePaths(year int, startCount int) {
	month := time.Month(startMonth)
	start := time.Date(year, month, startDate, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, time.Month(1), 0, 0, 0, 0, 0, time.UTC)

	for ; start.Before(end); start = start.AddDate(0, 0, 1) {
		if start.Weekday().String() == "Saturday" || start.Weekday().String() == "Sunday" {
			//do nothing
		} else {
			s := buildString(startCount, start.Month(), start.Day())
			time.Sleep(1 * time.Minute)
			wb := fetchPage(s)
			if wb == 200 {
				fmt.Println(wb)
			} else {
				fmt.Printf("%d, for %s\n", wb, s)
			}
			startCount = startCount + 1
		}
	}
}

func fetchPage(s string) int {
	url := base + s

	user_agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:85.0) Gecko/20100101 Firefox/85.0"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", user_agent)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Warn: can not get webpage for %s\n", url)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	writeBodyToFile(body, s)
	return resp.StatusCode
}

func writeBodyToFile(body []byte, name string) {
	filename := directory + name + ".txt"
	err := ioutil.WriteFile(filename, body, 0644)
	if err != nil {
		panic(err)
	}
}

func buildString(number int, month time.Month, date int) string {
	m := strings.ToLower(month.String())
	n := strconv.Itoa(number)
	d := strconv.Itoa(date)

	result := []string{n, m, d}
	return strings.Join(result, "-")
}
