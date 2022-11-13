import { Link } from "react-router-dom";

export function OverviewPage() {
  return (
    <div>
      <div className="tile is-ancestor">
        <div className="tile is-parent">
          <Link to="/raw" className="tile is-child box">
            <p className="title">Raw Data</p>
          </Link>
        </div>
        <div className="tile is-parent">
          <Link to="/speedometer" className="tile is-child box">
            <p className="title">Speedometer</p>
          </Link>
        </div>
        <div className="tile is-parent">
          <Link to="/temperature" className="tile is-child box">
            <p className="title">Temperature</p>
          </Link>
        </div>
      </div>
    </div>
  );
}
