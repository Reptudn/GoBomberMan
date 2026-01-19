export default function GameField({ fieldData = null, players = [], selfID = -1 } = {}) {
  if (!fieldData || !Array.isArray(players))
    return <p>No game field and player data available</p>;

  if (!fieldData.cells || !fieldData.width || !fieldData.height)
    return <p>Invalid game field data</p>;

  const getCellColor = (cell) => {
    const type = cell?.type ?? cell;
    if (type === undefined || type === null) return "lightgray";
    switch (String(type)) {
      case "0":
        return "lightgray"; // EMPTY
      case "1":
        return "saddlebrown"; // WALL_DESTRUCTABLE
      case "2":
        return "dimgray"; // WALL_INDESTRUCTABLE
      case "3": // BOMB
        return "red";
      case "4": // BOMB_EXPLOSION
        return "orange";
      case "5": // POWERUP_BOMB_COUNT_INCREASE
        return "green";
      default:
        return "lightgray";
    }
  };

  const renderCellContent = (cell, index) => {
    // cell may be an object with .type or a primitive
    const type = cell?.type ?? cell;
    // show short label for bombs/players etc.

    // check for player on this cell
    const x = index % fieldData.width;
    const y = Math.floor(index / fieldData.width);
    const playerHere = players.find(
      (p) => p?.pos && p.pos.x === x && p.pos.y === y && p?.alive,
    );
    if (playerHere) {
      if (selfID != -1 && playerHere.id === selfID) return "😼"
      return "👨";
    }

    if (String(type) === "3") return "💣";
    if (String(type) === "4") return "💥";
    if (String(type) === "5") return "⭐";
    if (String(type) === "6") return "🔷";

    // otherwise, nothing visible for empty/other types
    return "";
  };

  return (
    <>
      <p>Game Field:</p>
      <div
        className="game-field"
        style={{
          display: "grid",
          gridTemplateColumns: `repeat(${fieldData.width}, 50px)`,
          gridTemplateRows: `repeat(${fieldData.height}, 50px)`,
          gap: "2px",
        }}
      >
        {fieldData.cells.map((cell, index) => (
          <div
            key={index}
            className="game-cell"
            style={{
              width: "50px",
              height: "50px",
              backgroundColor: getCellColor(cell),
              border: "1px solid black",
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              fontSize: "24px",
              color: "white",
            }}
          >
            {renderCellContent(cell, index)}
          </div>
        ))}
      </div>
    </>
  );
}
