package main

import (
	"encoding/json"
	"log"
	"net/http"

	"nhooyr.io/websocket"
)

type Client struct{
	Nickname string
	conn *websocket.Conn
}

var(
	clients map[*Client]bool = make(map[*Client]bool)
)


func wsHandler(w http.ResponseWriter, r *http.Request){
	nickname := r.URL.Query().Get("nickname")
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil{
		log.Panic(err)
	}
	newClient := Client{
		Nickname: nickname,
		conn: conn,
	}
	clients[&newClient] = true
	for client := range clients{
		client.conn.Write(r.Context(), websocket.MessageText, []byte("O usuário " + nickname + " se conectou."))
	}

	for{
		_, data, err := newClient.conn.Read(r.Context())
		if err != nil{
			log.Print("Encerrando a conexão.")
			delete(clients, &newClient)
			return
		}
		log.Println(string(data))

		for client := range clients{
		client.conn.Write(r.Context(), websocket.MessageText, data)
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

	http.ListenAndServe(":8080", nil)
	log.Print("Starting server at localhost:8080")
}