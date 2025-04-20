import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import styled from 'styled-components';
import { useAuth } from '../context/AuthContext';
import Button from '../components/common/Button';
import { projectService } from '../services/projectService';

const Container = styled.div`
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
`;

const Title = styled.h1`
  font-size: 2rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 2rem;
`;

const FormGroup = styled.div`
  margin-bottom: 1.5rem;
`;

const Label = styled.label`
  display: block;
  font-weight: 600;
  margin-bottom: 0.5rem;
  color: #374151;
`;

const Input = styled.input`
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 1rem;
  
  &:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }
`;

const TextArea = styled.textarea`
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 1rem;
  min-height: 150px;
  resize: vertical;
  
  &:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
  }
`;

const ButtonContainer = styled.div`
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
`;

const ErrorMessage = styled.div`
  color: #ef4444;
  margin-top: 0.5rem;
  font-size: 0.875rem;
`;

const NewProject: React.FC = () => {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name.trim()) {
      setError('Project name is required');
      return;
    }
    
    try {
      setIsSubmitting(true);
      setError(null);
      
      const project = await projectService.createProject({
        name: name.trim(),
        description: description.trim()
      });
      
      // Navigate to the project edit page
      navigate(`/projects/${project.id}/edit`);
    } catch (err) {
      console.error('Error creating project:', err);
      setError('Failed to create project. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  };
  
  if (!isAuthenticated) {
    // Redirect to login if not authenticated
    navigate('/auth?mode=login');
    return null;
  }
  
  return (
    <Container>
      <Title>Create New Project</Title>
      
      <form onSubmit={handleSubmit}>
        <FormGroup>
          <Label htmlFor="name">Project Name</Label>
          <Input
            id="name"
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Enter a name for your project"
            required
          />
        </FormGroup>
        
        <FormGroup>
          <Label htmlFor="description">Description (Optional)</Label>
          <TextArea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Describe your project"
          />
        </FormGroup>
        
        {error && <ErrorMessage>{error}</ErrorMessage>}
        
        <ButtonContainer>
          <Button
            type="button"
            secondary
            onClick={() => navigate('/projects/')}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            primary
            disabled={isSubmitting}
          >
            {isSubmitting ? 'Creating...' : 'Create Project'}
          </Button>
        </ButtonContainer>
      </form>
    </Container>
  );
};

export default NewProject;
