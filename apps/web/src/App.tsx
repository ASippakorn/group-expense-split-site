import { FormEvent, useEffect, useState } from "react";
import { Navigate, Route, Routes, useNavigate } from "react-router-dom";
import { CircleDollarSign, LogOut, Plus, ReceiptText, UserPlus, UsersRound } from "lucide-react";
import * as api from "./api";
import type { Balance, Expense, Group, Participant, Settlement, SuggestedTransfer, Tag, User } from "./api";

export function App() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .getMe()
      .then(({ user }) => setUser(user))
      .catch(() => setUser(null))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="grid min-h-screen place-items-center bg-mist text-ink">Loading Splitr...</div>;
  }

  return (
    <Routes>
      <Route path="/login" element={<AuthPage mode="login" onUser={setUser} />} />
      <Route path="/register" element={<AuthPage mode="register" onUser={setUser} />} />
      <Route
        path="/"
        element={user ? <Dashboard user={user} onUser={setUser} /> : <Navigate to="/login" replace />}
      />
    </Routes>
  );
}

function AuthPage({ mode, onUser }: { mode: "login" | "register"; onUser: (user: User) => void }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const navigate = useNavigate();

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSubmitting(true);
    setError("");

    try {
      const result = mode === "login" ? await api.login(email, password) : await api.register(email, password);
      onUser(result.user);
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to sign in.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className="min-h-screen bg-mist px-4 py-10 text-ink">
      <section className="mx-auto grid min-h-[calc(100vh-5rem)] max-w-5xl items-center gap-8 md:grid-cols-[1.1fr_0.9fr]">
        <div>
          <p className="text-sm font-semibold uppercase tracking-wide text-leaf">Splitr</p>
          <h1 className="mt-4 max-w-xl text-4xl font-semibold leading-tight md:text-5xl">
            Keep group money clear from the first shared bill.
          </h1>
          <p className="mt-5 max-w-lg text-base leading-7 text-ink/70">
            Create groups, add expenses, and keep balances explainable with a ledger-first workflow.
          </p>
        </div>

        <form onSubmit={submit} className="rounded-lg bg-white p-6 shadow-panel" noValidate>
          <h2 className="text-2xl font-semibold">{mode === "login" ? "Sign in" : "Create account"}</h2>
          <div className="mt-6 space-y-4">
            <label className="block text-sm font-medium">
              Email
              <input
                className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                type="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                required
              />
            </label>
            <label className="block text-sm font-medium">
              Password
              <input
                className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                type="password"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
                minLength={8}
              />
            </label>
          </div>

          {error ? (
            <p className="mt-4 rounded-md border border-coral/30 bg-coral/10 px-3 py-2 text-sm text-coral" role="alert">
              {error}
            </p>
          ) : null}

          <button
            className="mt-6 inline-flex w-full items-center justify-center rounded-md bg-leaf px-4 py-2.5 font-semibold text-white outline-none hover:bg-leaf/90 focus:ring-2 focus:ring-leaf/30 disabled:opacity-60"
            disabled={submitting}
          >
            {submitting ? "Working..." : mode === "login" ? "Sign in" : "Create account"}
          </button>

          <button
            type="button"
            className="mt-4 text-sm font-medium text-leaf underline-offset-4 hover:underline"
            onClick={() => navigate(mode === "login" ? "/register" : "/login")}
          >
            {mode === "login" ? "Need an account?" : "Already have an account?"}
          </button>
        </form>
      </section>
    </main>
  );
}

function Dashboard({ user, onUser }: { user: User; onUser: (user: User | null) => void }) {
  const [groups, setGroups] = useState<Group[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [expenses, setExpenses] = useState<Expense[]>([]);
  const [balances, setBalances] = useState<Balance[]>([]);
  const [suggestedTransfers, setSuggestedTransfers] = useState<SuggestedTransfer[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [settlements, setSettlements] = useState<Settlement[]>([]);
  const [groupBalances, setGroupBalances] = useState<Record<string, Balance[]>>({});
  const [unavailableGroupBalances, setUnavailableGroupBalances] = useState<Set<string>>(new Set());
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [participantEmail, setParticipantEmail] = useState("");
  const [expenseDescription, setExpenseDescription] = useState("");
  const [expenseAmount, setExpenseAmount] = useState("");
  const [expenseDate, setExpenseDate] = useState("");
  const [payerParticipantId, setPayerParticipantId] = useState("");
  const [expenseParticipantIds, setExpenseParticipantIds] = useState<string[]>([]);
  const [splitType, setSplitType] = useState<"equal" | "manual_amount" | "percentage" | "tag">("equal");
  const [splitValues, setSplitValues] = useState<Record<string, string>>({});
  const [tagId, setTagId] = useState("");
  const [tagName, setTagName] = useState("");
  const [tagParticipantIds, setTagParticipantIds] = useState<string[]>([]);
  const [settlementPayerId, setSettlementPayerId] = useState("");
  const [settlementReceiverId, setSettlementReceiverId] = useState("");
  const [settlementAmount, setSettlementAmount] = useState("");
  const [settlementDate, setSettlementDate] = useState("");
  const [settlementNote, setSettlementNote] = useState("");
  const [settlementError, setSettlementError] = useState("");
  const [error, setError] = useState("");
  const [participantError, setParticipantError] = useState("");
  const [expenseError, setExpenseError] = useState("");
  const [balanceError, setBalanceError] = useState("");
  const [tagError, setTagError] = useState("");
  const [participantsLoading, setParticipantsLoading] = useState(false);
  const [expensesLoading, setExpensesLoading] = useState(false);
  const [balancesLoading, setBalancesLoading] = useState(false);

  useEffect(() => {
    api.listGroups().then(({ groups }) => {
      setGroups(groups);
      void Promise.all(
        groups.map(async (group) => {
          try {
            const { balances } = await api.listBalances(group.id);
            return { groupID: group.id, balances };
          } catch {
            return { groupID: group.id, balances: null };
          }
        }),
      ).then((results) => {
        const loadedBalances: Record<string, Balance[]> = {};
        for (const result of results) {
          if (result.balances) {
            loadedBalances[result.groupID] = result.balances;
          }
        }
        setGroupBalances(loadedBalances);
        setUnavailableGroupBalances(new Set(results.filter((result) => !result.balances).map((result) => result.groupID)));
      });
    });
  }, []);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    try {
      const { group } = await api.createGroup({ name, defaultCurrency: "THB", description });
      setGroups((current) => [group, ...current]);
      await openGroup(group);
      setName("");
      setDescription("");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Unable to create group.");
    }
  }

  async function signOut() {
    await api.logout();
    onUser(null);
  }

  async function openGroup(group: Group) {
    setSelectedGroup(group);
    setParticipantError("");
    setExpenseError("");
    setBalanceError("");
    setParticipantsLoading(true);
    setExpensesLoading(true);
    setBalancesLoading(true);
    try {
      const [{ participants }, { expenses }, { balances }, { suggestedTransfers }, { tags }, { settlements }] = await Promise.all([
        api.listParticipants(group.id),
        api.listExpenses(group.id),
        api.listBalances(group.id),
        api.listSuggestedTransfers(group.id),
        api.listTags(group.id),
        api.listSettlements(group.id),
      ]);
      setParticipants(participants);
      setExpenses(expenses);
      setBalances(balances);
      setSuggestedTransfers(suggestedTransfers);
      setTags(tags);
      setSettlements(settlements);
      setGroupBalances((current) => ({ ...current, [group.id]: balances }));
      setUnavailableGroupBalances((current) => {
        const next = new Set(current);
        next.delete(group.id);
        return next;
      });
      setPayerParticipantId(participants[0]?.id ?? "");
      setExpenseParticipantIds(participants.map((participant) => participant.id));
      setSplitValues({});
      setSplitType("equal");
      setTagId("");
      setTagParticipantIds(participants.map((participant) => participant.id));
      setSettlementPayerId(participants[1]?.id ?? "");
      setSettlementReceiverId(participants[0]?.id ?? "");
    } catch (err) {
      setParticipantError(err instanceof Error ? err.message : "Unable to load Group details.");
    } finally {
      setParticipantsLoading(false);
      setExpensesLoading(false);
      setBalancesLoading(false);
    }
  }

  async function submitTag(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedGroup || !tagName.trim() || tagParticipantIds.length === 0) {
      setTagError("Enter a Tag name and select at least one Participant.");
      return;
    }
    setTagError("");
    try {
      const { tag } = await api.createTag(selectedGroup.id, { name: tagName, participantIds: tagParticipantIds });
      setTags((current) => [...current, tag].sort((left, right) => left.name.localeCompare(right.name)));
      setTagName("");
    } catch (err) {
      setTagError(err instanceof Error ? err.message : "Unable to create Tag.");
    }
  }

  async function refreshBalances(groupID: string) {
    try {
      const [{ balances }, { suggestedTransfers }] = await Promise.all([
        api.listBalances(groupID),
        api.listSuggestedTransfers(groupID),
      ]);
      setBalances(balances);
      setSuggestedTransfers(suggestedTransfers);
      setBalanceError("");
      setGroupBalances((current) => ({ ...current, [groupID]: balances }));
      setUnavailableGroupBalances((current) => {
        const next = new Set(current);
        next.delete(groupID);
        return next;
      });
    } catch {
      setBalanceError("Unable to refresh Balances.");
      setUnavailableGroupBalances((current) => new Set(current).add(groupID));
    }
  }

  async function submitSettlement(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedGroup) return;
    const amountMinor = parseMoneyMinor(settlementAmount);
    if (amountMinor === null || amountMinor <= 0 || !settlementDate || !settlementPayerId || !settlementReceiverId || settlementPayerId === settlementReceiverId) {
      setSettlementError("Enter a positive amount, date, and two different Participants."); return;
    }
    try {
      const { settlement } = await api.createSettlement(selectedGroup.id, { payerParticipantId: settlementPayerId, receiverParticipantId: settlementReceiverId, amountMinor, currency: selectedGroup.defaultCurrency, settlementDate, note: settlementNote });
      setSettlements((current) => [settlement, ...current]); setSettlementAmount(""); setSettlementNote(""); setSettlementError(""); await refreshBalances(selectedGroup.id);
    } catch (err) { setSettlementError(err instanceof Error ? err.message : "Unable to record Settlement."); }
  }

  async function removeSettlement(settlementID: string) {
    if (!selectedGroup) return;
    try { await api.deleteSettlement(selectedGroup.id, settlementID); setSettlements((current) => current.filter((settlement) => settlement.id !== settlementID)); await refreshBalances(selectedGroup.id); }
    catch (err) { setSettlementError(err instanceof Error ? err.message : "Unable to delete Settlement."); }
  }

  async function submitParticipant(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedGroup) {
      return;
    }

    setParticipantError("");
    try {
      const { participant } = await api.addParticipant(selectedGroup.id, participantEmail);
      setParticipants((current) => [...current, participant]);
      setPayerParticipantId((current) => current || participant.id);
      setExpenseParticipantIds((current) => (current.includes(participant.id) ? current : [...current, participant.id]));
      setParticipantEmail("");
      await refreshBalances(selectedGroup.id);
    } catch (err) {
      setParticipantError(err instanceof Error ? err.message : "Unable to add Participant.");
    }
  }

  async function submitExpense(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedGroup) {
      return;
    }

    const amountMinor = parseMoneyMinor(expenseAmount);
    if (
      amountMinor === null ||
      amountMinor <= 0 ||
      !expenseDescription.trim() ||
      !expenseDate ||
      !payerParticipantId ||
      (splitType !== "tag" && (expenseParticipantIds.length === 0 || !expenseParticipantIds.includes(payerParticipantId))) ||
      (splitType === "tag" && (!tagId || !tags.find((tag) => tag.id === tagId)?.participants.some((participant) => participant.id === payerParticipantId)))
    ) {
      setExpenseError("Enter an Expense description, positive amount, date, Payer, and include the Payer in the Split or selected Tag.");
      return;
    }

    const splits = expenseParticipantIds.map((participantId) => ({
      participantId,
      amountMinor: splitType === "manual_amount" ? parseMoneyMinor(splitValues[participantId] ?? "") : undefined,
      percentage: splitType === "percentage" ? splitValues[participantId]?.trim() : undefined,
    }));
    if (splitType === "manual_amount" && (splits.some((split) => split.amountMinor === null) || splits.reduce((total, split) => total + (split.amountMinor ?? 0), 0) !== amountMinor)) {
      setExpenseError("Manual Split amounts must exactly equal the Expense amount.");
      return;
    }
    if (splitType === "percentage" && (!splits.every((split) => isValidPercentage(split.percentage ?? "")) || percentageBasisPoints(splits) !== 10000)) {
      setExpenseError("Percentage Splits must exactly total 100%.");
      return;
    }

    setExpenseError("");
    try {
      const { expense } = await api.createExpense(selectedGroup.id, {
        description: expenseDescription,
        amountMinor,
        currency: selectedGroup.defaultCurrency,
        expenseDate,
        payerParticipantId,
        participantIds: expenseParticipantIds,
        ...(splitType === "tag" ? { splitType, tagId } : splitType !== "equal" ? {
          splitType,
          splits: splits.map((split) => ({
            participantId: split.participantId,
            ...(split.amountMinor !== undefined ? { amountMinor: split.amountMinor ?? undefined } : {}),
            ...(split.percentage !== undefined ? { percentage: split.percentage } : {}),
          })),
        } : {}),
      });
      setExpenses((current) => [expense, ...current]);
      await refreshBalances(selectedGroup.id);
      setExpenseDescription("");
      setExpenseAmount("");
      setSplitValues({});
    } catch (err) {
      setExpenseError(err instanceof Error ? err.message : "Unable to add Expense.");
    }
  }

  function toggleExpenseParticipant(participantID: string) {
    setExpenseParticipantIds((current) =>
      current.includes(participantID)
        ? current.filter((selectedParticipantID) => selectedParticipantID !== participantID)
        : [...current, participantID],
    );
  }

  function changeSplitType(value: "equal" | "manual_amount" | "percentage" | "tag") {
    setSplitType(value);
    setSplitValues((current) => {
      if (Object.keys(current).length > 0 || value !== "percentage") return current;
      const next: Record<string, string> = {};
      expenseParticipantIds.forEach((participantID, index) => {
        const basisPoints = index === expenseParticipantIds.length - 1 ? 10000 - Math.floor(10000 / expenseParticipantIds.length) * index : Math.floor(10000 / expenseParticipantIds.length);
        next[participantID] = formatPercentageInput(basisPoints);
      });
      return next;
    });
  }

  return (
    <main className="min-h-screen bg-mist text-ink">
      <header className="border-b border-ink/10 bg-white">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4">
          <div className="flex items-center gap-3">
            <div className="grid h-10 w-10 place-items-center rounded-md bg-leaf text-white">
              <ReceiptText aria-hidden size={22} />
            </div>
            <div>
              <p className="text-xl font-semibold">Splitr</p>
              <p className="text-sm text-ink/60">{user.email}</p>
            </div>
          </div>
          <button
            className="inline-flex items-center gap-2 rounded-md border border-ink/15 px-3 py-2 text-sm font-medium hover:bg-mist focus:outline-none focus:ring-2 focus:ring-leaf/30"
            onClick={signOut}
          >
            <LogOut aria-hidden size={16} />
            Sign out
          </button>
        </div>
      </header>

      <section className="mx-auto grid max-w-6xl gap-6 px-4 py-8 lg:grid-cols-[360px_1fr]">
        <div className="space-y-6">
          <form onSubmit={submit} className="rounded-lg bg-white p-5 shadow-panel" noValidate>
            <div className="flex items-center gap-2">
              <Plus aria-hidden className="text-leaf" size={20} />
              <h1 className="text-xl font-semibold">Create group</h1>
            </div>
            <label className="mt-5 block text-sm font-medium">
              Name
              <input
                className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                value={name}
                onChange={(event) => setName(event.target.value)}
                required
              />
            </label>
            <label className="mt-4 block text-sm font-medium">
              Description
              <textarea
                className="mt-2 min-h-24 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                value={description}
                onChange={(event) => setDescription(event.target.value)}
              />
            </label>
            <p className="mt-3 text-sm text-ink/60">Default currency: THB</p>
            {error ? <p className="mt-4 text-sm text-coral">{error}</p> : null}
            <button className="mt-5 inline-flex w-full items-center justify-center rounded-md bg-leaf px-4 py-2.5 font-semibold text-white hover:bg-leaf/90 focus:outline-none focus:ring-2 focus:ring-leaf/30">
              Create group
            </button>
          </form>

          {selectedGroup ? (
            <>
              <section className="rounded-lg bg-white p-5 shadow-panel" aria-labelledby="participants-heading">
                <div className="flex items-center gap-2">
                  <UserPlus aria-hidden className="text-leaf" size={20} />
                  <h2 id="participants-heading" className="text-xl font-semibold">
                    Participants
                  </h2>
                </div>
                <p className="mt-1 text-sm text-ink/60">{selectedGroup.name}</p>

                {selectedGroup.ownerId === user.id ? (
                  <form onSubmit={submitParticipant} className="mt-5" noValidate>
                    <label className="block text-sm font-medium">
                      Participant email
                      <input
                        className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                        type="email"
                        value={participantEmail}
                        onChange={(event) => setParticipantEmail(event.target.value)}
                        required
                      />
                    </label>
                    <button className="mt-3 inline-flex w-full items-center justify-center rounded-md bg-leaf px-4 py-2.5 font-semibold text-white hover:bg-leaf/90 focus:outline-none focus:ring-2 focus:ring-leaf/30">
                      Add Participant
                    </button>
                  </form>
                ) : null}

                {participantError ? (
                  <p className="mt-4 rounded-md border border-coral/30 bg-coral/10 px-3 py-2 text-sm text-coral" role="alert">
                    {participantError}
                  </p>
                ) : null}

                <div className="mt-5 space-y-2">
                  {participantsLoading ? <p className="text-sm text-ink/60">Loading Participants...</p> : null}
                  {!participantsLoading && participants.length === 0 ? (
                    <p className="text-sm text-ink/60">No Participants yet.</p>
                  ) : null}
                  {participants.map((participant) => (
                    <div
                      key={participant.id}
                      className="flex items-center justify-between gap-3 rounded-md border border-ink/10 px-3 py-2 text-sm"
                    >
                      <span className="font-medium">{participant.user.email}</span>
                      <span className="rounded-md bg-mist px-2 py-1 capitalize text-ink/70">{participant.role}</span>
                    </div>
                  ))}
                </div>
              </section>

              <section className="rounded-lg bg-white p-5 shadow-panel" aria-labelledby="tags-heading">
                <h2 id="tags-heading" className="text-xl font-semibold">Tags</h2>
                <p className="mt-1 text-sm text-ink/60">Create a Group-level Tag for the Participants who share certain Expenses.</p>
                <form onSubmit={submitTag} className="mt-4 space-y-3" noValidate>
                  <label className="block text-sm font-medium">
                    Tag name
                    <input className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20" value={tagName} onChange={(event) => setTagName(event.target.value)} />
                  </label>
                  <fieldset className="space-y-2">
                    <legend className="text-sm font-medium">Tag Participants</legend>
                    {participants.map((participant) => (
                      <label key={participant.id} className="flex items-center gap-2 text-sm">
                        <input type="checkbox" checked={tagParticipantIds.includes(participant.id)} onChange={() => setTagParticipantIds((current) => current.includes(participant.id) ? current.filter((id) => id !== participant.id) : [...current, participant.id])} />
                        {participant.user.email}
                      </label>
                    ))}
                  </fieldset>
                  {tagError ? <p className="text-sm text-coral" role="alert">{tagError}</p> : null}
                  <button className="inline-flex rounded-md border border-ink/15 px-3 py-2 text-sm font-medium hover:bg-mist">Create Tag</button>
                </form>
                <div className="mt-4 space-y-2 text-sm">
                  {tags.length === 0 ? <p className="text-ink/60">No Tags yet.</p> : tags.map((tag) => <p key={tag.id}><span className="font-medium">{tag.name}</span>: {tag.participants.map((participant) => participant.user.email).join(", ")}</p>)}
                </div>
              </section>

              <section className="rounded-lg bg-white p-5 shadow-panel" aria-labelledby="expenses-heading">
                <div className="flex items-center gap-2">
                  <CircleDollarSign aria-hidden className="text-leaf" size={20} />
                  <h2 id="expenses-heading" className="text-xl font-semibold">
                    Expenses
                  </h2>
                </div>
                <p className="mt-1 text-sm text-ink/60">{selectedGroup.name}</p>

                {participants.length > 0 ? (
                  <form onSubmit={submitExpense} className="mt-5 space-y-4" noValidate>
                    <label className="block text-sm font-medium">
                      Expense description
                      <input
                        className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                        value={expenseDescription}
                        onChange={(event) => setExpenseDescription(event.target.value)}
                        required
                      />
                    </label>

                    <div className="grid gap-4 sm:grid-cols-2">
                      <label className="block text-sm font-medium">
                        Amount
                        <input
                          className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                          inputMode="decimal"
                          value={expenseAmount}
                          onChange={(event) => setExpenseAmount(event.target.value)}
                          required
                        />
                      </label>
                      <label className="block text-sm font-medium">
                        Date
                        <input
                          className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                          type="date"
                          value={expenseDate}
                          onChange={(event) => setExpenseDate(event.target.value)}
                          required
                        />
                      </label>
                    </div>

                    <label className="block text-sm font-medium">
                      Payer
                      <select
                        className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20"
                        value={payerParticipantId}
                        onChange={(event) => setPayerParticipantId(event.target.value)}
                        required
                      >
                        {participants.map((participant) => (
                          <option key={participant.id} value={participant.id}>
                            {participant.user.email}
                          </option>
                        ))}
                      </select>
                    </label>

                    <label className="block text-sm font-medium">
                      Split type
                      <select className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2 outline-none focus:border-leaf focus:ring-2 focus:ring-leaf/20" value={splitType} onChange={(event) => changeSplitType(event.target.value as "equal" | "manual_amount" | "percentage" | "tag")}>
                        <option value="equal">Equal</option>
                        <option value="manual_amount">Manual amount</option>
                        <option value="percentage">Percentage</option>
                        <option value="tag">Tag</option>
                      </select>
                    </label>

                    {splitType === "tag" ? (
                      <label className="block text-sm font-medium">
                        Tag
                        <select aria-label="Tag" className="mt-2 w-full rounded-md border border-ink/20 px-3 py-2" value={tagId} onChange={(event) => setTagId(event.target.value)}>
                          <option value="">Select a Tag</option>
                          {tags.map((tag) => <option key={tag.id} value={tag.id}>{tag.name}</option>)}
                        </select>
                      </label>
                    ) : (
                    <fieldset className="space-y-2">
                      <legend className="text-sm font-medium">{splitType === "equal" ? "Equal Split" : "Split Participants"}</legend>
                      {participants.map((participant) => (
                        <label key={participant.id} className="flex items-center gap-2 text-sm">
                          <input
                            className="h-4 w-4 rounded border-ink/30 text-leaf focus:ring-leaf/30"
                            type="checkbox"
                            checked={expenseParticipantIds.includes(participant.id)}
                            onChange={() => toggleExpenseParticipant(participant.id)}
                          />
                          <span>Split with {participant.user.email}</span>
                          {splitType === "manual_amount" ? (
                            <input aria-label={`Amount for ${participant.user.email}`} className="ml-auto w-28 rounded-md border border-ink/20 px-2 py-1" inputMode="decimal" value={splitValues[participant.id] ?? ""} onChange={(event) => setSplitValues((current) => ({ ...current, [participant.id]: event.target.value }))} />
                          ) : null}
                          {splitType === "percentage" ? (
                            <input aria-label={`Percentage for ${participant.user.email}`} className="ml-auto w-28 rounded-md border border-ink/20 px-2 py-1" inputMode="decimal" value={splitValues[participant.id] ?? ""} onChange={(event) => setSplitValues((current) => ({ ...current, [participant.id]: event.target.value }))} />
                          ) : null}
                        </label>
                      ))}
                    </fieldset>
                    )}

                    {expenseError ? (
                      <p className="rounded-md border border-coral/30 bg-coral/10 px-3 py-2 text-sm text-coral" role="alert">
                        {expenseError}
                      </p>
                    ) : null}

                    <button className="inline-flex w-full items-center justify-center rounded-md bg-leaf px-4 py-2.5 font-semibold text-white hover:bg-leaf/90 focus:outline-none focus:ring-2 focus:ring-leaf/30">
                      Add Expense
                    </button>
                  </form>
                ) : null}

                <div className="mt-5 space-y-3">
                  {expensesLoading ? <p className="text-sm text-ink/60">Loading Expenses...</p> : null}
                  {!expensesLoading && expenses.length === 0 ? <p className="text-sm text-ink/60">No Expenses yet.</p> : null}
                  {expenses.map((expense) => (
                    <article key={expense.id} className="rounded-md border border-ink/10 p-3 text-sm">
                      <div className="flex items-start justify-between gap-3">
                        <div>
                          <h3 className="font-semibold">{expense.description}</h3>
                          <p className="mt-1 text-ink/60">
                            {expense.payer.user.email} paid {formatMoney(expense.amountMinor, expense.currency)}
                          </p>
                          <p className="mt-1 text-ink/60">Split: {formatSplitType(expense.splitType)}</p>
                        </div>
                        <span className="rounded-md bg-mist px-2 py-1 text-ink/70">{expense.expenseDate}</span>
                      </div>
                      <ul className="mt-3 space-y-1 text-ink/70">
                        {expense.splits.map((split) => (
                          <li key={split.id}>
                            {split.participant.user.email}: {formatMoney(split.amountMinor, expense.currency)}{expense.splitType === "percentage" ? ` (${split.percentage}%)` : ""}
                          </li>
                        ))}
                      </ul>
                    </article>
                  ))}
                </div>
              </section>

              <section className="rounded-lg bg-white p-5 shadow-panel" aria-labelledby="settlements-heading">
                <h2 id="settlements-heading" className="text-xl font-semibold">Settlements</h2>
                <form onSubmit={submitSettlement} className="mt-4 space-y-3" noValidate>
                  <div className="grid gap-3 sm:grid-cols-2">
                    <label className="text-sm font-medium">Settlement payer<select aria-label="Settlement payer" className="mt-1 w-full rounded-md border border-ink/20 px-3 py-2" value={settlementPayerId} onChange={(event) => setSettlementPayerId(event.target.value)}>{participants.map((participant) => <option key={participant.id} value={participant.id}>{participant.user.email}</option>)}</select></label>
                    <label className="text-sm font-medium">Settlement receiver<select aria-label="Settlement receiver" className="mt-1 w-full rounded-md border border-ink/20 px-3 py-2" value={settlementReceiverId} onChange={(event) => setSettlementReceiverId(event.target.value)}>{participants.map((participant) => <option key={participant.id} value={participant.id}>{participant.user.email}</option>)}</select></label>
                    <label className="text-sm font-medium">Settlement amount<input aria-label="Settlement amount" className="mt-1 w-full rounded-md border border-ink/20 px-3 py-2" inputMode="decimal" value={settlementAmount} onChange={(event) => setSettlementAmount(event.target.value)} /></label>
                    <label className="text-sm font-medium">Settlement date<input aria-label="Settlement date" className="mt-1 w-full rounded-md border border-ink/20 px-3 py-2" type="date" value={settlementDate} onChange={(event) => setSettlementDate(event.target.value)} /></label>
                  </div>
                  <label className="block text-sm font-medium">Note<input aria-label="Settlement note" className="mt-1 w-full rounded-md border border-ink/20 px-3 py-2" value={settlementNote} onChange={(event) => setSettlementNote(event.target.value)} /></label>
                  {settlementError ? <p className="text-sm text-coral" role="alert">{settlementError}</p> : null}
                  <button className="rounded-md bg-leaf px-3 py-2 text-sm font-semibold text-white">Record Settlement</button>
                </form>
                <div className="mt-4 space-y-2">{settlements.length === 0 ? <p className="text-sm text-ink/60">No Settlements yet.</p> : settlements.map((settlement) => <article key={settlement.id} className="flex items-center justify-between gap-3 rounded-md border border-ink/10 p-3 text-sm"><p>{settlement.payer.user.email} repaid {settlement.receiver.user.email} {formatMoney(settlement.amountMinor, settlement.currency)}{settlement.note ? ` · ${settlement.note}` : ""}</p><button className="text-coral underline" onClick={() => void removeSettlement(settlement.id)}>Delete</button></article>)}</div>
              </section>

              <section className="rounded-lg bg-white p-5 shadow-panel" aria-labelledby="balances-heading">
                <div className="flex items-center gap-2">
                  <CircleDollarSign aria-hidden className="text-leaf" size={20} />
                  <h2 id="balances-heading" className="text-xl font-semibold">
                    Balances
                  </h2>
                </div>
                <p className="mt-1 text-sm text-ink/60">{selectedGroup.name}</p>
                <div className="mt-5 space-y-2">
                  {balancesLoading ? <p className="text-sm text-ink/60">Loading Balances...</p> : null}
                  {balanceError ? <p className="text-sm text-coral" role="alert">{balanceError}</p> : null}
                  {!balancesLoading && balances.length === 0 ? <p className="text-sm text-ink/60">Everyone is settled up.</p> : null}
                  {balances.map((balance) => (
                    <div
                      key={balance.participant.id}
                      className="flex items-center justify-between gap-3 rounded-md border border-ink/10 px-3 py-2 text-sm"
                    >
                      <p className="font-medium">
                        {balance.participant.user.email}: {formatBalance(balance.amountMinor, selectedGroup.defaultCurrency)}
                      </p>
                      <p className="text-ink/60">
                        Paid {formatMoney(balance.paidAmountMinor, selectedGroup.defaultCurrency)} · Owed {formatMoney(balance.owedAmountMinor, selectedGroup.defaultCurrency)}
                      </p>
                    </div>
                  ))}
                </div>
              </section>

              <section className="rounded-lg bg-white p-5 shadow-panel" aria-labelledby="suggested-transfers-heading">
                <h2 id="suggested-transfers-heading" className="text-xl font-semibold">Suggested Transfers</h2>
                <p className="mt-1 text-sm text-ink/60">A deterministic, minimum-transfer way to settle the current Balances.</p>
                <div className="mt-5 space-y-2">
                  {suggestedTransfers.length === 0 ? <p className="text-sm text-ink/60">Everyone is settled up.</p> : suggestedTransfers.map((transfer, index) => (
                    <p key={`${transfer.payerParticipantId}-${transfer.receiverParticipantId}-${index}`} className="rounded-md border border-ink/10 px-3 py-2 text-sm">
                      {transfer.payer.user.email} pays {transfer.receiver.user.email} {formatMoney(transfer.amountMinor, selectedGroup.defaultCurrency)}
                    </p>
                  ))}
                </div>
              </section>
            </>
          ) : null}
        </div>

        <div>
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-semibold">Your groups</h2>
            <span className="rounded-md bg-white px-3 py-1 text-sm text-ink/70">{groups.length} total</span>
          </div>

          <div className="mt-4 grid gap-3">
            {groups.length === 0 ? (
              <div className="rounded-lg border border-dashed border-ink/25 bg-white p-8 text-center">
                <UsersRound aria-hidden className="mx-auto text-gold" size={34} />
                <p className="mt-3 font-medium">No groups yet</p>
                <p className="mt-1 text-sm text-ink/60">Create your first shared expense space.</p>
              </div>
            ) : (
              groups.map((group) => (
                <article key={group.id} className="rounded-lg bg-white p-5 shadow-panel">
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <h3 className="text-lg font-semibold">{group.name}</h3>
                      <p className="mt-1 text-sm text-ink/60">{group.description || "No description"}</p>
                    </div>
                    <span className="rounded-md bg-mist px-2 py-1 text-sm font-medium">{group.defaultCurrency}</span>
                  </div>
                  <dl className="mt-5 grid gap-3 text-sm">
                    <Metric
                      label="Your balance"
                      value={
                        unavailableGroupBalances.has(group.id)
                          ? "Balance unavailable"
                          : formatBalance(
                              groupBalances[group.id]?.find((balance) => balance.participant.user.id === user.id)?.amountMinor ?? 0,
                              group.defaultCurrency,
                            )
                      }
                    />
                  </dl>
                  <button
                    className="mt-5 inline-flex items-center justify-center rounded-md border border-ink/15 px-3 py-2 text-sm font-medium hover:bg-mist focus:outline-none focus:ring-2 focus:ring-leaf/30"
                    onClick={() => void openGroup(group)}
                  >
                    Open {group.name}
                  </button>
                </article>
              ))
            )}
          </div>
        </div>
      </section>
    </main>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md bg-mist p-3">
      <dt className="text-ink/60">{label}</dt>
      <dd className="mt-1 font-semibold">{value}</dd>
    </div>
  );
}

function parseMoneyMinor(value: string) {
  const trimmed = value.trim();
  if (!/^\d+(\.\d{1,2})?$/.test(trimmed)) {
    return null;
  }

  const [whole, cents = ""] = trimmed.split(".");
  return Number(whole) * 100 + Number(cents.padEnd(2, "0"));
}

function isValidPercentage(value: string) {
  return /^\d+(\.\d{1,2})?$/.test(value.trim());
}

function percentageBasisPoints(splits: { percentage?: string }[]) {
  return splits.reduce((total, split) => {
    const [whole, fraction = ""] = (split.percentage ?? "").trim().split(".");
    return total + Number(whole) * 100 + Number(fraction.padEnd(2, "0"));
  }, 0);
}

function formatPercentageInput(basisPoints: number) {
  return (basisPoints / 100).toFixed(2).replace(/\.00$/, "");
}

function formatSplitType(splitType: Expense["splitType"]) {
  if (splitType === "manual_amount") return "Manual amount";
  if (splitType === "percentage") return "Percentage";
  return "Equal";
}

function formatMoney(amountMinor: number, currency: string) {
  const sign = amountMinor < 0 ? "-" : "";
  const absolute = Math.abs(amountMinor);
  const whole = Math.floor(absolute / 100);
  const cents = String(absolute % 100).padStart(2, "0");
  return `${sign}${currency} ${whole}.${cents}`;
}

function formatBalance(amountMinor: number, currency: string) {
  if (amountMinor > 0) {
    return `You are owed ${formatMoney(amountMinor, currency)}`;
  }
  if (amountMinor < 0) {
    return `You owe ${formatMoney(-amountMinor, currency)}`;
  }
  return "Settled up";
}
