import React from "react";
import { HashRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
// import Chat from "./pages/Chat";
import SignUp from "./pages/SignUp";
import Home from "./pages/Home";
import Contacts from "./pages/Contacts.jsx";
import Profile from "./pages/Profile.jsx";
import Layout from "./components/Layout";
import { Navigate, Outlet } from "react-router-dom";
import { useState } from "react";
const App = () => {
  const [isLoggedIn, setIsLoggedIn] = useState(false);

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
          {/* <Route path="/chat" element={<Chat />} /> */}
          <Route path="/contact" element={<Contacts />} />
          <Route path="/profile" element={<Profile />} />
        </Route>

        <Route path="/Login" element={<Navigate to="/" />} />
        <Route path="/signup" element={<Navigate to="/" />} />
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </Router>
  );
};

export default App;
