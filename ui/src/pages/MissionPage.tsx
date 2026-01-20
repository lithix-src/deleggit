import { motion } from "framer-motion";
import { Brain, MessageSquare, Terminal, Activity } from "lucide-react";
import { useState } from "react";
import { Button } from "../components/ui/button";

export function MissionPage() {
    const [messages, setMessages] = useState([
        { role: "system", content: "Catalyst Vibe Engine Online. How can I assist?" }
    ]);
    const [input, setInput] = useState("");

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

            {/* Header */}
            <header className="absolute top-0 left-0 right-0 p-6 flex justify-between items-center z-10">
                <div className="flex items-center gap-3">
                    <Brain className="text-emerald-500 h-6 w-6" />
                    <h1 className="text-xl font-bold tracking-widest text-slate-900">CATALYST <span className="text-indigo-600">CONTROL</span></h1>
                </div>
                <div className="flex gap-4 text-xs font-mono text-slate-500">
                    <span className="flex items-center gap-2"><Activity className="h-3 w-3" /> CORE: ONLINE</span>
                    <span className="flex items-center gap-2"><Terminal className="h-3 w-3" /> SWARM: IDLE</span>
                </div>
            </header>

            {/* Main Interface */}
            <main className="w-full max-w-3xl z-10 flex flex-col gap-8">

                {/* Vibe Visualizer (Placeholder) */}
                <div className="h-32 flex items-center justify-center">
                    <div className="relative">
                        <div className="absolute -inset-4 bg-emerald-500/20 blur-xl rounded-full animate-pulse" />
                        <Brain className="h-16 w-16 text-emerald-400 relative z-10 opacity-80" />
                    </div>
                </div>

                {/* Chat Console */}
                <div className="bg-white border border-slate-300 rounded-lg shadow-2xl shadow-slate-200 overflow-hidden flex flex-col h-[500px]">

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

            {/* Footer Navigation */}
            <footer className="absolute bottom-0 left-0 right-0 p-6 flex justify-center gap-8 z-10">
                <a href="/admin" className="text-xs font-mono text-slate-600 hover:text-slate-400 transition-colors">ACCESS ADMIN CONSOLE</a>
            </footer>
        </div>
    );
}
