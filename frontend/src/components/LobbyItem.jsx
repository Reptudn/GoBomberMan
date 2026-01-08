import { useNavigate } from "react-router-dom";
import { API_URL } from "../App.jsx";
import "./LobbyItem.css";

export default function LobbyItem({ name, id, currentPlayers, maxPlayers }) {
  const navigate = useNavigate();

  const handleLobbyClick = async () => {
    console.log(`Clicked on lobby ${id}`);

    const res = await fetch(`${API_URL}/join-game/${id}`, {
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

    navigate("/game", {
      state: { gameId: gameId, url: url },
    });
  };

  return (
    <div onClick={() => handleLobbyClick()} className="lobby-item">
      <h2>{name || "N/A"}</h2>
      <p>ID: {id || "N/A"}</p>
      <p>
        Players: {currentPlayers || "N/A"}/{maxPlayers || "N/A"}
      </p>
    </div>
  );
}
