# Tmi Client

The following code is an example on how this irc bot client can be used:

```golang
package main

import (
  "log"
  "os"
  
  "github.com/dfalgout/tmiclient"
)

var (
  nick  = os.Getenv("TMI_NICK")
  token = os.Getenv("TMI_TOKEN")
)

func msgHandler(tmi *tmiclient.TMI, result *tmiclient.Result) {
  msg := result.Value.(*tmiclient.Message)
  log.Println("Channel:", msg.Channel)
  log.Println("Badges:", msg.Badges)
  log.Println("Color:", msg.Color)
  log.Println("DisplayName:", msg.DisplayName)
  log.Println("Emotes:", msg.Emotes)
  log.Println("Id:", msg.Id)
  log.Println("IsMod:", msg.IsMod)
  log.Println("IsSubscriber:", msg.IsSubscriber)
  log.Println("IsTurbo:", msg.IsTurbo)
  log.Println("RoomId:", msg.RoomId)
  log.Println("SentTime:", msg.SentTime)
  log.Println("TmiSentTime:", msg.TmiSentTime)  
  log.Println("UserId:", msg.UserId)
  log.Println("UserType:", msg.UserType)
  log.Println("Payload:", msg.Payload)
}

func joinHandler(tmi *tmiclient.TMI, result *tmiclient.Result) {
  event := result.Value.(*tmiclient.Event)
  log.Println("Channel:", event.Channel)
  log.Println("DisplayName:", event.DisplayName)
}

func partHandler(tmi *tmiclient.TMI, result *tmiclient.Result) {
  event := result.Value.(*tmiclient.Event)
  log.Println("EVENT -> PART")
  log.Println("Channel:", event.Channel)
  log.Println("DisplayName:", event.DisplayName)
}

func main() {
  chans := []string{"you_channel_name", "another_channel_name"}

  tmi := tmiclient.NewTMI(nick, token, chans)
  tmi.RegisterHandler("msg", msgHandler)
  tmi.RegisterHandler("join", joinHandler)
  tmi.RegisterHandler("part", partHandler)
  tmi.Connect()
}
```

Handlers have the following prototype:

```golang
func someHandler(*tmiclient.TMI, *tmiclient.Result) {}
```

The following events are currently available for subscription and the objects that are injected into the handler:
*Note: Objects are injected into the Result.Value property*

| Event | Object |
|-------|--------|
|banned | \*tmiclient.Banned |
|error  | \*tmiclient.Error  |
|join   | \*tmiclient.Event  |
|part   | \*tmiclient.Event  |
|msg    | \*tmiclient.Message|
|mod:add| \*tmiclient.Event  |
|mod:sub| \*tmiclient.Event  |
|names  | \*tmiclient.Names  |
|notice | \*tmiclient.Notice |
|roomstate | \*tmiclient.RoomState |
|globaluserstate | \*tmiclient.UserState |
