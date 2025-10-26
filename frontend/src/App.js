// src/App.jsx
import React, { useState } from "react";

function App() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  const handleRegister = async () => {
    try {
      const res = await fetch("http://localhost:8080/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include", // cookie будут уходить
      });
      const data = await res.json();
      setMessage(data.status || data.error);
    } catch {
      setMessage("Registration failed");
    }
  };

  const handleLogin = async () => {
    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
        credentials: "include",
      });
      const data = await res.json();
      setMessage(data.status || data.error);

      if (res.status === 200 && data.status === "logged in") {
        setIsLoggedIn(true);
      } else {
        setIsLoggedIn(false);
      }
    } catch {
      setMessage("Login failed");
      setIsLoggedIn(false);
    }
  };

  const handleLogout = async () => {
    try {
      const res = await fetch("http://localhost:8080/logout", {
        method: "POST",
        credentials: "include",
      });
      const data = await res.json();
      setMessage(data.status || data.error);
      setIsLoggedIn(false);
    } catch {
      setMessage("Logout failed");
    }
  };

  const handleRefresh = async () => {
    try {
      const res = await fetch("http://localhost:8080/refresh", {
        method: "POST",
        credentials: "include",
      });
      const data = await res.json();
      setMessage(data.status || data.error);
    } catch {
      setMessage("Refresh failed");
    }
  };

  const handleCheckAuth = async () => {
    try {
      const res = await fetch("http://localhost:8080/healf", {
        method: "GET",
        credentials: "include",
      });
      const data = await res.json();
      if (res.ok && data.authenticated) {
        setMessage("User is authenticated ✅");
        setIsLoggedIn(true);
      } else {
        setMessage("User not authenticated ❌");
        setIsLoggedIn(false);
      }
    } catch {
      setMessage("Check auth failed");
      setIsLoggedIn(false);
    }
  };

  return (
    <div style={{ padding: "2rem", fontFamily: "sans-serif" }}>
      <h1>Auth Demo</h1>
      <p>Status: {isLoggedIn ? "Logged In ✅" : "Logged Out ❌"}</p>

      <input
        type="text"
        placeholder="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        style={{ marginRight: "0.5rem" }}
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        style={{ marginRight: "0.5rem" }}
      />

      <div style={{ marginTop: "1rem" }}>
        <button onClick={handleRegister} style={{ marginRight: "0.5rem" }}>
          Register
        </button>
        <button onClick={handleLogin} style={{ marginRight: "0.5rem" }}>
          Login
        </button>
        <button onClick={handleLogout} style={{ marginRight: "0.5rem" }}>
          Logout
        </button>
        <button onClick={handleRefresh} style={{ marginRight: "0.5rem" }}>
          Refresh Token
        </button>
        <button onClick={handleCheckAuth}>Check Auth</button>
      </div>

      <div style={{ marginTop: "1rem" }}>
        <strong>Response:</strong> {message}
      </div>
    </div>
  );
}

export default App;
