import { api } from './api';

interface MosaicGenerationRequest {
  project_id?: number;
  main_image_id: string;
  tile_image_ids: string[];
  tile_size: number;
  tile_density: number;
  overlay_ratio: number;
  style: 'classic' | 'random' | 'flowing';
  color_correction: boolean;
}

interface MosaicGenerationResponse {
  id: string;
  status: string;
  created_at: string;
}

interface MosaicGenerationStatus {
  id: string;
  status: string;
  progress: number;
  created_at: string;
  updated_at: string;
  result_url?: string;
}

class MosaicService {

  async uploadMainImage(file: File, projectId?: number): Promise<any> {
    try {
      const formData = new FormData();
      formData.append('image', file);

      if (projectId) {
        formData.append('project_id', projectId.toString());
      }

      const response = await api.upload('/images/main', formData);
      console.log('MosaicService: API response:', response); // Debugging line
      return response;
    } catch (error) {
      console.error('MosaicService: Error uploading main image:', error); // Ensure this is logged
      throw error; // Propagate the error
    }
  }

  async uploadTileImages(files: File[], projectId?: number, collectionId?: number): Promise<any> {
    try {
      const formData = new FormData();
      
      // Append each file to the form data
      files.forEach(file => {
        formData.append('images[]', file);
      });
      
      if (projectId) {
        formData.append('project_id', projectId.toString());
      }
      
      if (collectionId) {
        formData.append('collection_id', collectionId.toString());
      }
      
      return await api.upload('/images/tiles', formData);
    } catch (error) {
      console.error('Error uploading tile images:', error);
      throw error;
    }
  }

  async getTileCollections(): Promise<any> {
    try {
      return await api.get('/images/tiles');
    } catch (error) {
      console.error('Error fetching tile collections:', error);
      throw error;
    }
  }

  async generateMosaic(data: MosaicGenerationRequest): Promise<MosaicGenerationResponse> {
    try {
      return await api.post<MosaicGenerationResponse>('/generate', data);
    } catch (error) {
      console.error('Error generating mosaic:', error);
      throw error;
    }
  }

  async getGenerationStatus(id: string): Promise<MosaicGenerationStatus> {
    try {
      return await api.get<MosaicGenerationStatus>(`/generate/${id}/status`);
    } catch (error) {
      console.error(`Error fetching generation status for ${id}:`, error);
      throw error;
    }
  }
}

export const mosaicService = new MosaicService();
