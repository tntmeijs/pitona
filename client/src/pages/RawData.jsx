import { useState } from "react";
import { WebSocketWrapper } from "../components/WebSocketWrapper";

export function RawData() {
  const [pids, setPids] = useState([]);

  const onMessage = (msg) => {
    setPids(msg.Pids);
  };

  return (
    <WebSocketWrapper onMessage={onMessage}>
      <h1>Data from WebSocket connection:</h1>
      <table className="table">
        <thead>
          <tr>
            <th>Mode</th>
            <th>PID</th>
            <th>Data</th>
          </tr>
        </thead>
        <tbody>
          {pids.map((pid, index) => {
            return (
              <tr key={index}>
                <td>{pid.Mode}</td>
                <td>{pid.Pid}</td>
                <td>{pid.Data}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </WebSocketWrapper>
  );
}
