import { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useStateDispatch } from "../Context";
import { useWebSocket } from "../hooks/UseWebSocket";
import { Modal } from "./Modal";

const PAGE_LOADER_DELAY_MS = 750;

export function WebSocketWrapper({ onMessage, children }) {
  const navigate = useNavigate();
  const dispatch = useStateDispatch();

  // TODO: move this to the page that uses it - this wrapper should be removed
  const [isLoading, isReady, webSocket] = useWebSocket();

  useEffect(() => {
    setTimeout(
      () =>
        dispatch({
          type: "update",
          id: "showPageLoader",
          value: isLoading,
        }),
      PAGE_LOADER_DELAY_MS
    );
  }, [dispatch, isLoading]);

  if (isLoading) {
    return null;
  }

  if (!isReady) {
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
          webSocket.send(
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
