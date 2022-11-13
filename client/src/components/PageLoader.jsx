import "bulma-pageloader";
import { useStateContext } from "../Context";

export function PageLoader({ text, show }) {
  const state = useStateContext();

  return (
    <div className={`pageloader ${state.showPageLoader ? "is-active" : ""}`}>
      {text && <span className="title">{text}</span>}
    </div>
  );
}
