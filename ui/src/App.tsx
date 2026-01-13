import { useState } from "react";
import { SensorGrid } from "./features/dashboard/SensorGrid";
import { Sidebar } from "./components/layout/Sidebar";
import { ActiveAgents } from "./features/dashboard/ActiveAgents";
import { RepoFeed } from "./features/dashboard/RepoFeed";

function App() {
  const [currentView, setCurrentView] = useState("dashboard");

  return (
    <div className="flex min-h-screen bg-black text-zinc-100 font-sans selection:bg-emerald-500/30">
      <Sidebar currentView={currentView} setView={setCurrentView} />

      <div className="flex-1 flex flex-col h-screen overflow-hidden">
        <header className="border-b border-zinc-800 bg-zinc-950/50 backdrop-blur h-14 shrink-0 px-6 flex items-center justify-between">
          <div className="flex items-center gap-2 md:hidden">
            <div className="w-3 h-3 bg-emerald-500 rounded-full animate-pulse" />
            <h1 className="font-bold tracking-tight text-lg">DELEGGIT</h1>
          </div>

          <div className="flex items-center gap-4 ml-auto">
            <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-zinc-900 border border-zinc-800">
              <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
              <span className="text-xs font-mono text-zinc-400">SYSTEM ONLINE</span>
            </div>
          </div>
        </header>

        <main className="flex-1 p-6 overflow-y-auto bg-black/50">
          <div className="max-w-7xl mx-auto space-y-6">

            {currentView === "dashboard" && (
              <>
                {/* Row 1: Hardware Metrics */}
                <div>
                  <h2 className="text-sm font-mono text-zinc-500 mb-2 uppercase tracking-wider">Hardware Telemetry</h2>
                  <SensorGrid />
                </div>

                {/* Row 2: Agent Swarm & External Events */}
                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-[400px]">
                  <div className="lg:col-span-2 h-full flex flex-col">
                    <h2 className="text-sm font-mono text-zinc-500 mb-2 uppercase tracking-wider">Agent Operations</h2>
                    <ActiveAgents />
                  </div>
                  <div className="lg:col-span-1 h-full flex flex-col">
                    <h2 className="text-sm font-mono text-zinc-500 mb-2 uppercase tracking-wider">Repository Watch</h2>
                    <RepoFeed />
                  </div>
                </div>
              </>
            )}

            {currentView === "workflows" && (
              <div className="flex items-center justify-center h-[500px] border border-dashed border-zinc-800 rounded bg-zinc-900/20">
                <div className="text-center">
                  <h2 className="text-xl font-mono text-zinc-500 mb-2">WORKFLOW CANVAS</h2>
                  <p className="text-zinc-700 text-sm">Waiting for Phase 2...</p>
                </div>
              </div>
            )}

            {currentView === "hardware" && (
              <div className="flex items-center justify-center h-[500px] border border-dashed border-zinc-800 rounded bg-zinc-900/20">
                <div className="text-center">
                  <h2 className="text-xl font-mono text-zinc-500 mb-2">HARDWARE BRIDGE CONFIG</h2>
                  <p className="text-zinc-700 text-sm">Waiting for Phase 2...</p>
                </div>
              </div>
            )}

            {/* Fallback for others */}
            {(currentView === "agents" || currentView === "settings") && (
              <div className="flex items-center justify-center h-[500px] border border-dashed border-zinc-800 rounded bg-zinc-900/20">
                <div className="text-center">
                  <h2 className="text-xl font-mono text-zinc-500 mb-2 uppercase">{currentView}</h2>
                  <p className="text-zinc-700 text-sm">Under Construction</p>
                </div>
              </div>
            )}

          </div>
        </main>
      </div>
    </div>
  );
}

export default App;
