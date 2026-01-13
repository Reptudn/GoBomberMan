import { useNavigate } from "react-router-dom";
import "./HomePage.css";
import { API_URL } from "../App";
import { useEffect, useState } from "react";

function HomePage() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [loadingMessage, setLoadingMessage] = useState(false);
  const [serverOnline, setServerOnline] = useState(false);

  const handleCreateGame = async () => {
    try {
      const res = await fetch(`${API_URL}/create-game`, {
        method: "POST",
      });

      if (!res.ok) {
        alert("Failed to create a game");
        setLoading(false);
        setLoadingMessage("");
        return;
      }

      const data = await res.json();

      const url = data.url;
      const gameId = data.gameId;

      console.log("Game created successfully");

      navigate("/game", {
        state: { gameId: gameId, url: url },
      });
    } catch (error) {
      console.error("Error creating game:", error);
      alert("Failed to create a game");
      setLoading(false);
      setLoadingMessage("");
    }
  };

  const handleJoinGame = async () => {
    const gameIdInput = prompt("Enter Game ID:");

    if (!gameIdInput) {
      setLoading(false);
      setLoadingMessage("");
      return;
    }

    const res = await fetch(`${API_URL}/join-game/${gameIdInput}`, {
      method: "POST",
    });

    if (!res.ok) {
      const data = await res.json();
      alert("Failed to join the game: ", data.status);
      return;
    }

    const data = await res.json();

    const url = data.url;
    const gameId = data.gameId;

    setTimeout(() => {
      navigate(`/game`, {
        state: { gameId: gameId, url: url },
      });
    }, 500);
  };

  // useEffect(() => {
  //   const checkServerStatus = async () => {
  //     try {
  //       const response = await fetch(`${API_URL}/ping`);
  //       if (response.ok) {
  //         setServerOnline(true);
  //       } else {
  //         setServerOnline(false);
  //       }
  //     } catch {
  //       setServerOnline(false);
  //     }
  //   };

  //   setInterval(async () => {
  //     if (!serverOnline) await checkServerStatus();
  //   }, 1000);
  // }, [serverOnline]);

  return (
    <div className="menu">
      <h1>B💣mberman</h1>
      {serverOnline ? <p>Server is online</p> : <p>Server is offline</p>}
      {loading ? (
        <p>{loadingMessage}</p>
      ) : (
        <div className="menu-buttons">
          <button
            onClick={async () => {
              setLoading(true);
              setLoadingMessage("Creating game...");
              await handleCreateGame();
            }}
          >
            Create Game
          </button>
          <button
            onClick={() => {
              setLoading(true);
              setLoadingMessage("Joining game...");
              handleJoinGame();
            }}
          >
            Join Game
          </button>
          <button onClick={() => navigate("/lobbies")}>Browse Lobbies</button>
          <p>by jkauker</p>
        </div>
      )}
    </div>
  );
}

export default HomePage;
