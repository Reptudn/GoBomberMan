import { BrowserRouter, Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import LobbiesPage from "./pages/LobbiesPage";
import GamePage from "./pages/GamePage";
import AboutPage from "./pages/AboutPage";
import "./App.css";

export const API_URL = "http://localhost:8080";

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/lobbies" element={<LobbiesPage />} />
        <Route path="/game" element={<GamePage />} />
        <Route path="/about" element={<AboutPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
