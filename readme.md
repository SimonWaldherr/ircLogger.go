#ircLogger.go

a [IRC](http://en.wikipedia.org/wiki/Internet_Relay_Chat) logger in [golang](http://golang.org).  

currently it logs to a [tsv](http://en.wikipedia.org/wiki/Tab-separated_values) file, in the next version [mySQL](http://en.wikipedia.org/wiki/MySQL) will also be possible.

```sh
go run ircLogger.go irc.freenode.net 6667 channel nickname
```

