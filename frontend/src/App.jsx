import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import LobbiesPage from "./pages/LobbiesPage";
import GamePage from "./pages/GamePage";
import "./App.css";

export const API_URL = `http://${window.location.host.split(":")[0]}:8080`;
// export const API_URL = `http://172.17.252.15:8080`;

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/lobbies" element={<LobbiesPage />} />
        <Route path="/game" element={<GamePage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
