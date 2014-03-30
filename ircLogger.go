package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"os"
	"regexp"
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
		channel: "#" + os.Args[3],
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
	if len(os.Args) != 5 {
		log.Fatal("this application needs 4 parameter to start\nExample: irc.freenode.net 6667 channel nickname\n")
	} else {
		var msgstr []string
		msgexp, _ := regexp.Compile(":([^!]+)!~([^\\s]+) ([^\\s]+) ([^\\s]+) :(.*)")
		ircbot := NewIRC()
		conn, _ := ircbot.Connect()
		conn.Write([]byte("USER " + ircbot.nick + " 8 * :" + ircbot.nick + "\r\n"))
		conn.Write([]byte("NICK " + ircbot.nick + "\r\n"))
		conn.Write([]byte("JOIN " + ircbot.channel + "\r\n"))
		defer conn.Close()

		f, _ := os.OpenFile("./log.tsv", os.O_APPEND|os.O_WRONLY, 0666)
		defer f.Close()
		f.WriteString("User	Channel	Message\n")

		reader := bufio.NewReader(conn)
		tp := textproto.NewReader(reader)
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
				fmt.Print("User: ")
				fmt.Print(msgstr[1])
				fmt.Print(" Channel: ")
				fmt.Print(msgstr[4])
				fmt.Print(" Message: ")
				fmt.Print(msgstr[5])
				fmt.Println()

				f.WriteString(msgstr[1] + "\t" + msgstr[4] + "\t" + msgstr[5] + "\n")
			}
		}
	}
}
