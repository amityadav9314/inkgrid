import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';

// Default API configuration
const BASE_URL = process.env.REACT_APP_BASE_URL || 'http://localhost:8034';
const API_PATH = process.env.REACT_APP_API_PATH || '/goinkgrid/api';
const API_URL = `${BASE_URL}${API_PATH}`;
const AUTH_URL = `${BASE_URL}${API_PATH}/auth`;

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      withCredentials: true,
    });

    console.log('API client initialized with baseURL:', API_URL);

    // Add response interceptor for debugging CORS issues
    this.client.interceptors.response.use(
        response => {
          return response;
        },
        error => {
          if (error.response) {
            console.error('API Error Response:', {
              status: error.response.status,
              headers: error.response.headers,
              data: error.response.data
            });
          } else if (error.request) {
            console.error('API Error Request:', error.request);
          } else {
            console.error('API Error:', error.message);
          }
          return Promise.reject(error);
        }
    );

    // Add request interceptor to include auth token
    this.client.interceptors.request.use(
        (config) => {
          const token = localStorage.getItem('token');
          if (token) {
            // Ensure the Authorization header is properly formatted as 'Bearer {token}'
            config.headers['Authorization'] = `Bearer ${token}`;
            console.log('Adding auth token to request:', config.url);
          } else {
            console.log('No auth token available for request:', config.url);
          }
          return config;
        },
        (error) => {
          console.error('Request interceptor error:', error);
          return Promise.reject(error);
        }
    );

    // Add response interceptor to handle common errors
    this.client.interceptors.response.use(
        (response) => response,
        (error) => {
          // Handle 401 Unauthorized errors (token expired)
          if (error.response && error.response.status === 401) {
            localStorage.removeItem('token');
            localStorage.removeItem('user');
            window.location.href = '/auth';
          }
          return Promise.reject(error);
        }
    );
  }

  // Generic GET request
  public async get<T = any>(
      url: string,
      config?: AxiosRequestConfig
  ): Promise<T> {
    const response: AxiosResponse<T> = await this.client.get(url, config);
    return response.data;
  }

  // Generic POST request
  public async post<T = any>(
      url: string,
      data?: any,
      config?: AxiosRequestConfig
  ): Promise<T> {
    const response: AxiosResponse<T> = await this.client.post(url, data, config);
    return response.data;
  }

  // Generic PUT request
  public async put<T = any>(
      url: string,
      data?: any,
      config?: AxiosRequestConfig
  ): Promise<T> {
    const response: AxiosResponse<T> = await this.client.put(url, data, config);
    return response.data;
  }

  // Generic DELETE request
  public async delete<T = any>(
      url: string,
      config?: AxiosRequestConfig
  ): Promise<T> {
    const response: AxiosResponse<T> = await this.client.delete(url, config);
    return response.data;
  }

  // Upload file(s)
  public async upload<T = any>(
      url: string,
      formData: FormData,
      config?: AxiosRequestConfig
  ): Promise<T> {
    const uploadConfig = {
      ...config,
      headers: {
        ...config?.headers,
        'Content-Type': 'multipart/form-data',
      },
    };

    try {
      const response: AxiosResponse<T> = await this.client.post(url, formData, uploadConfig);
      return response.data; // Return response.data instead of response
    } catch (error) {
      console.error("Upload error:", error);
      throw error;
    }
  }
}

export const api = new ApiClient();