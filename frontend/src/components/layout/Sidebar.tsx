import React from 'react';
import { NavLink } from 'react-router-dom';
import styled from 'styled-components';

interface SidebarProps {
  isOpen: boolean;
  toggleSidebar: () => void;
}

const SidebarContainer = styled.aside<{ isOpen: boolean }>`
  position: fixed;
  top: 0;
  left: 0;
  width: 250px;
  height: 100vh;
  background-color: #1f2937;
  color: white;
  transform: translateX(${props => (props.isOpen ? '0' : '-100%')});
  transition: transform 0.3s ease-in-out;
  z-index: 1000;
  
  @media (min-width: 768px) {
    position: sticky;
    transform: translateX(0);
  }
`;

const SidebarHeader = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem;
  border-bottom: 1px solid #374151;
`;

const SidebarLogo = styled.div`
  display: flex;
  align-items: center;
  font-size: 1.25rem;
  font-weight: 700;
  
  img {
    width: 2rem;
    height: 2rem;
    margin-right: 0.75rem;
  }
`;

const CloseButton = styled.button`
  background: none;
  border: none;
  color: #9ca3af;
  cursor: pointer;
  
  &:hover {
    color: white;
  }
  
  svg {
    width: 1.5rem;
    height: 1.5rem;
  }
  
  @media (min-width: 768px) {
    display: none;
  }
`;

const SidebarNav = styled.nav`
  padding: 1.5rem 0;
`;

const SidebarNavItem = styled(NavLink)`
  display: flex;
  align-items: center;
  padding: 0.75rem 1.5rem;
  color: #9ca3af;
  text-decoration: none;
  transition: all 0.2s ease-in-out;
  
  svg {
    width: 1.25rem;
    height: 1.25rem;
    margin-right: 0.75rem;
  }
  
  &:hover {
    color: white;
    background-color: #374151;
  }
  
  &.active {
    color: white;
    background-color: #3b82f6;
    
    &:hover {
      background-color: #2563eb;
    }
  }
`;

const SidebarFooter = styled.div`
  position: absolute;
  bottom: 0;
  width: 100%;
  padding: 1.5rem;
  border-top: 1px solid #374151;
  font-size: 0.875rem;
  color: #9ca3af;
`;

const Sidebar: React.FC<SidebarProps> = ({ isOpen, toggleSidebar }) => {
  return (
    <SidebarContainer isOpen={isOpen}>
      <SidebarHeader>
        <SidebarLogo>
          <img src="/logo-white.svg" alt="InkGrid Logo" />
          InkGrid
        </SidebarLogo>
        <CloseButton onClick={toggleSidebar}>
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </CloseButton>
      </SidebarHeader>
      
      <SidebarNav>
        <SidebarNavItem to="/dashboard" end>
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Dashboard
        </SidebarNavItem>
        
        <SidebarNavItem to="/create">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
          </svg>
          Create Mosaic
        </SidebarNavItem>
        
        <SidebarNavItem to="/projects">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          My Projects
        </SidebarNavItem>
        
        <SidebarNavItem to="/gallery">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          Gallery
        </SidebarNavItem>
        
        <SidebarNavItem to="/settings">
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          Settings
        </SidebarNavItem>
      </SidebarNav>
      
      <SidebarFooter>
        &copy; {new Date().getFullYear()} InkGrid
      </SidebarFooter>
    </SidebarContainer>
  );
};

export default Sidebar;
