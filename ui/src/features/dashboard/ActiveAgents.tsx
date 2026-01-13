import { useState, useRef, useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Terminal, Search, Zap, Database } from "lucide-react";
import { cn } from "@/lib/utils";

interface LogEntry {
    id: string;
    agent: string;
    message: string;
    time: string;
    variant: "scout" | "analyst" | "runner" | "default";
}

const AGENT_VARIANTS: Record<string, { color: string; icon: any }> = {
    "TrendScout": { color: "text-purple-400", icon: Search },
    "GapAnalyst": { color: "text-blue-400", icon: Database },
    "CodeRunner": { color: "text-amber-400", icon: Zap },
    "default": { color: "text-zinc-400", icon: Terminal },
};

export function ActiveAgents() {
    const [logs, setLogs] = useState<LogEntry[]>([]);
    const bottomRef = useRef<HTMLDivElement>(null);

    useEventSubscription("agent/+/log", (event: CloudEvent) => {
        const agentName = event.data.agent || "Unknown";
        let variant: LogEntry["variant"] = "default";
        if (agentName === "TrendScout") variant = "scout";
        if (agentName === "GapAnalyst") variant = "analyst";
        if (agentName === "CodeRunner") variant = "runner";

        const entry: LogEntry = {
            id: event.id,
            agent: agentName,
            message: event.data.message,
            time: new Date(event.time).toLocaleTimeString(),
            variant,
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
        <Card className="bg-zinc-950 border-zinc-800 col-span-1 md:col-span-2 flex flex-col h-full min-h-[300px]">
            <CardHeader className="flex flex-row items-center justify-between pb-2 shrink-0">
                <CardTitle className="text-sm font-mono text-zinc-400 flex items-center gap-2">
                    <Terminal className="h-4 w-4 text-emerald-500" />
                    AGENT SWARM ACTIVITY
                </CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0 p-0 relative">
                <ScrollArea className="h-[300px] w-full bg-zinc-900/50 p-4 font-mono text-xs">
                    <div className="space-y-1">
                        {logs.length === 0 && <div className="text-zinc-600 italic">Waiting for agent signals...</div>}
                        {logs.map((log) => {
                            const style = AGENT_VARIANTS[log.agent] || AGENT_VARIANTS["default"];
                            const Icon = style.icon;
                            return (
                                <div key={log.id} className="flex gap-3 items-start hover:bg-zinc-900/80 p-1 rounded transition-colors group">
                                    <span className="text-zinc-700 w-16 mobile-hide shrink-0">[{log.time}]</span>
                                    <div className={cn("flex items-center gap-2 font-bold shrink-0 w-28", style.color)}>
                                        <Icon className="h-3 w-3" />
                                        {log.agent}
                                    </div>
                                    <span className="text-zinc-300 break-all group-hover:text-zinc-100 transition-colors">{log.message}</span>
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
