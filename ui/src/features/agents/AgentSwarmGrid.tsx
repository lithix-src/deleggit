import { useEffect, useState } from "react";
import { api, Agent } from "../../lib/api";
import { Terminal, Shield, AlertTriangle } from "lucide-react";

export function AgentSwarmGrid() {
    const [agents, setAgents] = useState<Agent[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        api.getAgents()
            .then(setAgents)
            .catch(err => {
                console.error(err);
                setError("Failed to load swarm data.");
            })
            .finally(() => setIsLoading(false));
    }, []);

    if (isLoading) {
        return <div className="text-center py-20 animate-pulse text-indigo-500 font-mono">Loading Swarm Data...</div>;
    }

    if (error) {
        return (
            <div className="flex items-center justify-center p-8 text-red-500 gap-2">
                <AlertTriangle className="h-5 w-5" />
                <span>{error}</span>
            </div>
        );
    }

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {agents.map(agent => (
                <div key={agent.id} className="bg-white border border-slate-200 rounded-lg shadow-sm hover:shadow-md transition-shadow p-6 flex flex-col gap-4">
                    <div className="flex justify-between items-start">
                        <div className="p-3 bg-indigo-50 rounded-lg text-indigo-600">
                            <Terminal className="h-6 w-6" />
                        </div>
                        <span className="px-2 py-1 bg-emerald-100 text-emerald-700 text-xs font-bold rounded uppercase tracking-wide">Active</span>
                    </div>

                    <div>
                        <h3 className="font-bold text-lg text-slate-800">{agent.service}</h3>
                        <p className="text-sm text-slate-500 font-mono">{agent.role}</p>
                    </div>

                    <div className="mt-auto pt-4 border-t border-slate-100">
                        <div className="flex items-center justify-between text-xs text-slate-400 font-mono">
                            <span>ID: {agent.id.substring(0, 8)}...</span>
                            <div className="flex items-center gap-1 text-slate-500">
                                <Shield className="h-3 w-3" /> Secure
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </div>
    );
}
