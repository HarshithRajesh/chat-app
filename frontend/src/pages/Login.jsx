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

          const userData = {
            id: parseInt(data.id) || parseInt(data.user_id) || parseInt(data.userId),
            email: data.email || email,
            name: data.name || ""
          };
          localStorage.setItem('user', JSON.stringify(userData));
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
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">Welcome Back</h1>
        <p className="auth-subtitle">Sign in to continue to your account</p>
        
        <form className="auth-form" onSubmit={handleLogin}>
          <div className="input-group">
            <label className="input-label">Email</label>
            <input 
              className="input-field"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>
          
          <div className="input-group">
            <label className="input-label">Password</label>
            <input 
              className="input-field"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          
          {error && <div className="error-message">{error}</div>}
          
          <button className="auth-button" type="submit" disabled={isLoading}>
            {isLoading ? 'Signing in...' : 'Sign In'}
          </button>
        </form>
        
        <div className="auth-link">
          Don't have an account? <Link to="/signup">Sign up</Link>
        </div>
      </div>
    </div>
  );
};

export default Login;
