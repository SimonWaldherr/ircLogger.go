package main

import (
	"bufio"
	"github.com/mxk/go-sqlite/sqlite3"
	"log"
	"net"
	"net/textproto"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type IRC struct {
	server        string
	port          string
	nick          string
	user          string
	channel       string
	pread, pwrite chan string
	conn          net.Conn
}

func NewIRC() *IRC {
	return &IRC{
		server:  os.Args[1],
		port:    os.Args[2],
		nick:    os.Args[4],
		channel: os.Args[3],
		user:    os.Args[4],
	}
}

func (bot *IRC) Connect() (conn net.Conn, err error) {
	conn, err = net.Dial("tcp", bot.server+":"+bot.port)
	if err != nil {
		log.Fatal("unable to connect to IRC server ", err)
	}
	bot.conn = conn
	log.Printf("Connected to IRC server %s (%s)\n", bot.server, bot.conn.RemoteAddr())
	return bot.conn, nil
}

func main() {
	if len(os.Args) < 5 {
		log.Fatal("this application needs 5 parameter to start\nExample: irc.freenode.net 6667 channel nickname\n")
	} else {
		var msgstr []string
		msgexp, _ := regexp.Compile(":([^!]+)!~([^\\s]+) ([^\\s]+) ([^\\s]+) :(.*)")
		ircbot := NewIRC()
		conn, _ := ircbot.Connect()
		channels := strings.Split(ircbot.channel, ",")
		conn.Write([]byte("USER " + ircbot.nick + " 8 * :" + ircbot.nick + "\r\n"))
		conn.Write([]byte("NICK " + ircbot.nick + "\r\n"))

		for i := 0; i < len(channels); i++ {
			conn.Write([]byte("JOIN #" + channels[i] + "\r\n"))
		}

		defer conn.Close()

		reader := bufio.NewReader(conn)
		tp := textproto.NewReader(reader)

		sql, _ := sqlite3.Open("./irclog.sqlite3")

		for {
			line, err := tp.ReadLine()
			if err != nil {
				break
			}
			msgstr = msgexp.FindStringSubmatch(line)
			if line[0:4] == "PING" {
				conn.Write([]byte("PONG " + line[5:] + "\r\n"))
			}
			if len(msgstr) == 6 {
				sql.Exec("INSERT INTO log (channel, user, message, timestamp) VALUES('" + msgstr[4] + "', '" + msgstr[1] + "', '" + msgstr[5] + "', '" + strconv.FormatInt(time.Now().Unix(), 10) + "')")
			}
		}
	}
}
