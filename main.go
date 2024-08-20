package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)
	return localAddress.IP
}

func getClientIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	ip, _, _ := net.SplitHostPort(remoteAddr)
	return ip
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		logMessage := fmt.Sprintf("[%s] : %s", currentTime, err)

		color.Yellow("Local IP: %s", GetLocalIP())
		color.Red(logMessage)
		return
	}
	defer conn.Close()

	clientIP := getClientIP(r)
	color.Yellow("Client connected from IP: %s", clientIP)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			logMessage := fmt.Sprintf("[%s][%s] : %s", currentTime, clientIP, err)

			color.Red(logMessage)
			break
		}

		currentTime := time.Now().Format("2006-01-02 15:04:05")
		logMessage := fmt.Sprintf("[%s][%s] : %s", currentTime, clientIP, message)

		color.Blue(logMessage)

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			// Error occurred while writing message
			currentTime := time.Now().Format("2006-01-02 15:04:05")
			logMessage := fmt.Sprintf("[%s][%s] : %s", currentTime, clientIP, err)

			color.Red(logMessage)
			break
		}
	}
}

func main() {
	localIP := GetLocalIP().String()
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	fullURL := fmt.Sprintf("ws://%s%s", localIP, port)

	fmt.Printf("[Starting server on %s ...]\n", fullURL)
	http.HandleFunc("/ws", websocketHandler)
	duration := time.Duration(2) * time.Second
	time.Sleep(duration)
	color.Green("Server is online at %s", fullURL)
	log.Fatal(http.ListenAndServe(port, nil))
}
