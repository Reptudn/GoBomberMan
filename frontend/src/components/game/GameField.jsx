export default function GameField({ fieldData, players }) {
  // TODO: also render players on the field
  return (
    <div
      className="game-field"
      style={{
        display: "grid",
        gridTemplateColumns: `repeat(${fieldData[0].length}, 50px)`,
        gridTemplateRows: `repeat(${fieldData.length}, 50px)`,
        gap: "2px",
      }}
    >
      {fieldData.flat().map((cell, index) => (
        <div
          key={index}
          className="game-cell"
          style={{
            width: "50px",
            height: "50px",
            backgroundColor: cell === 0 ? "lightgray" : "steelblue",
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
              player.pos.x === index % fieldData[0].length &&
              player.pos.y === Math.floor(index / fieldData[0].length)
            ) {
              return "P";
            }
          })}
        </div>
      ))}
    </div>
  );
}
