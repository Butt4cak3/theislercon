/*
Copyright (C) 2025  Marius Becker

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package theislercon

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// The Client type contains methods for all RCON commands.
type Client struct {
	conn  net.Conn
	mutex sync.Mutex
}

// Connect tries to connect to the specified address.
func Connect(addr string) (*Client, error) {
	client := new(Client)
	err := client.connect(addr)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (client *Client) connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

// Auth tries to authenticate with the gameserver.
//
// If the reason for failure is an incorrect password, the function
// will return [ErrIncorrectPassword].
func (client *Client) Auth(password string) error {
	msg := make([]byte, 0, 100)
	msg = append(msg, Auth)
	msg = append(msg, password...)

	client.mutex.Lock()
	defer client.mutex.Unlock()

	err := client.send(msg)
	if err != nil {
		return err
	}

	res, err := client.recv()
	if err != nil {
		return err
	}

	response := string(res)

	if response == "Password Accepted" {
		return nil
	} else {
		return ErrIncorrectPassword
	}
}

// GetPlayerList returns a list of all connected players.
//
// The list contains both players that are playing and players that are
// connected, but not playing (i.e. players that are in the class selection
// screen). However, the list only contains IDs and names.
func (client *Client) GetPlayerList() ([]Player, error) {
	msg, err := client.ExecCommand(GetPlayerList)
	if err != nil {
		return nil, err
	}

	response, err := parseResponse(msg, "PlayerList")
	if err != nil {
		return nil, err
	}

	msg = response.Content

	lines := strings.Split(msg, ",")
	players := make([]Player, 0, len(lines)/3)

	for i := 0; i < len(lines); i += 3 {
		playerID := strings.TrimSpace(lines[i])
		if playerID == "" {
			break
		}
		name := strings.TrimSpace(lines[i+1])
		players = append(players, Player{ID: playerID, Name: name})
	}

	return players, nil
}

// Announce sends a message to all currently connected players.
//
// The message will be displayed in a large text box at the top of the screen.
func (client *Client) Announce(message string) error {
	_, err := client.ExecCommand(Announce, message)
	return err
}

// SendDirectMessage sends an announcement message to one specific user.
//
// The message will be shown in the same way as a regular announcement.
func (client *Client) SendDirectMessage(playerID, message string) error {
	_, err := client.ExecCommand(DirectMessage, playerID, message)
	return err
}

// GetPlayerData returns a list of all players that have spawned in.
//
// In contrast to [Client.GetPlayerList], this function returns complete [Player]
// structs and not just IDs and names. However, it excludes players that
// are still in the class selection screen.
func (client *Client) GetPlayerData() ([]Player, error) {
	msg, err := client.ExecCommand(GetPlayerData)
	if err != nil {
		return nil, err
	}

	return parsePlayerDataMessage(msg)
}

// GetServerDetails returns some information about the server.
func (client *Client) GetServerDetails() (*ServerDetails, error) {
	msg, err := client.ExecCommand(GetServerDetails)
	if err != nil {
		return nil, err
	}

	result, err := parseServerDetails(msg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// WipeCorpses removes all dead entities from the map.
func (client *Client) WipeCorpses() error {
	_, err := client.ExecCommand(WipeCorpses)
	return err
}

// UpdatePlayables defines the list of playable classes.
func (client *Client) UpdatePlayables(classes []Class) error {
	classNames := make([]string, len(classes))
	for i, class := range classes {
		classNames[i] = string(class)
	}
	_, err := client.ExecCommand(UpdatePlayables, classNames...)
	return err
}

// KickPlayer kicks the player from the server.
func (client *Client) KickPlayer(playerID, reason string) error {
	_, err := client.ExecCommand(KickPlayer, playerID, reason)
	return err
}

// Save saves the current state of the map.
func (client *Client) Save() error {
	_, err := client.ExecCommand(Save)
	return err
}

// ToggleWhitelist turns the whielist on or off and returns true, if the new state is on.
func (client *Client) ToggleWhitelist() (bool, error) {
	res, err := client.ExecCommand(ToggleWhitelist)
	if err != nil {
		return false, err
	}
	return onOff(res[len(res)-2:]), nil
}

// AddWhitelistID adds one or more PlayerIDs to the whitelist.
func (client *Client) AddWhitelistID(playerID ...string) error {
	if len(playerID) > 0 {
		_, err := client.ExecCommand(AddWhitelistID, playerID...)
		return err
	} else {
		return nil
	}
}

// RemoveWhitelistID removes one or more PlayerIDs from the whitelist.
func (client *Client) RemoveWhitelistID(playerID ...string) error {
	if len(playerID) > 0 {
		_, err := client.ExecCommand(RemoveWhitelistID, playerID...)
		return err
	} else {
		return nil
	}
}

// ToggleGlobalState turns on or off the global chat feature on the server.
//
// If you need to know whether global chat is already enabled,
// use [Client.GetServerDetails].
func (client *Client) ToggleGlobalChat() (bool, error) {
	res, err := client.ExecCommand(ToggleGlobalChat)
	if err != nil {
		return false, err
	}
	return onOff(res[len(res)-2:]), nil
}

// ToggleHumans turns on or off the humans feature in the game.
func (client *Client) ToggleHumans() (bool, error) {
	res, err := client.ExecCommand(ToggleHumans)
	if err != nil {
		return false, err
	}
	return onOff(res[len(res)-2:]), nil
}

// ToggleAI turns the spawning of AI on or off.
func (client *Client) ToggleAI() (bool, error) {
	res, err := client.ExecCommand(ToggleAI)
	if err != nil {
		return false, err
	}
	return onOff(res[len(res)-2:]), nil
}

// DisableAIClasses defines the list of AI classes that cannot spawn.
func (client *Client) DisableAIClasses(classes []AIClass) error {
	classNames := make([]string, len(classes))
	for i, class := range classes {
		classNames[i] = string(class)
	}
	_, err := client.ExecCommand(DisableAIClasses, classNames...)
	return err
}

// SetAIDensity sets the AI density that can also be defined in Game.ini.
func (client *Client) SetAIDensity(density float32) error {
	_, err := client.ExecCommand(SetAIDensity, fmt.Sprintf("%.3f", density))
	return err
}

// Close closes the connection to the server.
func (client *Client) Close() error {
	return client.conn.Close()
}

// ExecCommand formats a command and sends it to the server.
//
// The client should have a corresponding method for every supported command.
// However, there may be undocumented features in the RCON protocol or this
// library may become out-of-date. In those cases, you can use ExecCommand
// directly to send commands that the client does not support (yet).
//
// While command can be any byte, you should normally use one of the
// [MessageType] constants.
//
// This function returns the server's response as a string.
func (client *Client) ExecCommand(command byte, params ...string) (string, error) {
	cmd := make([]byte, 0, 1024)
	cmd = append(cmd, ExecCommand)
	cmd = append(cmd, command)
	for i, param := range params {
		if i > 0 {
			cmd = append(cmd, ',')
		}
		cmd = append(cmd, param...)
	}

	// We only ever want to send one request over this connection at a time.
	client.mutex.Lock()
	defer client.mutex.Unlock()

	err := client.send(cmd)
	if err != nil {
		return "", err
	}

	res, err := client.recv()
	if err != nil {
		return "", err
	}

	return string(res), nil
}

func (client *Client) send(msg []byte) error {
	_, err := client.conn.Write(msg)
	return err
}

func (client *Client) recv() ([]byte, error) {
	buf := make([]byte, 1024*10)
	client.conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
	n, err := client.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func onOff(s string) bool {
	return s == "On"
}
