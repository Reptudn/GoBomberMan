import { useRef, useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";

function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();

  const socketRef = useRef(null);
  const [isConnected, setIsConnected] = useState("connecting");
  const [error, setError] = useState(null);
  const [retryCount, setRetryCount] = useState(0);
  const maxRetries = 10;

  const { url, gameId } = location.state || {};

  let wsUrl = `ws://${url}/ws`;

  useEffect(() => {
    if (!wsUrl || !gameId) {
      alert("Missing wsUrl or gameId");
      navigate("/");
      return;
    }

    console.log("Game ID: ", gameId);
    console.log(`Connecting to game with ID: ${gameId}`);

    const socket = new WebSocket(wsUrl);

    const connectionTimeout = setTimeout(() => {
      if (socket.readyState !== WebSocket.OPEN) {
        console.log("Connection timeout, will retry...");
        socket.close();
      }
    }, 3000);

    socket.onclose = (event) => {
      clearTimeout(connectionTimeout);
      console.log("Socket closed: ", event);

      if (retryCount < maxRetries && event.code !== 1000) {
        console.log(
          `Retrying connection in 2 seconds... (${retryCount + 1}/${maxRetries})`,
        );
        setIsConnected("disconnected");
        setError(
          `Connection failed. Retrying... (${retryCount + 1}/${maxRetries})`,
        );

        setTimeout(() => {
          setRetryCount(retryCount + 1);
        }, 2000);
      } else {
        setIsConnected("disconnected");
        setError(
          `Failed to connect after ${maxRetries} attempts. Game server may not be ready yet.`,
        );
      }
    };

    socket.onerror = (error) => {
      clearTimeout(connectionTimeout);
      console.log("Socket error: ", error);
      setError("WebSocket connection error. Server may still be starting...");
    };

    socket.onmessage = (event) => {
      console.log("Socket message: ", event.data);
    };

    socket.onopen = () => {
      console.log("Socket connected");
      setIsConnected("connected");
      setRetryCount(0);
    };

    socketRef.current = socket;

    return () => {
      clearTimeout(connectionTimeout);
      console.log("Cleaning up socket");
      if (socket.readyState === WebSocket.OPEN) {
        socket.close(1000); // Normal closure
      }
    };
  }, [wsUrl, gameId]);

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
      <h1>Game Page</h1>
      <h3>WebSocket URL: {wsUrl || "Not Available"}</h3>
      <h3>
        WebSocket Status:{" "}
        {isConnected === "connected"
          ? "🟢 Connected"
          : isConnected === "disconnected"
            ? "🔴 Not Connected"
            : isConnected === "connecting"
              ? "🟡 Connecting"
              : "Unknown Status"}
      </h3>
      {error && (
        <div
          style={{
            color: "orange",
            padding: "10px",
            border: "1px solid orange",
            margin: "10px 0",
          }}
        >
          <strong>Status:</strong> {error}
        </div>
      )}
      <p>Joining Game ID: {gameId}</p>
      <button
        onClick={() => {
          if (socketRef.current.state === WebSocket.OPEN) {
            socketRef.current.close(1000);
          }
          navigate("/");
        }}
      >
        Go Home
      </button>
      <button onClick={() => socketRef.current?.send("Hello")}>
        Send Message
      </button>
    </div>
  );
}

export default GamePage;
