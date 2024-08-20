package main

import (
	"context"
	"log"

	"github.com/godbus/dbus/v5"

	"github.com/2opremio/go-mpris"
)

func main() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	names, err := mpris.List(context.Background(), conn)
	if err != nil {
		panic(err)
	}
	if len(names) == 0 {
		log.Fatal("No media player found.")
	}

	name := names[0]
	log.Println("Found media player:", name)

	player := mpris.New(conn, name)

	_, id := player.GetIdentity()
	log.Println("Media player identity:", id)

	player.Raise()
}
