import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import styled from 'styled-components';

interface ImageUploadProps {
  onUpload: (files: File[]) => void;
  multiple?: boolean;
  accept?: string;
  maxSize?: number;
  className?: string;
  children?: React.ReactNode;
}

const DropzoneContainer = styled.div`
  border: 2px dashed #d1d5db;
  border-radius: 0.5rem;
  padding: 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  background-color: #f9fafb;
  
  &:hover {
    border-color: #3b82f6;
    background-color: #f3f4f6;
  }
  
  &.active {
    border-color: #3b82f6;
    background-color: #eff6ff;
  }
`;

const UploadIcon = styled.div`
  margin-bottom: 1rem;
  color: #6b7280;
  
  svg {
    width: 2.5rem;
    height: 2.5rem;
  }
`;

const UploadText = styled.div`
  color: #374151;
  font-size: 1rem;
  
  p {
    margin: 0.5rem 0;
  }
  
  .primary {
    font-weight: 600;
    color: #3b82f6;
  }
  
  .secondary {
    font-size: 0.875rem;
    color: #6b7280;
  }
`;

const ErrorMessage = styled.div`
  color: #ef4444;
  margin-top: 0.5rem;
  font-size: 0.875rem;
`;

const ImageUpload: React.FC<ImageUploadProps> = ({
  onUpload,
  multiple = false,
  accept = 'image/jpeg, image/png, image/webp',
  maxSize = 10 * 1024 * 1024, // 10MB
  className,
  children,
}) => {
  const onDrop = useCallback((acceptedFiles: File[]) => {
    onUpload(acceptedFiles);
  }, [onUpload]);

  const {
    getRootProps,
    getInputProps,
    isDragActive,
    isDragReject,
    fileRejections,
  } = useDropzone({
    onDrop,
    accept: { 'image/*': ['.jpeg', '.jpg', '.png', '.webp', '.heic'] },
    multiple,
    maxSize,
  });

  const fileRejectionItems = fileRejections.map(({ file, errors }) => (
    <li key={file.name}>
      {file.name} - {errors.map(e => e.message).join(', ')}
    </li>
  ));

  return (
    <div className={className}>
      <DropzoneContainer
        {...getRootProps()}
        className={`${isDragActive ? 'active' : ''}`}
      >
        <input {...getInputProps()} />
        {children || (
          <>
            <UploadIcon>
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </UploadIcon>
            <UploadText>
              {isDragActive ? (
                <p className="primary">Drop the files here...</p>
              ) : (
                <>
                  <p className="primary">Drag & drop {multiple ? 'images' : 'an image'} here, or click to select {multiple ? 'files' : 'a file'}</p>
                  <p className="secondary">Supports JPG, PNG, WEBP, and HEIC formats</p>
                </>
              )}
            </UploadText>
          </>
        )}
      </DropzoneContainer>
      
      {fileRejectionItems.length > 0 && (
        <ErrorMessage>
          <ul>
            {fileRejectionItems}
          </ul>
        </ErrorMessage>
      )}
    </div>
  );
};

export default ImageUpload;
