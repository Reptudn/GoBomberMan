package shared

type GameState int

const (
    GameStateWaiting GameState = iota  // 0
    GameStatePlaying                    // 1
    GameStatePaused                     // 2
    GameStateFinished                   // 3
)