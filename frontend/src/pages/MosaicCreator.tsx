import React, { useState } from 'react';
import styled from 'styled-components';
import { useMosaic } from '../context/MosaicContext';
import MainImageSelector from '../components/mosaic/MainImageSelector';
import TileImageSelector from '../components/mosaic/TileImageSelector';
import Button from '../components/common/Button';

const Container = styled.div`
  max-width: 1200px;
  margin: 0 auto;
`;

const Title = styled.h1`
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 1.5rem;
  color: #1f2937;
`;

const StepContainer = styled.div`
  margin-bottom: 2rem;
`;

const StepIndicator = styled.div`
  display: flex;
  margin-bottom: 2rem;
`;

const Step = styled.div<{ active: boolean; completed: boolean }>`
  flex: 1;
  text-align: center;
  padding: 1rem;
  position: relative;
  
  &:not(:last-child)::after {
    content: '';
    position: absolute;
    top: 50%;
    right: 0;
    transform: translate(50%, -50%);
    width: 100%;
    height: 2px;
    background-color: ${props =>
    props.completed ? '#3b82f6' : '#e5e7eb'};
    z-index: 0;
  }
  
  .step-number {
    width: 2.5rem;
    height: 2.5rem;
    border-radius: 50%;
    background-color: ${props =>
    props.completed ? '#3b82f6' :
        props.active ? '#eff6ff' : '#f3f4f6'};
    border: 2px solid ${props =>
    props.completed ? '#3b82f6' :
        props.active ? '#3b82f6' : '#e5e7eb'};
    color: ${props =>
    props.completed ? 'white' :
        props.active ? '#3b82f6' : '#9ca3af'};
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto 0.5rem;
    font-weight: 600;
    position: relative;
    z-index: 1;
  }
  
  .step-label {
    font-size: 0.875rem;
    font-weight: 500;
    color: ${props =>
    props.active || props.completed ? '#1f2937' : '#6b7280'};
  }
`;

const ActionButtons = styled.div`
  display: flex;
  justify-content: space-between;
  margin-top: 2rem;
`;

const MosaicCreator: React.FC = () => {
  const { mainImage, tileImages } = useMosaic();
  const [currentStep, setCurrentStep] = useState(1);

  const steps = [
    { number: 1, label: 'Upload Main Image' },
    { number: 2, label: 'Select Tile Images' },
    { number: 3, label: 'Configure Settings' },
    { number: 4, label: 'Generate & Export' },
  ];

  const handleNext = () => {
    if (currentStep < steps.length) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handleBack = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
    }
  };

  const renderStepContent = () => {
    switch (currentStep) {
      case 1:
        return <MainImageSelector />;
      case 2:
        return <TileImageSelector />;
      case 3:
        return <div>Mosaic Settings (Coming Soon)</div>;
      case 4:
        return <div>Generate & Export (Coming Soon)</div>;
      default:
        return null;
    }
  };

  const isNextDisabled = () => {
    if (currentStep === 1) return !mainImage;
    if (currentStep === 2) return !tileImages || tileImages.length < 5; // Require at least 5 tile images
    return false;
  };

  return (
      <Container>
        <Title>Create Your Mosaic</Title>

        <StepIndicator>
          {steps.map(step => (
              <Step
                  key={step.number}
                  active={step.number === currentStep}
                  completed={step.number < currentStep}
              >
                <div className="step-number">
                  {step.number < currentStep ? (
                      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" width="16" height="16">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                  ) : (
                      step.number
                  )}
                </div>
                <div className="step-label">{step.label}</div>
              </Step>
          ))}
        </StepIndicator>

        <StepContainer>
          {renderStepContent()}
        </StepContainer>

        <ActionButtons>
          <Button
              secondary
              onClick={handleBack}
              disabled={currentStep === 1}
          >
            Back
          </Button>

          <Button
              primary
              onClick={handleNext}
              disabled={isNextDisabled()}
          >
            {currentStep === steps.length ? 'Finish' : 'Next'}
          </Button>
        </ActionButtons>
      </Container>
  );
};

export default MosaicCreator;