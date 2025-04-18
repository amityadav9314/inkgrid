import { api } from './api';

interface Image {
  id: string;
  user_id: number;
  project_id?: number;
  type: 'main' | 'tile';
  path: string;
  filename: string;
  width: number;
  height: number;
  format: string;
}

interface Project {
  id: number;
  user_id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
  settings: {
    tile_size: number;
    tile_density: number;
    overlay_ratio: number;
    style: string;
    [key: string]: any;
  };
  status: string;
  main_image?: Image;
}

interface CreateProjectData {
  name: string;
  description?: string;
  main_image_id?: string;
  settings?: {
    [key: string]: any;
  };
}

interface UpdateProjectData {
  name?: string;
  description?: string;
  main_image_id?: string;
  settings?: {
    [key: string]: any;
  };
}

class ProjectService {
  async getProjects(): Promise<Project[]> {
    try {
      const response = await api.get<{ projects: Project[] }>('/api/projects');
      return response.projects;
    } catch (error) {
      console.error('Error fetching projects:', error);
      throw error;
    }
  }

  async getProject(id: number): Promise<Project> {
    try {
      const response = await api.get<Project>(`/api/projects/${id}`);
      return response;
    } catch (error) {
      console.error(`Error fetching project ${id}:`, error);
      throw error;
    }
  }

  async createProject(data: CreateProjectData): Promise<Project> {
    try {
      const response = await api.post<Project>('/api/projects', data);
      return response;
    } catch (error) {
      console.error('Error creating project:', error);
      throw error;
    }
  }

  async updateProject(id: number, data: UpdateProjectData): Promise<Project> {
    try {
      const response = await api.put<Project>(`/api/projects/${id}`, data);
      return response;
    } catch (error) {
      console.error(`Error updating project ${id}:`, error);
      throw error;
    }
  }

  async deleteProject(id: number): Promise<void> {
    try {
      await api.delete(`/api/projects/${id}`);
    } catch (error) {
      console.error(`Error deleting project ${id}:`, error);
      throw error;
    }
  }
}

export const projectService = new ProjectService();
