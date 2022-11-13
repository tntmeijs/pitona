export function Modal({ show, closable, children, onDismissed }) {
  const onClose = () => {
    if (onDismissed) {
      onDismissed();
    }
  };

  return (
    <div className={`modal ${show ? "is-active" : ""}`}>
      <div
        className="modal-background"
        onClick={() => {
          if (closable) {
            onClose();
          }
        }}
      ></div>
      <div className="modal-content">{children}</div>
      {closable && (
        <button className="modal-close is-large" onClick={onClose}></button>
      )}
    </div>
  );
}
