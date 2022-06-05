import { useEffect } from 'react'
import {Routes, Route, Navigate} from 'react-router-dom'
import Home from './pages/home'
import Job from './pages/job'

function RedirectHome() {
  return <div></div>
}

function App() {

  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/job" element={<Job />} />
        <Route path="*" element={<Navigate replace to="/" />} />
      </Routes>
    </div>
  )
}

export default App
