import { useState } from "react";
import { SensorGrid } from "../features/dashboard/SensorGrid";
import { Sidebar } from "../components/layout/Sidebar";
import { ActiveAgents } from "../features/dashboard/ActiveAgents";
import { ContainerGrid } from "../features/dashboard/ContainerGrid";
import { GenericSensorGrid } from "../features/dashboard/GenericSensorGrid";
import { RepoFeed } from "../features/dashboard/RepoFeed";

export function AdminPage() {
    const [currentView, setCurrentView] = useState("dashboard");

    return (
        <div className="flex min-h-screen bg-slate-100 text-slate-900 font-sans selection:bg-emerald-500/30">
            <Sidebar currentView={currentView} setView={setCurrentView} />

            <div className="flex-1 flex flex-col h-screen overflow-hidden">
                <header className="border-b border-slate-200 bg-white/80 backdrop-blur h-14 shrink-0 px-6 flex items-center justify-between">
                    <div className="flex items-center gap-2 md:hidden">
                        <div className="w-3 h-3 bg-emerald-500 rounded-full animate-pulse" />
                        <h1 className="font-bold tracking-tight text-lg">CATALYST ADMIN</h1>
                    </div>

                    <div className="flex items-center gap-4 ml-auto">
                        <div className="flex items-center gap-2 px-3 py-1 rounded-full bg-white border border-slate-200 shadow-sm">
                            <div className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></div>
                            <span className="text-xs font-mono text-slate-500">ADMIN ACCESS</span>
                        </div>
                    </div>
                </header>

                <main className="flex-1 p-6 overflow-y-auto bg-slate-100">
                    <div className="max-w-7xl mx-auto space-y-6">

                        {currentView === "dashboard" && (
                            <>
                                {/* Row 1: Hardware Metrics */}
                                <div>
                                    <h2 className="text-sm font-mono text-slate-500 mb-2 uppercase tracking-wider">Hardware Telemetry</h2>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6 pb-6">
                                        <SensorGrid />
                                        <ContainerGrid />
                                    </div>

                                    {/* Dynamic Grid (V2) */}
                                    <GenericSensorGrid />
                                </div>

                                {/* Row 2: Agent Swarm & External Events */}
                                <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 h-[400px]">
                                    <div className="lg:col-span-2 h-full flex flex-col">
                                        <h2 className="text-sm font-mono text-slate-500 mb-2 uppercase tracking-wider">Agent Operations</h2>
                                        <ActiveAgents />
                                    </div>
                                    <div className="lg:col-span-1 h-full flex flex-col">
                                        <h2 className="text-sm font-mono text-slate-500 mb-2 uppercase tracking-wider">Repository Watch</h2>
                                        <RepoFeed />
                                    </div>
                                </div>
                            </>
                        )}

                        {currentView === "workflows" && (
                            <div className="flex items-center justify-center h-[500px] border border-dashed border-slate-300 rounded bg-slate-50">
                                <div className="text-center">
                                    <h2 className="text-xl font-mono text-slate-500 mb-2">WORKFLOW CANVAS</h2>
                                    <p className="text-slate-700 text-sm">Waiting for Phase 2...</p>
                                </div>
                            </div>
                        )}

                        {currentView === "hardware" && (
                            <div className="flex items-center justify-center h-[500px] border border-dashed border-slate-300 rounded bg-slate-50">
                                <div className="text-center">
                                    <h2 className="text-xl font-mono text-slate-500 mb-2">HARDWARE BRIDGE CONFIG</h2>
                                    <p className="text-slate-700 text-sm">Waiting for Phase 2...</p>
                                </div>
                            </div>
                        )}

                        {/* Fallback for others */}
                        {(currentView === "agents" || currentView === "settings") && (
                            <div className="flex items-center justify-center h-[500px] border border-dashed border-slate-300 rounded bg-slate-50">
                                <div className="text-center">
                                    <h2 className="text-xl font-mono text-slate-500 mb-2 uppercase">{currentView}</h2>
                                    <p className="text-slate-700 text-sm">Under Construction</p>
                                </div>
                            </div>
                        )}

                    </div>
                </main>
            </div>
        </div>
    );
}
