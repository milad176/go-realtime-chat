import { useState, useRef, useEffect } from "react";
import EmojiPicker from "emoji-picker-react";

function App() {
  const [username, setUsername] = useState(
    localStorage.getItem("username") || ""
  );

  const [room, setRoom] = useState("general");

  const [connected, setConnected] = useState(false);
  const [socket, setSocket] = useState(null);

  const [messages, setMessages] = useState([]);
  const [messageInput, setMessageInput] = useState("");

  const [showEmojiPicker, setShowEmojiPicker] = useState(false);

  const [onlineUsers, setOnlineUsers] = useState([]);
  const [typingUser, setTypingUser] = useState("");

  const messagesEndRef = useRef(null);
  const typingTimeoutRef = useRef(null);
  const typingTimeout = useRef(null);
  const isTyping = useRef(false);

  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({
      behavior: "smooth",
    });
  }, [messages]);

  function onEmojiClick(emojiData) {
    setMessageInput((prev) => prev + emojiData.emoji);
  }

  async function connect() {
    if (!username.trim()) {
      alert("Please enter a username");
      return;
    }

    localStorage.setItem("username", username);

    const response = await fetch(
      `http://localhost:8080/api/messages?room=${room}`
    );

    const history = await response.json();

    setMessages(history);

    const ws = new WebSocket(
      "ws://localhost:8080/api/ws"
    );

    ws.onopen = () => {
      ws.send(
        JSON.stringify({
          type: "join_room",
          roomId: room,
          username: username,
        })
      );

      setConnected(true);
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);

      switch (message.type) {
        case "online_users":
          setOnlineUsers(message.users || []);
          break;

        case "typing":
          setTypingUser(message.username);

          setTimeout(() => {
            setTypingUser("");
          }, 1500);
          break;

        case "chat_message":
          setMessages((prev) => [...prev, message]);
          break;

        default:
          break;
      }
    };

    ws.onclose = () => {
      setConnected(false);
    };

    setSocket(ws);
  }

  function sendMessage() {
    if (!socket) return;

    if (!messageInput.trim()) return;

    socket.send(
      JSON.stringify({
        type: "chat_message",
        username: username,
        content: messageInput,
      })
    );

    setMessageInput("");
    setShowEmojiPicker(false);
  }

  function handleKeyDown(event) {
    if (event.key === "Enter") {
      sendMessage();
    }
  }

function handleTyping(value) {
    setMessageInput(value);

    if (!socket) return;

    if (!isTyping.current) {
        socket.send(
            JSON.stringify({
                type: "typing",
                username,
                roomId: room,
            })
        );

        isTyping.current = true;
    }

    clearTimeout(typingTimeout.current);

    typingTimeout.current = setTimeout(() => {
        socket.send(
            JSON.stringify({
                type: "stop_typing",
                username,
                roomId: room,
            })
        );

        isTyping.current = false;
    }, 2000);
}

  return (
    <div
      style={{
        minHeight: "100vh",
        background: "#121212",
        color: "white",
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        padding: "20px",
        fontFamily: "Arial",
      }}
    >
      <div
        style={{
          width: "900px",
          background: "#1e1e1e",
          borderRadius: "12px",
          padding: "20px",
          boxShadow: "0 0 20px rgba(0,0,0,0.4)",
        }}
      >
        <h1
          style={{
            textAlign: "center",
            marginBottom: "20px",
          }}
        >
          Go Realtime Chat 🚀
        </h1>

        {!connected && (
          <>
            <input
              placeholder="Username"
              value={username}
              onChange={(e) =>
                setUsername(e.target.value)
              }
              style={{
                width: "100%",
                padding: "12px",
                marginBottom: "10px",
                borderRadius: "8px",
                border: "none",
              }}
            />

            <input
              placeholder="Room"
              value={room}
              onChange={(e) =>
                setRoom(e.target.value)
              }
              style={{
                width: "100%",
                padding: "12px",
                marginBottom: "20px",
                borderRadius: "8px",
                border: "none",
              }}
            />

            <button
              onClick={connect}
              style={{
                width: "100%",
                padding: "12px",
                border: "none",
                borderRadius: "8px",
                cursor: "pointer",
                fontSize: "16px",
              }}
            >
              Connect
            </button>
          </>
        )}

        {connected && (
          <>
            <div
              style={{
                display: "flex",
                justifyContent: "space-between",
                marginBottom: "15px",
              }}
            >
              <div>
                👤 {username}
              </div>

              <div>
                💬 {room}
              </div>

              <div>
                🟢 {onlineUsers.length} online
              </div>
            </div>

            <div
              style={{
                marginBottom: "10px",
                color: "#aaa",
                fontSize: "14px",
              }}
            >
              {onlineUsers.join(", ")}
            </div>

            <div
              style={{
                height: "450px",
                overflowY: "auto",
                border: "1px solid #333",
                borderRadius: "10px",
                padding: "15px",
                marginBottom: "15px",
                background: "#181818",
              }}
            >
              {messages.map((msg, index) => {
                const isMine =
                  msg.username === username;

                const previous =
                  messages[index - 1];

                const grouped =
                  previous &&
                  previous.username ===
                    msg.username;

                return (
                  <div
                    key={index}
                    style={{
                      display: "flex",
                      justifyContent: isMine
                        ? "flex-end"
                        : "flex-start",
                      marginBottom: grouped
                        ? "4px"
                        : "12px",
                    }}
                  >
                    <div
                      style={{
                        maxWidth: "70%",
                        background: isMine
                          ? "#2d5fff"
                          : "#2b2b2b",
                        padding: "10px",
                        borderRadius: "12px",
                      }}
                    >
                      {!grouped && (
                        <div
                          style={{
                            fontWeight: "bold",
                            marginBottom: "5px",
                          }}
                        >
                          {msg.username}
                        </div>
                      )}

                      <div>
                        {msg.content}
                      </div>

                      <div
                        style={{
                          fontSize: "11px",
                          opacity: 0.7,
                          marginTop: "5px",
                          textAlign: "right",
                        }}
                      >
                        {msg.created_at
                          ? new Date(
                              msg.created_at
                            ).toLocaleTimeString()
                          : ""}
                      </div>
                    </div>
                  </div>
                );
              })}

              <div
                style={{
                  color: "#aaa",
                  fontStyle: "italic",
                  height: "20px",
                }}
              >
                {typingUser &&
                  typingUser !== username &&
                  `${typingUser} is typing...`}
              </div>

              <div ref={messagesEndRef} />
            </div>

            <div
              style={{
                position: "relative",
              }}
            >
              {showEmojiPicker && (
                <div
                  style={{
                    position: "absolute",
                    bottom: "60px",
                    right: "0",
                    zIndex: 1000,
                  }}
                >
                  <EmojiPicker
                    onEmojiClick={onEmojiClick}
                  />
                </div>
              )}

              <div
                style={{
                  display: "flex",
                  gap: "10px",
                }}
              >
                <input
                  value={messageInput}
                  onChange={(e) => {
                    handleTyping(e.target.value)
                  }}
                  onKeyDown={handleKeyDown}
                  placeholder="Type message..."
                  style={{
                    flex: 1,
                    padding: "12px",
                    borderRadius: "8px",
                    border: "none",
                  }}
                />

                <button
                  onClick={() =>
                    setShowEmojiPicker(
                      (prev) => !prev
                    )
                  }
                  style={{
                    padding: "12px",
                    border: "none",
                    borderRadius: "8px",
                    cursor: "pointer",
                  }}
                >
                  😀
                </button>

                <button
                  onClick={sendMessage}
                  style={{
                    padding: "12px 20px",
                    border: "none",
                    borderRadius: "8px",
                    cursor: "pointer",
                  }}
                >
                  Send
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
}

export default App;