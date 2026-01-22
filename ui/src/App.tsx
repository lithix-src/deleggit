import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AdminPage } from "./pages/AdminPage";
import { MissionPage } from "./pages/MissionPage";
import { AgentsPage } from "./pages/AgentsPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MissionPage />} />
        <Route path="/admin" element={<AdminPage />} />
        <Route path="/agents" element={<AgentsPage />} />
      </Routes>
    </Router>
  );
}

export default App;
