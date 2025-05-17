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

// Package parser contains functions that help with parsing the
// RCON server responses.
package parser

import (
	"fmt"
	"strconv"
)

func Tag(msg, substr string) (string, error) {
	if msg[:len(substr)] == substr {
		return msg[len(substr):], nil
	} else {
		return "", fmt.Errorf("expected \"%s\"", substr)
	}
}

func ParseDigits(msg string) (string, string, error) {
	pos := 0
	for pos < len(msg) && IsAsciiDigit(msg[pos]) {
		pos++
	}
	if pos > 0 {
		return msg[:pos], msg[pos:], nil
	} else {
		return "", "", fmt.Errorf("expected digits")
	}
}

func ParseInt64(msg string) (int64, string, error) {
	digits, msg, err := ParseDigits(msg)
	if err != nil {
		return 0, "", err
	}
	num, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return 0, "", err
	}
	return num, msg, nil
}

func ParseInt(msg string) (int, string, error) {
	digits, msg, err := ParseDigits(msg)
	if err != nil {
		return 0, "", err
	}
	num, err := strconv.Atoi(digits)
	if err != nil {
		return 0, "", err
	}
	return num, msg, nil
}

func ParseFloat64(msg string) (float64, string, error) {
	pos := 0
	if pos < len(msg) && msg[pos] == '-' {
		pos++
	}

	for pos < len(msg) && IsAsciiDigit(msg[pos]) {
		pos++
	}

	if pos < len(msg) && msg[pos] == '.' {
		pos++
		for pos < len(msg) && IsAsciiDigit(msg[pos]) {
			pos++
		}
	}

	if pos > 0 && msg[pos-1] != '-' {
		f, err := strconv.ParseFloat(msg[:pos], 64)
		if err != nil {
			return 0, "", err
		}
		return f, msg[pos:], nil
	} else {
		return 0, "", fmt.Errorf("expected float at %s", msg[:10])
	}
}

func IsAsciiDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func Take(msg string, n int) (string, string) {
	n = min(n, len(msg))
	return msg[:n], msg[n:]
}

func SkipWhitespace(msg string) string {
	pos := 0
	for pos < len(msg) && (msg[pos] == ' ' || msg[pos] == '\n' || msg[pos] == '\r' || msg[pos] == '\t') {
		pos++
	}
	return msg[pos:]
}

func IsAsciiLetter(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func ParseTimestamp(msg string) (string, string, error) {
	year, msg, err := ParseDigits(msg)
	if err != nil {
		return "", "", err
	}

	msg, err = Tag(msg, ".")
	if err != nil {
		return "", "", err
	}

	month, msg, err := ParseDigits(msg)
	if err != nil {
		return "", "", err
	}

	msg, err = Tag(msg, ".")
	if err != nil {
		return "", "", err
	}

	day, msg, err := ParseDigits(msg)
	if err != nil {
		return "", "", err
	}

	msg, err = Tag(msg, "-")
	if err != nil {
		return "", "", err
	}

	hour, msg, err := ParseDigits(msg)
	if err != nil {
		return "", "", err
	}

	msg, err = Tag(msg, ".")
	if err != nil {
		return "", "", err
	}

	minute, msg, err := ParseDigits(msg)
	if err != nil {
		return "", "", err
	}

	msg, err = Tag(msg, ".")
	if err != nil {
		return "", "", err
	}

	second, msg, err := ParseDigits(msg)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%s.%s.%s-%s.%s.%s", year, month, day, hour, minute, second), msg, nil
}
