package main

import (
	jar "github.com/tsauzeau/authIntra/intra-auth/epiJar"
)

func main() {
	jar := jar.New("thomas.sauzeau@epitech.eu", "tmp")
	jar.Auth()
}
