import { useState, useEffect } from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

interface SensorValue {
    time: string;
    value: number;
}

export function SensorGrid() {
    const [cpuTemps, setCpuTemps] = useState<SensorValue[]>([]);
    const [currentTemp, setCurrentTemp] = useState<number>(0);
    const [memUsage, setMemUsage] = useState<SensorValue[]>([]);
    const [currentMem, setCurrentMem] = useState<number>(0);

    useEventSubscription("sensor/cpu/temp", (event: CloudEvent) => {
        // Expected event.data = { value: 65.4 }
        const val = event.data.value;
        const time = new Date(event.time).toLocaleTimeString();

        setCurrentTemp(val);
        setCpuTemps((prev) => {
            const nw = [...prev, { time, value: val }];
            return nw.slice(-20); // Keep last 20 points
        });
    });

    useEventSubscription("sensor/memory/usage", (event: CloudEvent) => {
        const val = event.data.value;
        const time = new Date(event.time).toLocaleTimeString();

        setCurrentMem(val);
        setMemUsage((prev) => {
            const nw = [...prev, { time, value: val }];
            return nw.slice(-20);
        });
    });

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-4">
            <Card className="bg-slate-900 border-slate-800 text-slate-100 shadow-xl shadow-slate-950/50">
                <CardHeader>
                    <CardTitle className="text-sm font-mono text-slate-400">CPU LOAD</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-4xl font-bold mb-4 font-mono text-emerald-400">
                        {currentTemp.toFixed(1)}%
                    </div>
                    <div className="h-[100px] w-full">
                        <ResponsiveContainer width="100%" height="100%">
                            <LineChart data={cpuTemps}>
                                <Line
                                    type="monotone"
                                    dataKey="value"
                                    stroke="#34d399"
                                    strokeWidth={2}
                                    dot={false}
                                    isAnimationActive={false}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </div>
                </CardContent>
            </Card>

            <Card className="bg-slate-900 border-slate-800 text-slate-100 shadow-xl shadow-slate-950/50">
                <CardHeader>
                    <CardTitle className="text-sm font-mono text-slate-400">MEMORY USAGE</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="text-4xl font-bold mb-4 font-mono text-purple-400">
                        {currentMem.toFixed(1)}%
                    </div>
                    <div className="h-[100px] w-full">
                        <ResponsiveContainer width="100%" height="100%">
                            <LineChart data={memUsage}>
                                <Line
                                    type="monotone"
                                    dataKey="value"
                                    stroke="#c084fc"
                                    strokeWidth={2}
                                    dot={false}
                                    isAnimationActive={false}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </div>
                </CardContent>
            </Card>
        </div>
    );
}
