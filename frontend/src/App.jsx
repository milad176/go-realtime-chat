import { useEffect, useRef } from "react";

function App() {
  const socketRef = useRef(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/api/ws");

    socketRef.current = socket;

    socket.onopen = () => {
      console.log("connected");

      socket.send(
        JSON.stringify({
          type: "join_room",
          room_id: "general",
        })
      );
    };

    socket.onmessage = (event) => {
      console.log("message:", event.data);
    };

    socket.onclose = () => {
      console.log("disconnected");
    };

    return () => {
      socket.close();
    };
  }, []);

  function sendMessage() {
    if (!socketRef.current) return;

    socketRef.current.send(
      JSON.stringify({
        type: "chat_message",
        content: "hello from react",
      })
    );
  }

  return (
    <div>
      <h1>Go Realtime Chat</h1>

      <button onClick={sendMessage}>
        Send Message
      </button>
    </div>
  );
}

export default App;