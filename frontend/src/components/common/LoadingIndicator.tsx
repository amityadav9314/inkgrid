import React from 'react';
import styled, { keyframes } from 'styled-components';

interface LoadingIndicatorProps {
  size?: 'small' | 'medium' | 'large';
  color?: string;
  fullScreen?: boolean;
  text?: string;
  className?: string;
}

const spin = keyframes`
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
`;

const FullScreenContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.8);
  z-index: 9999;
`;

const Container = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`;

const Spinner = styled.div<{ size: string; color: string }>`
  border-radius: 50%;
  animation: ${spin} 1s linear infinite;
  border-style: solid;
  border-color: ${props => `${props.color} transparent ${props.color} transparent`};
  
  ${props => {
    switch (props.size) {
      case 'small':
        return `
          width: 1.5rem;
          height: 1.5rem;
          border-width: 2px;
        `;
      case 'large':
        return `
          width: 3rem;
          height: 3rem;
          border-width: 4px;
        `;
      default: // medium
        return `
          width: 2rem;
          height: 2rem;
          border-width: 3px;
        `;
    }
  }}
`;

const LoadingText = styled.p<{ size: string }>`
  margin-top: 1rem;
  color: #4b5563;
  
  ${props => {
    switch (props.size) {
      case 'small':
        return 'font-size: 0.75rem;';
      case 'large':
        return 'font-size: 1.125rem;';
      default: // medium
        return 'font-size: 0.875rem;';
    }
  }}
`;

const LoadingIndicator: React.FC<LoadingIndicatorProps> = ({
  size = 'medium',
  color = '#3b82f6',
  fullScreen = false,
  text,
  className,
}) => {
  const content = (
    <Container className={className}>
      <Spinner size={size} color={color} />
      {text && <LoadingText size={size}>{text}</LoadingText>}
    </Container>
  );

  if (fullScreen) {
    return <FullScreenContainer>{content}</FullScreenContainer>;
  }

  return content;
};

export default LoadingIndicator;
