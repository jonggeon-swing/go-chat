package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	sessionManager := &SessionManager{
		sess: &sync.Map{},
	}
	messageChannel := make(chan string, 10)
	go listenFromUser(sessionManager, messageChannel)
	go listenFromMessageChannel(sessionManager, messageChannel)

	// 메인 고루틴이 종료되지 않도록 무한 대기
	log.Println("서버가 실행 중입니다. Ctrl+C로 종료할 수 있습니다.")
	WaitForStop()
}

func listenFromUser(manager *SessionManager, messageChannel chan<- string) {
	if listener, err := net.ListenTCP( //socket 연결부분
		"tcp",
		&net.TCPAddr{Port: 8080},
	); err != nil {
		log.Println("ERROR ...", err)
		panic(err)
	} else {
		for {
			log.Println("CONN ...")
			if conn, _ := listener.AcceptTCP(); conn != nil {
				_ = conn.SetKeepAlive(true)
				_ = conn.SetKeepAlivePeriod(4 * time.Minute)
				manager.sess.Store(conn.RemoteAddr().String(), Session{Socket: conn})
				go handleConnection(conn, messageChannel)
			}
		}
	}
}

func handleConnection(conn net.Conn, messageChannel chan<- string) {
	// Handle the connection
	// For example, read data from the connection and process it
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			// Handle error
			panic(err)
		}
	}(conn)

	// Simulate some processing
	for {
		buffer := make([]byte, 10)
		n, err := conn.Read(buffer)
		if err != nil {
			return
		}
		msg := string(buffer[:n])
		messageChannel <- msg
	}
}

func listenFromMessageChannel(manager *SessionManager, messageChannel <-chan string) {
	for {
		select {
		case msg := <-messageChannel:
			manager.sess.Range(func(key, value interface{}) bool {
				session := value.(Session)
				_, err := session.Socket.Write([]byte("SEND: " + msg))
				if err != nil {
					// Handle error
					return false
				}
				return true
			})
		}
	}
}

func WaitForStop() {
	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-termChan // Blocks here until interrupted
	os.Exit(1)
}
