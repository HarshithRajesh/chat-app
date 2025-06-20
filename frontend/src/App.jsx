import React from "react";
import './App.css'; // Add this import
import { HashRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import SignUp from "./pages/SignUp";
import Home from "./pages/Home";
import Contact from "./pages/Contacts.jsx";
import Profile from "./pages/Profile.jsx";
import Chat from "./pages/Chat.jsx"; // Make sure this import is correct
import Layout from "./components/Layout";
import { Navigate } from "react-router-dom";
import { useState } from "react";


const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(true); // Set to true for testing

  if (!isLoggedIn) {
    return (
      <Router>
        <Routes>
          <Route
            path="/Login"
            element={<Login setIsLoggedIn={setIsLoggedIn} />}
          />
          <Route
            path="/signup"
            element={<SignUp setIsLoggedIn={setIsLoggedIn} />}
          />
          <Route path="*" element={<Navigate to="/Login" />} />
        </Routes>
      </Router>
    );
  }

  return (
    <Router>
      <Routes>
        <Route element={<Layout setIsLoggedIn={setIsLoggedIn} />}>
          <Route path="/" element={<Home />} />
          <Route path="/contact" element={<Contact />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/chat" element={<Chat />} />
        </Route>

        <Route path="/Login" element={<Navigate to="/" />} />
        <Route path="/signup" element={<Navigate to="/" />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
};

export default App;
