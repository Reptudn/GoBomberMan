import { useRef, useEffect, useCallback } from "react";

export function useWebSocket(
  url,
  { protocols, onOpen, onMessage, onClose, onError } = {},
) {
  const wsRef = useRef(null);
  const handlersRef = useRef({
    onOpen,
    onMessage,
    onClose,
    onError,
  });

  useEffect(() => {
    handlersRef.current = {
      onOpen,
      onMessage,
      onClose,
      onError,
    };
  }, [onOpen, onMessage, onClose, onError]);

  useEffect(() => {
    if (!url) return;
    const ws = new WebSocket(url, protocols);
    wsRef.current = ws;

    ws.onopen = (event) => {
      handlersRef.current.onOpen?.(event);
    };

    ws.onmessage = (event) => {
      handlersRef.current.onMessage?.(event);
    };

    ws.onclose = (event) => {
      handlersRef.current.onClose?.(event);
    };

    ws.onerror = (event) => {
      handlersRef.current.onError?.(event);
    };

    return () => {
      if (
        ws.readyState === WebSocket.OPEN ||
        ws.readyState === WebSocket.CONNECTING
      ) {
        wsRef.current.close(1000);
      }
      wsRef.current = null;
    };
  }, [url, protocols]);

  const sendSocketMessage = useCallback(
    (data) => {
      if (!wsRef.current) {
        console.error("WebSocket is not connected");
        return;
      }

      if (wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.send(data);
      }
    },
    [wsRef],
  );

  const closeSocket = useCallback(
    (code) => {
      if (!wsRef.current) return;

      if (
        wsRef.current.readyState === WebSocket.OPEN ||
        wsRef.current.readyState === WebSocket.CONNECTING
      ) {
        wsRef.current.close(code);
      }
      wsRef.current = null;
    },
    [wsRef],
  );

  return { sendSocketMessage, closeSocket, socket: wsRef };
}
