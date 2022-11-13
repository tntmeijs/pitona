import { useEffect, useRef, useState } from "react";

const WEBSOCKET_URL = "ws://localhost:8080/ws";

export function useWebSocket() {
  const webSocket = useRef(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    if (webSocket.current === null) {
      setIsLoading(true);

      const socket = new WebSocket(WEBSOCKET_URL);

      socket.onopen = () => {
        setIsLoading(false);
        setIsReady(true);
      };

      socket.onclose = () => setIsReady(false);

      socket.onerror = () => {
        setIsLoading(false);
        socket.close();
      };

      webSocket.current = socket;
    }

    const webSocketCurrent = webSocket.current;

    return () => {
      if (webSocketCurrent.readyState === WebSocket.OPEN) {
        webSocketCurrent.close();
      }
    };
  }, []);

  return [isLoading, isReady, webSocket.current];
}
