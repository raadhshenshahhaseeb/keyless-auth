import React from "react"

import Demo from "./views/Demo"
import LandingPage from "./views/LandingPage"
import Wallet from "./views/Wallet"

import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
const App:React.FC = () => {
  

  return (
    <>
      <Router>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/wallet" element={<Wallet />} />
          <Route path="/demo" element={<Demo />} />
        </Routes>
    </Router>
    </>
  )
}

export default App
