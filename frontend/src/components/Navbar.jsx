import React from "react";
import { useNavigate } from "react-router-dom";
import { LogoutOutlined, SearchOutlined } from '@ant-design/icons';
import "./Navbar.css";

const Navbar = () => {
  const navigate = useNavigate();

  const handleLogout = () => {
    sessionStorage.removeItem("token");
    navigate("/");
  };

  return (
    <nav className="navbar">
      <div className="navbar-left">
        <div className="navbar-logo">
          <span className="navbar-title">EquiTrack</span>
        </div>
      </div>

      <div className="navbar-center">
        <div className="navbar-search">
          <input
            type="text"
            placeholder="Search stocks"
            className="search-input"
          />
          <button className="search-button">
            <span className="search-icon"><SearchOutlined/></span>
          </button>
        </div>
      </div>

      <div className="navbar-right">
        <div className="user-menu">
          <span className="logout-icon" onClick={handleLogout}>
            <LogoutOutlined /> Logout
          </span>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
