import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { App } from "./App";
import * as api from "./api";

vi.mock("./api", () => ({
  getMe: vi.fn(),
  listGroups: vi.fn(),
  createGroup: vi.fn(),
  listParticipants: vi.fn(),
  addParticipant: vi.fn(),
  login: vi.fn(),
  logout: vi.fn(),
  register: vi.fn(),
}));

describe("App", () => {
  beforeEach(() => {
    vi.mocked(api.getMe).mockRejectedValue(new Error("not signed in"));
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [] });
    vi.mocked(api.addParticipant).mockResolvedValue({
      participant: {
        id: "participant-2",
        user: { id: "user-2", email: "friend@example.com" },
        role: "participant",
        active: true,
      },
    });
    vi.mocked(api.createGroup).mockReset();
    vi.mocked(api.login).mockReset();
    vi.mocked(api.logout).mockReset();
    vi.mocked(api.register).mockReset();
  });

  afterEach(() => {
    cleanup();
  });

  it("shows the sign in screen when no session exists", async () => {
    render(
      <MemoryRouter initialEntries={["/login"]}>
        <App />
      </MemoryRouter>,
    );

    expect(await screen.findByRole("heading", { name: "Sign in" })).toBeInTheDocument();
  });

  it("adds a Participant by email from a Group detail view", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({
      groups: [
        {
          id: "group-1",
          name: "Bangkok Food Crawl",
          defaultCurrency: "THB",
          description: "Dinner",
          ownerId: "owner-1",
        },
      ],
    });
    vi.mocked(api.listParticipants).mockResolvedValue({
      participants: [
        {
          id: "participant-1",
          user: { id: "owner-1", email: "owner@example.com" },
          role: "owner",
          active: true,
        },
      ],
    });

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    );

    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    expect(await screen.findByText("owner@example.com")).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText("Participant email"), {
      target: { value: "friend@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Add Participant" }));

    await waitFor(() => expect(api.addParticipant).toHaveBeenCalledWith("group-1", "friend@example.com"));
    expect(await screen.findByText("friend@example.com")).toBeInTheDocument();
  });

  it("shows Participant add errors", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({
      groups: [
        {
          id: "group-1",
          name: "Bangkok Food Crawl",
          defaultCurrency: "THB",
          description: "Dinner",
          ownerId: "owner-1",
        },
      ],
    });
    vi.mocked(api.addParticipant).mockRejectedValue(new Error("No registered User has that email."));

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    );

    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    fireEvent.change(await screen.findByLabelText("Participant email"), {
      target: { value: "missing@example.com" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Add Participant" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("No registered User has that email.");
  });
});
