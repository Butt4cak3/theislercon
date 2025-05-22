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

// Classes are the different types of dinosaurs that
// players can choose to play.
type DinoClass string

const (
	Beipiaosaurus      DinoClass = "Beipiaosaurus"
	Carnotaurus        DinoClass = "Carnotaurus"
	Ceratosaurus       DinoClass = "Ceratosaurus"
	Deinosuchus        DinoClass = "Deinosuchus"
	Diabloceratops     DinoClass = "Diabloceratops"
	Dilophosaurus      DinoClass = "Dilophosaurus"
	Dryosaurus         DinoClass = "Dryosaurus"
	Gallimimus         DinoClass = "Gallimimus"
	Herrerasaurus      DinoClass = "Herrerasaurus"
	Hypsilophodon      DinoClass = "Hypsilophodon"
	Maiasaura          DinoClass = "Maiasaura"
	Omniraptor         DinoClass = "Omniraptor"
	Pachycephalosaurus DinoClass = "Pachycephalosaurus"
	Pteranodon         DinoClass = "Pteranodon"
	Stegosaurus        DinoClass = "Stegosaurus"
	Tenontosaurus      DinoClass = "Tenontosaurus"
	Troodon            DinoClass = "Troodon"
)

var AllClasses = [17]DinoClass{Beipiaosaurus, Carnotaurus, Ceratosaurus, Deinosuchus, Diabloceratops, Dilophosaurus, Dryosaurus, Gallimimus, Herrerasaurus, Hypsilophodon, Maiasaura, Omniraptor, Pachycephalosaurus, Pteranodon, Stegosaurus, Tenontosaurus, Troodon}

func IsClass(s string) bool {
	c := DinoClass(s)
	for _, class := range AllClasses {
		if c == class {
			return true
		}
	}
	return false
}

// Name returns a human-readable name for this class.
//
// In case a class is not yet supported by this library, the name as
// reported by the server will be used.
func (c DinoClass) Name() string {
	return string(c)
}

type AIClass string

const (
	Boar          AIClass = "Boar"
	Compsognathus AIClass = "Compsognathus"
	Deer          AIClass = "Deer"
	Goat          AIClass = "Goat"
	Pterodactylus AIClass = "Pterodactylus"
	Seaturtle     AIClass = "Seaturtle"
)

var AllAIClasses = [6]AIClass{Boar, Compsognathus, Deer, Goat, Pterodactylus, Seaturtle}

func IsAIClass(s string) bool {
	c := AIClass(s)
	for _, class := range AllAIClasses {
		if c == class {
			return true
		}
	}
	return false
}
