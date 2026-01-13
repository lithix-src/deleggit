import mqtt from "mqtt";
import { useEffect, useState } from "react";

// CloudEvent Schema Configuration
export interface CloudEvent {
    id: string;
    source: string;
    type: string;
    data: any;
    time: string;
}

const BROKER_URL = "ws://localhost:9001"; // WebSockets port

export const client = mqtt.connect(BROKER_URL, {
    keepalive: 60,
    protocolId: "MQTT",
    protocolVersion: 4,
    clean: true,
    reconnectPeriod: 1000,
    connectTimeout: 30 * 1000,
});

client.on("connect", () => {
    console.log("Connected to MQTT Broker via WebSockets");
});

client.on("error", (err) => {
    console.error("Connection error: ", err);
});

// React Hook for easy subscription
export function useEventSubscription(topic: string, handler: (payload: CloudEvent) => void) {
    useEffect(() => {
        const handleMessage = (chkTopic: string, message: Buffer) => {
            // Basic topic matching (exact or simple wildcard logic if needed)
            // For now, accept if topic matches
            if (chkTopic === topic || (topic.endsWith("#") && chkTopic.startsWith(topic.slice(0, -1)))) {
                try {
                    const payload = JSON.parse(message.toString());
                    handler(payload);
                } catch (e) {
                    console.error("Failed to parse message", e);
                }
            }
        };

        client.subscribe(topic);
        client.on("message", handleMessage);

        return () => {
            client.unsubscribe(topic);
            client.removeListener("message", handleMessage);
        };
    }, [topic, handler]);
}
