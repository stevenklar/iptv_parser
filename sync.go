package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "net/http"
    "encoding/json"
    "io/ioutil"

    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/sheets/v4"
    "github.com/stevenklar/iptv_parser/pkg"
)

func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
    tokFile := "token.json"
    tok, err := tokenFromFile(tokFile)
    if err != nil {
        tok = getTokenFromWeb(config)
        saveToken(tokFile, tok)
    }
    return config.Client(ctx, tok), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()
    tok := &oauth2.Token{}
    err = json.NewDecoder(f).Decode(tok)
    return tok, err
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
    authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
    fmt.Printf("Go to the following link in your browser then type the "+
        "authorization code: \n%v\n", authURL)

    var authCode string
    if _, err := fmt.Scan(&authCode); err != nil {
        log.Fatalf("Unable to read authorization code: %v", err)
    }

    tok, err := config.Exchange(context.Background(), authCode)
    if err != nil {
        log.Fatalf("Unable to retrieve token from web: %v", err)
    }
    return tok
}

func saveToken(file string, token *oauth2.Token) {
    fmt.Printf("Saving credential file to: %s\n", file)
    f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
    if err != nil {
        log.Fatalf("Unable to cache oauth token: %v", err)
    }
    defer f.Close()
    json.NewEncoder(f).Encode(token)
}

func getSheetByID(srv *sheets.Service, spreadsheetID string) (*sheets.Spreadsheet, error) {
    // Construct the request.
    spreadsheet, err := srv.Spreadsheets.Get(spreadsheetID).Do()
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve spreadsheet: %v", err)
    }
    return spreadsheet, nil
}

func updateSheetCell(srv *sheets.Service, spreadsheetID string, sheetName string, row int, col int, newValue interface{}) error {
    // Define the range of the cell to update.
    rangeToUpdate := fmt.Sprintf("%s!%s%d", sheetName, columnToLetter(col), row)

    // Create the value range object with the new value.
    valueRange := &sheets.ValueRange{
        Values: [][]interface{}{{newValue}},
    }

    // Call the Google Sheets API to update the value in the cell.
    _, err := srv.Spreadsheets.Values.Update(spreadsheetID, rangeToUpdate, valueRange).ValueInputOption("USER_ENTERED").Do()
    if err != nil {
        return fmt.Errorf("unable to update cell: %v", err)
    }

    return nil
}

// Helper function to convert a column number to a letter (e.g. 1 -> A, 2 -> B, etc.).
func columnToLetter(col int) string {
    letter := ""
    for col > 0 {
        col--
        letter = string('A'+col%26) + letter
        col /= 26
    }
    return letter
}

func readSheetData(srv *sheets.Service, spreadsheetID string, sheetName string) ([][]string, error) {
    // Define the range of cells to read.
    readRange := sheetName + "!A2:D"

    // Call the Google Sheets API to retrieve the values in the range.
    resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
    if err != nil {
        return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
    }

    // Convert the values to a 2D slice of strings.
    var values [][]string
    for _, row := range resp.Values {
        var rowValues []string
        for _, col := range row[:4] {
            s, ok := col.(string)
            if !ok {
                s = fmt.Sprint(col)
            }
            rowValues = append(rowValues, s)
        }
        values = append(values, rowValues)
    }

    return values, nil
}

func main() {
    if len(os.Args) != 4 {
      fmt.Printf("Usage: %s <host:port> <spreadsheetID> <sheetName>\n", os.Args[0])
      os.Exit(1)
    }

    host := os.Args[1]
    spreadsheetID := os.Args[2]
    sheetName := os.Args[3]

    iptv := &pkg.IPTV{
      CsvFile: "",
      Host:    host,
    }

    b, err := ioutil.ReadFile("credentials.json")
    if err != nil {
        log.Fatalf("Unable to read client secret file: %v", err)
    }

    // If modifying these scopes, delete your previously saved token.json.
    config, err := google.ConfigFromJSON(b, sheets.SpreadsheetsScope)
    if err != nil {
        log.Fatalf("Unable to parse client secret file to config: %v", err)
    }
    client, err := getClient(context.Background(), config)
    if err != nil {
        log.Fatalf("Unable to get client: %v", err)
    }

    srv, err := sheets.New(client)

    spreadsheet, err := getSheetByID(srv, spreadsheetID)
    if err != nil {
        log.Fatalf("Unable to retrieve spreadsheet: %v", err)
    }
    fmt.Printf("Spreadsheet title: %s\n", spreadsheet.Properties.Title)

    values, err := readSheetData(srv, spreadsheetID, sheetName)
    if err != nil {
        log.Fatalf("Unable to retrieve data from sheet: %v", err)
    }

    // Print the values in the first 3 columns of each row.
    for row, column := range values {
        realRow := row + 2

        user := pkg.User{
          Login:    column[0],
          Name:     column[1],
          Password: column[2],
        }

        if column[3] != "<nil>" {
            log.Printf("Skip user expiration (%v): %v", realRow, user.Login)
            continue
        }

        err := iptv.GetUserExpiration(&user)
        if err != nil {
			user.Error = err
            log.Printf("Unable to retrieve data from iptv host: %v", user.Error)
            continue
        }

        log.Printf("Update user expiration (%v): %v", realRow, user.Login)
        updateSheetCell(srv, spreadsheetID, sheetName, realRow, 4, user.Expires.Format("02.01.2006 15:04:05"))
    }
}
