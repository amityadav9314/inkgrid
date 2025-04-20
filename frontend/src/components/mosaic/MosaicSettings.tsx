import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { useMosaic } from '../../context/MosaicContext';
import { mosaicService } from '../../services/mosaicService';

const MosaicSettings: React.FC = () => {
  const { settings, updateSettings } = useMosaic();
  const [isLoading, setIsLoading] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [saveMessage, setSaveMessage] = useState<string | null>(null);
  const [localSettings, setLocalSettings] = useState({
    tileSize: settings.tileSize,
    tileDensity: settings.tileDensity,
    colorAdjustment: settings.colorAdjustment,
    style: 'classic' as 'classic' | 'random' | 'flowing',
  });

  useEffect(() => {
    const fetchSettings = async () => {
      setIsLoading(true);
      try {
        const response = await mosaicService.getMosaicSettings();
        if (response && response.settings) {
          // Map backend snake_case to frontend camelCase
          const fetchedSettings = {
            tileSize: response.settings.tile_size || settings.tileSize,
            tileDensity: response.settings.tile_density || settings.tileDensity,
            colorAdjustment: response.settings.color_adjustment || settings.colorAdjustment,
            style: response.settings.style || 'classic',
          };
          setLocalSettings(fetchedSettings);
          updateSettings(fetchedSettings);
        }
      } catch (error) {
        console.error('Error fetching settings:', error);
        // If we can't fetch settings, we'll use the defaults from context
      } finally {
        setIsLoading(false);
      }
    };

    fetchSettings();
  }, [updateSettings, settings]);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    
    // For range inputs, convert to number
    const numValue = e.target.type === 'range' ? parseInt(value, 10) : value;
    
    setLocalSettings({
      ...localSettings,
      [name]: numValue,
    });
  };

  const handleSaveSettings = async () => {
    setIsSaving(true);
    setSaveMessage(null);
    
    try {
      // Convert frontend camelCase to backend snake_case
      await mosaicService.saveMosaicSettings({
        tile_size: localSettings.tileSize,
        tile_density: localSettings.tileDensity,
        color_adjustment: localSettings.colorAdjustment,
        style: localSettings.style,
      });
      
      // Update the global context with the new settings
      updateSettings(localSettings);
      setSaveMessage('Settings saved successfully');
    } catch (error) {
      console.error('Error saving settings:', error);
      setSaveMessage('Failed to save settings. Please try again.');
    } finally {
      setIsSaving(false);
      
      // Clear the message after 3 seconds
      setTimeout(() => {
        setSaveMessage(null);
      }, 3000);
    }
  };

  if (isLoading) {
    return (
      <Container>
        <LoadingMessage>Loading settings...</LoadingMessage>
      </Container>
    );
  }

  return (
    <Container>
      <Title>Mosaic Settings</Title>
      <Description>
        Customize how your mosaic will be generated. These settings affect the
        appearance and quality of the final result.
      </Description>

      <SettingsForm>
        <SettingGroup>
          <SettingLabel>Tile Size</SettingLabel>
          <RangeContainer>
            <RangeInput
              type="range"
              name="tileSize"
              min="10"
              max="200"
              value={localSettings.tileSize}
              onChange={handleInputChange}
            />
            <RangeValue>{localSettings.tileSize}px</RangeValue>
          </RangeContainer>
          <SettingDescription>
            Controls the size of each tile in the mosaic. Smaller tiles create more detail but require more tile images.
          </SettingDescription>
        </SettingGroup>

        <SettingGroup>
          <SettingLabel>Tile Density</SettingLabel>
          <RangeContainer>
            <RangeInput
              type="range"
              name="tileDensity"
              min="1"
              max="100"
              value={localSettings.tileDensity}
              onChange={handleInputChange}
            />
            <RangeValue>{localSettings.tileDensity}%</RangeValue>
          </RangeContainer>
          <SettingDescription>
            Controls how densely packed the tiles are. Higher density means less gaps between tiles.
          </SettingDescription>
        </SettingGroup>

        <SettingGroup>
          <SettingLabel>Color Adjustment</SettingLabel>
          <RangeContainer>
            <RangeInput
              type="range"
              name="colorAdjustment"
              min="0"
              max="100"
              value={localSettings.colorAdjustment}
              onChange={handleInputChange}
            />
            <RangeValue>{localSettings.colorAdjustment}%</RangeValue>
          </RangeContainer>
          <SettingDescription>
            Controls how much the tile colors are adjusted to match the main image. Higher values create a more accurate representation.
          </SettingDescription>
        </SettingGroup>

        <SettingGroup>
          <SettingLabel>Style</SettingLabel>
          <SelectInput
            name="style"
            value={localSettings.style}
            onChange={handleInputChange}
          >
            <option value="classic">Classic</option>
            <option value="random">Random</option>
            <option value="flowing">Flowing</option>
          </SelectInput>
          <SettingDescription>
            Choose the arrangement style for your mosaic tiles.
          </SettingDescription>
        </SettingGroup>

        <SaveButtonContainer>
          <SaveButton onClick={handleSaveSettings} disabled={isSaving}>
            {isSaving ? 'Saving...' : 'Save Settings'}
          </SaveButton>
          {saveMessage && (
            <SaveMessage success={saveMessage.includes('success')}>
              {saveMessage}
            </SaveMessage>
          )}
        </SaveButtonContainer>
      </SettingsForm>
    </Container>
  );
};

const Container = styled.div`
  padding: 2rem;
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
`;

const Title = styled.h2`
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #333;
`;

const Description = styled.p`
  color: #666;
  margin-bottom: 2rem;
  line-height: 1.5;
`;

const SettingsForm = styled.div`
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const SettingGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
`;

const SettingLabel = styled.label`
  font-weight: 600;
  font-size: 1rem;
  color: #333;
`;

const SettingDescription = styled.p`
  font-size: 0.875rem;
  color: #666;
  margin-top: 0.25rem;
`;

const RangeContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 1rem;
`;

const RangeInput = styled.input`
  flex: 1;
  height: 6px;
  -webkit-appearance: none;
  appearance: none;
  width: 100%;
  background: #ddd;
  outline: none;
  border-radius: 3px;

  &::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    background: #3b82f6;
    border-radius: 50%;
    cursor: pointer;
  }

  &::-moz-range-thumb {
    width: 18px;
    height: 18px;
    background: #3b82f6;
    border-radius: 50%;
    cursor: pointer;
    border: none;
  }
`;

const RangeValue = styled.span`
  font-size: 0.875rem;
  font-weight: 600;
  color: #3b82f6;
  min-width: 60px;
  text-align: right;
`;

const SelectInput = styled.select`
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  background-color: white;
  font-size: 1rem;
  color: #333;
  width: 100%;
  max-width: 300px;

  &:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
  }
`;

const SaveButtonContainer = styled.div`
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0.5rem;
`;

const SaveButton = styled.button`
  background-color: #3b82f6;
  color: white;
  font-weight: 600;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;

  &:hover {
    background-color: #2563eb;
  }

  &:disabled {
    background-color: #93c5fd;
    cursor: not-allowed;
  }
`;

const SaveMessage = styled.div<{ success: boolean }>`
  font-size: 0.875rem;
  color: ${props => (props.success ? '#10b981' : '#ef4444')};
  font-weight: 500;
`;

const LoadingMessage = styled.div`
  text-align: center;
  padding: 2rem;
  color: #666;
  font-size: 1rem;
`;

export default MosaicSettings;
