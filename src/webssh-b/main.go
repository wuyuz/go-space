package main

import (
	"bytes"
	"log"
	ssh "webssh-b/ssh2ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	app := fiber.New()

	app.Static("/", "./")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		client, _ := ssh.NewSshClient()
		defer client.Close()
		ssConn, _ := ssh.NewSshConn(120, 32, client)
		defer ssConn.Close()

		quitChan := make(chan bool, 3)

		var logBuff = new(bytes.Buffer)

		// most messages are ssh output, not webSocket input
		go ssConn.ReceiveWsMsg(c, logBuff, quitChan)
		go ssConn.SendComboOutput(c, quitChan)
		go ssConn.SessionWait(quitChan)
		<-quitChan
	}))

	log.Fatal(app.Listen(":3000"))
}
