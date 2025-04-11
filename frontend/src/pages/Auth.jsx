import React, { useState, useEffect } from 'react';
import { message } from 'antd';
import { useNavigate } from 'react-router-dom'; 
import { EyeOutlined, EyeInvisibleOutlined } from '@ant-design/icons'; 
import './Auth.css';

const Auth = () => {
  const [isSignIn, setIsSignIn] = useState(true);
  const [userName, setUserName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false); 
  const [showConfirmPassword, setShowConfirmPassword] = useState(false); 
  const navigate = useNavigate();

  useEffect(() => {
    setUserName("user123");
    setPassword("user");
  }, []);

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword);
  };

  const toggleConfirmPasswordVisibility = () => {
    setShowConfirmPassword(!showConfirmPassword);
  };

  const handleSignInSubmit = async (event) => {
    event.preventDefault(); 
    setLoading(true); 
    try {
      const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userName, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Login failed');
      }

      const token = data['token'];
      const userId = data['id']; 

      sessionStorage.setItem('token', token);
      sessionStorage.setItem('userId', userId);
      sessionStorage.setItem('userName', userName);
      sessionStorage.setItem('password', password);

      navigate('/home'); 
    } catch (error) {
      message.error(`Error: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const handleRegisterSubmit = async (event) => {
    event.preventDefault(); 
    if (password !== confirmPassword) {
      message.error('Passwords do not match!');
      return;
    }

    setLoading(true); 
    try {
      const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ userName, email, password }),
      });

      if (!response.ok) {
        throw new Error('Registration failed');
      }
      message.success('Registration successful');
      setIsSignIn(true); 
    } catch (error) {
      message.error(`Error: ${error.message}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-container">
      <div className="auth-card">
        <h1 className="auth-title">{isSignIn ? 'Sign In' : 'Register'}</h1>

        <div className="auth-toggle">
          <button
            className={`auth-toggle-button ${isSignIn ? 'active' : ''}`}
            onClick={() => setIsSignIn(true)}
          >
            Sign In
          </button>
          <button
            className={`auth-toggle-button ${!isSignIn ? 'active' : ''}`}
            onClick={() => setIsSignIn(false)}
          >
            Register
          </button>
        </div>

        {isSignIn ? (
          <form onSubmit={handleSignInSubmit}>
            <div className="input-container">
              <input
                type="text"
                placeholder="User Name"
                className="input-field"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                required
              />
            </div>
            <div className="input-container">
              <div className="password-field">
                <input
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Password"
                  className="input-field"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <span className="eye-icon" onClick={togglePasswordVisibility}>
                  {showPassword ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                </span>
              </div>
            </div>
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? 'Loading...' : 'Login'}
            </button>
          </form>
        ) : (
          <form onSubmit={handleRegisterSubmit}>
            <div className="input-container">
              <input
                type="text"
                placeholder="UserName"
                className="input-field"
                value={userName}
                onChange={(e) => setUserName(e.target.value)}
                required
              />
            </div>
            <div className="input-container">
              <input
                type="email"
                placeholder="Email"
                className="input-field"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>
            <div className="input-container">
              <div className="password-field">
                <input
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Password"
                  className="input-field"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <span className="eye-icon" onClick={togglePasswordVisibility}>
                  {showPassword ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                </span>
              </div>
            </div>
            <div className="input-container">
              <div className="password-field">
                <input
                  type={showConfirmPassword ? 'text' : 'password'}
                  placeholder="Confirm Password"
                  className="input-field"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
                <span className="eye-icon" onClick={toggleConfirmPasswordVisibility}>
                  {showConfirmPassword ? <EyeInvisibleOutlined /> : <EyeOutlined />}
                </span>
              </div>
            </div>
            {password && confirmPassword && password !== confirmPassword && (
              <p className="password-error">Passwords do not match</p>
            )}
            <button type="submit" className="auth-button" disabled={loading}>
              {loading ? 'Loading...' : 'Register'}
            </button>
          </form>
        )}

        <div className="footer-links">
          <a href="#" className="footer-link">Forgot email?</a>
          <a href="#" className="footer-link">Need help?</a>
        </div>
        <div className="create-account">
          {isSignIn ? (
            <p>Don't have an account? <a href="#" onClick={() => setIsSignIn(false)}>Create one</a></p>
          ) : (
            <p>Already have an account? <a href="#" onClick={() => setIsSignIn(true)}>Sign in</a></p>
          )}
        </div>
      </div>
    </div>
  );
};

export default Auth;
