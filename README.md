# InkGrid - Mosaic Image Generator

InkGrid is a powerful application that transforms a main image into an artistic mosaic composed of numerous smaller images serving as tiles. It provides an intuitive interface for uploading, arranging, and customizing mosaic parameters to create stunning visual compositions that can be saved and shared across platforms.

## Features

- Upload a main image and multiple tile images
- Intelligent tile placement based on color matching and pattern recognition
- Customizable tile size, density, and arrangement
- Multiple mosaic styles (classic grid, random, flowing patterns)
- Real-time preview of mosaic effect
- High-resolution export options
- Project management for saving and editing mosaics

## Technology Stack

### Backend
- Golang with Gin framework
- PostgreSQL database
- Docker containerization

### Frontend
- React with TypeScript
- Electron for desktop application

## Getting Started

### Prerequisites
- Go 1.16+
- Node.js 14+
- PostgreSQL 12+
- Docker and Docker Compose

### Installation

1. Clone the repository
   ```
   git clone https://github.com/yourusername/inkgrid.git
   cd inkgrid
   ```

2. Start the application using Docker Compose
   ```
   docker-compose up
   ```

3. Access the application
   - Web: http://localhost:3000
   - API: http://localhost:8034/goinkgrid

## Development

### Backend
```
cd backend
go mod download
go run cmd/server/main.go
```

### Frontend
```
cd frontend
npm install
npm start
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.