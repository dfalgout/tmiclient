package tmiclient

import (
	"errors"
	"net"
	"strings"
)

// event handler required prototypes
type handler func(*TMI, *Result)
type handlers map[string][]handler

// Specify individual event data structures
type Result struct {
	Value interface{}
}

type Message struct {
	Badges       string
	Channel      string
	Color        string
	DisplayName  string
	Emotes       string
	Id           string
	IsMod        bool
	Payload      string
	RoomId       string
	SentTime     int
	IsSubscriber bool
	TmiSentTime  int
	IsTurbo      bool
	UserId       int
	UserType     string
}

// @badges=premium/1;color=;display-name=u_lost;emotes=;message-id=18;thread-id=108938100_159925790;turbo=1;user-id=108938100;user-type= :u_lost!u_lost@u_lost.tmi.twitch.tv WHISPER streamwars_bot :yo

type Whisper struct {
	Badges    string
	Color     string
	From      string
	Emotes    string
	MessageID int
	ThreadID  string
	IsTurbo   bool
	UserID    int
	UserType  string
	To        string
	Payload   string
}

type Event struct {
	Channel     string
	DisplayName string
}

type RoomState struct {
	BroadcasterLanguage string
	Channel             string
	R9K                 bool
	Slow                int
	SubsOnly            bool
}

type UserState struct {
	Color       string
	DisplayName string
	EmoteSets   []int
	Turbo       bool
	UserID      int
	UserType    string
}

type Names struct {
	Channel string
	Viewers []string
}

type Banned struct {
	Channel     string
	DisplayName string
	BanDuration int
	BanReason   string
}

type Error struct {
	Command string
	Reason  string
}

type Notice struct {
	Channel   string
	MessageID string
	Message   string
}

type TMI struct {
	Nick     string
	Token    string
	Chans    []string
	handlers handlers
	conn     *net.TCPConn
}

func NewTMI(nick string, token string, chans []string) *TMI {
	return &TMI{
		Nick:     nick,
		Token:    token,
		Chans:    chans,
		handlers: make(handlers),
	}
}

func parseNotice(channel string, msgID string, message string) *Notice {
	fieldMap := getFieldMap(msgID)

	return &Notice{
		Channel:   strings.Trim(channel, "#"),
		MessageID: fieldMap[msgID],
		Message:   message,
	}
}

func parseError(command string, reason string) *Error {
	return &Error{
		Command: command,
		Reason:  reason,
	}
}

// Private Helper functions
func parseEvent(channel string, viewer string) *Event {
	return &Event{
		Channel:     strings.Trim(channel, "#"),
		DisplayName: viewer,
	}
}

func parseClearChat(channel string, displayName string, data string) *Banned {
	fieldMap := getFieldMap(data)

	return &Banned{
		Channel:     strings.TrimPrefix(channel, "#"),
		DisplayName: displayName,
		BanDuration: getIntFromString(fieldMap[banDuration]),
		BanReason:   fieldMap[banReason],
	}
}

func parseNames(channel string, names []string) *Names {
	names[0] = strings.TrimPrefix(names[0], ":")
	return &Names{
		Channel: strings.TrimLeft(channel, "#"),
		Viewers: names,
	}
}

func parseGlobalUserState(channel string, state string) *UserState {
	fieldMap := getFieldMap(state)

	return &UserState{
		Color:       fieldMap[color],
		DisplayName: fieldMap[displayName],
		EmoteSets:   getIntListFromString(fieldMap[emoteSets]),
		Turbo:       getBoolFromString(fieldMap[turbo]),
		UserID:      getIntFromString(fieldMap[userID]),
		UserType:    fieldMap[userType],
	}
}

func parseRoomStateFields(channel string, state string) *RoomState {
	fieldMap := getFieldMap(state)

	return &RoomState{
		BroadcasterLanguage: fieldMap[broadcasterLang],
		Channel:             strings.Trim(channel, "#"),
		R9K:                 getBoolFromString(fieldMap[r9k]),
		Slow:                getIntFromString(fieldMap[slow]),
		SubsOnly:            getBoolFromString(fieldMap[subsOnly]),
	}
}

func parseWhisper(channel string, to string, message string, payload string) *Whisper {
	fieldMap := getFieldMap(message)

	return &Whisper{
		Badges:    fieldMap[badges],
		Color:     fieldMap[color],
		From:      fieldMap[displayName],
		Emotes:    fieldMap[emotes],
		IsTurbo:   getBoolFromString(fieldMap[turbo]),
		MessageID: getIntFromString(fieldMap[messageID]),
		Payload:   payload,
		ThreadID:  fieldMap[threadID],
		To:        to,
	}
}

func parseMessageFields(chn string, msg string, payload string) *Message {
	fieldMap := getFieldMap(msg)

	return &Message{
		Badges:       fieldMap[badges],
		Channel:      strings.Trim(chn, "#"),
		Color:        fieldMap[color],
		DisplayName:  fieldMap[displayName],
		Emotes:       fieldMap[emotes],
		Id:           fieldMap[id],
		IsMod:        getBoolFromString(fieldMap[mod]),
		Payload:      payload,
		RoomId:       fieldMap[roomID],
		SentTime:     getIntFromString(fieldMap[sentTS]),
		IsSubscriber: getBoolFromString(fieldMap[subscriber]),
		TmiSentTime:  getIntFromString(fieldMap[tmiSentTS]),
		IsTurbo:      getBoolFromString(fieldMap[turbo]),
		UserId:       getIntFromString(fieldMap[userID]),
		UserType:     fieldMap[userType],
	}
}

func (t *TMI) invokeHandlers(event string, result interface{}) {
	event = strings.ToLower(event)

	for _, v := range t.handlers[event] {
		v(t, &Result{
			Value: result,
		})
	}
}

func (t *TMI) login() error {
	if t.Nick == "" {
		return errors.New("Nickname Required")
	}

	if t.Token == "" {
		return errors.New("Token Required")
	}

	nickCmd := "NICK " + t.Nick
	passCmd := "PASS oauth:" + t.Token

	_, err := t.bulkSend([]string{passCmd, nickCmd})
	err = t.JoinChannels(t.Chans)
	return err
}

func (t *TMI) send(msg string) (int, error) {
	return t.conn.Write([]byte(msg + "\r\n"))
}

func (t *TMI) bulkSend(msg []string) (int, error) {
	var total = 0
	for _, m := range msg {
		_, err := t.send(m)
		if err != nil {
			return total, err
		}
		total++
	}

	return total, nil
}

// Public Helper Functions
func (t *TMI) JoinChannel(channel string) error {
	command := "JOIN #" + channel
	_, err := t.send(command)
	if err == nil {
		t.Chans = appendListUnique(t.Chans, channel)
	}
	return err
}

func (t *TMI) JoinChannels(channels []string) error {
	for _, v := range channels {
		err := t.JoinChannel(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TMI) LeaveChannel(channel string) error {
	command := "PART #" + channel
	_, err := t.send(command)
	if err == nil {
		// Also have to remove the channel from the TMI client channel list
		for k, v := range t.Chans {
			if v == channel {
				t.Chans = append(t.Chans[:k], t.Chans[k+1:]...)
			}
		}
	}
	return err
}

func (t *TMI) LeaveChannels(channels []string) error {
	for _, v := range channels {
		err := t.LeaveChannel(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TMI) SendMsg(channel string, msg string) error {
	command := "PRIVMSG #" + channel + " :" + msg
	_, err := t.send(command)
	return err
}

func (t *TMI) SendWhisper(viewer string, msg string) error {
	command := "/w " + viewer + " " + msg
	return t.SendMsg("jtv", command)
}

func (t *TMI) GiveMod(channel string, viewer string) error {
	command := "/mod " + viewer
	return t.SendMsg(channel, command)
}

func (t *TMI) RemoveMod(channel string, viewer string) error {
	command := "/unmod " + viewer
	return t.SendMsg(channel, command)
}

func (t *TMI) RegisterHandler(event string, h handler) error {
	event = strings.ToLower(event)

	if ok := isin(event, validEvents); ok {
		t.handlers[event] = append(t.handlers[event], h)
	} else {
		return errors.New(event + " is not a valid event")
	}

	return nil
}

func (t *TMI) SubscribeAll(events []string) error {
	for _, v := range events {
		err := t.Subscribe(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TMI) Subscribe(event string) error {
	command := "CAP REQ :twitch.tv/" + event
	_, err := t.send(command)
	return err
}
