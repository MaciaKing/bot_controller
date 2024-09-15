package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	ip := flag.String("ip", "0.0.0.0", "King botnet IP direction")
	port := flag.String("port", "8080", "King botnet port")
	flag.Parse()

	// Url definition for WebSocket
	host := *ip + ":" + *port
	u := url.URL{Scheme: "ws", Host: host, Path: "/ws"}
	log.Printf("Connecting to %s", u.String())

	// Establecer la conexión
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	// Configurar un canal para manejar las interrupciones del sistema (ej. Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Canal para mensajes de salida
	done := make(chan struct{})

	// Goroutine para leer mensajes desde el servidor
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			log.Printf("Message from bot king: %s", message)
		}
	}()

	// Bucle para esperar mensajes o interrupciones
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection...")

			// Intentar cerrar la conexión limpiamente
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Close error:", err)
				return
			}
		}
	}
}
