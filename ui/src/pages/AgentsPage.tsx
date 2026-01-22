import { Brain, Activity } from "lucide-react";
import { AgentSwarmGrid } from "../features/agents/AgentSwarmGrid";

export function AgentsPage() {
    return (
        <div className="min-h-screen bg-slate-100 text-slate-900 font-sans p-8">
            <header className="mb-8 flex items-center justify-between">
                <div>
                    <h1 className="text-2xl font-bold text-slate-900 flex items-center gap-3">
                        <Brain className="text-emerald-500 h-8 w-8" />
                        AGENT <span className="text-indigo-600">SWARM</span>
                    </h1>
                    <p className="text-slate-500 font-mono text-sm mt-1">Autonomous Entities Management Console</p>
                </div>
                <a href="/" className="text-sm font-medium text-slate-600 hover:text-indigo-600 flex items-center gap-2">
                    <Activity className="h-4 w-4" /> Return to Mission Control
                </a>
            </header>

            <AgentSwarmGrid />
        </div>
    );
}
