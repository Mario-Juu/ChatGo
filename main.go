package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type Client struct{
	Nickname string
	conn *websocket.Conn
	ctx context.Context
}

type Message struct{
	From string `json:"from"`
	Content string `json:"content"`
	SentAt string `json:"sentAt"`
}

var(
	clients map[*Client]bool = make(map[*Client]bool)
	joinCh chan *Client = make(chan *Client)
	broadcastCh chan Message = make(chan Message)
)


func wsHandler(w http.ResponseWriter, r *http.Request){
	nickname := r.URL.Query().Get("nickname")
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil{
		log.Panic(err)
	}
	
	
	go joiner()
	go broadcast()

	newClient := Client{
		Nickname: nickname,
		conn: conn,
		ctx: r.Context(),
	}


	joinCh <- &newClient

	
	for{
		_, data, err := newClient.conn.Read(r.Context())
		if err != nil{
			log.Print("Encerrando a conexão.")
			delete(clients, &newClient)
			broadcastCh <- Message{From: newClient.Nickname, Content: newClient.Nickname + " saiu do chat.", SentAt: time.Now().Format("2006-01-02 15:04:05")}
			break
		}
		log.Println(string(data))

		
		var msgRec Message 
		json.Unmarshal(data, &msgRec)

		broadcastCh <- Message{From: msgRec.From, Content: msgRec.Content, SentAt: time.Now().Format("2006-01-02 15:04:05")}
	}
	}


func broadcast(){
	for newMsg := range broadcastCh{
		
	for client := range clients{
		msg, err := json.Marshal(newMsg)
		if err != nil{
			log.Print(err)
		}
		client.conn.Write(client.ctx, websocket.MessageText, msg)
	}

	}
}


func joiner(){
	for newClient := range joinCh{
		clients[newClient] = true

	for client := range clients{
		broadcastCh <- Message{From: client.Nickname, Content: client.Nickname + " entrou no chat.", SentAt: time.Now().Format("2006-01-02 15:04:05")}
	}
	}
}

func main(){

	http.HandleFunc("/ws", wsHandler)

	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/clients", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		var res []*Client
		for c := range clients{
			res = append(res, c)
		}
		json.NewEncoder(w).Encode(res)
	})

	log.Print("Starting server at localhost:8080")
	http.ListenAndServeTLS(":8080", "server.crt", "server.key", nil)
}