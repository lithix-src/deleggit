import { useEffect, useState } from "react";
import { api, Agent } from "../../lib/api";
import { Terminal, Shield, AlertTriangle, Settings, Eye } from "lucide-react";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { AgentDetailDialog } from "./AgentDetailDialog";

export function AgentSwarmGrid() {
    const [agents, setAgents] = useState<Agent[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [isAdmin, setIsAdmin] = useState(false);
    const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null);

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
        <div className="space-y-6">
            {/* Admin Controls Toolbar */}
            <div className="flex justify-end items-center gap-3 pb-4 border-b border-slate-200/60">
                <span className="text-xs font-mono text-slate-500 uppercase">Admin Mode</span>
                <Switch
                    checked={isAdmin}
                    onCheckedChange={setIsAdmin}
                />
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {agents.map(agent => (
                    <div
                        key={agent.id}
                        onClick={() => setSelectedAgent(agent)}
                        className="bg-white border border-slate-200 rounded-lg shadow-sm hover:shadow-md hover:border-indigo-200 transition-all cursor-pointer p-6 flex flex-col gap-4 group relative"
                    >
                        <div className="flex justify-between items-start">
                            <div className="p-3 bg-indigo-50 group-hover:bg-indigo-100 transition-colors rounded-lg text-indigo-600">
                                <Terminal className="h-6 w-6" />
                            </div>
                            <span className="px-2 py-1 bg-emerald-100 text-emerald-700 text-xs font-bold rounded uppercase tracking-wide">Active</span>
                        </div>

                        <div>
                            <h3 className="font-bold text-lg text-slate-800 group-hover:text-indigo-700 transition-colors">{agent.service}</h3>
                            <p className="text-sm text-slate-500 font-mono">{agent.role}</p>
                        </div>

                        <div className="mt-auto pt-4 border-t border-slate-100 flex items-center justify-between">
                            <div className="flex items-center gap-1 text-xs text-slate-400 font-mono">
                                <Shield className="h-3 w-3" /> Secure
                            </div>

                            {/* Hover Actions */}
                            <div className="opacity-0 group-hover:opacity-100 transition-opacity flex gap-2">
                                <Button size="icon" variant="ghost" className="h-6 w-6 text-slate-400 hover:text-indigo-600">
                                    <Eye className="h-4 w-4" />
                                </Button>
                                {isAdmin && (
                                    <Button size="icon" variant="ghost" className="h-6 w-6 text-slate-400 hover:text-indigo-600">
                                        <Settings className="h-4 w-4" />
                                    </Button>
                                )}
                            </div>
                        </div>
                    </div>
                ))}
            </div>

            <AgentDetailDialog
                agent={selectedAgent}
                isOpen={!!selectedAgent}
                onClose={() => setSelectedAgent(null)}
                isAdmin={isAdmin}
            />
        </div>
    );
}
