import { useRef, useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";

function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();

  const socketRef = useRef(null);
  const [isConnected, setIsConnected] = useState(false);

  const { wsUrl, gameId } = location.state || {};

  useEffect(() => {
    if (!wsUrl || !gameId) {
      console.error("Missing wsUrl or gameId");
      return;
    }

    console.log("Game ID: ", gameId);
    console.log(`Connecting to game with ID: ${gameId}`);

    const socket = new WebSocket(wsUrl);

    socket.onclose = (event) => {
      console.log("Socket closed: ", event);
      setIsConnected(false);
      navigate("/");
    };

    socket.onerror = (error) => {
      console.log("Socket error: ", error);
      alert("Socket error occurred");
    };

    socket.onmessage = (event) => {
      console.log("Socket message: ", event.data);
    };

    socket.onopen = () => {
      console.log("Socket connected");
      setIsConnected(true);
    };

    socketRef.current = socket;

    return () => {
      console.log("Cleaning up socket");
      socket.close();
    };
  }, [wsUrl, gameId, navigate]);

  if (!gameId) {
    return (
      <div>
        <h1>Game Page</h1>
        <p>No Game Id</p>
        <button onClick={() => navigate("/")}>Go Home</button>
      </div>
    );
  }

  return (
    <div>
      <h1>Game Page {isConnected ? "Connected" : "Not Connected"}</h1>
      <p>Joining Game ID: {gameId}</p>
    </div>
  );
}

export default GamePage;
