import { useState } from "react";

const SignUp = ({ setIsLoggedIn }) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSignUp = async (e) => {
    e.preventDefault();

    // Check if passwords match
    if (password !== confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      // STEP 1: Send signup data to your backend
      const response = await fetch("/signup", {
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
        // STEP 2: Signup successful - automatically log them in
        setIsLoggedIn(true);

        // Optional: Store user info
        // const data = await response.json();
        // localStorage.setItem('userToken', data.token);
        // localStorage.setItem('userEmail', email);
      } else {
        // Signup failed
        const errorData = await response.json();
        setError(errorData.message || "Signup failed");
      }
    } catch (error) {
      setError("Something went wrong. Please try again.");
      console.error("Signup error:", error);
    }

    setIsLoading(false);
  };

  return (
    <div>
      <h2>Create Your Account</h2>

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
          {isLoading ? "Creating Account..." : "Sign Up"}
        </button>
      </form>

      <p>
        Already have an account?
        <a href="/Login">Login here</a>
      </p>
    </div>
  );
};

export default SignUp;
