import React, { createContext, useState, useContext } from 'react';

interface Image {
  id: string;
  path: string;
  filename: string;
  width: number;
  height: number;
  format: string;
}

interface MosaicSettings {
  tileSize: number;
  tileDensity: number;
  overlayRatio: number;
  style: 'classic' | 'random' | 'flowing';
  colorCorrection: boolean;
}

interface MosaicContextType {
  mainImage: Image | null;
  tileImages: Image[];
  settings: MosaicSettings;
  generationId: string | null;
  generationStatus: string | null;
  setMainImage: (image: Image | null) => void;
  addTileImage: (image: Image) => void;
  removeTileImage: (id: string) => void;
  clearTileImages: () => void;
  updateSettings: (settings: Partial<MosaicSettings>) => void;
  resetSettings: () => void;
  setGenerationId: (id: string | null) => void;
  setGenerationStatus: (status: string | null) => void;
}

const defaultSettings: MosaicSettings = {
  tileSize: 50,
  tileDensity: 80,
  overlayRatio: 0.7,
  style: 'classic',
  colorCorrection: true
};

const MosaicContext = createContext<MosaicContextType | undefined>(undefined);

export const useMosaic = () => {
  const context = useContext(MosaicContext);
  if (context === undefined) {
    throw new Error('useMosaic must be used within a MosaicProvider');
  }
  return context;
};

export const MosaicProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [mainImage, setMainImage] = useState<Image | null>(null);
  const [tileImages, setTileImages] = useState<Image[]>([]);
  const [settings, setSettings] = useState<MosaicSettings>(defaultSettings);
  const [generationId, setGenerationId] = useState<string | null>(null);
  const [generationStatus, setGenerationStatus] = useState<string | null>(null);

  const addTileImage = (image: Image) => {
    setTileImages(prevImages => [...prevImages, image]);
  };

  const removeTileImage = (id: string) => {
    setTileImages(prevImages => prevImages.filter(image => image.id !== id));
  };

  const clearTileImages = () => {
    setTileImages([]);
  };

  const updateSettings = (newSettings: Partial<MosaicSettings>) => {
    setSettings(prevSettings => ({ ...prevSettings, ...newSettings }));
  };

  const resetSettings = () => {
    setSettings(defaultSettings);
  };

  const value = {
    mainImage,
    tileImages,
    settings,
    generationId,
    generationStatus,
    setMainImage,
    addTileImage,
    removeTileImage,
    clearTileImages,
    updateSettings,
    resetSettings,
    setGenerationId,
    setGenerationStatus
  };

  return (
    <MosaicContext.Provider value={value}>
      {children}
    </MosaicContext.Provider>
  );
};
