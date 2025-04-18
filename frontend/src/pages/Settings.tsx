import React, { useState } from 'react';
import styled from 'styled-components';
import { useAuth } from '../context/AuthContext';
import Button from '../components/common/Button';

const Container = styled.div`
  max-width: 800px;
  margin: 0 auto;
`;

const Title = styled.h1`
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 2rem;
  color: #1f2937;
`;

const Card = styled.div`
  background-color: white;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  padding: 1.5rem;
  margin-bottom: 2rem;
`;

const CardTitle = styled.h2`
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: #1f2937;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid #e5e7eb;
`;

const Form = styled.form`
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
`;

const FormGroup = styled.div`
  display: flex;
  flex-direction: column;
`;

const Label = styled.label`
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 0.5rem;
  color: #4b5563;
`;

const Input = styled.input`
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  font-size: 1rem;
  
  &:focus {
    outline: none;
    border-color: #3b82f6;
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
  }
`;

const SuccessMessage = styled.div`
  background-color: #d1fae5;
  color: #065f46;
  padding: 0.75rem;
  border-radius: 0.375rem;
  margin-bottom: 1.5rem;
`;

const ErrorMessage = styled.div`
  background-color: #fee2e2;
  color: #b91c1c;
  padding: 0.75rem;
  border-radius: 0.375rem;
  margin-bottom: 1.5rem;
`;

const ToggleGroup = styled.div`
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
`;

const ToggleItem = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
`;

const ToggleLabel = styled.div`
  display: flex;
  flex-direction: column;
`;

const ToggleTitle = styled.span`
  font-weight: 500;
  color: #1f2937;
`;

const ToggleDescription = styled.span`
  font-size: 0.875rem;
  color: #6b7280;
`;

const Toggle = styled.label`
  position: relative;
  display: inline-block;
  width: 3rem;
  height: 1.5rem;
`;

const ToggleInput = styled.input`
  opacity: 0;
  width: 0;
  height: 0;
  
  &:checked + span {
    background-color: #3b82f6;
  }
  
  &:checked + span:before {
    transform: translateX(1.5rem);
  }
`;

const ToggleSlider = styled.span`
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: #e5e7eb;
  transition: 0.4s;
  border-radius: 1.5rem;
  
  &:before {
    position: absolute;
    content: "";
    height: 1.25rem;
    width: 1.25rem;
    left: 0.125rem;
    bottom: 0.125rem;
    background-color: white;
    transition: 0.4s;
    border-radius: 50%;
  }
`;

const Settings: React.FC = () => {
  const { user, isAuthenticated } = useAuth();
  
  const [profileForm, setProfileForm] = useState({
    name: user?.name || '',
    email: user?.email || '',
  });
  
  const [passwordForm, setPasswordForm] = useState({
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });
  
  const [notifications, setNotifications] = useState({
    email: true,
    projectUpdates: true,
    marketing: false,
  });
  
  const [profileSuccess, setProfileSuccess] = useState<string | null>(null);
  const [profileError, setProfileError] = useState<string | null>(null);
  
  const [passwordSuccess, setPasswordSuccess] = useState<string | null>(null);
  const [passwordError, setPasswordError] = useState<string | null>(null);
  
  const handleProfileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setProfileForm(prev => ({ ...prev, [name]: value }));
  };
  
  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setPasswordForm(prev => ({ ...prev, [name]: value }));
  };
  
  const handleNotificationChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, checked } = e.target;
    setNotifications(prev => ({ ...prev, [name]: checked }));
  };
  
  const handleProfileSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Implement profile update
    setProfileSuccess('Profile updated successfully!');
    setProfileError(null);
    
    // Reset success message after 3 seconds
    setTimeout(() => {
      setProfileSuccess(null);
    }, 3000);
  };
  
  const handlePasswordSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      setPasswordError('New passwords do not match.');
      setPasswordSuccess(null);
      return;
    }
    
    // TODO: Implement password update
    setPasswordSuccess('Password updated successfully!');
    setPasswordError(null);
    
    // Reset form and success message
    setPasswordForm({
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    });
    
    setTimeout(() => {
      setPasswordSuccess(null);
    }, 3000);
  };
  
  if (!isAuthenticated) {
    return (
      <Container>
        <Title>Settings</Title>
        <Card>
          <div style={{ textAlign: 'center', padding: '2rem' }}>
            <p>Please log in to access your settings.</p>
          </div>
        </Card>
      </Container>
    );
  }
  
  return (
    <Container>
      <Title>Settings</Title>
      
      <Card>
        <CardTitle>Profile Information</CardTitle>
        {profileSuccess && <SuccessMessage>{profileSuccess}</SuccessMessage>}
        {profileError && <ErrorMessage>{profileError}</ErrorMessage>}
        
        <Form onSubmit={handleProfileSubmit}>
          <FormGroup>
            <Label htmlFor="name">Name</Label>
            <Input
              type="text"
              id="name"
              name="name"
              value={profileForm.name}
              onChange={handleProfileChange}
            />
          </FormGroup>
          
          <FormGroup>
            <Label htmlFor="email">Email</Label>
            <Input
              type="email"
              id="email"
              name="email"
              value={profileForm.email}
              onChange={handleProfileChange}
            />
          </FormGroup>
          
          <div>
            <Button primary type="submit">
              Save Changes
            </Button>
          </div>
        </Form>
      </Card>
      
      <Card>
        <CardTitle>Change Password</CardTitle>
        {passwordSuccess && <SuccessMessage>{passwordSuccess}</SuccessMessage>}
        {passwordError && <ErrorMessage>{passwordError}</ErrorMessage>}
        
        <Form onSubmit={handlePasswordSubmit}>
          <FormGroup>
            <Label htmlFor="currentPassword">Current Password</Label>
            <Input
              type="password"
              id="currentPassword"
              name="currentPassword"
              value={passwordForm.currentPassword}
              onChange={handlePasswordChange}
            />
          </FormGroup>
          
          <FormGroup>
            <Label htmlFor="newPassword">New Password</Label>
            <Input
              type="password"
              id="newPassword"
              name="newPassword"
              value={passwordForm.newPassword}
              onChange={handlePasswordChange}
            />
          </FormGroup>
          
          <FormGroup>
            <Label htmlFor="confirmPassword">Confirm New Password</Label>
            <Input
              type="password"
              id="confirmPassword"
              name="confirmPassword"
              value={passwordForm.confirmPassword}
              onChange={handlePasswordChange}
            />
          </FormGroup>
          
          <div>
            <Button primary type="submit">
              Update Password
            </Button>
          </div>
        </Form>
      </Card>
      
      <Card>
        <CardTitle>Notification Preferences</CardTitle>
        
        <ToggleGroup>
          <ToggleItem>
            <ToggleLabel>
              <ToggleTitle>Email Notifications</ToggleTitle>
              <ToggleDescription>Receive notifications via email</ToggleDescription>
            </ToggleLabel>
            <Toggle>
              <ToggleInput
                type="checkbox"
                name="email"
                checked={notifications.email}
                onChange={handleNotificationChange}
              />
              <ToggleSlider />
            </Toggle>
          </ToggleItem>
          
          <ToggleItem>
            <ToggleLabel>
              <ToggleTitle>Project Updates</ToggleTitle>
              <ToggleDescription>Get notified when your projects are completed</ToggleDescription>
            </ToggleLabel>
            <Toggle>
              <ToggleInput
                type="checkbox"
                name="projectUpdates"
                checked={notifications.projectUpdates}
                onChange={handleNotificationChange}
              />
              <ToggleSlider />
            </Toggle>
          </ToggleItem>
          
          <ToggleItem>
            <ToggleLabel>
              <ToggleTitle>Marketing</ToggleTitle>
              <ToggleDescription>Receive marketing and promotional emails</ToggleDescription>
            </ToggleLabel>
            <Toggle>
              <ToggleInput
                type="checkbox"
                name="marketing"
                checked={notifications.marketing}
                onChange={handleNotificationChange}
              />
              <ToggleSlider />
            </Toggle>
          </ToggleItem>
        </ToggleGroup>
        
        <div style={{ marginTop: '1.5rem' }}>
          <Button primary onClick={() => alert('Preferences saved!')}>
            Save Preferences
          </Button>
        </div>
      </Card>
    </Container>
  );
};

export default Settings;
