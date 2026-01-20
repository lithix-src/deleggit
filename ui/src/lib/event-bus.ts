import mqtt from "mqtt";
import { useEffect } from "react";

// CloudEvent Schema Configuration
export interface CloudEvent {
    id: string;
    source: string;
    type: string;
    data: any;
    time: string;
}

const BROKER_URL = "ws://localhost:30002"; // WebSockets port (NodePort)

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
            // Regex based matching for + and # wildcards
            const pattern = "^" + topic.replace(/\+/g, "[^/]+").replace(/#/g, ".*") + "$";
            const regex = new RegExp(pattern);

            if (regex.test(chkTopic)) {
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
