import { useMemo } from "react";
import { useRef, useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import GameField from "../components/game/GameField";
import Chat from "../components/game/Chat";

function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();

  const socketRef = useRef(null);
  const [isConnected, setIsConnected] = useState("connecting");
  const [gameRunning, setGameRunning] = useState(false);
  const [gameData, setGameData] = useState(null);
  const [error, setError] = useState(null);
  const [retryCount, setRetryCount] = useState(0);
  const maxRetries = 10;

  const { url, gameId } = location.state || {};
  const wsUrl = useMemo(() => {
    if (!url) return null;
    return `ws://${url}/ws`;
  }, [url]);

  useEffect(() => {
    if (!url || !gameId) {
      alert("Missing wsUrl or gameId");
      navigate("/");
      return;
    }

    if (
      socketRef.current?.readyState === WebSocket.OPEN ||
      socketRef.current?.readyState === WebSocket.CONNECTING
    ) {
      console.log("Socket already exists, skipping connection");
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
      navigate("/");
    };

    socket.onmessage = (event) => {
      console.log("Socket message: ", event.data);
      const data = JSON.parse(event.data);

      switch (data.type) {
        case "game_state":
          console.log("game state update", data);
          setGameData(data.message);
          break;
        case "game_start":
          setGameRunning(true);
          setError(null);
          break;
        case "game_end":
          setGameRunning(false);
          break;
        case "chat":
          console.log("Chat message: ", data.message);
          break;
        case "success":
          console.log("Success message: ", data.message);
          break;
        case "error":
          setError(data.message || "Unknown error from server");
          break;
        default:
          console.log("Unknown message type: ", data.type);
      }
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
        socket.close(1000);
      }
    };
  }, [wsUrl, gameId, navigate]);

  useEffect(() => {
    if (!gameRunning) return;

    const handleKeyPress = (e) => {
      switch (e.key) {
        case "ArrowUp":
          socketRef.current?.send(
            JSON.stringify({ type: "move", direction: "up" }),
          );
          break;
        case "ArrowDown":
          socketRef.current?.send(
            JSON.stringify({ type: "move", message: { direction: "down" } }),
          );
          break;
        case "ArrowLeft":
          socketRef.current?.send(
            JSON.stringify({ type: "move", message: { direction: "left" } }),
          );
          break;
        case "ArrowRight":
          socketRef.current?.send(
            JSON.stringify({ type: "move", message: { direction: "right" } }),
          );
          break;
        case "Space":
          socketRef.current?.send(JSON.stringify({ type: "place_bomb" }));
          break;
        default:
          break;
      }
    };
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [gameRunning]);

  if (!gameId) {
    return (
      <div>
        <h1>Game Page</h1>
        <p>No Game Id</p>
        <button
          onClick={() => {
            socketRef.current?.close(1000);
            navigate("/");
          }}
        >
          Go Home
        </button>
      </div>
    );
  }

  return (
    <div>
      <h1>Game Page</h1>
      <h3>WebSocket URL: {url || "Not Available"}</h3>
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
      {gameRunning ? (
        <GameField
          fieldData={gameData.field || {}}
          players={gameData.players || {}}
        />
      ) : (
        <p>Waiting for game to start...</p>
      )}
      {/* <Chat socket={socketRef} /> */}
      <p>Joining Game ID: {gameId}</p>
      <button
        onClick={() => {
          socketRef.current?.close(1000);
          navigate("/");
        }}
      >
        Go Home
      </button>
      <button
        onClick={() => {
          const msg = prompt("Enter message");
          socketRef.current?.send(msg || "Hallo");
        }}
      >
        Send Message
      </button>
      <button
        onClick={() => {
          socketRef.current?.send(JSON.stringify({ type: "start_game" }));
        }}
      >
        Start Game
      </button>
    </div>
  );
}

export default GamePage;
