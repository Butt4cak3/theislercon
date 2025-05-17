# The Isle RCON Client

This is a RCON client library for The Isle Evrima dedicated servers.

## The Client type

Client is the main interface of this library. It contains methods for all known RCON commands. To create a new client, use `Connect(addr)`. Make sure to call `Client.Close()` when you're done with it.

These are the methods that are currently supported:

- AddWhitelistID
- Announce
- Auth
- DisableAIClasses
- GetPlayerData
- GetPlayerList
- GetServerDetails
- KickPlayer
- RemoveWhitelistID
- Save
- SendDirectMessage
- SetAIDensity
- ToggleAI
- ToggleGlobalChat
- ToggleHumans
- ToggleWhitelist
- UpdatePlayables
- WipeCorpses

## Example: Get a list of connected players

```go
import (
    "fmt"
    rcon "github.com/butt4cak3/theislercon"
)

func main() {
    client, err := rcon.Connect("127.0.0.1:8888")
    if err != nil {
        fmt.Println("Could not connect")
        return
    }
    defer client.Close()

    err = client.Auth("YourSecurePasswordHere")
    if err != nil {
        fmt.Println("Could not authenticate")
        return
    }

    players, err := client.GetPlayerList()
    if err != nil {
        fmt.Println("Failed to get player list")
        return
    }

    for _, player := range players {
        fmt.Printf("%s\n", player.Name)
    }
}
```

## The RCON protocol

What follows is a somewhat technical description of the underlying protocol. It may contain errors or misconceptions, because (apart from the command table below) it was mostly reverse engineered.

The communication between client and server happens in a request-response-format. The client sends a request to the server and the server then responds to that request. The server never sends any data without first receiving a request and it only responds at most once to every request.

### Authentication

The first message that must be sent to the server after the connection has been established is the authentication request. It consists of the byte 0x01 and the RCON password as defined in the server's Game.ini config file.

On successful authentication, the server will respond with the string "Password Accepted". Otherwise, it will respond with "Incorrect Password".

### Sending commands

A request frame consists of three parts: The ExecCommand byte, a command byte from the list below, and arguments. The command byte can be one of the values in the table below. The arguments are a list of zero or more strings, separated by commas. The command does _not_ need to be terminated by anything, although some client libraries append a NULL byte (0x00).

### Server responses

A response usually consists of three parts:

1. A timestamp, enclosed in brackets (`[` and `]`)
2. The name of the type of response as a string
3. Data, depending on the type of response

For some reason, _all_ of these components are optional. Most, but not all of the times, there is a timestamp. Sometimes the type of response is there, sometimes not.

If present, the closing bracket of the timestamp is followed by a space. However, there is no space between the response type name and the data.

> [!NOTE]
> The authentication response is special and does not contain a timestamp nor the response type name.

### Commands

#### Special bytes

| Command       | Byte | Arguments                  |
| ------------- | ---- | -------------------------- |
| Auth          | 0x01 | RCON password              |
| ExecCommand   | 0x02 | Command byte, arguments... |
| ResponseValue | 0x03 | Unknown                    |

#### Commands

| Command           | Byte | Response name     | Arguments         |
| ----------------- | ---- | ----------------- | ----------------- |
| Announce          | 0x10 | Announce          | Message           |
| DirectMessage     | 0x11 | DirectMessage     | PlayerID, Message |
| GetServerDetails  | 0x12 | ServerDetails     |                   |
| WipeCorpses       | 0x13 | WipeCorpses       |                   |
| UpdatePlayables   | 0x15 | UpdatePlayables   |                   |
| BanPlayer         | 0x20 | BanPlayer         | Unknown           |
| KickPlayer        | 0x30 | KickPlayer        | PlayerID          |
| GetPlayerList     | 0x40 | PlayerList        |                   |
| Save              | 0x50 | Save              |                   |
| GetPlayerData     | 0x77 | PlayerData        |                   |
| ToggleWhitelist   | 0x81 | ToggleWhitelist   |                   |
| AddWhitelistID    | 0x82 | AddWhitelistID    | PlayerID          |
| RemoveWhitelistID | 0x83 | RemoveWhitelistID | PlayerID          |
| ToggleGlobalChat  | 0x84 | ToggleGlobalChat  |                   |
| ToggleHumans      | 0x86 | ToggleHumans      |                   |
| ToggleAI          | 0x90 | ToggleAI          |                   |
| DisableAIClasses  | 0x91 | DisableAIClasses  | AIClasses         |
| AIDensity         | 0x92 | AIDensity         | Density           |
