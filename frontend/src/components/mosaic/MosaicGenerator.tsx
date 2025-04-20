import React, { useState, useEffect, useRef } from 'react';
import styled from 'styled-components';
import { useMosaic } from '../../context/MosaicContext';
import Button from '../common/Button';
import LoadingIndicator from '../common/LoadingIndicator';
import { mosaicService } from '../../services/mosaicService';
import { TransformWrapper, TransformComponent } from 'react-zoom-pan-pinch';

interface MosaicGeneratorProps {
  projectId?: number;
}

// Update the settings interface to include style
interface MosaicSettings {
  tileSize: number;
  tileDensity: number;
  colorAdjustment: number;
  style?: "classic" | "random" | "flowing";
}

// Update the MosaicGenerationStatus interface to include error property
interface MosaicGenerationStatus {
  id: string;
  status: 'pending' | 'processing' | 'completed' | 'failed';
  progress: number;
  sd_url?: string;
  hd_url?: string;
  error?: string;
  created_at: string;
  updated_at: string;
  tile_size?: number;
  tile_density?: number;
  style?: string;
}

// API response type for mosaic status
interface ApiMosaicResponse {
  id: string;
  status: string;
  progress: number;
  sd_url?: string;
  hd_url?: string;
  error?: string;
  created_at: string;
  updated_at: string;
  tile_size?: number;
  tile_density?: number;
  style?: string;
  [key: string]: any;
}

const Container = styled.div`
  background-color: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  padding: 1.5rem;
  margin-bottom: 2rem;
`;

const Title = styled.h2`
  font-size: 1.25rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 1rem;
`;

const Description = styled.p`
  color: #6b7280;
  margin-bottom: 1.5rem;
`;

const GenerateButton = styled(Button)`
  margin-bottom: 1.5rem;
`;

const MosaicPreviewContainer = styled.div`
  margin-top: 1.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  overflow: hidden;
  position: relative;
`;

const ZoomControls = styled.div`
  position: absolute;
  top: 1rem;
  right: 1rem;
  z-index: 10;
  display: flex;
  gap: 0.5rem;
  background-color: rgba(255, 255, 255, 0.8);
  padding: 0.5rem;
  border-radius: 0.375rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
`;

const ZoomButton = styled.button`
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 1.25rem;
  font-weight: bold;

  &:hover {
    background-color: #f3f4f6;
  }

  &:active {
    background-color: #e5e7eb;
  }
`;

const QualityToggle = styled.div`
  display: flex;
  margin-bottom: 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  overflow: hidden;
`;

const QualityOption = styled.button<{ active: boolean }>`
  flex: 1;
  padding: 0.75rem;
  background-color: ${props => props.active ? '#3b82f6' : 'white'};
  color: ${props => props.active ? 'white' : '#4b5563'};
  border: none;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;

  &:hover {
    background-color: ${props => props.active ? '#3b82f6' : '#f3f4f6'};
  }
`;

const MosaicList = styled.div`
  margin-top: 2rem;
`;

const MosaicListTitle = styled.h3`
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 1rem;
`;

const MosaicItem = styled.div`
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  padding: 1rem;
  margin-bottom: 1rem;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: #3b82f6;
    box-shadow: 0 1px 3px 0 rgba(59, 130, 246, 0.2);
  }
`;

const MosaicItemHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
`;

const MosaicItemTitle = styled.h4`
  font-size: 1rem;
  font-weight: 500;
  color: #1f2937;
`;

const MosaicItemDate = styled.span`
  font-size: 0.875rem;
  color: #6b7280;
`;

const MosaicItemPreview = styled.img`
  width: 100%;
  height: 150px;
  object-fit: cover;
  border-radius: 0.25rem;
  margin-top: 0.5rem;
`;

const ErrorMessage = styled.div`
  color: #ef4444;
  padding: 1rem;
  background-color: #fee2e2;
  border-radius: 0.375rem;
  margin-bottom: 1rem;
`;

const MosaicGenerator: React.FC<MosaicGeneratorProps> = ({ projectId }) => {
  const { mainImage, tileImages, settings } = useMosaic();
  const [isGenerating, setIsGenerating] = useState(false);
  const [generationId, setGenerationId] = useState<string | null>(null);
  const [generationStatus, setGenerationStatus] = useState<MosaicGenerationStatus | null>(null);
  const [generationError, setGenerationError] = useState<string | null>(null);
  const [quality, setQuality] = useState<'sd' | 'hd'>('sd');
  const [previousMosaics, setPreviousMosaics] = useState<MosaicGenerationStatus[]>([]);
  const [isLoadingPrevious, setIsLoadingPrevious] = useState(false);
  const statusCheckInterval = useRef<NodeJS.Timeout | null>(null);

  // Fetch previous mosaics when component mounts
  useEffect(() => {
    if (projectId) {
      fetchPreviousMosaics();
    }

    return () => {
      // Clear interval when component unmounts
      if (statusCheckInterval.current) {
        clearInterval(statusCheckInterval.current);
      }
    };
  }, [projectId]);

  // Check generation status periodically
  useEffect(() => {
    if (generationId && isGenerating) {
      // Start checking status immediately
      checkGenerationStatus();

      // Set up interval to check status every 2 seconds
      statusCheckInterval.current = setInterval(checkGenerationStatus, 2000);

      // Clean up interval when generation is complete or component unmounts
      return () => {
        if (statusCheckInterval.current) {
          clearInterval(statusCheckInterval.current);
          statusCheckInterval.current = null;
        }
      };
    }
  }, [generationId, isGenerating]);

  const fetchPreviousMosaics = async () => {
    if (!projectId) return;

    try {
      setIsLoadingPrevious(true);
      const response = await mosaicService.getProjectMosaics(projectId);

      // Type cast each mosaic in the response
      const typedMosaics: MosaicGenerationStatus[] = (response.mosaics || []).map((mosaic: ApiMosaicResponse) => ({
        ...mosaic,
        status: mosaic.status as 'pending' | 'processing' | 'completed' | 'failed',
      }));

      setPreviousMosaics(typedMosaics);
    } catch (error) {
      console.error('Error fetching previous mosaics:', error);
    } finally {
      setIsLoadingPrevious(false);
    }
  };

  const checkGenerationStatus = async () => {
    if (!generationId) return;

    try {
      const response = await mosaicService.getGenerationStatus(generationId);

      // Type cast the response to match your interface
      const status: MosaicGenerationStatus = {
        ...response,
        // Ensure the status is one of the allowed string literals
        status: response.status as 'pending' | 'processing' | 'completed' | 'failed',
      };

      setGenerationStatus(status);

      // If generation is complete or failed, stop checking
      if (status.status === 'completed' || status.status === 'failed') {
        if (statusCheckInterval.current) {
          clearInterval(statusCheckInterval.current);
          statusCheckInterval.current = null;
        }

        setIsGenerating(false);

        // If failed, show error
        if (status.status === 'failed' && status.error) {
          setGenerationError(status.error);
        }

        // If completed, refresh previous mosaics list
        if (status.status === 'completed') {
          fetchPreviousMosaics();
        }
      }
    } catch (error) {
      console.error('Error checking generation status:', error);
      setGenerationError('Failed to check generation status. Please try again.');
      setIsGenerating(false);

      if (statusCheckInterval.current) {
        clearInterval(statusCheckInterval.current);
        statusCheckInterval.current = null;
      }
    }
  };

  const handleGenerateMosaic = async () => {
    if (!projectId || !mainImage || !tileImages || tileImages.length === 0) {
      setGenerationError('Missing required data for mosaic generation. Please make sure you have selected a main image and at least one tile image.');
      return;
    }

    try {
      setIsGenerating(true);
      setGenerationError(null);
      setGenerationStatus(null);

      // Prepare tile image IDs
      const tileImageIDs = tileImages.map(img => img.id);

      // Generate mosaic - fixed property names to match API naming
      const response = await mosaicService.generateMosaic({
        project_id: projectId,
        main_image_id: mainImage.id,
        tile_image_ids: tileImageIDs,
        tile_size: settings?.tileSize || 50,
        tile_density: settings?.tileDensity || 80,
        overlay_ratio: (settings?.colorAdjustment || 50) / 100,
        style: (settings?.style as "classic" | "random" | "flowing") || 'classic',
        color_correction: true
      });

      setGenerationId(response.id);
    } catch (error) {
      console.error('Error generating mosaic:', error);
      setGenerationError('Failed to start mosaic generation. Please try again.');
      setIsGenerating(false);
    }
  };

  const handleSelectMosaic = (mosaic: ApiMosaicResponse) => {
    // Type cast the mosaic to match your interface
    const typedMosaic: MosaicGenerationStatus = {
      ...mosaic,
      status: mosaic.status as 'pending' | 'processing' | 'completed' | 'failed',
    };

    setGenerationStatus(typedMosaic);
    setGenerationId(mosaic.id);
    setGenerationError(null);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
  };

  const renderMosaicPreview = () => {
    if (!generationStatus || generationStatus.status !== 'completed') {
      return null;
    }

    const imageUrl = quality === 'sd' ? generationStatus.sd_url : generationStatus.hd_url;

    if (!imageUrl) {
      return <div>No image available for the selected quality.</div>;
    }

    return (
        <TransformWrapper
            initialScale={1}
            minScale={0.5}
            maxScale={4}
            wheel={{ step: 0.1 }}
        >
          {({ zoomIn, zoomOut, resetTransform }: {
            zoomIn: () => void;
            zoomOut: () => void;
            resetTransform: () => void;
          }) => (
              <>
                <ZoomControls>
                  <ZoomButton onClick={() => zoomIn()}>+</ZoomButton>
                  <ZoomButton onClick={() => zoomOut()}>-</ZoomButton>
                  <ZoomButton onClick={() => resetTransform()}>↺</ZoomButton>
                </ZoomControls>
                <TransformComponent>
                  <img
                      src={imageUrl}
                      alt="Generated Mosaic"
                      style={{ width: '100%', display: 'block' }}
                  />
                </TransformComponent>
              </>
          )}
        </TransformWrapper>
    );
  };

  return (
      <Container>
        <Title>Generate & Export Mosaic</Title>
        <Description>
          Generate your mosaic using the main image, tile images, and settings you've configured.
          You can choose between standard definition (SD) and high definition (HD) quality.
        </Description>

        {generationError && (
            <ErrorMessage>{generationError}</ErrorMessage>
        )}

        {!generationStatus || generationStatus.status !== 'completed' ? (
            <GenerateButton
                onClick={handleGenerateMosaic}
                disabled={isGenerating}
            >
              {isGenerating ? 'Generating...' : 'Generate Mosaic'}
            </GenerateButton>
        ) : (
            <>
              <QualityToggle>
                <QualityOption
                    active={quality === 'sd'}
                    onClick={() => setQuality('sd')}
                >
                  Standard Definition
                </QualityOption>
                <QualityOption
                    active={quality === 'hd'}
                    onClick={() => setQuality('hd')}
                >
                  High Definition
                </QualityOption>
              </QualityToggle>

              <Button
                  onClick={() => {
                    const imageUrl = quality === 'sd' ? generationStatus.sd_url : generationStatus.hd_url;
                    if (imageUrl) {
                      window.open(imageUrl, '_blank');
                    }
                  }}
              >
                Download {quality.toUpperCase()} Image
              </Button>
            </>
        )}

        {isGenerating && (
            <div style={{ marginTop: '1.5rem' }}>
              <LoadingIndicator
                  text={`Generating mosaic... ${generationStatus ? generationStatus.progress + '%' : ''}`}
              />
            </div>
        )}

        {generationStatus && generationStatus.status === 'completed' && (
            <MosaicPreviewContainer>
              {renderMosaicPreview()}
            </MosaicPreviewContainer>
        )}

        {projectId && (
            <MosaicList>
              <MosaicListTitle>Previous Mosaics</MosaicListTitle>

              {isLoadingPrevious ? (
                  <LoadingIndicator text="Loading previous mosaics..." />
              ) : previousMosaics.length > 0 ? (
                  previousMosaics.map((mosaic) => (
                      <MosaicItem
                          key={mosaic.id}
                          onClick={() => handleSelectMosaic(mosaic)}
                      >
                        <MosaicItemHeader>
                          <MosaicItemTitle>
                            Mosaic {mosaic.id}
                          </MosaicItemTitle>
                          <MosaicItemDate>
                            {formatDate(mosaic.created_at)}
                          </MosaicItemDate>
                        </MosaicItemHeader>
                        <div>
                          Tile Size: {mosaic.tile_size},
                          Density: {mosaic.tile_density},
                          Style: {mosaic.style}
                        </div>
                        {mosaic.sd_url && (
                            <MosaicItemPreview src={mosaic.sd_url} alt={`Mosaic ${mosaic.id}`} />
                        )}
                      </MosaicItem>
                  ))
              ) : (
                  <p>No previous mosaics found.</p>
              )}
            </MosaicList>
        )}
      </Container>
  );
};

export default MosaicGenerator;