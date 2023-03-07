package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
  "io/ioutil"
  "strconv"
)

// IPTV struct to hold the data
type IPTV struct {
	csvFile string
	host    string
}

// printUsers function to read csv file and print user info
func (iptv *IPTV) printUsers() {
	usersFile, err := os.Open(iptv.csvFile)
	if err != nil {
		panic(err)
	}
	defer usersFile.Close()

	reader := csv.NewReader(usersFile)
	reader.Comma = '\t'
  _, err = reader.Read()
  if err != nil {
    panic(err)
  }

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		user := map[string]string{
			"login":    record[0],
			"name":     record[1],
			"password": record[2],
		}
		iptv.printUser(user)
	}
}

// printUser function to print a single user's info
func (iptv *IPTV) printUser(user map[string]string) {
	fmt.Println(user["login"])
	fmt.Printf("\033[34m%s\033[0m\n", user["name"])

	if user["password"] == "" {
		fmt.Printf("\033[33m%s\033[0m\n", "MISSING PASSWORD")
	} else {
		url := fmt.Sprintf("http://%s/player_api.php?username=%s&password=%s", iptv.host, user["login"], user["password"])
		fmt.Println(url)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.47 Safari/537.36")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var values map[string]interface{}
    err = json.Unmarshal(data, &values)
    if err != nil {
        panic(err)
    }

		// var values map[string]interface{}
		// err = json.NewDecoder(resp.Body).Decode(&values)
		// if err != nil {
		// 	panic(err)
		// }

		userInfo := values["user_info"].(map[string]interface{})
		if userInfo["auth"].(float64) == 0 {
			fmt.Printf("\033[31m%s\033[0m\n", "INVALID PASSWORD")
			return
		}

    tsStr := userInfo["exp_date"].(string)
    ts, err := strconv.ParseInt(tsStr, 10, 64)
    if err != nil {
      panic(err)
    }
    fmt.Println(time.Unix(ts, 0).UTC().Format("02.01.2006 15:04:05"))
	}

	fmt.Println("======================")
}

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <host:port> <users_csv_file_path>\n", os.Args[0])
		os.Exit(1)
	}

	iptv := &IPTV{
		csvFile: os.Args[2],
		host:    os.Args[1],
	}

	iptv.printUsers()
}

