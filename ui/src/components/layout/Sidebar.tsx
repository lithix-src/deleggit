import { LayoutDashboard, Server, Settings, GitGraph, FileCode } from "lucide-react";
import { cn } from "@/lib/utils";

interface SidebarProps extends React.HTMLAttributes<HTMLDivElement> {
    currentView: string;
    setView: (view: string) => void;
}

export function Sidebar({ className, currentView, setView }: SidebarProps) {
    const navItems = [
        { name: "Dashboard", id: "dashboard", icon: LayoutDashboard },
        { name: "Workflows", id: "workflows", icon: GitGraph },
        { name: "Hardware", id: "hardware", icon: Server },
        { name: "Agents", id: "agents", icon: FileCode },
        { name: "Settings", id: "settings", icon: Settings },
    ];

    return (
        <div className={cn("pb-12 w-64 border-r border-slate-800 h-[calc(100vh-3.5rem)] bg-slate-900 hidden md:block", className)}>
            <div className="space-y-4 py-4">
                <div className="px-3 py-2">
                    <h2 className="mb-2 px-4 text-xs font-semibold tracking-tight text-slate-500 font-mono uppercase">
                        Platform
                    </h2>
                    <div className="space-y-1">
                        {navItems.map((item) => (
                            <button
                                key={item.id}
                                onClick={() => setView(item.id)}
                                className={cn(
                                    "w-full justify-start flex items-center gap-3 px-4 py-2 text-sm font-medium rounded-md transition-colors",
                                    currentView === item.id
                                        ? "bg-slate-800 text-slate-100 font-bold"
                                        : "text-slate-400 hover:bg-slate-800/50 hover:text-slate-200"
                                )}
                            >
                                <item.icon className="h-4 w-4" />
                                {item.name}
                            </button>
                        ))}
                    </div>
                </div>

                <div className="px-3 py-2">
                    <div className="px-4 py-4 rounded border border-slate-800 bg-slate-950/50 mx-2">
                        <div className="flex items-center gap-2 text-xs font-mono text-emerald-500 mb-2">
                            <div className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                            STATUS: OPTIMAL
                        </div>
                        <div className="text-[10px] text-slate-500 font-mono">
                            Uptime: 04:22:19<br />
                            Swarm: 3 Agents
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
}
