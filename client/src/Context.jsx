import { createContext, useContext, useReducer } from "react";

const StateContext = createContext({});
const StateDispatch = createContext(null);

export function useStateContext() {
  return useContext(StateContext);
}

export function useStateDispatch() {
  return useContext(StateDispatch);
}

export function ContextProvider({ children }) {
  const [state, dispatch] = useReducer((state, action) => {
    switch (action.type) {
      case "update": {
        return { ...state, [action.id]: action.value };
      }

      case "delete": {
        delete state[action.id];
        return state;
      }

      default: {
        throw Error("Unknown action: " + action.type);
      }
    }
  }, {});

  return (
    <StateContext.Provider value={state}>
      <StateDispatch.Provider value={dispatch}>
        {children}
      </StateDispatch.Provider>
    </StateContext.Provider>
  );
}
