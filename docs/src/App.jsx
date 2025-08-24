import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './pages/Home';
import QuickStart from './pages/QuickStart';
import Installation from './pages/Installation';
import Configuration from './pages/Configuration';
import Features from './pages/Features';
import API from './pages/API';
import Troubleshooting from './pages/Troubleshooting';
import Terms from './pages/Terms';
import './App.css';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/quick-start" element={<QuickStart />} />
          <Route path="/installation" element={<Installation />} />
          <Route path="/configuration" element={<Configuration />} />
          <Route path="/features" element={<Features />} />
          <Route path="/api" element={<API />} />
          <Route path="/troubleshooting" element={<Troubleshooting />} />
          <Route path="/terms" element={<Terms />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
