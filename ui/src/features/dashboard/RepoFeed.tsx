import { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { GitPullRequest, AlertCircle } from "lucide-react";

interface RepoEvent {
    id: string;
    repo: string;
    title: string;
    time: string;
}

export function RepoFeed() {
    const [events, setEvents] = useState<RepoEvent[]>([]);

    useEventSubscription("repo/+/issue", (ev: CloudEvent) => {
        setEvents((prev) => [
            {
                id: ev.id,
                repo: ev.data.repo,
                title: ev.data.title,
                time: new Date(ev.time).toLocaleTimeString(),
            },
            ...prev.slice(0, 5),
        ]);
    });

    return (
        <Card className="bg-zinc-950 border-zinc-800">
            <CardHeader className="pb-2">
                <CardTitle className="text-sm font-mono text-zinc-400 flex items-center gap-2">
                    <GitPullRequest className="h-4 w-4 text-purple-500" />
                    REPO EVENTS
                </CardTitle>
            </CardHeader>
            <CardContent>
                <div className="space-y-4">
                    {events.length === 0 && <div className="text-sm text-zinc-600 italic px-2">No recent events</div>}
                    {events.map((e) => (
                        <div key={e.id} className="flex items-start gap-3 border-b border-zinc-800/50 pb-2 last:border-0 hover:bg-zinc-900/40 p-2 rounded transition-colors cursor-pointer">
                            <AlertCircle className="h-4 w-4 text-amber-500 mt-0.5" />
                            <div>
                                <div className="text-sm font-medium text-zinc-200">{e.title}</div>
                                <div className="text-xs text-zinc-500">{e.repo} • {e.time}</div>
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}
