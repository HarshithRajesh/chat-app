import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { login } from "../services/api";

function Login() {

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');

    if (!email.trim()) {
      setError('Email required');
      return;
    }
    if (!password.trim()) {
      setError('Password is required');
      return;
    }
    setLoading(true);

    try {
      const responseData = await login({ email, password });
      console.log('Loging in with', responseData);
      navigate('/chat');
    } catch (err) {
      setError('Loging in failed');
      console.log('Log in Failed ', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="Email">Email:</label>
          <input type="email" id="email" value={email} onChange={(e) => setEmail(e.target.value)} disabled={loading} />
        </div>
        <div>
          <label htmlFor="Password">Password:</label>
          <input type="password" id="password" value={password} onChange={(e) => setPassword(e.target.value)} disabled={loading} />
        </div>
        <button type="submit" disabled={loading}>{loading ? 'Loging In...' : 'Log In'}</button>
        {error && <p style={{ color: 'red' }}>{error}</p>}
      </form>
      <p>Placeholder for login page</p>
      <Link to='/chat'>Go to chat</Link>
      <p>Dont have an account<Link to='/signup'>Sing Up</Link></p>
    </div>
  );
}

export default Login;
