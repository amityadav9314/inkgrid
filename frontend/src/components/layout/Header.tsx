import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { useAuth } from '../../context/AuthContext';

const HeaderContainer = styled.header`
  background-color: #ffffff;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

const Logo = styled(Link)`
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #1f2937;
  font-size: 1.5rem;
  font-weight: 700;
  
  img {
    height: 2.5rem;
    margin-right: 0.75rem;
  }
`;

const Nav = styled.nav`
  display: flex;
  align-items: center;
`;

const NavLink = styled(Link)`
  color: #4b5563;
  text-decoration: none;
  padding: 0.5rem 1rem;
  font-weight: 500;
  transition: color 0.2s ease-in-out;
  
  &:hover {
    color: #3b82f6;
  }
  
  &.active {
    color: #3b82f6;
  }
`;

const AuthButtons = styled.div`
  display: flex;
  gap: 1rem;
`;

const LoginButton = styled(Link)`
  color: #3b82f6;
  text-decoration: none;
  padding: 0.5rem 1rem;
  font-weight: 500;
  border-radius: 0.375rem;
  transition: all 0.2s ease-in-out;
  
  &:hover {
    background-color: #eff6ff;
  }
`;

const SignupButton = styled(Link)`
  background-color: #3b82f6;
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  font-weight: 500;
  border-radius: 0.375rem;
  transition: all 0.2s ease-in-out;
  
  &:hover {
    background-color: #2563eb;
  }
`;

const UserMenu = styled.div`
  position: relative;
`;

const UserButton = styled.button`
  display: flex;
  align-items: center;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.375rem;
  
  &:hover {
    background-color: #f3f4f6;
  }
  
  img {
    width: 2rem;
    height: 2rem;
    border-radius: 50%;
    margin-right: 0.5rem;
  }
  
  span {
    margin-right: 0.5rem;
  }
  
  svg {
    width: 1rem;
    height: 1rem;
  }
`;

const DropdownMenu = styled.div<{ isOpen: boolean }>`
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.5rem;
  width: 12rem;
  background-color: white;
  border-radius: 0.375rem;
  box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  z-index: 10;
  overflow: hidden;
  display: ${props => (props.isOpen ? 'block' : 'none')};
`;

const DropdownItem = styled.button`
  width: 100%;
  text-align: left;
  padding: 0.75rem 1rem;
  background: none;
  border: none;
  cursor: pointer;
  color: #4b5563;
  font-size: 0.875rem;
  
  &:hover {
    background-color: #f3f4f6;
    color: #1f2937;
  }
  
  &.logout {
    color: #ef4444;
    
    &:hover {
      background-color: #fef2f2;
    }
  }
`;

const MobileMenuButton = styled.button`
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.5rem;
  
  @media (max-width: 768px) {
    display: block;
  }
  
  svg {
    width: 1.5rem;
    height: 1.5rem;
    color: #4b5563;
  }
`;

const MobileMenu = styled.div<{ isOpen: boolean }>`
  display: none;
  
  @media (max-width: 768px) {
    display: ${props => (props.isOpen ? 'block' : 'none')};
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    background-color: white;
    padding: 1rem;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
    z-index: 20;
  }
`;

const MobileNavLink = styled(Link)`
  display: block;
  padding: 0.75rem 0;
  color: #4b5563;
  text-decoration: none;
  font-weight: 500;
  border-bottom: 1px solid #e5e7eb;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    color: #3b82f6;
  }
`;

const Header: React.FC = () => {
  const { user, isAuthenticated, logout } = useAuth();
  const navigate = useNavigate();
  const [userMenuOpen, setUserMenuOpen] = useState(false);
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  
  const handleLogout = () => {
    logout();
    setUserMenuOpen(false);
    navigate('/');
  };
  
  return (
    <HeaderContainer>
      <Logo to="/">
        <img src="/logo.svg" alt="InkGrid Logo" />
        InkGrid
      </Logo>
      
      <MobileMenuButton onClick={() => setMobileMenuOpen(!mobileMenuOpen)}>
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
        </svg>
      </MobileMenuButton>
      
      <Nav>
        <NavLink to="/">Home</NavLink>
        {isAuthenticated && (
          <>
            <NavLink to="/projects/">My Projects</NavLink>
          </>
        )}
      </Nav>
      
      {isAuthenticated ? (
        <UserMenu>
          <UserButton onClick={() => setUserMenuOpen(!userMenuOpen)}>
            <img src={`https://ui-avatars.com/api/?name=${user?.name || 'User'}&background=random`} alt="User" />
            <span>{user?.name}</span>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
            </svg>
          </UserButton>
          
          <DropdownMenu isOpen={userMenuOpen}>
            <DropdownItem onClick={() => {
              setUserMenuOpen(false);
              navigate('/settings');
            }}>
              Settings
            </DropdownItem>
            <DropdownItem className="logout" onClick={handleLogout}>
              Logout
            </DropdownItem>
          </DropdownMenu>
        </UserMenu>
      ) : (
        <AuthButtons>
          <LoginButton to="/auth?mode=login">Log in</LoginButton>
          <SignupButton to="/auth?mode=register">Sign up</SignupButton>
        </AuthButtons>
      )}
      
      <MobileMenu isOpen={mobileMenuOpen}>
        <MobileNavLink to="/" onClick={() => setMobileMenuOpen(false)}>Home</MobileNavLink>
        {isAuthenticated ? (
          <>
            <MobileNavLink to="/projects/" onClick={() => setMobileMenuOpen(false)}>My Projects</MobileNavLink>
            <MobileNavLink to="/settings" onClick={() => setMobileMenuOpen(false)}>Settings</MobileNavLink>
            <MobileNavLink to="#" onClick={() => {
              setMobileMenuOpen(false);
              handleLogout();
            }}>Logout</MobileNavLink>
          </>
        ) : (
          <>
            <MobileNavLink to="/auth?mode=login" onClick={() => setMobileMenuOpen(false)}>Log in</MobileNavLink>
            <MobileNavLink to="/auth?mode=register" onClick={() => setMobileMenuOpen(false)}>Sign up</MobileNavLink>
          </>
        )}
      </MobileMenu>
    </HeaderContainer>
  );
};

export default Header;
