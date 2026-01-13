import { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { GitPullRequest, AlertCircle, GitCommit, GitMerge } from "lucide-react";

interface RepoEvent {
    id: string;
    type: string;
    repo: string;
    title: string;
    author?: string;
    time: string;
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
            },
            ...prev.slice(0, 5),
        ]);
    });

    const getIcon = (type: string) => {
        switch (type) {
            case "repo.push": return <GitCommit className="h-4 w-4 text-blue-500 mt-0.5" />;
            case "repo.pr": return <GitMerge className="h-4 w-4 text-purple-500 mt-0.5" />;
            default: return <AlertCircle className="h-4 w-4 text-amber-500 mt-0.5" />;
        }
    };

    return (
        <Card className="bg-slate-900 border-slate-800 shadow-xl shadow-slate-950/50">
            <CardHeader className="pb-2">
                <CardTitle className="text-sm font-mono text-slate-400 flex items-center gap-2">
                    <GitPullRequest className="h-4 w-4 text-purple-500" />
                    REPO STREAM
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {events.length === 0 && <div className="text-sm text-slate-600 italic px-2">Waiting for changes...</div>}
                    {events.map((e) => (
                        <div key={e.id} className="flex items-start gap-3 border-b border-slate-800/50 pb-2 last:border-0 hover:bg-slate-800/40 p-2 rounded transition-colors cursor-pointer group">
                            {getIcon(e.type)}
                            <div>
                                <div className="text-sm font-medium text-slate-200 group-hover:text-blue-400 transition-colors">
                                    {e.title}
                                </div>
                                <div className="text-xs text-slate-500 flex items-center gap-2">
                                    <span>{e.repo}</span>
                                    {e.author && <span className="text-slate-600">• @{e.author}</span>}
                                    <span className="text-slate-700">• {e.time}</span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}
