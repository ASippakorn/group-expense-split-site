import { cleanup, fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import { MemoryRouter } from "react-router-dom";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { App } from "./App";
import * as api from "./api";

vi.mock("./api", () => ({
  getMe: vi.fn(),
  listGroups: vi.fn(),
  createGroup: vi.fn(),
  listParticipants: vi.fn(),
  listBalances: vi.fn(),
  listSuggestedTransfers: vi.fn(),
  listSettlements: vi.fn(),
  createSettlement: vi.fn(),
  deleteSettlement: vi.fn(),
  addParticipant: vi.fn(),
  listExpenses: vi.fn(),
  createExpense: vi.fn(),
  listTags: vi.fn(),
  createTag: vi.fn(),
  login: vi.fn(),
  logout: vi.fn(),
  register: vi.fn(),
}));

describe("App", () => {
  beforeEach(() => {
    vi.mocked(api.getMe).mockRejectedValue(new Error("not signed in"));
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [] });
	vi.mocked(api.listBalances).mockResolvedValue({ balances: [] });
    vi.mocked(api.listSuggestedTransfers).mockResolvedValue({ suggestedTransfers: [] });
    vi.mocked(api.listSettlements).mockResolvedValue({ settlements: [] });
    vi.mocked(api.listExpenses).mockResolvedValue({ expenses: [] });
    vi.mocked(api.listTags).mockResolvedValue({ tags: [] });
    vi.mocked(api.addParticipant).mockResolvedValue({
      participant: {
        id: "participant-2",
        user: { id: "user-2", email: "friend@example.com" },
        role: "participant",
        active: true,
      },
    });
    vi.mocked(api.createGroup).mockReset();
    vi.mocked(api.createExpense).mockReset();
    vi.mocked(api.createTag).mockReset();
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
    await waitFor(() => expect(screen.getAllByText("friend@example.com").length).toBeGreaterThan(0));
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

  it("creates an equal Split Expense from a Group detail view", async () => {
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
        {
          id: "participant-2",
          user: { id: "user-2", email: "friend@example.com" },
          role: "participant",
          active: true,
        },
      ],
    });
    vi.mocked(api.createExpense).mockResolvedValue({
      expense: {
        id: "expense-1",
        description: "Noodles",
        amountMinor: 10001,
        currency: "THB",
        expenseDate: "2026-08-14",
        splitType: "equal",
        payerParticipantId: "participant-1",
        payer: {
          id: "participant-1",
          user: { id: "owner-1", email: "owner@example.com" },
          role: "owner",
          active: true,
        },
        splits: [
          {
            id: "split-1",
            participantId: "participant-1",
            participant: {
              id: "participant-1",
              user: { id: "owner-1", email: "owner@example.com" },
              role: "owner",
              active: true,
            },
            amountMinor: 5001,
          },
          {
            id: "split-2",
            participantId: "participant-2",
            participant: {
              id: "participant-2",
              user: { id: "user-2", email: "friend@example.com" },
              role: "participant",
              active: true,
            },
            amountMinor: 5000,
          },
        ],
      },
    });

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    );

    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    fireEvent.change(await screen.findByLabelText("Expense description"), {
      target: { value: "Noodles" },
    });
    fireEvent.change(screen.getByLabelText("Amount"), {
      target: { value: "100.01" },
    });
    fireEvent.change(screen.getByLabelText("Date"), {
      target: { value: "2026-08-14" },
    });
    fireEvent.change(screen.getByLabelText("Payer"), {
      target: { value: "participant-1" },
    });
    fireEvent.click(screen.getByRole("button", { name: "Add Expense" }));

    await waitFor(() =>
      expect(api.createExpense).toHaveBeenCalledWith("group-1", {
        description: "Noodles",
        amountMinor: 10001,
        currency: "THB",
        expenseDate: "2026-08-14",
        payerParticipantId: "participant-1",
        participantIds: ["participant-1", "participant-2"],
      }),
    );
    expect(await screen.findByText("Noodles")).toBeInTheDocument();
    expect(screen.getByText("Split: Equal")).toBeInTheDocument();
    expect(screen.getByText("owner@example.com: THB 50.01")).toBeInTheDocument();
    expect(screen.getByText("friend@example.com: THB 50.00")).toBeInTheDocument();
  });

  it("renders Balance summaries on Group cards and in Group detail", async () => {
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
        { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
        { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
      ],
    });
    vi.mocked(api.listBalances).mockResolvedValue({
      balances: [
        {
          participant: { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
          paidAmountMinor: 10001,
          owedAmountMinor: 5001,
          amountMinor: 5000,
        },
        {
          participant: { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
          paidAmountMinor: 0,
          owedAmountMinor: 5000,
          amountMinor: -5000,
        },
      ],
    });

    render(
      <MemoryRouter>
        <App />
      </MemoryRouter>,
    );

    expect(await screen.findByText("You are owed THB 50.00")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: /open Bangkok Food Crawl/i }));

    expect(await screen.findByRole("heading", { name: "Balances" })).toBeInTheDocument();
    expect(screen.getByText("owner@example.com: You are owed THB 50.00")).toBeInTheDocument();
    expect(screen.getByText("friend@example.com: You owe THB 50.00")).toBeInTheDocument();
  });

  it("renders Suggested Transfers and the settled-up empty state in Group detail", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [{ id: "group-1", name: "Bangkok Food Crawl", defaultCurrency: "THB", description: "Dinner", ownerId: "owner-1" }] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [
      { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
      { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
    ] });
    vi.mocked(api.listSuggestedTransfers).mockResolvedValue({ suggestedTransfers: [{
      payerParticipantId: "participant-2", payer: { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
      receiverParticipantId: "participant-1", receiver: { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true }, amountMinor: 5000,
    }] });

    render(<MemoryRouter><App /></MemoryRouter>);
    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));

    expect(await screen.findByRole("heading", { name: "Suggested Transfers" })).toBeInTheDocument();
    expect(screen.getByText("friend@example.com pays owner@example.com THB 50.00")).toBeInTheDocument();
  });

  it("shows the Suggested Transfers empty state when everyone is settled up", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [{ id: "group-1", name: "Bangkok Food Crawl", defaultCurrency: "THB", description: "Dinner", ownerId: "owner-1" }] });

    render(<MemoryRouter><App /></MemoryRouter>);
    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));

    expect(await screen.findByRole("heading", { name: "Suggested Transfers" })).toBeInTheDocument();
    expect(within(screen.getByRole("heading", { name: "Suggested Transfers" }).closest("section")!).getByText("Everyone is settled up.")).toBeInTheDocument();
  });

  it("records and deletes a Settlement from the Group detail view", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [{ id: "group-1", name: "Bangkok Food Crawl", defaultCurrency: "THB", description: "Dinner", ownerId: "owner-1" }] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [{ id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true }, { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true }] });
    vi.mocked(api.createSettlement).mockResolvedValue({ settlement: { id: "settlement-1", payerParticipantId: "participant-2", payer: { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true }, receiverParticipantId: "participant-1", receiver: { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true }, amountMinor: 5000, currency: "THB", settlementDate: "2026-08-30", note: "Cash" } });
    vi.mocked(api.deleteSettlement).mockResolvedValue();
    render(<MemoryRouter><App /></MemoryRouter>);
    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    fireEvent.change(await screen.findByLabelText("Settlement amount"), { target: { value: "50" } }); fireEvent.change(screen.getByLabelText("Settlement date"), { target: { value: "2026-08-30" } }); fireEvent.change(screen.getByLabelText("Settlement note"), { target: { value: "Cash" } }); fireEvent.click(screen.getByRole("button", { name: "Record Settlement" }));
    await waitFor(() => expect(api.createSettlement).toHaveBeenCalledWith("group-1", expect.objectContaining({ amountMinor: 5000, note: "Cash" })));
    expect(await screen.findByText(/friend@example.com repaid owner@example.com THB 50.00/)).toBeInTheDocument(); fireEvent.click(screen.getByRole("button", { name: "Delete" })); await waitFor(() => expect(api.deleteSettlement).toHaveBeenCalledWith("group-1", "settlement-1"));
  });

  it("switches to manual amount Splits and shows client validation errors", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [{ id: "group-1", name: "Bangkok Food Crawl", defaultCurrency: "THB", description: "Dinner", ownerId: "owner-1" }] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [
      { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
      { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
    ] });

    render(<MemoryRouter><App /></MemoryRouter>);
    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    fireEvent.change(await screen.findByLabelText("Expense description"), { target: { value: "Dinner" } });
    fireEvent.change(screen.getByLabelText("Amount"), { target: { value: "100" } });
    fireEvent.change(screen.getByLabelText("Date"), { target: { value: "2026-08-14" } });
    fireEvent.change(screen.getByLabelText("Split type"), { target: { value: "manual_amount" } });
    fireEvent.change(screen.getByLabelText("Amount for owner@example.com"), { target: { value: "20" } });
    fireEvent.change(screen.getByLabelText("Amount for friend@example.com"), { target: { value: "70" } });
    fireEvent.click(screen.getByRole("button", { name: "Add Expense" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Manual Split amounts must exactly equal the Expense amount.");
    expect(api.createExpense).not.toHaveBeenCalled();
  });

  it("sends percentage Split values from the Group detail form", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [{ id: "group-1", name: "Bangkok Food Crawl", defaultCurrency: "THB", description: "Dinner", ownerId: "owner-1" }] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [
      { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
      { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
    ] });
    vi.mocked(api.createExpense).mockResolvedValue({ expense: {
      id: "expense-2", description: "Dinner", amountMinor: 10000, currency: "THB", expenseDate: "2026-08-14", splitType: "percentage", payerParticipantId: "participant-1",
      payer: { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
      splits: [],
    } });

    render(<MemoryRouter><App /></MemoryRouter>);
    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    fireEvent.change(await screen.findByLabelText("Expense description"), { target: { value: "Dinner" } });
    fireEvent.change(screen.getByLabelText("Amount"), { target: { value: "100" } });
    fireEvent.change(screen.getByLabelText("Date"), { target: { value: "2026-08-14" } });
    fireEvent.change(screen.getByLabelText("Split type"), { target: { value: "percentage" } });
    fireEvent.change(screen.getByLabelText("Percentage for owner@example.com"), { target: { value: "40" } });
    fireEvent.change(screen.getByLabelText("Percentage for friend@example.com"), { target: { value: "59" } });
    fireEvent.click(screen.getByRole("button", { name: "Add Expense" }));
    expect(await screen.findByRole("alert")).toHaveTextContent("Percentage Splits must exactly total 100%.");
    expect(api.createExpense).not.toHaveBeenCalled();
    fireEvent.change(screen.getByLabelText("Percentage for friend@example.com"), { target: { value: "60" } });
    fireEvent.click(screen.getByRole("button", { name: "Add Expense" }));

    await waitFor(() => expect(api.createExpense).toHaveBeenCalledWith("group-1", expect.objectContaining({
      splitType: "percentage",
      splits: [{ participantId: "participant-1", percentage: "40" }, { participantId: "participant-2", percentage: "60" }],
    })));
  });

  it("creates a Tag and uses it for a tagged Expense", async () => {
    vi.mocked(api.getMe).mockResolvedValue({ user: { id: "owner-1", email: "owner@example.com" } });
    vi.mocked(api.listGroups).mockResolvedValue({ groups: [{ id: "group-1", name: "Bangkok Food Crawl", defaultCurrency: "THB", description: "Dinner", ownerId: "owner-1" }] });
    vi.mocked(api.listParticipants).mockResolvedValue({ participants: [
      { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
      { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
    ] });
    vi.mocked(api.createTag).mockResolvedValue({ tag: { id: "tag-1", name: "Alcohol", participants: [
      { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true },
      { id: "participant-2", user: { id: "user-2", email: "friend@example.com" }, role: "participant", active: true },
    ] } });
    vi.mocked(api.createExpense).mockResolvedValue({ expense: {
      id: "expense-3", description: "Wine", amountMinor: 9000, currency: "THB", expenseDate: "2026-08-14", splitType: "tag", tagId: "tag-1", payerParticipantId: "participant-1",
      payer: { id: "participant-1", user: { id: "owner-1", email: "owner@example.com" }, role: "owner", active: true }, splits: [],
    } });

    render(<MemoryRouter><App /></MemoryRouter>);
    fireEvent.click(await screen.findByRole("button", { name: /open Bangkok Food Crawl/i }));
    fireEvent.change(await screen.findByLabelText("Tag name"), { target: { value: "Alcohol" } });
    fireEvent.click(screen.getByRole("button", { name: "Create Tag" }));
    await waitFor(() => expect(api.createTag).toHaveBeenCalledWith("group-1", { name: "Alcohol", participantIds: ["participant-1", "participant-2"] }));

    fireEvent.change(screen.getByLabelText("Expense description"), { target: { value: "Wine" } });
    fireEvent.change(screen.getByLabelText("Amount"), { target: { value: "90" } });
    fireEvent.change(screen.getByLabelText("Date"), { target: { value: "2026-08-14" } });
    fireEvent.change(screen.getByLabelText("Split type"), { target: { value: "tag" } });
    fireEvent.change(screen.getByLabelText("Tag"), { target: { value: "tag-1" } });
    fireEvent.click(screen.getByRole("button", { name: "Add Expense" }));

    await waitFor(() => expect(api.createExpense).toHaveBeenCalledWith("group-1", expect.objectContaining({ splitType: "tag", tagId: "tag-1" })));
  });
});
