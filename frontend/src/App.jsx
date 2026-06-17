import { useState } from "react";

function App() {
  const [username, setUsername] = useState("");
  const [room, setRoom] = useState("general");

  const [connected, setConnected] = useState(false);

  const [socket, setSocket] = useState(null);

  const [messages, setMessages] = useState([]);

  const [messageInput, setMessageInput] = useState("");

  async function connect() {

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
        })
      );

      setConnected(true);
    };

    ws.onmessage = (event) => {

      const message = JSON.parse(event.data);

      setMessages(prev => [...prev, message]);
    };

    ws.onclose = () => {
      setConnected(false);
    };

    setSocket(ws);
  }

  function sendMessage() {

    if (!socket) return;

    socket.send(
      JSON.stringify({
        type: "chat_message",
        username,
        content: messageInput,
      })
    );

    setMessageInput("");
  }

  return (
    <div style={{
      maxWidth: "800px",
      margin: "40px auto",
      fontFamily: "Arial"
    }}>
      <h1>Go Realtime Chat</h1>

      {!connected && (
        <>
          <input
            placeholder="Username"
            value={username}
            onChange={(e) =>
              setUsername(e.target.value)
            }
          />

          <br />
          <br />

          <input
            placeholder="Room"
            value={room}
            onChange={(e) =>
              setRoom(e.target.value)
            }
          />

          <br />
          <br />

          <button onClick={connect}>
            Connect
          </button>
        </>
      )}

      {connected && (
        <>
          <h3>
            Room: {room}
          </h3>

          <div
            style={{
              border: "1px solid #ddd",
              height: "400px",
              overflowY: "auto",
              padding: "10px",
              marginBottom: "10px"
            }}
          >
            {messages.map((msg, index) => (
              <div key={index}>
                <strong>
                  {msg.username}
                </strong>
                : {msg.content}
              </div>
            ))}
          </div>

          <input
            value={messageInput}
            onChange={(e) =>
              setMessageInput(e.target.value)
            }
            placeholder="Type message..."
          />

          <button onClick={sendMessage}>
            Send
          </button>
        </>
      )}
    </div>
  );
}

export default App;