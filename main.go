package main

import "bitbucket.org/lewington/erosai-server/server"

func main() {
	s := server.New()
	s.SpinUp()
}
