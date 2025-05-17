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

type Player struct {
	ID       string // Steam or EOS ID
	Name     string
	Location Location
	Class    Class
	Growth   int8 // Growth as percentage. 75% means fully grown in the current game version.
	Health   int8 // Health as percentage
	Stamina  int8 // Stamina as percentage
	Hunger   int8 // Hunger as percentage
	Thirst   int8 // Thirst as percentage
}

type Location struct {
	X float64 // Latitude
	Y float64 // Longitude
	Z float64 // Altitude
}
