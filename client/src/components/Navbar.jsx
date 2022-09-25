import { useState } from "react";
import { Link } from "react-router-dom";

export function Navbar() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <nav className="navbar" role="navigation" aria-label="main navigation">
      <div className="navbar-brand">
        <Link className="navbar-item" to="/">
          <img
            alt="logo"
            src="https://bulma.io/images/bulma-logo.png"
            width="112"
            height="28"
          />
        </Link>

        <a
          role="button"
          className={`navbar-burger ${isOpen ? "is-active" : ""}`}
          aria-label="menu"
          aria-expanded="false"
          data-target="navbarBasicExample"
          onClick={() => setIsOpen(!isOpen)}
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div className={`navbar-menu ${isOpen ? "is-active" : ""}`}>
        <div className="navbar-start">
          <Link className="navbar-item" to="/">
            Overview
          </Link>

          <Link className="navbar-item" to="/system">
            System Info
          </Link>

          <div className="navbar-item has-dropdown is-hoverable">
            <a className="navbar-link">More</a>

            <div className="navbar-dropdown">
              <Link className="navbar-item" to="/debugging">
                Debugging
              </Link>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
