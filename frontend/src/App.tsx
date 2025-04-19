import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { AuthProvider } from './context/AuthContext';
import { MosaicProvider } from './context/MosaicContext';
import Header from './components/layout/Header';
import Footer from './components/layout/Footer';
import Home from './pages/Home';
import Auth from './pages/Auth';
import MosaicCreator from './pages/MosaicCreator';
import Projects from './pages/Projects';
import Settings from './pages/Settings';
import './styles/global.css';
import Sidebar from 'components/layout/Sidebar';

const App: React.FC = () => {
  return (
    <Router>
      <AuthProvider>
        <MosaicProvider>
          <div className="app">
            <Header />
            {/* <Sidebar isOpen={false} toggleSidebar={function (): void { */}
              {/* throw new Error('Function not implemented.'); */}
            {/* } } /> */}
            <main className="main-content">
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/auth" element={<Auth />} />
                <Route path="/create" element={<MosaicCreator />} />
                <Route path="/projects" element={<Projects />} />
                <Route path="/settings" element={<Settings />} />
              </Routes>
            </main>
            <Footer />
          </div>
        </MosaicProvider>
      </AuthProvider>
    </Router>
  );
};

export default App;
