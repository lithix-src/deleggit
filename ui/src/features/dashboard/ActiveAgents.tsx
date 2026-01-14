import { useState, useRef, useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Terminal, Search, Zap, Database, Activity, Filter } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuCheckboxItem,
    DropdownMenuContent,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils";

interface LogEntry {
    id: string;
    agent: string;
    message: string;
    time: string;
    variant: "scout" | "analyst" | "runner" | "infra" | "ui" | "sim" | "qa" | "default";
}

const AGENT_CONFIG: Record<string, { color: string; icon: any; badge: string }> = {
    // Standard Roles
    "Interface": { color: "text-pink-400", badge: "bg-pink-950/50 text-pink-400 border-pink-800 hover:bg-pink-900/50", icon: Zap },
    "Orchestrator": { color: "text-emerald-400", badge: "bg-emerald-950/50 text-emerald-400 border-emerald-800 hover:bg-emerald-900/50", icon: Database },
    "Infrastructure": { color: "text-blue-400", badge: "bg-blue-950/50 text-blue-400 border-blue-800 hover:bg-blue-900/50", icon: Terminal },
    "Compliance": { color: "text-cyan-400", badge: "bg-cyan-950/50 text-cyan-400 border-cyan-800 hover:bg-cyan-900/50", icon: Search },
    "Simulation": { color: "text-zinc-400", badge: "bg-zinc-900 text-zinc-400 border-zinc-700 hover:bg-zinc-800", icon: Activity },
    "default": { color: "text-slate-400", badge: "bg-slate-900 text-slate-400 border-slate-700 hover:bg-slate-800", icon: Terminal },
};

export function ActiveAgents() {
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const [filters, setFilters] = useState<string[]>(["Interface", "Orchestrator", "Infrastructure"]); // Default active filters
    const bottomRef = useRef<HTMLDivElement>(null);

    useEventSubscription("agent/+/log", (event: CloudEvent) => {
        const agentName = event.data.agent || "Unknown";

        const entry: LogEntry = {
            id: event.id || Math.random().toString(),
            agent: agentName,
            message: event.data.message,
            time: new Date(event.time).toLocaleTimeString(),
            variant: "default",
        };
        setLogs((prev) => {
            const newLogs = [...prev, entry];
            if (newLogs.length > 50) return newLogs.slice(-50);
            return newLogs;
        });
    });

    useEffect(() => {
        if (bottomRef.current) {
            bottomRef.current.scrollIntoView({ behavior: "smooth" });
        }
    }, [logs]);

    return (
        <Card className="bg-slate-900 border-slate-800 col-span-1 md:col-span-2 flex flex-col h-full min-h-[300px] shadow-xl shadow-slate-950/50">
            <CardHeader className="flex flex-row items-center justify-between pb-2 shrink-0">
                <CardTitle className="text-sm font-mono text-slate-400 flex items-center gap-2">
                    <Terminal className="h-4 w-4 text-emerald-500" />
                    AGENT SWARM ACTIVITY
                </CardTitle>
                <div className="flex items-center gap-2">
                    <Badge variant="outline" className="border-slate-800 text-slate-500 font-mono text-[10px]">
                        ACTIVE: {new Set(logs.map(l => l.agent)).size}
                    </Badge>
                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant="ghost" size="icon" className="h-6 w-6 text-slate-400 hover:text-white">
                                <Filter className="h-4 w-4" />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" className="bg-slate-900 border-slate-800 text-slate-200">
                            {Object.keys(AGENT_CONFIG).filter(k => k !== 'default').map(agent => (
                                <DropdownMenuCheckboxItem
                                    key={agent}
                                    checked={filters.includes(agent)}
                                    onCheckedChange={() => {
                                        setFilters(prev =>
                                            prev.includes(agent) ? prev.filter(f => f !== agent) : [...prev, agent]
                                        );
                                    }}
                                    className="text-xs"
                                >
                                    <span className={cn("mr-2", AGENT_CONFIG[agent].color)}>●</span> {agent}
                                </DropdownMenuCheckboxItem>
                            ))}
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            </CardHeader>
            <CardContent className="flex-1 min-h-0 p-0 relative">
                <ScrollArea className="h-[300px] w-full bg-slate-950/30 font-mono text-xs">
                    <div className="flex flex-col w-full">
                        {logs.length === 0 && <div className="text-slate-600 italic p-4">Waiting for agent signals...</div>}
                        {logs.filter(log => filters.includes(log.agent)).map((log) => {
                            const config = AGENT_CONFIG[log.agent] || AGENT_CONFIG["default"];
                            return (
                                <div key={log.id} className="flex gap-3 items-center hover:bg-slate-800/40 px-4 py-2 border-b border-slate-900/50 last:border-0 group transition-colors">
                                    <span className="text-slate-600 w-16 shrink-0 text-[10px]">{log.time}</span>

                                    <div className="w-24 shrink-0">
                                        <Badge variant="outline" className={cn("text-[10px] h-5 px-1.5 font-normal w-full justify-center", config.badge)}>
                                            {log.agent}
                                        </Badge>
                                    </div>

                                    <span className="text-slate-300 break-all group-hover:text-slate-100 transition-colors flex-1">
                                        {log.message}
                                    </span>
                                </div>
                            );
                        })}
                    </div>
                    <div ref={bottomRef} className="h-1" />
                </ScrollArea>
            </CardContent>
        </Card>
    );
}
