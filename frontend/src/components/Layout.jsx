import { Outlet, Link } from "react-router-dom";

const Layout = ({ setIsLoggedIn }) => {
  // Function to handle logout
  const handleLogout = () => {
    setIsLoggedIn(false);
  };

  return (
    <div>
      <nav style={{ padding: "1rem", backgroundColor: "#f0f0f0" }}>
        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
          }}
        >
          <div>
            <Link to="/" style={{ marginRight: "1rem" }}>
              Home
            </Link>
            <Link to="/contact" style={{ marginRight: "1rem" }}>
              Contact
            </Link>
            <Link to="/profile" style={{ marginRight: "1rem" }}>
              Profile
            </Link>
          </div>

          <div>
            <button onClick={handleLogout}>Logout</button>
          </div>
        </div>
      </nav>

      <main style={{ padding: "1rem" }}>
        <Outlet />
      </main>
    </div>
  );
};

export default Layout;
