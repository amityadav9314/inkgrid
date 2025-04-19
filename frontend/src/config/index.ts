// frontend/src/config/index.ts
export const config = {
    BASE_URL: process.env.REACT_APP_BASE_URL || 'http://localhost:8034',
    API_PATH: process.env.REACT_APP_API_PATH || '/goinkgrid/api',

    // Computed properties
    get API_URL() {
        return `${this.BASE_URL}${this.API_PATH}`;
    },

    // Helper for image paths
    getImageUrl(path: string): string {
        if (!path) return '';
        if (path.startsWith('http')) return path;
        // Remove leading slash if present to avoid double slashes
        const formattedPath = path.startsWith('/') ? path.substring(1) : path;
        return `${this.BASE_URL}/${formattedPath}`;
    }
};