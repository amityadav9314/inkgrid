import React, { createContext, useContext, useState, ReactNode } from 'react';

// Define the image info type
export interface ImageInfo {
  id: string;
  path: string;
  filename: string;
  width: number;
  height: number;
  format: string;
}

// Define the context shape
export interface MosaicContextType {
  mainImage: ImageInfo | null;
  setMainImage: (image: ImageInfo | null) => void;
  tileImages: ImageInfo[];
  setTileImages: (images: ImageInfo[]) => void;
  settings: {
    tileSize: number;
    tileDensity: number;
    colorAdjustment: number;
  };
  updateSettings: (newSettings: Partial<typeof defaultSettings>) => void;
  generationStatus: string | null;
  setGenerationStatus: (status: string | null) => void;
  generationId: string | null;
  setGenerationId: (id: string | null) => void;
}

const defaultSettings = {
  tileSize: 50,
  tileDensity: 80,
  colorAdjustment: 50,
};

// Create the context
const MosaicContext = createContext<MosaicContextType>({
  mainImage: null,
  setMainImage: () => {},
  tileImages: [],
  setTileImages: () => {},
  settings: defaultSettings,
  updateSettings: () => {},
  generationStatus: null,
  setGenerationStatus: () => {},
  generationId: null,
  setGenerationId: () => {},
});

export const MosaicProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [mainImage, setMainImage] = useState<ImageInfo | null>(null);
  const [tileImages, setTileImages] = useState<ImageInfo[]>([]);
  const [settings, setSettings] = useState(defaultSettings);
  const [generationStatus, setGenerationStatus] = useState<string | null>(null);
  const [generationId, setGenerationId] = useState<string | null>(null);

  const updateSettings = (newSettings: Partial<typeof defaultSettings>) => {
    setSettings(prev => ({ ...prev, ...newSettings }));
  };

  const value: MosaicContextType = {
    mainImage,
    setMainImage,
    tileImages,
    setTileImages,
    settings,
    updateSettings,
    generationStatus,
    setGenerationStatus,
    generationId,
    setGenerationId
  };

  return (
      <MosaicContext.Provider value={value}>
        {children}
      </MosaicContext.Provider>
  );
};

export const useMosaic = () => useContext(MosaicContext);