import { useState } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { Box, Layers, PlayCircle, StopCircle, RefreshCw } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

interface Container {
    id: string;
    names: string;
    image: string;
    state: string; // "running", "exited"
    status: string;
}

export function ContainerGrid() {
    const [containers, setContainers] = useState<Container[]>([]);

    useEventSubscription("infra/docker/state", (ev: CloudEvent) => {
        if (ev.data && Array.isArray(ev.data)) {
            setContainers(ev.data);
        }
    });

    const getStatusColor = (state: string) => {
        if (state === "running") return "text-emerald-500 bg-emerald-950/30 border-emerald-800";
        if (state === "exited") return "text-slate-500 bg-slate-900 border-slate-700";
        return "text-amber-500 bg-amber-950/30 border-amber-800";
    };

    return (
        <Card className="bg-slate-900 border-slate-800 shadow-xl shadow-slate-950/50">
            <CardHeader className="pb-2 flex flex-row items-center justify-between">
                <CardTitle className="text-sm font-mono text-slate-400 flex items-center gap-2">
                    <Layers className="h-4 w-4 text-blue-500" />
                    CONTAINER RUNTIME
                </CardTitle>
                <div className="flex gap-2">
                    <a href="http://localhost:30030" target="_blank" rel="noopener noreferrer" className="text-[10px] text-slate-500 hover:text-blue-400 font-mono flex items-center gap-1 transition-colors">
                        GRAFANA ↗
                    </a>
                </div>
            </CardHeader>
            <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                    {containers.length === 0 && (
                        <div className="col-span-full text-center p-4 text-slate-600 italic border border-dashed border-slate-800 rounded">
                            <RefreshCw className="h-4 w-4 animate-spin mx-auto mb-2" />
                            Connecting to Runtime...
                        </div>
                    )}
                    {containers.map((c) => (
                        <div key={c.id} className="bg-slate-950/40 border border-slate-800 rounded p-3 flex flex-col gap-2 hover:border-slate-700 transition-colors group">
                            <div className="flex justify-between items-start">
                                <div className="flex items-center gap-2">
                                    <Box className="h-3.5 w-3.5 text-slate-500 group-hover:text-slate-300" />
                                    <span className="font-mono text-xs font-bold text-slate-300 truncate max-w-[120px]" title={c.names}>
                                        {c.names.replace("/", "")}
                                    </span>
                                </div>
                                <Badge variant="outline" className={cn("text-[10px] px-1.5 h-5", getStatusColor(c.state))}>
                                    {c.state.toUpperCase()}
                                </Badge>
                            </div>

                            <div className="flex items-center gap-2 text-[10px] text-slate-500 font-mono pl-0.5">
                                {c.state === "running" ? (
                                    <PlayCircle className="h-3 w-3 text-emerald-600" />
                                ) : (
                                    <StopCircle className="h-3 w-3 text-slate-700" />
                                )}
                                <span className="truncate" title={c.status}>{c.status}</span>
                            </div>

                            <div className="text-[10px] text-slate-600 truncate mt-1">
                                {c.image}
                            </div>
                        </div>
                    ))}
                </div>
            </CardContent>
        </Card>
    );
}
