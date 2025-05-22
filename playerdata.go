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
	"strings"

	"github.com/butt4cak3/theislercon/internal/parser"
)

func parsePlayerDataMessage(message string) ([]Player, error) {
	response, err := parseResponse(message, "PlayerData")
	if err != nil {
		return nil, err
	}

	message = response.Content

	lines := strings.Split(message, "\n")

	result := make([]Player, 0)

	for _, line := range lines {
		msg := parser.SkipWhitespace(line)
		if msg == "" {
			continue
		}

		msg, err := parser.Tag(msg, "Name:")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		pidIndex := strings.LastIndex(msg, ", PlayerID:")
		name, msg := parser.Take(msg, pidIndex)
		_, msg = parser.Take(msg, 1)

		msg = parser.SkipWhitespace(msg)

		msg, err = parser.Tag(msg, "PlayerID:")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		playerId, msg, err := parser.ParseDigits(msg)
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		msg, err = parser.Tag(msg, "Location:")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		x, msg, err := parseLocationField(msg, "X")
		if err != nil {
			return nil, err
		}
		y, msg, err := parseLocationField(msg, "Y")
		if err != nil {
			return nil, err
		}
		z, msg, err := parseLocationField(msg, "Z")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		msg, err = parser.Tag(msg, "Class:")
		if err != nil {
			return nil, err
		}

		msg = parser.SkipWhitespace(msg)

		class, msg, err := parseClassName(msg)
		if err != nil {
			return nil, err
		}

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		growth, msg, err := parsePercentageField(msg, "Growth")
		if err != nil {
			return nil, err
		}

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		health, msg, err := parsePercentageField(msg, "Health")
		if err != nil {
			return nil, err
		}

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		stamina, msg, err := parsePercentageField(msg, "Stamina")
		if err != nil {
			return nil, err
		}

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		hunger, msg, err := parsePercentageField(msg, "Hunger")
		if err != nil {
			return nil, err
		}

		msg, err = parser.Tag(msg, ",")
		if err != nil {
			return nil, err
		}

		thirst, _, err := parsePercentageField(msg, "Thirst")
		if err != nil {
			return nil, err
		}

		class = class[3 : len(class)-2]

		player := Player{playerId, name, Location{x, y, z}, DinoClass(class), growth, health, stamina, hunger, thirst}
		result = append(result, player)
	}

	return result, nil
}

func parseLocationField(msg, name string) (float64, string, error) {
	msg = parser.SkipWhitespace(msg)

	msg, err := parser.Tag(msg, name)
	if err != nil {
		return 0, "", err
	}

	msg, err = parser.Tag(msg, "=")
	if err != nil {
		return 0, "", err
	}

	v, msg, err := parser.ParseFloat64(msg)
	if err != nil {
		return 0, "", err
	}
	return v, msg, nil
}

func parsePercentageField(msg, name string) (int8, string, error) {
	msg = parser.SkipWhitespace(msg)

	msg, err := parser.Tag(msg, name)
	if err != nil {
		return 0, "", err
	}
	msg, err = parser.Tag(msg, ":")
	if err != nil {
		return 0, "", err
	}
	msg = parser.SkipWhitespace(msg)
	f, msg, err := parser.ParseFloat64(msg)
	if err != nil {
		return 0, "", err
	}
	return int8(f * 100), msg, nil
}

func parseClassName(msg string) (string, string, error) {
	pos := 0
	for pos < len(msg) && (msg[pos] == '_' || parser.IsAsciiLetter(msg[pos])) {
		pos++
	}
	if pos > 0 {
		return msg[:pos], msg[pos:], nil
	} else {
		return "", "", fmt.Errorf("expected class name")

	}
}
