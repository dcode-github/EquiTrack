import React, { createContext, useContext, useRef } from "react";

const WebSocketContext = createContext();

export const WebSocketProvider = ({ children }) => {
  const wsRef = useRef(null);

  const connect = (url) => {
    const socket = new WebSocket(url);
    wsRef.current = socket;
    return socket;
  };

  const disconnect = () => {
    if (wsRef.current) {
      wsRef.current.close();
      wsRef.current = null;
    }
  };

  return (
    <WebSocketContext.Provider value={{ connect, disconnect, wsRef }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => useContext(WebSocketContext);
