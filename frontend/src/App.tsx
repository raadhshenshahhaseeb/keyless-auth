import LandingPage from "./views/LandingPage"
import Wallet from "./views/Wallet"
import React from "react"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
const App:React.FC = () => {
  

  return (
    <>
      <Router>
        <Routes>
          <Route path="/" element={<LandingPage />} />
          <Route path="/wallet" element={<Wallet />} />
        </Routes>
    </Router>
    </>
  )
}

export default App
