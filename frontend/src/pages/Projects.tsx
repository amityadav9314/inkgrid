import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import styled from 'styled-components';
import { useAuth } from '../context/AuthContext';
import Button from '../components/common/Button';
import LoadingIndicator from '../components/common/LoadingIndicator';
import { projectService } from '../services/projectService';

interface Project {
  id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
  status: string;
  main_image?: {
    path: string;
    filename: string;
  };
}

const Container = styled.div`
  max-width: 1200px;
  margin: 0 auto;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
`;

const Title = styled.h1`
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
`;

const EmptyState = styled.div`
  text-align: center;
  padding: 4rem 2rem;
  background-color: #f9fafb;
  border-radius: 0.5rem;
  border: 2px dashed #e5e7eb;
`;

const EmptyStateIcon = styled.div`
  margin-bottom: 1.5rem;
  color: #9ca3af;
  
  svg {
    width: 4rem;
    height: 4rem;
  }
`;

const EmptyStateTitle = styled.h2`
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #1f2937;
`;

const EmptyStateText = styled.p`
  color: #6b7280;
  margin-bottom: 1.5rem;
  max-width: 500px;
  margin-left: auto;
  margin-right: auto;
`;

const ProjectGrid = styled.div`
  display: grid;
  grid-template-columns: repeat(1, 1fr);
  gap: 2rem;
  
  @media (min-width: 640px) {
    grid-template-columns: repeat(2, 1fr);
  }
  
  @media (min-width: 1024px) {
    grid-template-columns: repeat(3, 1fr);
  }
`;

const ProjectCard = styled.div`
  background-color: white;
  border-radius: 0.5rem;
  overflow: hidden;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
  
  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1), 0 4px 6px -2px rgba(0, 0, 0, 0.05);
  }
`;

const ProjectImage = styled.div<{ imagePath?: string }>`
  height: 180px;
  background-color: #f3f4f6;
  background-image: ${props => props.imagePath ? `url(${props.imagePath})` : 'none'};
  background-size: cover;
  background-position: center;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  
  svg {
    width: 3rem;
    height: 3rem;
  }
`;

const ProjectContent = styled.div`
  padding: 1.5rem;
`;

const ProjectTitle = styled.h3`
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #1f2937;
`;

const ProjectDescription = styled.p`
  color: #6b7280;
  font-size: 0.875rem;
  margin-bottom: 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
`;

const ProjectMeta = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.75rem;
  color: #9ca3af;
  margin-bottom: 1rem;
`;

const ProjectStatus = styled.span<{ status: string }>`
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
  
  ${props => {
    switch (props.status) {
      case 'completed':
        return `
          background-color: #d1fae5;
          color: #065f46;
        `;
      case 'in_progress':
        return `
          background-color: #eff6ff;
          color: #1e40af;
        `;
      case 'new':
        return `
          background-color: #f3f4f6;
          color: #1f2937;
        `;
      default:
        return `
          background-color: #f3f4f6;
          color: #1f2937;
        `;
    }
  }}
`;

const ProjectActions = styled.div`
  display: flex;
  gap: 0.5rem;
`;

const Projects: React.FC = () => {
  const { isAuthenticated, user, token } = useAuth();
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // Debug log for Projects component
  useEffect(() => {
    console.log('Projects component - Auth state:', { 
      isAuthenticated, 
      user, 
      token,
      localStorageToken: localStorage.getItem('token'),
      localStorageUser: localStorage.getItem('user')
    });
  }, [isAuthenticated, user, token]);
  
  useEffect(() => {
    const fetchProjects = async () => {
      console.log('fetchProjects called, isAuthenticated:', isAuthenticated);
      if (!isAuthenticated) {
        console.log('Not authenticated, skipping fetch');
        return;
      }
      
      try {
        setIsLoading(true);
        console.log('Fetching projects...');
        const data = await projectService.getProjects();
        console.log('Projects data received:', data);
        setProjects(data);
        setError(null);
      } catch (err) {
        console.error('Error fetching projects:', err);
        setError('Failed to load projects. Please try again later.');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchProjects();
  }, [isAuthenticated]);
  
  if (!isAuthenticated) {
    return (
      <Container>
        <EmptyState>
          <EmptyStateIcon>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
            </svg>
          </EmptyStateIcon>
          <EmptyStateTitle>Authentication Required</EmptyStateTitle>
          <EmptyStateText>
            Please log in or create an account to view and manage your projects.
          </EmptyStateText>
          <Link to="/auth?mode=login">
            <Button primary>Log In</Button>
          </Link>
        </EmptyState>
      </Container>
    );
  }
  
  if (isLoading) {
    return (
      <Container>
        <div style={{ display: 'flex', justifyContent: 'center', padding: '4rem' }}>
          <LoadingIndicator size="large" text="Loading projects..." />
        </div>
      </Container>
    );
  }
  
  if (error) {
    return (
      <Container>
        <Header>
          <Title>My Projects</Title>
          <Link to="/projects/new">
            <Button primary>Create New</Button>
          </Link>
        </Header>
        <div style={{ textAlign: 'center', padding: '2rem', color: '#ef4444' }}>
          {error}
        </div>
      </Container>
    );
  }
  
  if (projects.length === 0) {
    return (
      <Container>
        <Header>
          <Title>My Projects</Title>
          <Link to="/projects/new">
            <Button primary>Create New</Button>
          </Link>
        </Header>
        <EmptyState>
          <EmptyStateIcon>
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </EmptyStateIcon>
          <EmptyStateTitle>No Projects Yet</EmptyStateTitle>
          <EmptyStateText>
            You haven't created any mosaic projects yet. Start by creating your first masterpiece!
          </EmptyStateText>
          <Link to="/projects/new">
            <Button primary>Create Your First Mosaic</Button>
          </Link>
        </EmptyState>
      </Container>
    );
  }
  
  return (
    <Container>
      <Header>
        <Title>My Projects</Title>
        <Link to="/projects/new">
          <Button primary>Create New</Button>
        </Link>
      </Header>
      
      <ProjectGrid>
        {projects.map(project => (
          <ProjectCard key={project.id}>
            <ProjectImage imagePath={project.main_image?.path}>
              {!project.main_image && (
                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              )}
            </ProjectImage>
            
            <ProjectContent>
              <ProjectTitle>{project.name}</ProjectTitle>
              <ProjectDescription>{project.description || 'No description'}</ProjectDescription>
              
              <ProjectMeta>
                <span>Created: {new Date(project.created_at).toLocaleDateString()}</span>
                <ProjectStatus status={project.status}>
                  {project.status === 'completed' ? 'Completed' : 
                   project.status === 'in_progress' ? 'In Progress' : 'New'}
                </ProjectStatus>
              </ProjectMeta>
              
              <ProjectActions>
                <Link to={`/projects/${project.id}/edit`} style={{ flex: 1 }}>
                  <Button primary fullWidth>Edit</Button>
                </Link>
                <Link to={`/projects/${project.id}/view`} style={{ flex: 1 }}>
                  <Button secondary fullWidth>View</Button>
                </Link>
              </ProjectActions>
            </ProjectContent>
          </ProjectCard>
        ))}
      </ProjectGrid>
    </Container>
  );
};

export default Projects;
