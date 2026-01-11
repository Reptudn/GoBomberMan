import { useEffect, useState } from "react";
import LobbyItem from "../components/LobbyItem";
import { useNavigate } from "react-router-dom";

const API_URL = "http://localhost:8080";

export default function LobbyPage() {
  const [lobbies, setLobbies] = useState([]);
  const navigate = useNavigate();

  const getGames = () => {
    fetch(`${API_URL}/list-games`)
      .then((response) => response.json())
      .then((data) => {
        setLobbies(data.games || []);
        console.log(data);
      })
      .catch((error) => alert("Error fetching lobbies:", error));
  };

  useEffect(() => {
    getGames();
  }, []);

  return (
    <div>
      <div
        style={{
          position: "sticky",
          top: 0,
          zIndex: 10,
          backgroundColor: "#242424",
          padding: "1rem",
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
        }}
      >
        <h1>Lobbies Page</h1>
        <button onClick={() => navigate("/")}>Go Back</button>
        <button onClick={() => getGames()}>Refresh</button>
      </div>
      {lobbies.length > 0 ? (
        lobbies.map((lobby) => (
          <LobbyItem
            key={lobby.id || lobby.gameId}
            name={lobby.name}
            id={lobby.id || lobby.gameId}
            currentPlayers={lobby.currentPlayers}
            maxPlayers={lobby.maxPlayers}
          />
        ))
      ) : (
        <div>No lobbies available</div>
      )}
    </div>
  );
}
