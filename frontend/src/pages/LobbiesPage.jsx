import { useEffect, useState } from "react";
import LobbyItem from "../components/LobbyItem";
import { useNavigate } from "react-router-dom";

export default function LobbyPage() {
  const [lobbies, setLobbies] = useState([]);
  const navigate = useNavigate();

  useEffect(() => {
    fetch("/api/list-games")
      .then((response) => response.json())
      .then((data) => setLobbies(data.games || []))
      .catch((error) => console.error("Error fetching lobbies:", error));
  }, []);

  return (
    <div>
      <h1>Lobbies Page</h1>
      {lobbies.length > 0 ? (
        lobbies.map((lobby) => (
          <LobbyItem
            key={lobby.gameId || lobby.id}
            name={lobby.name}
            id={lobby.gameId || lobby.id}
            currentPlayer={lobby.currentPlayer}
            maxPlayers={lobby.maxPlayers}
          />
        ))
      ) : (
        <div>No lobbies available</div>
      )}
      <button onClick={() => navigate("/")}>Go Back</button>
    </div>
  );
}
