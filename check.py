#!python3
import csv
import json
import urllib.request
import sys
from datetime import datetime

# local imports
from helpers import color

class IPTV:
    """IPTV Helper class"""

    def __init__(self, host, file):
        self.csv = file
        self.host = host


    def print_users(self):
        with open(self.csv, newline='') as usersfile:
            reader = csv.DictReader(usersfile, delimiter="\t")
            for row in reader:
                self.print_user(row)

    def print_user(self, user):
        print(user['login'])
        print(color.OKBLUE + user['name'] + color.ENDC)

        if user['password'] == '':
            print(color.WARNING + "MISSING PASSWORD" + color.ENDC)
        else:
            url = "http://%s/player_api.php?username=%s&password=%s" % (self.host, user['login'], user['password'])
            response = urllib.request.urlopen(url)
            data = response.read()
            values = json.loads(data)
            user_info = values['user_info']
            if (user_info['auth'] == 0):
                print(color.FAIL + "INVALID PASSWORD" + color.ENDC)
                return

            ts = int(user_info['exp_date'])
            print(datetime.utcfromtimestamp(ts).strftime('%d.%m.%Y %H:%M:%S'))
        print("======================")

if len(sys.argv) != 3:
    exit("Usage: %s <host:port> <users_csv_file_path>" % sys.argv[0])


iptv = IPTV(sys.argv[1], sys.argv[2])
iptv.print_users()
