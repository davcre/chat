# Chat
Chat app with Websockets in Go and Vue.js.

## Installation
Follow the instructions on [this](https://golang.org/doc/install) site to install go.

For saving and restoring messages, you need to install and configure Redis.

Initialize the Go modules with your GitHub repository address:

```bash
go mod init github.com/<your GitHub username>/<project name>
```

Fetch the Go modules:

```bash
go get -u github.com/gomodule/redigo/redis
```

```bash
go get github.com/gorilla/websocket
```

```bash
go get github.com/joho/godotenv
```

Create the .env file to pass the port number to the server and Redis variables.

Compile the server:

```bash
go build
```

## Usage

Start the chat with the following command:

```bash
go run ./
```