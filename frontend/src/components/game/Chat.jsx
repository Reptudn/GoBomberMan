import { useState } from "react";

export default function Chat({ messages, sendChatMessage }) {
  return (
    <div className="chat" style={{ width: "300px", border: "1px solid black" }}>
      <div className="chat-messages">
        {messages.map((msg, index) => (
          <ChatMessage key={index} sender={msg.sender} message={msg.message} />
        ))}
      </div>
      <ChatInput socket={sendChatMessage} />
    </div>
  );
}

function ChatMessage({ sender, message }) {
  return (
    <div className="chat-message">
      {sender}: {message}
    </div>
  );
}

function ChatInput({ sendChatMessage }) {
  const [message, setMessage] = useState("");

  const sendText = () => {
    const chatMessage = { type: "chat_message", content: message };
    socket.send(JSON.stringify(chatMessage));
    setMessage("");
  };

  return (
    <input
      type="text"
      className="chat-input"
      placeholder="Type a message..."
      onSubmit={(e) => {
        e.preventDefault();
        sendText();
      }}
      onChange={(e) => sendChatMessage(e.target.value)}
    />
  );
}
