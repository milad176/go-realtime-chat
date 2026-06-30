# 🚀 Go Realtime Chat

A modern real-time chat application built with **Go**, **React**,
**WebSockets**, and **PostgreSQL**.

This project demonstrates how to build a scalable backend capable of
handling real-time communication while persisting chat history to a
relational database. It combines REST APIs for historical data with
WebSockets for instant messaging, following a clean architecture that
separates API handlers, repositories, WebSocket infrastructure, and
database access.

This project is designed to demonstrate:

- concurrent WebSocket client handling with goroutines
- channel-based message broadcasting
- room-based real-time communication
- PostgreSQL message persistence
- graceful connection cleanup
- full-stack communication between Go backend and React frontend

------------------------------------------------------------------------

## ✨ Features

### Real-time Chat

-   ⚡ Instant messaging using WebSockets
-   💬 Room-based chat
-   📜 Persistent message history
-   😀 Emoji support
-   ⌨️ Typing indicator
-   🟢 Online users list
-   ⏰ Message timestamps
-   🔄 Automatic message history loading

### Backend

-   REST API
-   WebSocket server
-   PostgreSQL persistence
-   Repository Pattern
-   Graceful shutdown
-   SQL migrations
-   CORS support
-   Room management
-   Connection management

### Frontend

-   React
-   Vite
-   JavaScript
-   Emoji Picker React
-   Modern dark UI
-   Auto-scroll
-   Persistent username (LocalStorage)
-   Enter-to-send
-   Responsive layout

------------------------------------------------------------------------

# 🛠 Tech Stack

## Backend

-   Go
-   Gorilla WebSocket
-   PostgreSQL
-   pgx
-   UUID

## Frontend

-   React
-   Vite
-   JavaScript

------------------------------------------------------------------------

# 🏗 Architecture Overview

``` text
                +----------------------+
                |      React UI        |
                +----------+-----------+
                           |
                 REST API  |  WebSocket
                           |
          +----------------v----------------+
          |          Go HTTP Server         |
          +----------------+----------------+
                           |
          +----------------v----------------+
          |              Hub               |
          |--------------------------------|
          | Client Management             |
          | Room Management               |
          | Broadcasting                  |
          | Online Users                  |
          | Typing Events                 |
          +----------------+--------------+
                           |
          +----------------v----------------+
          |        Repository Layer         |
          +----------------+----------------+
                           |
          +----------------v----------------+
          |          PostgreSQL             |
          +---------------------------------+
```

------------------------------------------------------------------------

# 📂 Project Structure

```text
go-realtime-chat/
│
├── backend/
│   ├── .env
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── api/
│   │   │   ├── middleware.go
│   │   │   └── server.go
│   │   │
│   │   ├── config/
│   │   │   └── config.go
│   │   │
│   │   ├── db/
│   │   │   ├── migrate.go
│   │   │   └── postgres.go
│   │   │
│   │   ├── handler/
│   │   │   ├── user_handler.go
│   │   │   ├── room_handler.go
│   │   │   └── message_handler.go
│   │   │
│   │   ├── models/
│   │   │   ├── user.go
│   │   │   ├── room.go
│   │   │   └── message.go
│   │   │
│   │   ├── repository/
│   │   │   ├── user_repository.go
│   │   │   ├── room_repository.go
│   │   │   └── message_repository.go
│   │   │
│   │   └── ws/
│   │       ├── client.go
│   │       ├── hub.go
│   │       ├── handler.go
│   │       ├── message.go
│   │       └── broadcast_message.go
│   │
│   └── migrations/
│       ├── 001_init.sql
│       └── 001_second.sql
│
├── frontend/
│   ├── src/
│   ├── public/
│   ├── package.json
│   ├── package-lock.json
│   └── vite.config.js
│
└── README.md
```

------------------------------------------------------------------------

# ⚙️ Running the Project

## Clone

``` bash
git clone https://github.com/yourusername/go-realtime-chat.git
cd go-realtime-chat
```

## Backend

``` bash
cd backend
go mod tidy
go run ./cmd/server/main.go
```

Backend runs on:

``` text
http://localhost:8080
```

## Frontend

``` bash
cd frontend
npm install
npm run dev
```

Frontend runs on:

``` text
http://localhost:5173
```

------------------------------------------------------------------------

# 🌐 REST API

  Method   Endpoint
  -------- ------------------------------
  GET      `/api/health`
  GET      `/api/messages?room=general`
  POST     `/api/users`
  POST     `/api/rooms`
  GET      `/api/rooms`

------------------------------------------------------------------------

# 🔌 WebSocket Protocol

Connect:

``` text
ws://localhost:8080/api/ws
```

Join Room

``` json
{
  "type": "join_room",
  "roomId": "general",
  "username": "milad"
}
```

Chat Message

``` json
{
  "type": "chat_message",
  "username": "milad",
  "content": "Hello everyone!"
}
```

Typing

``` json
{
  "type": "typing",
  "username": "milad",
  "roomId": "general"
}
```

Stop Typing

``` json
{
  "type": "stop_typing",
  "username": "milad",
  "roomId": "general"
}
```

------------------------------------------------------------------------

# 🗄 Database

Current entities:

-   Users
-   Rooms
-   Messages

Messages are stored in PostgreSQL and automatically loaded when users
join a room.

------------------------------------------------------------------------

# 🚀 Roadmap

-   JWT Authentication
-   User Registration & Login
-   User Avatars
-   Private Messaging
-   File Uploads
-   Read Receipts
-   Edit/Delete Messages
-   Docker & Docker Compose
-   Redis Pub/Sub
-   Kubernetes Deployment
-   CI/CD with GitHub Actions
-   Unit & Integration Tests

------------------------------------------------------------------------

# 📚 What I Learned

-   WebSocket communication in Go
-   Goroutines and channels
-   Repository pattern
-   PostgreSQL integration
-   React state management
-   Real-time application architecture
-   REST + WebSocket integration
-   Graceful shutdown
-   Concurrent client management

------------------------------------------------------------------------

# 📄 License

This project is licensed under the MIT License.
