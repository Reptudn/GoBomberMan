import { useNavigate } from "react-router-dom";

export default function LobbyItem(name, id, currentPlayer, maxPlayers) {
  const navigate = useNavigate();

  const handleLobbyClick = async () => {
    console.log(`Clicked on lobby ${id}`);

    const res = await fetch(`/api/join-game/${id}`);

    if (!res.ok) {
      alert("Failed to join the game");
      return;
    }

    const data = await res.json();

    navigate("/game", {
      state: { gameId: data.gameId, socketUrl: data.socketUrl },
    });
  };

  return (
    <div onClick={() => handleLobbyClick()}>
      <h2>{name}</h2>
      <p>ID: {id}</p>
      <p>
        Players: {currentPlayer}/{maxPlayers}
      </p>
    </div>
  );
}
