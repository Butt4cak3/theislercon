/*
theislercon - A library for communicating with The Isle servers
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

type ServerDetails struct {
	Name                           string
	Password                       string
	Map                            string
	MaxPlayers                     int
	CurrentPlayers                 int
	EnableMutations                bool
	EnableHumans                   bool
	HasPassword                    bool
	QueueEnabled                   bool
	Whitelist                      bool
	SpawnAI                        bool
	AllowRecordingGameplay         bool
	UseRegionSpawning              bool
	UseRegionSpawnCooldown         bool
	RegionSpawnCooldownTimeSeconds int
	DayLengthMinutes               int
	NightLengthMinutes             int
	EnableGlobalChat               bool
}

func parseServerDetails(msg string) (*ServerDetails, error) {
	response, err := parseResponse(msg, "ServerDetails")
	if err != nil {
		return nil, err
	}

	msg = response.Content

	msg = parser.SkipWhitespace(msg)

	details := new(ServerDetails)

	for len(msg) > 0 {
		m := msg
		key, m, err := parseKey(m)
		if err != nil {
			return nil, err
		}

		m = parser.SkipWhitespace(m)

		switch key {
		case "ServerName":
			details.Name, m = parseStringValue(m)
		case "ServerPassword":
			details.Password, m = parseStringValue(m)
		case "ServerMap":
			details.Map, m = parseStringValue(m)
		case "ServerMaxPlayers":
			details.MaxPlayers, m, err = parseIntValue(m)
		case "ServerCurrentPlayers":
			details.CurrentPlayers, m, err = parseIntValue(m)
		case "bEnableMutations":
			details.EnableMutations, m, err = parseBoolValue(m)
		case "bEnableHumans":
			details.EnableHumans, m, err = parseBoolValue(m)
		case "bServerPassword":
			details.HasPassword, m, err = parseBoolValue(m)
		case "bQueueEnabled":
			details.QueueEnabled, m, err = parseBoolValue(m)
		case "bServerWhitelist":
			details.Whitelist, m, err = parseBoolValue(m)
		case "bSpawnAI":
			details.SpawnAI, m, err = parseBoolValue(m)
		case "bAllowRecordingReplay":
			details.AllowRecordingGameplay, m, err = parseBoolValue(m)
		case "bUseRegionSpawning":
			details.UseRegionSpawning, m, err = parseBoolValue(m)
		case "bUseRegionSpawnCooldown":
			details.UseRegionSpawnCooldown, m, err = parseBoolValue(m)
		case "RegionSpawnCooldownTimeSeconds":
			details.RegionSpawnCooldownTimeSeconds, m, err = parseIntValue(m)
		case "ServerDayLengthMinutes":
			details.DayLengthMinutes, m, err = parseIntValue(m)
		case "ServerNightLengthMinutes":
			details.NightLengthMinutes, m, err = parseIntValue(m)
		case "bEnableGlobalChat":
			details.EnableGlobalChat, m, err = parseBoolValue(m)
		default:
			return nil, fmt.Errorf("unknown key %s", key)
		}

		if err != nil {
			return nil, err
		}

		m = parser.SkipWhitespace(m)

		if len(m) > 0 && m[0] == ',' {
			m = m[1:]
		}

		m = parser.SkipWhitespace(m)

		msg = m
	}

	return details, nil
}

func parseKey(msg string) (string, string, error) {
	pos := 0
	for pos < len(msg) && parser.IsAsciiLetter(msg[pos]) {
		pos++
	}
	msg = parser.SkipWhitespace(msg)

	key, msg := parser.Take(msg, pos)

	msg, err := parser.Tag(msg, ":")
	if err != nil {
		return "", "", err
	}

	return key, msg, nil
}

func parseStringValue(msg string) (string, string) {
	pos := 0
	for len(msg) > 0 && msg[pos] != ',' {
		pos++
	}
	return msg[:pos], msg[pos:]
}

func parseIntValue(msg string) (int, string, error) {
	return parser.ParseInt(msg)
}

func parseBoolValue(msg string) (bool, string, error) {
	if msg[0:4] == "true" {
		return true, msg[4:], nil
	} else if msg[0:5] == "false" {
		return false, msg[5:], nil
	} else {
		return false, "", fmt.Errorf("expected bool")
	}
}
