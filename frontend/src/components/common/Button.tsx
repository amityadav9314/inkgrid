import React from 'react';
import styled from 'styled-components';

interface ButtonProps {
  primary?: boolean;
  secondary?: boolean;
  danger?: boolean;
  fullWidth?: boolean;
  disabled?: boolean;
  type?: 'button' | 'submit' | 'reset';
  onClick?: (event: React.MouseEvent<HTMLButtonElement>) => void;
  children: React.ReactNode;
}

const StyledButton = styled.button<ButtonProps>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
  font-weight: 600;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  border: none;
  outline: none;
  
  ${(props) => props.fullWidth && `
    width: 100%;
  `}
  
  ${(props) => props.primary && `
    background-color: #3b82f6;
    color: white;
    
    &:hover:not(:disabled) {
      background-color: #2563eb;
    }
  `}
  
  ${(props) => props.secondary && `
    background-color: #e5e7eb;
    color: #1f2937;
    
    &:hover:not(:disabled) {
      background-color: #d1d5db;
    }
  `}
  
  ${(props) => props.danger && `
    background-color: #ef4444;
    color: white;
    
    &:hover:not(:disabled) {
      background-color: #dc2626;
    }
  `}
  
  &:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
`;

const Button: React.FC<ButtonProps> = ({
  primary = true,
  secondary = false,
  danger = false,
  fullWidth = false,
  disabled = false,
  type = 'button',
  onClick,
  children,
}) => {
  return (
    <StyledButton
      primary={primary}
      secondary={secondary}
      danger={danger}
      fullWidth={fullWidth}
      disabled={disabled}
      type={type}
      onClick={onClick}
    >
      {children}
    </StyledButton>
  );
};

export default Button;
