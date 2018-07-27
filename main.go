package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	ch := make(chan string)

	//open and read a file of bill numbers
	file, err := os.Open("./bills.txt")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer file.Close()

	var bills []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		bills = append(bills, scanner.Text())
		fmt.Println(scanner.Text())
	}

	fmt.Println("Length: " + strconv.Itoa(len(bills)))
	start := time.Now()
	for i := 0; i < len(bills); i++ {
		fmt.Println(i)
		go getTMWin(i, bills[i], ch)
		time.Sleep(2000 * time.Millisecond)
		fmt.Println(time.Since(start))
	}
	for i := 0; i < 2000; i++ {
		fmt.Println(<-ch)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

}
func getTMWin(count int, billNumber string, ch chan<- string) {
	// Make HTTP GET request
	response, err := http.Get("https://mydaylightupgrd.dylt.com/trace/external_bill_viewer.msw?trace_type=~PTLORDER&search_value=" + billNumber)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Copy data from the response to standard output
	n, err := io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Number of bytes copied to STDOUT:", n)
}

func getToken() string {
	//Consumer Key: x5Vxusddiy2pYqwpZytwxqkG0lW7Z6a5
	//Consumer Secret: ThzO25vxF0RDuA2U
	body := strings.NewReader(`client_secret=ThzO25vxF0RDuA2U&grant_type=client_credentials&client_id=x5Vxusddiy2pYqwpZytwxqkG0lW7Z6a5`)
	req, err := http.NewRequest("POST", "https://api.dylt.com/oauth/client_credential/accesstoken?grant_type=client_credentials", body)
	if err != nil {
		// handle err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var data map[string]string
	json.Unmarshal(token, &data)
	//fmt.Println(data)
	//fmt.Println(data["access_token"])
	return data["access_token"]
}
