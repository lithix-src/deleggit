import { motion } from "framer-motion";
import { Brain, MessageSquare, Terminal, Activity, ChevronDown, Plus } from "lucide-react";
import { useState, useEffect } from "react";
import { Button } from "../components/ui/button";
import { api, Context, Repo } from "../lib/api";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "../components/ui/dropdown-menu";

export function MissionPage() {
    const [messages, setMessages] = useState([
        { role: "system", content: "Catalyst Vibe Engine Online. How can I assist?" }
    ]);
    const [input, setInput] = useState("");
    const [context, setContext] = useState<Context | null>(null);
    const [repos, setRepos] = useState<Repo[]>([]);
    const [isLoading, setIsLoading] = useState(false);

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            const [ctx, r] = await Promise.all([api.getContext(), api.getRepos()]);
            setContext(ctx);
            setRepos(r);
        } catch (e) {
            console.error(e);
        }
    };

    const handleSwitchRepo = async (repoId: string, defaultBranch: string) => {
        if (!context || context.active_repo_id === repoId) return;
        setIsLoading(true);
        try {
            await api.setContext(repoId, defaultBranch);
            await loadData(); // Reload to confirm
        } catch (e) {
            console.error("Context Switch Failed", e);
        } finally {
            setIsLoading(false);
        }
    };

    const handleSend = () => {
        if (!input.trim()) return;
        // Mock send
        setMessages([...messages, { role: "user", content: input }]);
        setInput("");
        // Simulate response
        setTimeout(() => {
            setMessages(prev => [...prev, { role: "system", content: "Processing your request via Liaison..." }]);
        }, 800);
    };

    return (
        <div className="min-h-screen bg-slate-100 text-slate-900 font-sans flex flex-col items-center justify-center p-4 relative overflow-hidden">

            {/* Ambient Background - Subtle */}
            <div className="absolute inset-0 bg-gradient-to-br from-slate-200/50 via-slate-100 to-slate-200/50 pointer-events-none" />

            {/* Command Bar Header */}
            <header className="absolute top-0 left-0 right-0 p-6 flex justify-between items-center z-10 pointer-events-none">
                <div className="flex items-center gap-3 pointer-events-auto">
                    <Brain className="text-emerald-500 h-6 w-6" />
                    <h1 className="text-xl font-bold tracking-widest text-slate-900">CATALYST <span className="text-indigo-600">CONTROL</span></h1>
                </div>

                {/* Navigation Actions */}
                <div className="flex items-center gap-4 pointer-events-auto">
                    {/* System Status Indicators (now prominent) */}
                    <div className="hidden md:flex gap-4 text-xs font-mono text-slate-500 mr-4 border-r border-slate-300 pr-4">
                        <span className="flex items-center gap-2"><Activity className="h-3 w-3 text-emerald-500" /> CORE: ONLINE</span>
                        <span className="flex items-center gap-2"><Terminal className="h-3 w-3 text-blue-500" /> SWARM: IDLE</span>
                        {isLoading && <span className="flex items-center gap-2 text-indigo-500 animate-pulse">SWITCHING CONTEXT...</span>}
                    </div>

                    <Button variant="ghost" size="sm" className="text-slate-600 hover:text-indigo-600">
                        <MessageSquare className="h-4 w-4 mr-2" />
                        Docs
                    </Button>

                    <a href="/admin">
                        <Button variant="outline" size="sm" className="bg-white border-slate-300 shadow-sm hover:border-indigo-400 hover:text-indigo-600 hover:shadow-md transition-all">
                            <Terminal className="h-4 w-4 mr-2" />
                            Admin Console
                        </Button>
                    </a>
                </div>
            </header>

            {/* Main Interface */}
            <main className="w-full max-w-4xl z-10 flex flex-col gap-6 mt-20">

                {/* Project Context HUD - Interactive */}
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">

                    {/* Active Repo Selector */}
                    <div className="bg-white/80 backdrop-blur border border-slate-200 p-4 rounded-lg shadow-sm flex items-center gap-3">
                        <div className="p-2 bg-indigo-50 rounded-md text-indigo-600">
                            <Activity className="h-5 w-5" />
                        </div>
                        <div className="flex-1">
                            <p className="text-[10px] font-mono text-slate-500 uppercase tracking-wider">Active Repository</p>
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <button className="flex items-center gap-1 text-sm font-semibold text-slate-800 hover:text-indigo-600 outline-none w-full">
                                        <span className="truncate">{context ? `${context.org}/${context.name}` : "Loading..."}</span>
                                        <ChevronDown className="h-3 w-3 opacity-50 ml-auto" />
                                    </button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="start" className="w-[200px]">
                                    {repos.map(r => (
                                        <DropdownMenuItem key={r.id} onClick={() => handleSwitchRepo(r.id, r.default_branch)}>
                                            <span className={context?.active_repo_id === r.id ? "font-bold text-indigo-600" : ""}>
                                                {r.org}/{r.name}
                                            </span>
                                        </DropdownMenuItem>
                                    ))}
                                    <DropdownMenuItem className="text-slate-400 cursor-not-allowed">
                                        <Plus className="h-3 w-3 mr-2" /> Add Repo...
                                    </DropdownMenuItem>
                                </DropdownMenuContent>
                            </DropdownMenu>
                        </div>
                    </div>

                    {/* Branch Indicator (Read Only for now) */}
                    <div className="bg-white/80 backdrop-blur border border-slate-200 p-4 rounded-lg shadow-sm flex items-center gap-3">
                        <div className="p-2 bg-emerald-50 rounded-md text-emerald-600">
                            <Terminal className="h-5 w-5" />
                        </div>
                        <div>
                            <p className="text-[10px] font-mono text-slate-500 uppercase tracking-wider">Current Branch</p>
                            <div className="flex items-center gap-2">
                                <span className="text-sm font-semibold text-slate-800">{context ? context.active_branch : "..."}</span>
                                <span className="px-1.5 py-0.5 rounded-full bg-emerald-100 text-emerald-700 text-[10px] font-medium">Synced</span>
                            </div>
                        </div>
                    </div>

                    {/* Environment (Read Only) */}
                    <div className="bg-white/80 backdrop-blur border border-slate-200 p-4 rounded-lg shadow-sm flex items-center gap-3">
                        <div className="p-2 bg-slate-100 rounded-md text-slate-600">
                            <Brain className="h-5 w-5" />
                        </div>
                        <div>
                            <p className="text-[10px] font-mono text-slate-500 uppercase tracking-wider">Workspace Path</p>
                            <p className="text-xs font-mono text-slate-700 truncate max-w-[200px]" title={context?.local_path || "Local"}>
                                {context?.local_path || "Local Development"}
                            </p>
                        </div>
                    </div>
                </div>

                {/* Chat Console */}
                <div className="bg-white border border-slate-200 rounded-lg shadow-xl shadow-slate-200/50 overflow-hidden flex flex-col h-[500px]">

                    {/* Message Stream */}
                    <div className="flex-1 p-6 space-y-4 overflow-y-auto">
                        {messages.map((m, i) => (
                            <motion.div
                                key={i}
                                initial={{ opacity: 0, y: 10 }}
                                animate={{ opacity: 1, y: 0 }}
                                className={`flex ${m.role === 'user' ? 'justify-end' : 'justify-start'}`}
                            >
                                <div className={`max-w-[80%] p-3 rounded text-sm ${m.role === 'user'
                                    ? 'bg-emerald-50 border border-emerald-200 text-emerald-900'
                                    : 'bg-slate-50 border border-slate-200 text-slate-700'
                                    }`}>
                                    {m.content}
                                </div>
                            </motion.div>
                        ))}
                    </div>

                    {/* Input Area */}
                    <div className="p-4 border-t border-slate-200 bg-slate-50/80 flex gap-2">
                        <input
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            onKeyDown={(e) => e.key === 'Enter' && handleSend()}
                            className="flex-1 bg-white border border-slate-300 rounded px-4 py-2 text-sm focus:outline-none focus:border-indigo-500 transition-colors text-slate-800 placeholder:text-slate-400"
                            placeholder="State your intent..."
                        />
                        <Button onClick={handleSend} size="sm" variant="ghost" className="text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50">
                            <MessageSquare className="h-4 w-4" />
                        </Button>
                    </div>
                </div>

            </main>

            {/* Footer - Copyright Only */}
            <footer className="absolute bottom-0 left-0 right-0 p-6 flex justify-center z-10 pointer-events-none">
                <span className="text-[10px] font-mono text-slate-400">CATALYST &copy; 2026 // OPEN SOURCE INTELLIGENCE</span>
            </footer>
        </div>
    );
}
