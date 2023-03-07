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
  var output string

  output += user["login"] + "\t"
  output += user["name"] + "\t"

	if user["password"] == "" {
		output += "MISSING PASSWORD"
	} else {
		url := fmt.Sprintf("http://%s/player_api.php?username=%s&password=%s", iptv.host, user["login"], user["password"])
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

    if resp.ContentLength == 0 {
      output += "INVALID_CONTENT_API"
      fmt.Println(output)
			return
    }

    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      output += "INVALID_DATA_API"
      fmt.Println(output)
			return
    }

    var values map[string]interface{}
    err = json.Unmarshal(data, &values)
    if err != nil {
      output += "INVALID_PASSWORD_API"
      fmt.Println(output)
			return
    }

		// var values map[string]interface{}
		// err = json.NewDecoder(resp.Body).Decode(&values)
		// if err != nil {
		// 	panic(err)
		// }

		userInfo := values["user_info"].(map[string]interface{})
		if userInfo["auth"].(float64) == 0 {
			output += "INVALID PASSWORD"
      fmt.Println(output)
			return
		}

    tsStr := userInfo["exp_date"].(string)
    ts, err := strconv.ParseInt(tsStr, 10, 64)
    if err != nil {
      panic(err)
    }
    output += time.Unix(ts, 0).UTC().Format("02.01.2006 15:04:05") + "\t"
	}

	fmt.Println(output)
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

