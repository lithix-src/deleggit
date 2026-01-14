import { create } from 'zustand';

// Universal Sensor Data Point
export interface SensorPoint {
    time: string;
    value: number;
}

// Universal Sensor "Card" State
export interface SensorState {
    id: string; // e.g. "sensor/cpu/temp"
    label: string; // "CPU Load"
    unit: string; // "%"
    currentValue: number;
    history: SensorPoint[];
}

interface SensorStore {
    sensors: Record<string, SensorState>;
    updateSensor: (id: string, value: number, unit?: string, label?: string) => void;
}

export const useSensorStore = create<SensorStore>((set) => ({
    sensors: {},
    updateSensor: (id, value, unit = "", label = "") =>
        set((state) => {
            const existing = state.sensors[id];
            const time = new Date().toLocaleTimeString();

            if (existing) {
                // Update existing sensor
                return {
                    sensors: {
                        ...state.sensors,
                        [id]: {
                            ...existing,
                            currentValue: value,
                            history: [...existing.history.slice(-19), { time, value }], // Keep last 20
                            // Update label/unit if provided and different (optional, but good for self-healing)
                            label: label || existing.label,
                            unit: unit || existing.unit,
                        },
                    },
                };
            } else {
                // Discovery: New Sensor!
                return {
                    sensors: {
                        ...state.sensors,
                        [id]: {
                            id,
                            label: label || id, // Fallback to ID if no label
                            unit,
                            currentValue: value,
                            history: [{ time, value }],
                        },
                    },
                };
            }
        }),
}));
