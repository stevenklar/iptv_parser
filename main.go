package main

import (
	"fmt"
	"os"
  "text/tabwriter"
  "github.com/fatih/color"
  "github.com/stevenklar/iptv_parser/pkg"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <host:port> <users_csv_file_path>\n", os.Args[0])
		os.Exit(1)
	}

	iptv := &pkg.IPTV{
		CsvFile: os.Args[2],
		Host:    os.Args[1],
	}

	users, err := iptv.GetUsers()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

  spacerColor := color.New(color.Bold, color.FgYellow).SprintFunc()

  fmt.Println(spacerColor("=================================================="))
  printUsersPretty(users)
  fmt.Println(spacerColor("=================================================="))
  printUsers(users)
  fmt.Println(spacerColor("=================================================="))
}

func printUsers(users []pkg.User) {
	fmt.Println("Login\tName\tExpires")
	for _, user := range users {
		fmt.Printf("%s\t%s\t", user.Login, user.Name)
		if user.Error != nil {
			fmt.Printf("%s\n", user.Error)
		} else {
			fmt.Printf("%s\n", user.Expires.Format("02.01.2006 15:04:05"))
		}
	}
}

func printUsersPretty(users []pkg.User) {
	// Set up a tab writer to align columns of text.
	w := tabwriter.NewWriter(os.Stdout, 4, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)

  // Define some colors for highlighting the header row.
  headerColor := color.New(color.Bold, color.FgHiGreen).SprintFunc()

  // Define some colors for highlighting errors.
  errorColor := color.New(color.Bold, color.FgRed).SprintFunc()

	// Write the header row.
	fmt.Fprintf(w, "%-30s | %-21s | %-25s\n", headerColor("Login"), headerColor("Name"), headerColor("Expires"))
	fmt.Fprintf(w, "%-19s | %-10s | %-25s\n", "-----", "-----", "-----")


	// Write each user's information to the tab writer.
	for _, user := range users {
		// Write the login and name columns.
		fmt.Fprintf(w, "%-19s | %-10s | ", user.Login, user.Name)

		// Write the expiration date or error message.
		if user.Error != nil {
			fmt.Fprintf(w, "%-25s\n", errorColor(user.Error))
		} else {
			fmt.Fprintf(w, "%-25s\n", user.Expires.Format("02.01.2006 15:04:05"))
		}
	}

	// Flush the tab writer to stdout.
	w.Flush()
}

