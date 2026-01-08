cd ./scripts

echo "Building frontend..."
chmod +x build_frontend.sh
./build_frontend.sh
echo "Frontend built successfully!"

echo "Building game server..."
chmod +x build_game_server.sh
./build_game_server.sh
echo "Game server built successfully!"

echo "Building backend..."
chmod +x build_backend.sh
./build_backend.sh
echo "Backend built successfully!"
