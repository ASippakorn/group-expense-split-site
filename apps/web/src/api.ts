export type User = {
  id: string;
  email: string;
};

export type Group = {
  id: string;
  name: string;
  defaultCurrency: string;
  description: string;
  ownerId: string;
};

export type Participant = {
  id: string;
  user: User;
  role: "owner" | "participant";
  active: boolean;
};

type ApiError = {
  error: {
    code: string;
    message: string;
    fields?: Record<string, string>;
  };
};

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "/api/v1";

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...init?.headers,
    },
    ...init,
  });

  if (!response.ok) {
    let message = "Something went wrong.";
    try {
      const body = (await response.json()) as ApiError;
      message = body.error.message;
    } catch {
      message = response.statusText;
    }
    throw new Error(message);
  }

  if (response.status === 204) {
    return undefined as T;
  }

  return (await response.json()) as T;
}

export function register(email: string, password: string) {
  return request<{ user: User }>("/auth/register", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export function login(email: string, password: string) {
  return request<{ user: User }>("/auth/login", {
    method: "POST",
    body: JSON.stringify({ email, password }),
  });
}

export function logout() {
  return request<void>("/auth/logout", { method: "POST" });
}

export function getMe() {
  return request<{ user: User }>("/me");
}

export function listGroups() {
  return request<{ groups: Group[] }>("/groups");
}

export function createGroup(input: { name: string; defaultCurrency: string; description: string }) {
  return request<{ group: Group }>("/groups", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function listParticipants(groupId: string) {
  return request<{ participants: Participant[] }>(`/groups/${groupId}/participants`);
}

export function addParticipant(groupId: string, email: string) {
  return request<{ participant: Participant }>(`/groups/${groupId}/participants`, {
    method: "POST",
    body: JSON.stringify({ email }),
  });
}
