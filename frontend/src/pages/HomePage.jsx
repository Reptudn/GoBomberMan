import { useNavigate } from "react-router-dom";
import "./HomePage.css";

function HomePage() {
  const navigate = useNavigate();

  const handleCreateGame = async () => {
    const res = await fetch("/api/create-game", {
      method: "POST",
    });

    if (!res.ok) {
      alert("Failed to create a game");
      return;
    }

    const data = await res.json();

    navigate(`/game/${data.gameId}`, {
      state: { gameId: data.gameId, wsUrl: data.wsUrl, port: data.port },
    });
  };

  const handleJoinGame = () => {
    const gameId = prompt("Enter Game ID:");
    if (gameId) {
      navigate(`/game/${gameId}`);
    }
  };

  return (
    <div className="menu">
      <h1>B💣mberman</h1>
      <div className="menu-buttons">
        <button onClick={() => handleCreateGame()}>Create Game</button>
        <button onClick={() => handleJoinGame()}>Join Game</button>
        <button onClick={() => navigate("/lobbies")}>Browse Lobbies</button>
        <button onClick={() => navigate("/about")}>About</button>
      </div>
    </div>
  );
}

export default HomePage;
