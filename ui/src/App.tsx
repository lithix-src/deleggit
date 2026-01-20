import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AdminPage } from "./pages/AdminPage";
import { MissionPage } from "./pages/MissionPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MissionPage />} />
        <Route path="/admin" element={<AdminPage />} />
      </Routes>
    </Router>
  );
}

export default App;
