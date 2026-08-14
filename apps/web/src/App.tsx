import { FormEvent, useEffect, useState } from "react";
import { Navigate, Route, Routes, useNavigate } from "react-router-dom";
import { LogOut, Plus, ReceiptText, UserPlus, UsersRound } from "lucide-react";
import * as api from "./api";
import type { Group, Participant, User } from "./api";

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
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [participantEmail, setParticipantEmail] = useState("");
  const [error, setError] = useState("");
  const [participantError, setParticipantError] = useState("");
  const [participantsLoading, setParticipantsLoading] = useState(false);

  useEffect(() => {
    api.listGroups().then(({ groups }) => setGroups(groups));
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
    setParticipantsLoading(true);
    try {
      const { participants } = await api.listParticipants(group.id);
      setParticipants(participants);
    } catch (err) {
      setParticipantError(err instanceof Error ? err.message : "Unable to load Participants.");
    } finally {
      setParticipantsLoading(false);
    }
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
      setParticipantEmail("");
    } catch (err) {
      setParticipantError(err instanceof Error ? err.message : "Unable to add Participant.");
    }
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
                  <dl className="mt-5 grid grid-cols-2 gap-3 text-sm sm:grid-cols-3">
                    <Metric label="Net balance" value="THB 0.00" />
                    <Metric label="Expenses" value="0" />
                    <Metric label="Settlements" value="0" />
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
