## This project is there for me to learn Kubernetes Deployment and Golang.

## Game
- Min Players is 2
- So far any player can start a game (no admin roles yet)
- players can move with the arrow keys and place bombs with the spacebar
- players can collect powerups
- players can die

- Backend and Game Server written in Go
- Frontent UI written in React

### HOW TO RUN?
If you have kubectl installed and running you can start the backend with:
```bash
make deploy-backend
```

If you want to run the frontend you can use:
(Currently only in Dev mode)
```bash
cd frontend
npm run dev
```

> PS: The devcontainer is not yet tested and might not work at all :D

Have fun lol!

---

https://echo.labstack.com/docs
