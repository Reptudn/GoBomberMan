import { useMemo } from "react";
import { useState, useEffect } from "react";
import { useNavigate, useLocation } from "react-router-dom";
import GameField from "../components/game/GameField";
import { useWebSocket } from "../hooks/useWebSocket";

function GamePage() {
  const navigate = useNavigate();
  const location = useLocation();

  const [isConnected, setIsConnected] = useState("connecting");
  const [gameRunning, setGameRunning] = useState(false);
  const [fieldData, setFieldData] = useState({});
  const [playersData, setPlayersData] = useState({});
  const [error, setError] = useState(null);

  const { url, gameId } = location.state || {};
  const wsUrl = useMemo(() => {
    if (!url) return null;
    const wsProtocol = window.location.protocol === "https:" ? "wss" : "ws";
    return `${wsProtocol}://${window.location.host.split(":")[0]}:${url.split(":")[1]}/ws`;
  }, [url]);

  const { sendSocketMessage, closeSocket, socket } = useWebSocket(wsUrl, {
    onOpen: (event) => {
      console.log("Socket opened", event);
      setIsConnected("connected");
    },
    onMessage: (event) => {
      const data = JSON.parse(event.data);
      console.log("Socket message: ", data);

      switch (data.type) {
        case "game_state":
          console.log("game state field", data.message.field);
          console.log("game state players", data.message.players);
          setFieldData(data.message.field);
          setPlayersData(data.message.players);
          break;
        case "game_start":
          setGameRunning(true);
          setError(null);
          break;
        case "game_over":
          alert("Game Over!", data.message);
          setGameRunning(false);
          navigate("/");
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
    },
    onClose: (event) => {
      console.log("Socket closed", event);
      setIsConnected("disconnected");
      if (!error) {
        setError("WebSocket connection closed by server");
      }
    },
    onError: (error) => {
      console.error("WS ERROR: ", error);
      setError("WebSocket connection error. Server may still be starting...");
      setIsConnected("disconnected");
      // navigate("/");
    },
  });

  console.log("Socket connected", socket);

  useEffect(() => {
    if (!gameRunning) return;

    const handleKeyPress = (e) => {
      switch (e.key) {
        case "ArrowUp":
          sendSocketMessage(
            JSON.stringify({ type: "move", data: { direction: "up" } }),
          );
          break;
        case "ArrowDown":
          sendSocketMessage(
            JSON.stringify({ type: "move", data: { direction: "down" } }),
          );
          break;
        case "ArrowLeft":
          sendSocketMessage(
            JSON.stringify({ type: "move", data: { direction: "left" } }),
          );
          break;
        case "ArrowRight":
          sendSocketMessage(
            JSON.stringify({ type: "move", data: { direction: "right" } }),
          );
          break;
        case " ":
          console.log("Placing bomb");
          sendSocketMessage(JSON.stringify({ type: "place_bomb" }));
          break;
        default:
          break;
      }
    };
    window.addEventListener("keydown", handleKeyPress);

    return () => {
      window.removeEventListener("keydown", handleKeyPress);
    };
  }, [gameRunning, sendSocketMessage]);

  if (!gameId) {
    return (
      <div>
        <h1>Game Page</h1>
        <p>No Game Id</p>
        <button
          onClick={() => {
            closeSocket(1000);
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
          <button
            onClick={() => {
              setError(null);
              setIsConnected("connecting");
              // trigger reconnect by re-creating the ws url: you could force a re-mount or toggle wsUrl
              // simplest: call close then rely on effect to recreate if url didn't change
              closeSocket(1000);
            }}
          >
            Retry / Reconnect
          </button>
        </div>
      )}
      {gameRunning ? (
        <>
          <GameField
            fieldData={fieldData || {}}
            players={
              Array.isArray(playersData)
                ? playersData
                : Object.values(playersData || {})
            }
          />
          <p>Game ID: {gameId}</p>
        </>
      ) : (
        <p>Waiting for game to start...</p>
      )}
      {/* <Chat socket={socketRef} /> */}
      <p>Joining Game ID: {gameId}</p>
      <button
        onClick={() => {
          navigate("/");
        }}
      >
        Go Home
      </button>
      <button
        onClick={() => {
          const msg = prompt("Enter message");
          sendSocketMessage(msg || "Hallo");
        }}
      >
        Send Message
      </button>
      {!gameRunning && (
        <button
          onClick={() => {
            sendSocketMessage(JSON.stringify({ type: "start_game" }));
          }}
        >
          Start Game
        </button>
      )}
    </div>
  );
}

export default GamePage;
