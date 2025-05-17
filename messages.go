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

	"github.com/butt4cak3/theislercon/internal/parser"
)

type MessageType = byte

const (
	Auth              MessageType = 0x01
	ExecCommand       MessageType = 0x02
	ResponseValue     MessageType = 0x03
	Announce          MessageType = 0x10
	DirectMessage     MessageType = 0x11
	GetServerDetails  MessageType = 0x12
	WipeCorpses       MessageType = 0x13
	UpdatePlayables   MessageType = 0x15
	BanPlayer         MessageType = 0x20
	KickPlayer        MessageType = 0x30
	GetPlayerList     MessageType = 0x40
	Save              MessageType = 0x50
	GetPlayerData     MessageType = 0x77
	ToggleWhitelist   MessageType = 0x81
	AddWhitelistID    MessageType = 0x82
	RemoveWhitelistID MessageType = 0x83
	ToggleGlobalChat  MessageType = 0x84
	ToggleHumans      MessageType = 0x86
	ToggleAI          MessageType = 0x90
	DisableAIClasses  MessageType = 0x91
	SetAIDensity      MessageType = 0x92
)

type Response struct {
	Timestamp string
	Type      string
	Content   string
}

func parseResponse(msg string, responseType string) (*Response, error) {
	var timestamp string
	var err error

	if msg[0] == '[' {
		msg, err = parser.Tag(msg, "[")
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMalformedResponse, err)
		}

		timestamp, msg, err = parser.ParseTimestamp(msg)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMalformedResponse, err)
		}

		msg, err = parser.Tag(msg, "]")
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrMalformedResponse, err)
		}

		msg = parser.SkipWhitespace(msg)
	}

	msg, err = parser.Tag(msg, responseType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrMalformedResponse, err)
	}

	return &Response{timestamp, responseType, msg}, nil
}
