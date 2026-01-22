import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Agent } from "@/lib/api";
import { Terminal, Shield, Cpu, Activity, Lock } from "lucide-react";

interface AgentDetailDialogProps {
    agent: Agent | null;
    isOpen: boolean;
    onClose: () => void;
    isAdmin: boolean;
}

export function AgentDetailDialog({ agent, isOpen, onClose, isAdmin }: AgentDetailDialogProps) {
    if (!agent) return null;

    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="bg-white sm:max-w-2xl">
                <DialogHeader>
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-indigo-50 rounded text-indigo-600">
                            <Terminal className="h-6 w-6" />
                        </div>
                        <div>
                            <DialogTitle className="text-xl font-bold flex items-center gap-2">
                                {agent.service}
                                <Badge variant="secondary" className="bg-emerald-100 text-emerald-700 hover:bg-emerald-100">
                                    Active
                                </Badge>
                            </DialogTitle>
                            <DialogDescription className="font-mono text-xs mt-1">
                                ID: {agent.id}
                            </DialogDescription>
                        </div>
                    </div>
                </DialogHeader>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-4">
                    {/* Status Column */}
                    <div className="space-y-4">
                        <div className="bg-slate-50 p-4 rounded-lg border border-slate-100">
                            <h4 className="text-xs font-bold text-slate-500 uppercase mb-3 flex items-center gap-2">
                                <Activity className="h-3 w-3" /> Runtime Status
                            </h4>
                            <div className="space-y-2 text-sm">
                                <div className="flex justify-between">
                                    <span className="text-slate-500">Role</span>
                                    <span className="font-medium font-mono">{agent.role}</span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-slate-500">Security Level</span>
                                    <span className="font-medium flex items-center gap-1 text-emerald-600">
                                        <Shield className="h-3 w-3" /> Standard
                                    </span>
                                </div>
                                <div className="flex justify-between">
                                    <span className="text-slate-500">Uptime</span>
                                    <span className="font-medium font-mono">14m 22s</span>
                                </div>
                            </div>
                        </div>

                        <div className="bg-slate-50 p-4 rounded-lg border border-slate-100">
                            <h4 className="text-xs font-bold text-slate-500 uppercase mb-3 flex items-center gap-2">
                                <Cpu className="h-3 w-3" /> Capabilities
                            </h4>
                            <div className="flex flex-wrap gap-2">
                                <Badge variant="outline" className="text-slate-500 bg-white">Code Analysis</Badge>
                                <Badge variant="outline" className="text-slate-500 bg-white">File I/O</Badge>
                                <Badge variant="outline" className="text-slate-500 bg-white">Command Exec</Badge>
                            </div>
                        </div>
                    </div>

                    {/* Config Column */}
                    <div className="flex flex-col h-full">
                        <h4 className="text-xs font-bold text-slate-500 uppercase mb-3 flex items-center justify-between">
                            <span>Configuration</span>
                            {isAdmin && (
                                <Badge className="bg-indigo-600 hover:bg-indigo-700 cursor-pointer text-[10px]">
                                    Edit Config
                                </Badge>
                            )}
                            {!isAdmin && (
                                <span className="flex items-center gap-1 text-[10px] text-slate-400">
                                    <Lock className="h-3 w-3" /> Read Only
                                </span>
                            )}
                        </h4>

                        <ScrollArea className="flex-1 w-full bg-slate-900 rounded-lg p-4 h-[200px] text-xs font-mono text-emerald-400 shadow-inner">
                            <pre>
                                {JSON.stringify(agent.config || {}, null, 2)}
                            </pre>
                        </ScrollArea>
                        <p className="text-[10px] text-slate-400 mt-2 text-right">
                            Configuration Source: Database
                        </p>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    );
}
