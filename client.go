package tmiclient

import (
	"bufio"
	"net"
	"net/textproto"
	"strings"
)

// Accepted events to register handlers to
var (
	validEvents = []string{
		"banned",
		"error",
		"join",
		"part",
		"msg",
		"mod:add",
		"mod:sub",
		"names",
		"notice",
		"roomstate",
		"globaluserstate",
	}
)

func handleMessage(t *TMI, msg string) {
	// split the msg by spaces
	msgParts := strings.Split(msg, " ")
	// log.Println(msgParts)

	if msgParts[0] == "PING" {
		t.send("PONG " + msgParts[1])
		return
	}

	if len(msgParts) >= 2 {
		switch strings.ToLower(msgParts[1]) {
		case "join", "part":
			channel := msgParts[2]
			viewer := getUserName(msgParts[0])
			event := msgParts[1]

			t.invokeHandlers(event, parseEvent(channel, viewer))
			break
		case "mode":
			channel := msgParts[2]
			viewer := msgParts[4]
			mod := msgParts[3]
			event := "mod"

			if strings.Contains(mod, "+") {
				event += ":add"
			} else {
				event += ":sub"
			}

			t.invokeHandlers(event, parseEvent(channel, viewer))
			break
		case "353":
			channel := msgParts[4]
			names := msgParts[5:]

			t.invokeHandlers("names", parseNames(channel, names))
			break
		case "421":
			command := msgParts[3]
			msgParts[4] = strings.TrimPrefix(msgParts[4], ":")
			reason := strings.Join(msgParts[4:], " ")

			t.invokeHandlers("error", parseError(command, reason))
			break
		default:
			break
		}
	}

	if len(msgParts) >= 3 {
		switch strings.ToLower(msgParts[2]) {
		case "privmsg":
			channel := msgParts[3]
			message := msgParts[0]
			payload := strings.TrimPrefix(strings.Join(msgParts[4:], " "), ":")

			t.invokeHandlers("msg", parseMessageFields(channel, message, payload))
			break
		case "roomstate":
			channel := msgParts[3]
			state := msgParts[0]

			t.invokeHandlers("roomstate", parseRoomStateFields(channel, state))
			break
		case "globaluserstate":
			channel := msgParts[3]
			state := msgParts[0]

			t.invokeHandlers("globaluserstate", parseGlobalUserState(channel, state))
			break
		case "clearchat":
			channel := msgParts[3]
			displayName := msgParts[4]
			data := msgParts[0]

			t.invokeHandlers("banned", parseClearChat(channel, displayName, data))
			break
		case "notice":
			channel := msgParts[3]
			msgID := strings.TrimPrefix(msgParts[0], "@")
			payload := strings.TrimPrefix(strings.Join(msgParts[4:], " "), ":")

			t.invokeHandlers("notice", parseNotice(channel, msgID, payload))
			break
		default:
			break
		}
	}
}

func (t *TMI) Connect() {
	subscriptions := []string{
		membership,
		tags,
		commands,
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
	handleWriteError(err)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	handleWriteError(err)
	defer conn.Close()
	t.conn = conn

	err = t.login()
	handleWriteError(err)

	err = t.SubscribeAll(subscriptions)
	handleWriteError(err)

	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		msg, err := tp.ReadLine()
		handleWriteError(err)

		// Handle messages concurrently
		go handleMessage(t, msg)
	}
}
