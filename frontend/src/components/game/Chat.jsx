import { useState } from "react";

export default function Chat({ messages = [], sendChatMessage }) {
  return (
    <div className="chat" style={{ width: "300px", border: "1px solid black" }}>
      <div className="chat-messages">
        {messages.map((msg, index) => (
          <ChatMessage key={index} sender={msg.sender} message={msg.message} />
        ))}
      </div>
      <ChatInput sendChatMessage={sendChatMessage} />
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

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!message) return;
    const chatMessage = { type: "chat_message", content: message };
    // parent expects to receive a ready-to-send string or object; here we stringify
    if (typeof sendChatMessage === "function") {
      sendChatMessage(JSON.stringify(chatMessage));
    }
    setMessage("");
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        className="chat-input"
        placeholder="Type a message..."
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button type="submit">Send</button>
    </form>
  );
}
