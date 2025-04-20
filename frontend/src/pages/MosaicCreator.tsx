import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { useMosaic } from '../context/MosaicContext';
import MainImageSelector from '../components/mosaic/MainImageSelector';
import TileImageSelector from '../components/mosaic/TileImageSelector';
import MosaicSettings from '../components/mosaic/MosaicSettings';
import Button from '../components/common/Button';
import { projectService } from '../services/projectService';
import { useAuth } from '../context/AuthContext';
import LoadingIndicator from '../components/common/LoadingIndicator';
import { config } from '../config';

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

const ErrorMessage = styled.div`
  color: #ef4444;
  padding: 1rem;
  background-color: #fee2e2;
  border-radius: 0.5rem;
  margin-bottom: 1.5rem;
`;

const MosaicCreator: React.FC = () => {
  const { mainImage, tileImages, setMainImage, setTileImages } = useMosaic();
  const [currentStep, setCurrentStep] = useState(1);
  const [project, setProject] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLoadingImages, setIsLoadingImages] = useState(false);
  
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    // Redirect to login if not authenticated
    if (!isAuthenticated) {
      navigate('/auth');
      return;
    }

    // Load project data if editing an existing project
    if (id) {
      const fetchProject = async () => {
        try {
          setIsLoading(true);
          const projectData = await projectService.getProject(parseInt(id));
          setProject(projectData);
          setIsLoading(false);
          
          // Fetch project images
          await fetchProjectImages(parseInt(id));
        } catch (error) {
          console.error('Error fetching project:', error);
          setError('Failed to load project data. Please try again.');
          setIsLoading(false);
        }
      };

      fetchProject();
    } else {
      setIsLoading(false);
    }
  }, [id, navigate, isAuthenticated]);
  
  const fetchProjectImages = async (projectId: number) => {
    try {
      setIsLoadingImages(true);
      const response = await projectService.getProjectImages(projectId);
      
      // Set main image if available
      if (response.main_images && response.main_images.length > 0) {
        const mainImg = response.main_images[0];
        setMainImage({
          id: mainImg.id,
          path: mainImg.path,
          filename: mainImg.filename,
          width: mainImg.width,
          height: mainImg.height,
          format: mainImg.format
        });
      }
      
      // Set tile images if available
      if (response.tile_images && response.tile_images.length > 0) {
        const tileImgs = response.tile_images.map((img: any) => ({
          id: img.id,
          path: img.path,
          filename: img.filename,
          width: img.width,
          height: img.height,
          format: img.format
        }));
        setTileImages(tileImgs);
      }
      
      setIsLoadingImages(false);
    } catch (error) {
      console.error('Error fetching project images:', error);
      setIsLoadingImages(false);
    }
  };

  const steps = [
    { number: 1, label: 'Upload Main Image' },
    { number: 2, label: 'Select Tile Images' },
    { number: 3, label: 'Configure Settings' },
    { number: 4, label: 'Generate & Export' },
  ];

  const handleNext = () => {
    // Validate current step
    if (currentStep === 1 && !mainImage) {
      alert('Please upload a main image before proceeding.');
      return;
    }

    if (currentStep === 2 && (!tileImages || tileImages.length === 0)) {
      alert('Please upload at least one tile image before proceeding.');
      return;
    }

    if (currentStep < steps.length) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
    }
  };

  const renderStep = () => {
    if (isLoading) {
      return <LoadingIndicator text="Loading project data..." />;
    }

    if (error) {
      return <ErrorMessage>{error}</ErrorMessage>;
    }

    switch (currentStep) {
      case 1:
        return <MainImageSelector projectId={id ? parseInt(id) : undefined} />;
      case 2:
        return <TileImageSelector projectId={id ? parseInt(id) : undefined} />;
      case 3:
        return <MosaicSettings projectId={id ? parseInt(id) : undefined} />;
      case 4:
        return (
            <div>
              <h2>Generate & Export</h2>
              <p>This feature is coming soon!</p>
            </div>
        );
      default:
        return null;
    }
  };

  return (
      <Container>
        <Title>{project ? `Edit Project: ${project.name}` : 'Create New Mosaic'}</Title>

        <StepIndicator>
          {steps.map((step) => (
              <Step
                  key={step.number}
                  active={currentStep === step.number}
                  completed={currentStep > step.number}
              >
                <div className="step-number">{step.number}</div>
                <div className="step-label">{step.label}</div>
              </Step>
          ))}
        </StepIndicator>

        <StepContainer>
          {isLoadingImages ? (
            <LoadingIndicator text="Loading project images..." />
          ) : (
            renderStep()
          )}
        </StepContainer>

        <ActionButtons>
          <Button
              secondary
              onClick={handlePrevious}
              disabled={currentStep === 1}
          >
            Previous
          </Button>
          <Button
              onClick={handleNext}
              disabled={currentStep === steps.length}
          >
            {currentStep === steps.length ? 'Finish' : 'Next'}
          </Button>
        </ActionButtons>
      </Container>
  );
};

export default MosaicCreator;