import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";

const SignUp = ({ setIsLoggedIn }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();

  const handleSignUp = async (e) => {
    e.preventDefault();

    if (password !== confirmPassword) {
      setError("Passwords don't match");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      const response = await fetch("http://localhost:8080/signup", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          email: email,
          password: password,
        }),
      });

      if (response.ok) {
        const contentType = response.headers.get("content-type");

        if (contentType && contentType.includes("application/json")) {
          const data = await response.json();
          console.log("Signup successful (JSON):", data);

          // Store user data in localStorage
          const userData = {
            id: parseInt(data.id) || parseInt(data.user_id) || parseInt(data.userId),
            email: data.email || email,
            name: data.name || ""
          };
          localStorage.setItem("user", JSON.stringify(userData));
        } else {
          const textData = await response.text();
          console.log("Signup successful (Text):", textData);

          // Store basic user data - you'll need user ID from backend
          const userData = {
            id: null, // Backend should return user ID
            email: email,
            name: "",
          };
          localStorage.setItem("user", JSON.stringify(userData));
        }

        setIsLoggedIn(true);
        // Redirect to profile page instead of home
        navigate("/profile");
      } else {
        try {
          const errorData = await response.json();
          setError(errorData.message || "Signup failed");
        } catch {
          setError(`Signup failed: ${response.status} ${response.statusText}`);
        }
      }
    } catch (error) {
      console.error("Signup error:", error);
      setError("Something went wrong. Please try again");
    }

    setIsLoading(false);
  };

  return (
    <div>
      <h2>Sign Up for Chat App</h2>
      <form onSubmit={handleSignUp}>
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
        <div>
          <input
            type="password"
            placeholder="Confirm Password"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            required
          />
        </div>
        {error && <div style={{ color: "red" }}>{error}</div>}
        <button type="submit" disabled={isLoading}>
          {isLoading ? "Signing Up..." : "Sign Up"}
        </button>
      </form>
      <p>
        Already have an account? <Link to="/login">Login here</Link>
      </p>
    </div>
  );
};

export default SignUp;
