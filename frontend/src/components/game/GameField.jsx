export default function GameField({ fieldData, players }) {
  if (!fieldData || !players)
    return <p>No game field and player data available</p>;

  const getCellColor = (cell) => {
    if (!cell.type) return "lightgray";
    switch (cell.type) {
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

  // TODO: also render players on the field
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
            {cell !== 0 ? cell : ""}
            {players.forEach((player) => {
              if (
                player.pos.x === index % fieldData.width &&
                player.pos.y === Math.floor(index / fieldData.height)
              ) {
                return "P";
              }
            })}
          </div>
        ))}
      </div>
    </>
  );
}
