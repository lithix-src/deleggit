import { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { ScrollArea } from "@/components/ui/scroll-area";
import { GitPullRequest, AlertCircle, GitCommit, GitMerge } from "lucide-react";

interface RepoEvent {
    id: string;
    type: string;
    repo: string;
    title: string;
    author?: string;
    time: string;
    url?: string;
}

export function RepoFeed() {
    const [events, setEvents] = useState<RepoEvent[]>([]);

    useEventSubscription("repo/+/event", (ev: CloudEvent) => {
        if (!ev || !ev.data) return; // Safety check

        setEvents((prev) => [
            {
                id: ev.id || Math.random().toString(),
                type: ev.type || "repo.issue",
                repo: ev.data.repo || "unknown",
                title: ev.data.title || "Untitled Event",
                author: ev.data.author,
                time: ev.time ? new Date(ev.time).toLocaleTimeString() : new Date().toLocaleTimeString(),
                url: ev.data.url,
            },
            ...prev.slice(0, 50),
        ]);
    });

    const getIcon = (type: string) => {
        switch (type) {
            case "repo.push": return <GitCommit className="h-3 w-3 text-blue-500 mt-0.5" />;
            case "repo.pr": return <GitMerge className="h-3 w-3 text-purple-500 mt-0.5" />;
            default: return <AlertCircle className="h-3 w-3 text-amber-500 mt-0.5" />;
        }
    };

    const handleEventClick = (url?: string) => {
        if (url) {
            window.open(url, "_blank");
        }
    };

    return (
        <Card className="bg-slate-900 border-slate-800 shadow-xl shadow-slate-950/50 flex flex-col h-full min-h-[300px]">
            <CardHeader className="pb-2 shrink-0">
                <CardTitle className="text-sm font-mono text-slate-400 flex items-center gap-2">
                    <GitPullRequest className="h-4 w-4 text-purple-500" />
                    REPO STREAM
                </CardTitle>
            </CardHeader>
            <CardContent className="flex-1 min-h-0 p-0 relative">
                <ScrollArea className="h-[300px] w-full bg-slate-950/30 p-4 font-mono text-xs">
                    <div className="space-y-1">
                        {events.length === 0 && <div className="text-slate-600 italic">Waiting for changes...</div>}
                        {events.map((e) => (
                            <div
                                key={e.id}
                                onClick={() => handleEventClick(e.url)}
                                className="flex items-start gap-3 hover:bg-slate-800/80 p-1 rounded transition-colors cursor-pointer group"
                            >
                                {getIcon(e.type)}
                                <div className="flex-1">
                                    <div className="text-slate-200 group-hover:text-blue-400 transition-colors font-bold">
                                        {e.title}
                                    </div>
                                    <div className="text-slate-500 flex items-center gap-2 text-[10px] mt-0.5">
                                        <span>{e.repo}</span>
                                        {e.author && <span className="text-slate-600">• @{e.author}</span>}
                                        <span className="text-slate-700">• {e.time}</span>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                </ScrollArea>
            </CardContent>
        </Card>
    );
}
