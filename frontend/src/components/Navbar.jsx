import { Link } from "react-router-dom";
export function Navbar() {
  return (
    <>
      <Link to="/">
        <button>Home</button>
      </Link>
      <Link to="/Contact">
        <button>Contact</button>
      </Link>
      <Link to="/Profile">
        <button>Profile</button>
      </Link>
    </>
  );
}
