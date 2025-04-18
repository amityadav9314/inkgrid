import React, { useState } from 'react';
import styled from 'styled-components';
import { useMosaic } from '../../context/MosaicContext';
import ImageUpload from '../common/ImageUpload';
import Button from '../common/Button';
import LoadingIndicator from '../common/LoadingIndicator';
import { mosaicService } from '../../services/mosaicService';

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

const ImagePreviewContainer = styled.div`
  position: relative;
  margin-top: 1.5rem;
`;

const ImagePreview = styled.img`
  width: 100%;
  max-height: 400px;
  object-fit: contain;
  border-radius: 0.375rem;
  border: 1px solid #e5e7eb;
`;

const ImageActions = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 1rem;
  gap: 0.5rem;
`;

const MainImageSelector: React.FC = () => {
  const { mainImage, setMainImage } = useMosaic();
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  const handleUpload = async (files: File[]) => {
    if (files.length === 0) return;
    
    const file = files[0]; // Only use the first file
    setIsUploading(true);
    setUploadError(null);
    
    try {
      const response = await mosaicService.uploadMainImage(file);
      setMainImage({
        id: response.id,
        path: response.path,
        filename: response.filename,
        width: response.width,
        height: response.height,
        format: response.format
      });
    } catch (error) {
      console.error('Error uploading main image:', error);
      setUploadError('Failed to upload image. Please try again.');
    } finally {
      setIsUploading(false);
    }
  };

  const handleRemove = () => {
    setMainImage(null);
  };

  return (
    <Container>
      <Title>Main Image</Title>
      <Description>
        Upload the main image that will be transformed into a mosaic. 
        This will be the foundation of your mosaic artwork.
      </Description>
      
      {!mainImage ? (
        <>
          <ImageUpload 
            onUpload={handleUpload} 
            multiple={false}
          />
          {isUploading && (
            <div style={{ marginTop: '1rem', textAlign: 'center' }}>
              <LoadingIndicator size="small" text="Uploading image..." />
            </div>
          )}
          {uploadError && (
            <div style={{ marginTop: '1rem', color: '#ef4444' }}>
              {uploadError}
            </div>
          )}
        </>
      ) : (
        <ImagePreviewContainer>
          <ImagePreview 
            src={mainImage.path} 
            alt={mainImage.filename} 
          />
          <ImageActions>
            <Button 
              secondary 
              onClick={handleRemove}
            >
              Remove
            </Button>
          </ImageActions>
        </ImagePreviewContainer>
      )}
    </Container>
  );
};

export default MainImageSelector;
