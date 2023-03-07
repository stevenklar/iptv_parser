# IPTV Helper

The IPTV Helper is a command-line tool that helps you manage IPTV user accounts. Given a CSV file containing user login information, the tool checks each user's account status and prints the results to the console.

## Installation

To install the IPTV Helper, you'll need to have Go installed on your system. Once you have Go installed, you can download and install the tool using the following command:
```sh
go install github.com/stevenklar/iptv_parser
```

This will download and install the tool into your `$GOPATH/bin` directory.

## Usage

To use the IPTV Helper, you'll need to have a CSV file containing user login information. Each row in the CSV file should contain three columns, separated by tabs:

`<login>\t<name>\t<password>`

You can then run the IPTV Helper with the following command:

```
iptv_parser host:port <csv_file>
```

Here, `<host:port>` is the address of your IPTV server (e.g. `xtreme-provider:8080`), and `<csv-file>` is the path to your CSV file.

## Output

The IPTV Helper prints the status of each user's account to the console. If the user's account is valid, the tool prints the expiration date of the account. If the user's account is invalid (e.g. the password is incorrect), the tool prints a warning message.

In addition to printing output to the console, you can also redirect the output to a file using the `>` operator, like so:

```
iptv_parser host:port <csv_file> > output.txt
```

This will write the program's output to a file called `output.txt`, which you can then open in a text editor or import into a Google Sheet.

## Contributing

If you'd like to contribute to the IPTV Helper, feel free to fork the project and submit a pull request. Before submitting a pull request, please make sure that your changes are thoroughly tested and documented.

## License

The IPTV Helper is open source software licensed under the MIT License. See the LICENSE file for more details.

