import React, { useState } from 'react';
import styled from 'styled-components';
import { useMosaic } from '../../context/MosaicContext';
import { ImageInfo } from '../../context/MosaicContext'; // Import the ImageInfo type
import ImageUpload from '../common/ImageUpload';
import LoadingIndicator from '../common/LoadingIndicator';
import { mosaicService } from '../../services/mosaicService';
import { config } from '../../config';


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

const TileImagesGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 1rem;
  margin-top: 1.5rem;
`;

const TileImageItem = styled.div`
  position: relative;
  border-radius: 0.375rem;
  overflow: hidden;
  aspect-ratio: 1;

  &:hover .overlay {
    opacity: 1;
  }
`;

const TileImage = styled.img`
  width: 100%;
  height: 100%;
  object-fit: cover;
`;

const ImageOverlay = styled.div`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.2s;
  cursor: pointer;
`;

const RemoveButton = styled.button`
  background-color: #ef4444;
  color: white;
  border: none;
  border-radius: 9999px;
  width: 2rem;
  height: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;

  &:hover {
    background-color: #dc2626;
  }
`;

const StatsContainer = styled.div`
  margin-top: 1.5rem;
  padding: 1rem;
  background-color: #f9fafb;
  border-radius: 0.375rem;
  display: flex;
  justify-content: space-between;
`;

const Stat = styled.div`
  text-align: center;

  .value {
    font-size: 1.25rem;
    font-weight: 600;
    color: #1f2937;
  }

  .label {
    font-size: 0.875rem;
    color: #6b7280;
  }
`;

interface ImageResponse {
  id: string;
  path: string;
  filename: string;
  width: number;
  height: number;
  format: string;
}

interface ImagesArrayResponse {
  images: ImageResponse[];
}

type UploadResponse = ImageResponse | ImageResponse[] | ImagesArrayResponse;

const TileImageSelector: React.FC = () => {
  const { tileImages, setTileImages } = useMosaic();
  const [isUploading, setIsUploading] = useState(false);
  const [uploadError, setUploadError] = useState<string | null>(null);

  // Update the handleUpload function in TileImageSelector.tsx
  const handleUpload = async (files: File[]) => {
    if (files.length === 0) return;

    setIsUploading(true);
    setUploadError(null);

    try {
      const response = (await mosaicService.uploadTileImages(
          files
      )) as UploadResponse;

      // Debug the response structure
      console.log('Upload response:', response);

      // Handle different response formats
      let newImages: ImageInfo[] = [];

      if (Array.isArray(response)) {
        newImages = response.map((img: ImageResponse) => ({
          id: img.id,
          path: config.getImageUrl(img.path), // Use config helper
          filename: img.filename,
          width: img.width,
          height: img.height,
          format: img.format,
        }));
      } else if (response && typeof response === 'object') {
        // Handle object response with images property
        if ('images' in response && Array.isArray(response.images)) {
          newImages = response.images.map((img: ImageResponse) => ({
            id: img.id,
            path: config.getImageUrl(img.path), // Use config helper
            filename: img.filename,
            width: img.width,
            height: img.height,
            format: img.format,
          }));
        } else {
          // Handle single image object
          const img = response as ImageResponse;
          newImages = [
            {
              id: img.id,
              path: config.getImageUrl(img.path), // Use config helper
              filename: img.filename,
              width: img.width,
              height: img.height,
              format: img.format,
            },
          ];
        }
      }

      // Update the tile images array
      setTileImages([...(tileImages || []), ...newImages]);
    } catch (error) {
      console.error('Error uploading tile images:', error);
      setUploadError('Failed to upload images. Please try again.');
    } finally {
      setIsUploading(false);
    }
  };

  const handleRemove = (id: string) => {
    if (!tileImages) return;
    setTileImages(tileImages.filter((img) => img.id !== id));
  };

  return (
    <Container>
      <Title>Tile Images</Title>
      <Description>
        Upload multiple images that will be used as tiles in your mosaic. For
        best results, use images with diverse colors and patterns.
      </Description>

      <ImageUpload onUpload={handleUpload} multiple={true} accept='image/*' />

      {isUploading && (
        <div style={{ marginTop: '1rem', textAlign: 'center' }}>
          <LoadingIndicator size='small' text='Uploading images...' />
        </div>
      )}

      {uploadError && (
        <div style={{ marginTop: '1rem', color: '#ef4444' }}>{uploadError}</div>
      )}

      {tileImages && tileImages.length > 0 && (
        <>
          <StatsContainer>
            <Stat>
              <div className='value'>{tileImages.length}</div>
              <div className='label'>Total Images</div>
            </Stat>
            <Stat>
              <div className='value'>
                {Math.floor(tileImages.length / 10) * 10}+
              </div>
              <div className='label'>Recommended</div>
            </Stat>
          </StatsContainer>

          <TileImagesGrid>
            {tileImages.map((image) => (
              <TileImageItem key={image.id}>
                <TileImage src={image.path} alt={image.filename} />
                <ImageOverlay className='overlay'>
                  <RemoveButton onClick={() => handleRemove(image.id)}>
                    <svg
                      xmlns='http://www.w3.org/2000/svg'
                      width='16'
                      height='16'
                      viewBox='0 0 24 24'
                      fill='none'
                      stroke='currentColor'
                      strokeWidth='2'
                      strokeLinecap='round'
                      strokeLinejoin='round'
                    >
                      <line x1='18' y1='6' x2='6' y2='18'></line>
                      <line x1='6' y1='6' x2='18' y2='18'></line>
                    </svg>
                  </RemoveButton>
                </ImageOverlay>
              </TileImageItem>
            ))}
          </TileImagesGrid>
        </>
      )}
    </Container>
  );
};

export default TileImageSelector;
