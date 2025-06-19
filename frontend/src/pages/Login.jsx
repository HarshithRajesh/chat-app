import { useState } from "react";
import { Link } from "react-router-dom";

const Login = ({ setIsLoggedIn }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleLogin = async (e) => {
    e.preventDefault();

    setIsLoading(true);
    setError("");

    console.log("Starting login request...");

    try {
      console.log("Sending request to:", "http://localhost:8080/Login");
      console.log("Request body:",{email,password})
      const response = await fetch("http://localhost:8080/Login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          password: password,
        }),
      });

      console.log("Response received:",response);
      console.log("Response status:", response.status);

      if (response.ok) {
        
        const contentType = response.headers.get("content-type");
        console.log("Content-Type:", contentType);

        if (contentType && contentType.includes("application/json")) {
          const data = await response.json();
          console.log("Login successful (JSON):", data);
        } else {
         
          const textData = await response.text();
          console.log("Login successful (Text):", textData);
        }

        console.log("Setting isLoggedIn to true");
        setIsLoggedIn(true);
      } else {
        
        try {
          const errorData = await response.json();
          setError(errorData.message || "Login failed");
        } catch {
          
          setError(`Login failed: ${response.status} ${response.statusText}`);
        }
      }
    } catch (error) {
      console.error("Fetch error:", error);
      setError("Something went wrong. Please try again");
      
    }

    setIsLoading(false);
    console.log("Login request completed");
  };

  return (
    <div>
      <h2>Login to the Chat</h2>
      <form onSubmit={handleLogin}>
        <div>
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>

        <div>
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
          />
        </div>
        {error && <div style={{ color: "red" }}>{error}</div>}
        <button type="submit" disabled={isLoading}>
          {isLoading ? "Logging In..." : "Login"}
        </button>
      </form>

      <p>
        Don't have an account? <Link to="/signup">Sign Up here</Link>
      </p>
    </div>
  );
};

export default Login;
