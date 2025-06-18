import { Link } from "react-router-dom";
export function Navbar() {
  return (
    <>
      <Link to="/">
        <button>Home</button>
      </Link>
      <Link to="/"></Link>
    </>
  );
}
