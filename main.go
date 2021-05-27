package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
 
  "github.com/joho/godotenv"
  "github.com/gorilla/websocket"
  "github.com/go-redis/redis"
)

type Message struct {
  Email string `json:"email"`
  Username string `json:"username"`
  Message string `json:"message"`
}

var rdb *redis.Client

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader {
  CheckOrigin: func(r *http.Request) bool {
    return true
  },
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    log.Fatal(err)
  }
  
  defer ws.Close()
  clients[ws] = true

  if rdb.Exists("messages").Val() != 0 {
    sendPrevMsgs(ws)
  }
  
  for {
    var msg Message
    err := ws.ReadJSON(&msg)
    if err != nil {
      log.Println(err)
      delete(clients, ws)
      break
    }
    broadcast <- msg
  }
}

func sendPrevMsgs(ws *websocket.Conn) {
  messages, err := rdb.LRange("messages", 0, -1).Result()
  if err != nil {
    panic(err)
  }

  for _, message := range messages {
    var msg Message
    json.Unmarshal([]byte(message), &msg)
    messageClient(ws, msg)
  }
}

func handleMessages() {
  for {
    msg := <-broadcast
    saveInDb(msg)
    messageClients(msg)
  }
}

func saveInDb(msg Message) {
  json, err := json.Marshal(msg)
  if err != nil {
    panic(err)
  }

  if err := rdb.RPush("messages", json).Err(); err != nil {
    panic(err)
  }
}

func messageClients(msg Message) {
  for client := range clients {
    messageClient(client, msg)
  }
}

func messageClient(client *websocket.Conn, msg Message) {
  err := client.WriteJSON(msg)
  if err != nil {
    log.Println(err)
    client.Close()
    delete(clients, client)
  }
}

func main() {
  var appConfig map[string]string
  appConfig, err := godotenv.Read()

  if err != nil {
    log.Fatal("Error reading .env file")
  }

  port := appConfig["PORT"]
  redisURL := fmt.Sprintf("redis://:%s@%s:%s/1",
    appConfig["REDIS_PASSWD"],
    appConfig["REDIS_ADDR"],
    appConfig["REDIS_PORT"],
  )
  opt, err := redis.ParseURL(redisURL)
  if err != nil {
    log.Fatal("Error parsing redis URL")
  }
  rdb = redis.NewClient(opt) 

  log.Println("Staring development server at http://127.0.0.1:"+port)
  log.Println("Quit the server with CTRL-C") 

  http.HandleFunc("/ws", handleConnections)
  go handleMessages()

  fs := http.FileServer(http.Dir("./public"))
  http.Handle("/", fs)

  log.Fatal(http.ListenAndServe(":"+port, nil))
}
