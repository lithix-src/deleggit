export interface Repo {
    id: string;
    org: string;
    name: string;
    default_branch: string;
}

export interface Context {
    active_repo_id: string;
    active_branch: string;
    org: string;
    name: string;
}

export interface Agent {
    id: string;
    service: string;
    role: string;
    config: any;
}

const API_BASE = "http://localhost:8080/api";

export const api = {
    async getRepos(): Promise<Repo[]> {
        const res = await fetch(`${API_BASE}/repos`);
        if (!res.ok) throw new Error("Failed to fetch repos");
        return res.json();
    },

    async getContext(): Promise<Context> {
        const res = await fetch(`${API_BASE}/context`);
        if (!res.ok) throw new Error("Failed to fetch context");
        return res.json();
    },

    async setContext(repoId: string, branch: string): Promise<void> {
        const res = await fetch(`${API_BASE}/context`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ repo_id: repoId, branch }),
        });
        if (!res.ok) {
            const err = await res.text();
            throw new Error(`Failed to set context: ${err}`);
        }
    },

    async getAgents(): Promise<Agent[]> {
        const res = await fetch(`${API_BASE}/agents`);
        if (!res.ok) throw new Error("Failed to fetch agents");
        return res.json();
    }
};
