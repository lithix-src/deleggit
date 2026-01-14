import { useSensorStore } from "@/lib/store";
import { useEventSubscription, CloudEvent } from "@/lib/event-bus";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { ResponsiveContainer, LineChart, Line } from "recharts";
import { Activity } from "lucide-react";

export function GenericSensorGrid() {
    const { sensors, updateSensor } = useSensorStore();

    // Wildcard Subscription to ALL sensors
    useEventSubscription("sensor/+/+", (event: CloudEvent) => {
        // Topic breakdown: sensor/<domain>/<metric>
        // event.data expected: { value: N, unit: "X", label: "Y" }

        // Safety Clean
        if (!event.data || typeof event.data.value !== 'number') return;

        // Construct a unique ID from the CloudEvent topic (if available) or fallback
        // The event bus abstraction might not pass the exact topic string easily depending on implementation.
        // Assuming your event-bus callback logic allows topic extraction, or we map explicitly.
        // For now, let's assume we map "sensor/cpu/temp" manually if needed, OR 
        // we update the event bus to pass the topics.

        // HACK: For now, we infer ID from the event type or we rely on the Backend Standard
        // Backend Standard says type should be "sensor.cpu.usage"

        // Let's use the Event Type as the unique key for now, or the topic if we can get it.
        // Since useEventSubscription is a wrapper, let's look at `event.type`.
        const id = event.type;

        // Extract CSS (Catalyst Service Standard) fields
        const { value, unit, label } = event.data;

        updateSensor(id, value, unit, label);
    });

    const sensorList = Object.values(sensors);

    return (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 p-4 mb-8 border-t border-dashed border-slate-800 pt-8 relative">
            <div className="absolute top-0 left-0 -translate-y-1/2 bg-slate-950 px-2 text-xs font-mono text-slate-500 flex items-center gap-2 border border-slate-800 rounded-full ml-4">
                <Activity className="h-3 w-3 text-emerald-500" />
                DYNAMIC SENSOR GRID (v2)
            </div>

            {sensorList.length === 0 && (
                <div className="col-span-full text-center p-8 text-slate-600 italic font-mono text-sm">
                    No CSS-Compliant Sensors Discovered.
                </div>
            )}

            {sensorList.map((sensor) => (
                <Card key={sensor.id} className="bg-slate-900 border-slate-800 text-slate-100 shadow-xl shadow-slate-950/50">
                    <CardHeader className="pb-2">
                        <CardTitle className="text-sm font-mono text-slate-400 uppercase tracking-wider flex justify-between">
                            {sensor.label}
                            <span className="text-[10px] text-slate-600 normal-case">{sensor.id}</span>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-4xl font-bold mb-4 font-mono text-slate-200">
                            {sensor.currentValue.toFixed(1)}
                            <span className="text-lg text-slate-500 ml-1">{sensor.unit}</span>
                        </div>
                        <div className="h-[100px] w-full">
                            <ResponsiveContainer width="100%" height="100%">
                                <LineChart data={sensor.history}>
                                    <Line
                                        type="monotone"
                                        dataKey="value"
                                        stroke="#94a3b8" // Slate-400
                                        strokeWidth={2}
                                        dot={false}
                                        isAnimationActive={false}
                                    />
                                </LineChart>
                            </ResponsiveContainer>
                        </div>
                    </CardContent>
                </Card>
            ))}
        </div>
    );
}
