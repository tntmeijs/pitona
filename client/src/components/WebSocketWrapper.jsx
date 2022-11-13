import { useEffect, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useStateDispatch } from "../Context";
import { Modal } from "./Modal";

const PAGE_LOADER_DELAY_MS = 750;
const WEBSOCKET_URL = "ws://localhost:8080/ws";

export function WebSocketWrapper({ onMessage, children }) {
  const webSocket = useRef(null);
  const [isLoading, setIsLoading] = useState(true);
  const [webSocketReady, setWebSocketReady] = useState(false);

  const navigate = useNavigate();
  const dispatch = useStateDispatch();

  // Responsible for managing the WebSocket connection's life cycle
  useEffect(() => {
    if (webSocket.current === null) {
      dispatch({ type: "update", id: "showPageLoader", value: true });
      setIsLoading(true);

      const socket = new WebSocket(WEBSOCKET_URL);
      webSocket.current = socket;
    }

    const webSocketCurrent = webSocket.current;

    return () => {
      if (webSocketCurrent.readyState === WebSocket.OPEN) {
        console.log("Unmount - closing connection");
        webSocketCurrent.close();
      }
    };
  }, [dispatch]);

  useEffect(() => {
    if (webSocket.current === null) {
      return;
    }

    webSocket.current.onopen = () => {
      setTimeout(
        () =>
          dispatch({
            type: "update",
            id: "showPageLoader",
            value: false,
          }),
        PAGE_LOADER_DELAY_MS
      );

      setIsLoading(false);
      setWebSocketReady(true);
    };

    webSocket.current.onerror = () => {
      setWebSocketReady(false);
      setIsLoading(false);
    };

    webSocket.current.onclose = () => {
      setWebSocketReady(false);
    };

    webSocket.current.onmessage = (event) => onMessage(JSON.parse(event.data));
  }, [dispatch, onMessage]);

  if (isLoading) {
    return null;
  }

  if (!webSocketReady) {
    return (
      <Modal show={true}>
        <div className="notification is-danger is-light block">
          <p className="title">Unable to connect to the server 😢</p>
          <p className="mb-4">
            We could not establish a WebSocket connection with the server.
            Please return to the home page and try again.
          </p>
          <button className="button is-danger" onClick={() => navigate("/")}>
            Back
          </button>
        </div>
      </Modal>
    );
  }

  return (
    <div>
      <button
        onClick={() =>
          webSocket.current.send(
            JSON.stringify({
              pids: [
                {
                  mode: "01",
                  pid: "01",
                  responseSizeInBytes: 1,
                },
                {
                  mode: "01",
                  pid: "46",
                  responseSizeInBytes: 1,
                },
              ],
            })
          )
        }
      >
        ping server
      </button>
      {children}
    </div>
  );
}
