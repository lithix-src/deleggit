# Catalyst Architecture & Standards

> **Status**: Draft (Prop. 1.0)
> **Focus**: Standardization, Extensibility, and Decoupling.

## 1. The "Catalyst Service Standard" (CSS)

To ensure the system scales ("Sensor Net Growth") without rebuilding the core, all services must adhere to the **CSS**.

### 1.1. The Data Protocol (CloudEvent)
All entities (Sensors, Agents, Services) MUST communicate via **MQTT** using the **CloudEvent JSON** format.

**Schema**:
```json
{
  "specversion": "1.0",
  "id": "uuid-v4",
  "source": "service-name",
  "type": "domain.entity.action",
  "time": "ISO-8601",
  "data": { ... }
}
```

### 1.2. The "Universal Sensor" Contract
To allow the UI to *dynamically* render new sensors without code changes, all telemetry MUST follow this `data` schema:

**Topic**: `sensor/<domain>/<metric>`  (e.g., `sensor/kitchen/temp`)
**Payload (`data`)**:
```json
{
  "value": 24.5,
  "unit": "°C",
  "label": "Kitchen Temperature",
  "meta": {
    "status": "nominal", // nominal, warning, error
    "trend": "up"        // up, down, flat
  }
}
```

**UI Behavior**:
The Frontend `GenericSensorGrid` subscribes to `sensor/#`.
- On message: Check if Card exists.
- If No: Create Card using `label` and `unit`.
- If Yes: Update Value and Chart.

### 1.3. Service "Manifest" (Configuration)
Services should not hardcode ports. Configuration is injected via usage of a central `config.env` or K8s ConfigMap.

| Variable | Standard | Default |
| :--- | :--- | :--- |
| `BROKER_URL` | MQTT Connection String | `tcp://localhost:1883` |
| `LOG_LEVEL` | Logging Granularity | `INFO` |
| `PORT` | Service HTTP Port | *(Service Specific)* |

---

## 2. Infrastructure Patterns

### 2.1. The "Sidecar" Pattern for Legacy Data
If a source cannot speak MQTT/CloudEvent (e.g., a USB Serial Device), write a small **Adapter Service** (Golang) that strictly translates:
`Raw Serial -> [Adapter] -> Standard CloudEvent -> MQTT`

### 2.2. The "Unified Agent" Pattern for Platform Data
Do not run sidecars for things the platform can see.
- **Example**: `docker-watcher`.
- **Role**: It is the *sole authority* on Container State. It publishes to `infra/docker/state` and `infra/docker/metrics`.

---

## 3. UI Architecture: "The Thin Glass"

The UI should be **Metadata Driven**.
- **Avoid**: Hardcoded components like `KitchenTempCard.tsx`.
- **Prefer**: `WidgetGrid` iterating over a `useSensorStore` Map.
- **Config**: `dashboard.json` defines layout preferences (order, size), but *availability* is driven by live data.
