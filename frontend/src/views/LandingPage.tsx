import React from 'react';
import Header from '../components/Header';
import Hero from '../components/Hero';
import Features from '../components/Features';
import Footer from '../components/Footer';

const LandingPage: React.FC = () => {
  return (
    <>
      <Header />
      <Hero />
      <Features />
      <Footer />
    </>
  );
};

export default LandingPage;