import { useEffect } from "react";

function App() {

  useEffect(() => {

    const socket = new WebSocket("ws://localhost:8080/api/ws");

    socket.onopen = () => {
      console.log("connected");

      socket.send("hello from react frontend");
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

  return (
    <div>
      <h1>Go Realtime Chat</h1>
    </div>
  );
}

export default App;