import { useNavigate } from "react-router-dom";
import "./HomePage.css";
import { API_URL } from "../App";
import { useState } from "react";

function HomePage() {
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [loadingMessage, setLoadingMessage] = useState(false);

  const handleCreateGame = async () => {
    try {
      const res = await fetch(`${API_URL}/create-game`, {
        method: "POST",
      });

      if (!res.ok) {
        alert("Failed to create a game");
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

  return (
    <div className="menu">
      <h1>B💣mberman</h1>
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
          <button onClick={() => navigate("/about")}>About</button>
        </div>
      )}
    </div>
  );
}

export default HomePage;
