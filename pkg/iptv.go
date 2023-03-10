package pkg

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
  "strconv"
  "io"
)

// IPTV struct to hold the data
type IPTV struct {
	CsvFile string
	Host    string
}

type User struct {
	Login    string
	Name     string
	Password string
	Expires  time.Time
	Error    error
}

func (iptv *IPTV) GetUsers() ([]User, error) {
	usersFile, err := os.Open(iptv.CsvFile)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer usersFile.Close()

	reader := csv.NewReader(usersFile)
	reader.Comma = '\t'
	_, err = reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading csv file: %w", err)
	}

	var users []User
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, fmt.Errorf("error reading csv record: %w", err)
		}

		user := User{
			Login:    record[0],
			Name:     record[1],
			Password: record[2],
		}
		err = iptv.GetUserExpiration(&user)
		if err != nil {
			user.Error = err
		}
		users = append(users, user)
	}

	return users, nil
}

func (iptv *IPTV) GetUserExpiration(user *User) error {
	if user.Password == "" {
		return fmt.Errorf("missing password for user %s", user.Login)
	}

	url := fmt.Sprintf("http://%s/player_api.php?username=%s&password=%s", iptv.Host, user.Login, user.Password)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/35.0.1916.47 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.ContentLength == 0 {
		return fmt.Errorf("invalid content from API for user %s", user.Login)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var values map[string]interface{}
	err = json.Unmarshal(data, &values)
	if err != nil {
		return fmt.Errorf("unable to unmarshal data for user %s", user.Login)
	}

	userInfo := values["user_info"].(map[string]interface{})
	if userInfo["auth"].(float64) == 0 {
		return fmt.Errorf("unable to read user_info for user %s", user.Login)
	}

	tsStr := userInfo["exp_date"].(string)
	ts, err := strconv.ParseInt(tsStr, 10, 64)
  if err != nil {
    return fmt.Errorf("error parsing expiration date: %w", err)
  }

  user.Expires = time.Unix(ts, 0).UTC()
  return nil
}

