import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  Bot,
  Check,
  CircleDollarSign,
  Clipboard,
  Download,
  ExternalLink,
  FileJson,
  Github,
  Inbox,
  Link,
  MailPlus,
  Paperclip,
  Printer,
  RefreshCw,
  Search,
  Server,
  ShieldAlert,
  ShieldCheck,
  Trash2,
  Upload
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import {
  createApiClient,
  errorMessage,
  uploadEmlFile,
  type AssistTone,
  type Message
} from "@/api/client";
import { Toast, type ToastState } from "@/components/Toast";
import { demoMessages, fallbackDraft, filterMessages } from "@/features/mail/demo";
import {
  createSnapshot,
  decodeSnapshotHash,
  defaultUIState,
  encodeSnapshotHash,
  fileLooksLikeEML,
  fileResultId,
  looksLikeEMLText,
  maxShareUrlLength,
  messageListToCSV,
  messageListToJSON,
  parseSnapshotText,
  parseStoredUIState,
  serializeUIState,
  snapshotToJSON,
  stateFileName,
  summarizeBatch,
  textToEmlFile,
  type FileImportResult,
  type UIState
} from "@/features/mail/workbench";
import { copyText, downloadTextFile, readClipboardText } from "@/lib/browser";
import { buildInfo } from "@/lib/build";
import { fetchLatestRepoCommit } from "@/lib/github";
import { formatMailDate, relativeAge } from "@/lib/time";

const repoUrl = "https://github.com/baditaflorin/localhuman-mail";
const paypalUrl = "https://www.paypal.com/paypalme/florinbadita";
const defaultBackendUrl = import.meta.env.VITE_LOCALHUMAN_API_URL ?? "http://localhost:8080";
const uiStateKey = "localhuman-mail.ui-state";
const legacyBackendUrlKey = "localhuman-mail.backend-url";
const toneOptions: AssistTone[] = ["concise", "warm", "decisive"];

export function App() {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const stateInputRef = useRef<HTMLInputElement | null>(null);
  const [uiState, setUIState] = useState<UIState>(() => initialUIState());
  const [toast, setToast] = useState<ToastState>(null);
  const [importResults, setImportResults] = useState<FileImportResult[]>([]);
  const [isImporting, setIsImporting] = useState(false);
  const [dragActive, setDragActive] = useState(false);

  const api = useMemo(() => createApiClient(uiState.backendUrl), [uiState.backendUrl]);

  useEffect(() => {
    window.localStorage.setItem(uiStateKey, serializeUIState(uiState));
    window.localStorage.setItem(legacyBackendUrlKey, uiState.backendUrl);
  }, [uiState]);

  useEffect(() => {
    const snapshot = decodeSnapshotHash(window.location.hash);
    if (!snapshot) {
      return;
    }
    setUIState({ ...snapshot.ui, snapshotMessages: snapshot.messages });
    setToast({ kind: "success", message: "Restored shared localhuman-mail state" });
    window.history.replaceState(null, "", window.location.pathname + window.location.search);
  }, []);

  const patchUIState = (patch: Partial<UIState>) => {
    setUIState((current) => ({ ...current, ...patch }));
  };

  const versionQuery = useQuery({
    queryKey: ["version", uiState.backendUrl],
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/version");
      if (error || !data) throw new Error(errorMessage(error, "Backend version unavailable"));
      return data;
    },
    retry: false
  });

  const capabilitiesQuery = useQuery({
    queryKey: ["capabilities", uiState.backendUrl],
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/capabilities");
      if (error || !data) throw new Error(errorMessage(error, "Capabilities unavailable"));
      return data.capabilities;
    },
    enabled: versionQuery.isSuccess,
    retry: false
  });

  const latestCommitQuery = useQuery({
    queryKey: ["github-latest-commit"],
    queryFn: fetchLatestRepoCommit,
    retry: false,
    staleTime: 60_000
  });

  const messagesQuery = useQuery({
    queryKey: ["messages", uiState.backendUrl],
    queryFn: async () => {
      const { data, error } = await api.GET("/api/v1/messages", {
        params: { query: { limit: 100 } }
      });
      if (error || !data) throw new Error(errorMessage(error, "Messages unavailable"));
      return data.messages;
    },
    enabled: versionQuery.isSuccess,
    retry: false
  });

  const backendOnline = versionQuery.isSuccess;
  const backendMessages = messagesQuery.data ?? [];
  const sourceMessages = backendMessages.length
    ? backendMessages
    : uiState.snapshotMessages.length
      ? uiState.snapshotMessages
      : demoMessages;
  const visibleMessages = filterMessages(sourceMessages, uiState.query);
  const currentMessages = visibleMessages.length ? visibleMessages : sourceMessages;
  const selectedMessage =
    visibleMessages.find((message) => message.id === uiState.selectedId) ??
    visibleMessages[0] ??
    sourceMessages[0];
  const hasSnapshot = uiState.snapshotMessages.length > 0 && backendMessages.length === 0;
  const batchSummary = summarizeBatch(importResults);
  const debugEnabled = new URLSearchParams(window.location.search).get("debug") === "1";

  const importDemoMutation = useMutation({
    mutationFn: async () => {
      const { data, error } = await api.POST("/api/v1/import/demo");
      if (error || !data) throw new Error(errorMessage(error, "Demo import failed"));
      return data;
    },
    onSuccess: async (result) => {
      setToast({ kind: "success", message: `Imported ${result.imported} backend demo messages` });
      await queryClient.invalidateQueries({ queryKey: ["messages", uiState.backendUrl] });
    },
    onError: (error) => {
      setToast({
        kind: "error",
        message: errorMessage(error, "Start the backend to import demo mail")
      });
    }
  });

  const assistMutation = useMutation({
    mutationFn: async (message: Message) => {
      if (!backendOnline) {
        return {
          draft: fallbackDraft(message, uiState.tone),
          model: "browser-demo",
          source: "fallback"
        };
      }
      const { data, error } = await api.POST("/api/v1/assist/reply", {
        body: { messageId: message.id, tone: uiState.tone }
      });
      if (error || !data) throw new Error(errorMessage(error, "Reply assist failed"));
      return data;
    },
    onSuccess: (result) => {
      patchUIState({ draft: result.draft });
      setToast({ kind: "success", message: `Draft generated by ${result.model}` });
    },
    onError: (error) => {
      setToast({ kind: "error", message: errorMessage(error, "Reply assist failed") });
    }
  });

  const importFiles = async (files: File[]) => {
    if (!backendOnline) {
      setToast({ kind: "error", message: "Start the backend before importing mail." });
      return;
    }
    if (files.length === 0) {
      return;
    }
    const initialResults: FileImportResult[] = files.map((file, index) => ({
      id: fileResultId(file, index),
      name: file.name,
      size: file.size,
      status: "pending",
      imported: 0,
      skipped: 0
    }));
    setImportResults(initialResults);
    setIsImporting(true);
    try {
      for (const [index, file] of files.entries()) {
        const id = initialResults[index].id;
        try {
          if (!(await fileLooksLikeEML(file))) {
            updateImportResult(id, {
              status: "error",
              error: "This file does not look like an exported .eml message.",
              imported: 0,
              skipped: 0
            });
            continue;
          }
          const result = await uploadEmlFile(uiState.backendUrl, file);
          updateImportResult(id, {
            status: result.imported > 0 ? "imported" : "skipped",
            imported: result.imported,
            skipped: result.skipped
          });
        } catch (error) {
          updateImportResult(id, {
            status: "error",
            error: errorMessage(error, "Could not import this file"),
            imported: 0,
            skipped: 0
          });
        }
      }
      await queryClient.invalidateQueries({ queryKey: ["messages", uiState.backendUrl] });
      setToast({ kind: "success", message: "Batch import finished" });
    } finally {
      setIsImporting(false);
    }
  };

  const updateImportResult = (id: string, patch: Partial<FileImportResult>) => {
    setImportResults((current) =>
      current.map((result) => (result.id === id ? { ...result, ...patch } : result))
    );
  };

  const importPasteText = async () => {
    const text = uiState.pasteText.trim();
    if (!looksLikeEMLText(text)) {
      setToast({
        kind: "error",
        message: "Paste raw .eml text with headers such as From, To, Subject, or Date."
      });
      return;
    }
    await importFiles([textToEmlFile(text)]);
    patchUIState({ pasteText: "" });
  };

  const importClipboardText = async () => {
    try {
      const text = await readClipboardText();
      patchUIState({ pasteText: text });
      if (looksLikeEMLText(text)) {
        await importFiles([textToEmlFile(text)]);
        patchUIState({ pasteText: "" });
      } else {
        setToast({ kind: "error", message: "Clipboard text does not look like raw .eml content." });
      }
    } catch (error) {
      setToast({ kind: "error", message: errorMessage(error, "Paste into the text box instead.") });
    }
  };

  const copyValue = async (value: string, label: string) => {
    try {
      await copyText(value);
      setToast({ kind: "success", message: `${label} copied` });
    } catch (error) {
      setToast({ kind: "error", message: errorMessage(error, `Could not copy ${label}`) });
    }
  };

  const exportSnapshot = () => {
    const snapshot = createSnapshot(uiState, currentMessages, buildInfo);
    downloadTextFile(stateFileName, snapshotToJSON(snapshot), "application/json");
    setToast({ kind: "success", message: "State file downloaded" });
  };

  const importSnapshotFile = async (file: File) => {
    try {
      const snapshot = parseSnapshotText(await file.text());
      setUIState({ ...snapshot.ui, snapshotMessages: snapshot.messages });
      setToast({ kind: "success", message: `Restored ${snapshot.messages.length} messages` });
    } catch (error) {
      setToast({ kind: "error", message: errorMessage(error, "State file could not be imported") });
    }
  };

  const shareSnapshot = async () => {
    const snapshot = createSnapshot(uiState, sourceMessages, buildInfo);
    const url = new URL(window.location.href);
    url.hash = encodeSnapshotHash(snapshot);
    if (url.href.length > maxShareUrlLength) {
      setToast({
        kind: "error",
        message: "This state is too large for a share URL. Download a state file instead."
      });
      return;
    }
    window.history.replaceState(null, "", url);
    await copyValue(url.href, "Share URL");
  };

  const clearLocalState = () => {
    setUIState(defaultUIState(defaultBackendUrl));
    setImportResults([]);
    setToast({ kind: "success", message: "Local UI state cleared" });
  };

  return (
    <main
      className={`min-h-screen bg-paper text-ink ${dragActive ? "ring-4 ring-fern" : ""}`}
      onDragOver={(event) => {
        event.preventDefault();
        setDragActive(true);
      }}
      onDragLeave={() => setDragActive(false)}
      onDrop={(event) => {
        event.preventDefault();
        setDragActive(false);
        void importFiles(Array.from(event.dataTransfer.files));
      }}
    >
      <div className="mx-auto grid min-h-screen max-w-[1500px] grid-cols-1 lg:grid-cols-[280px_minmax(420px,520px)_1fr]">
        <aside className="border-b border-line bg-white px-5 py-5 lg:border-b-0 lg:border-r">
          <div className="flex items-center gap-3">
            <span className="grid h-10 w-10 place-items-center rounded-lg bg-ink text-white">
              <Inbox aria-hidden="true" className="h-5 w-5" />
            </span>
            <div>
              <h1 className="text-lg font-semibold">localhuman-mail</h1>
              <p className="text-xs uppercase tracking-[0.18em] text-steel">private inbox AI</p>
            </div>
          </div>

          <div className="mt-7 space-y-3 print:hidden">
            <label className="block text-xs font-semibold uppercase tracking-[0.14em] text-steel">
              Backend
            </label>
            <div className="flex items-center gap-2 rounded-lg border border-line bg-paper px-3 py-2">
              <Server aria-hidden="true" className="h-4 w-4 text-steel" />
              <input
                aria-label="Backend URL"
                className="min-w-0 flex-1 bg-transparent text-sm outline-none"
                value={uiState.backendUrl}
                onChange={(event) => patchUIState({ backendUrl: event.target.value })}
              />
            </div>
            <div className="flex items-center gap-2 text-sm">
              <span
                className={`h-2.5 w-2.5 rounded-full ${backendOnline ? "bg-fern" : "bg-coral"}`}
                aria-hidden="true"
              />
              <span>{backendOnline ? "Connected" : "Demo mode"}</span>
            </div>
          </div>

          <div className="mt-7 grid grid-cols-2 gap-2 print:hidden">
            <a
              className="inline-flex items-center justify-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
              href={repoUrl}
              target="_blank"
              rel="noreferrer"
            >
              <Github aria-hidden="true" className="h-4 w-4" />
              Star
            </a>
            <a
              className="inline-flex items-center justify-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
              href={paypalUrl}
              target="_blank"
              rel="noreferrer"
            >
              <CircleDollarSign aria-hidden="true" className="h-4 w-4" />
              PayPal
            </a>
          </div>

          <div className="mt-7 space-y-2 rounded-lg border border-line bg-paper p-3 print:hidden">
            <div className="flex items-center gap-2 text-sm font-semibold">
              <ShieldCheck aria-hidden="true" className="h-4 w-4 text-fern" />
              Local boundary
            </div>
            <p className="text-sm leading-6 text-steel">
              Mail stays in your backend unless you explicitly export or import a browser state
              file.
            </p>
          </div>

          <div className="mt-7 space-y-2 print:hidden">
            <button
              type="button"
              className="inline-flex w-full items-center justify-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
              onClick={() => stateInputRef.current?.click()}
            >
              <Upload aria-hidden="true" className="h-4 w-4" />
              Import state
            </button>
            <button
              type="button"
              className="inline-flex w-full items-center justify-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
              onClick={exportSnapshot}
            >
              <Download aria-hidden="true" className="h-4 w-4" />
              Export state
            </button>
            <button
              type="button"
              className="inline-flex w-full items-center justify-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
              onClick={clearLocalState}
            >
              <Trash2 aria-hidden="true" className="h-4 w-4" />
              Clear local
            </button>
            <input
              ref={stateInputRef}
              className="hidden"
              type="file"
              accept="application/json,.json"
              onChange={(event) => {
                const file = event.target.files?.[0];
                if (file) void importSnapshotFile(file);
                event.currentTarget.value = "";
              }}
            />
          </div>

          <dl className="mt-7 grid grid-cols-[auto_1fr] gap-x-3 gap-y-2 text-xs text-steel">
            <dt>UI</dt>
            <dd className="font-mono text-ink">v{buildInfo.version}</dd>
            <dt>Build</dt>
            <dd className="font-mono text-ink">{buildInfo.commit}</dd>
            {latestCommitQuery.data ? (
              <>
                <dt>Latest</dt>
                <dd>
                  <a
                    className="font-mono text-ink underline decoration-line underline-offset-4 hover:text-fern"
                    href={latestCommitQuery.data.url}
                    target="_blank"
                    rel="noreferrer"
                  >
                    {latestCommitQuery.data.sha}
                  </a>
                </dd>
              </>
            ) : null}
            {versionQuery.data ? (
              <>
                <dt>API</dt>
                <dd className="font-mono text-ink">v{versionQuery.data.version}</dd>
              </>
            ) : null}
          </dl>
        </aside>

        <section className="border-b border-line bg-white lg:border-b-0 lg:border-r print:hidden">
          <div className="sticky top-0 z-10 border-b border-line bg-white/95 px-4 py-4 backdrop-blur">
            <div className="flex items-center gap-2 rounded-lg border border-line bg-paper px-3 py-2">
              <Search aria-hidden="true" className="h-4 w-4 text-steel" />
              <input
                aria-label="Search messages"
                placeholder="Search mail"
                className="min-w-0 flex-1 bg-transparent text-sm outline-none"
                value={uiState.query}
                onChange={(event) => patchUIState({ query: event.target.value })}
              />
            </div>
            <div className="mt-3 flex flex-wrap items-center gap-2">
              <button
                className="inline-flex items-center gap-2 rounded-lg bg-ink px-3 py-2 text-sm font-semibold text-white hover:bg-slate-700 disabled:cursor-not-allowed disabled:opacity-60"
                type="button"
                onClick={() => importDemoMutation.mutate()}
                disabled={importDemoMutation.isPending}
              >
                <MailPlus aria-hidden="true" className="h-4 w-4" />
                Demo
              </button>
              <button
                className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper disabled:cursor-not-allowed disabled:opacity-60"
                type="button"
                onClick={() => fileInputRef.current?.click()}
                disabled={isImporting}
              >
                <Upload aria-hidden="true" className="h-4 w-4" />
                EML files
              </button>
              <button
                className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                type="button"
                onClick={() => void importClipboardText()}
              >
                <Clipboard aria-hidden="true" className="h-4 w-4" />
                Clipboard
              </button>
              <button
                className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                type="button"
                onClick={() => void queryClient.invalidateQueries()}
              >
                <RefreshCw aria-hidden="true" className="h-4 w-4" />
                Sync
              </button>
              <input
                ref={fileInputRef}
                className="hidden"
                type="file"
                accept=".eml,message/rfc822,text/plain"
                multiple
                onChange={(event) => {
                  void importFiles(Array.from(event.target.files ?? []));
                  event.currentTarget.value = "";
                }}
              />
            </div>
            <textarea
              aria-label="Paste raw EML text"
              className="mt-3 min-h-20 w-full resize-y rounded-lg border border-line bg-paper p-3 text-sm leading-6 outline-none focus:border-fern"
              placeholder="Paste raw .eml text here"
              value={uiState.pasteText}
              onChange={(event) => patchUIState({ pasteText: event.target.value })}
            />
            <div className="mt-2 flex flex-wrap items-center gap-2">
              <button
                type="button"
                className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper disabled:cursor-not-allowed disabled:opacity-60"
                disabled={!uiState.pasteText.trim() || isImporting}
                onClick={() => void importPasteText()}
              >
                <Upload aria-hidden="true" className="h-4 w-4" />
                Import paste
              </button>
              <span className="text-xs text-steel">
                Drop files anywhere. {hasSnapshot ? "Viewing imported state." : "Demo fills gaps."}
              </span>
            </div>
          </div>

          {importResults.length ? (
            <div className="border-b border-line bg-paper px-4 py-3 text-sm">
              <div className="flex flex-wrap items-center gap-2 font-semibold">
                <span>{isImporting ? "Importing" : "Import complete"}</span>
                <span className="text-steel">
                  {batchSummary.imported} imported, {batchSummary.skipped} skipped,{" "}
                  {batchSummary.failed} failed
                </span>
              </div>
              <div className="mt-2 grid gap-1">
                {importResults.map((result) => (
                  <div key={result.id} className="flex items-start justify-between gap-3 text-xs">
                    <span className="min-w-0 truncate">{result.name}</span>
                    <span
                      className={`shrink-0 ${
                        result.status === "error" ? "text-rose-700" : "text-steel"
                      }`}
                    >
                      {result.error ?? result.status}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          ) : null}

          <div className="divide-y divide-line" aria-label="Message list">
            {visibleMessages.map((message) => {
              const selected = selectedMessage?.id === message.id;
              return (
                <button
                  key={message.id}
                  type="button"
                  className={`grid w-full gap-2 px-4 py-4 text-left transition hover:bg-paper ${
                    selected ? "bg-emerald-50" : "bg-white"
                  }`}
                  onClick={() => {
                    patchUIState({ selectedId: message.id, draft: "" });
                  }}
                >
                  <span className="flex items-center justify-between gap-3">
                    <span className="truncate text-sm font-semibold">{message.from}</span>
                    <span className="shrink-0 text-xs text-steel">
                      {formatMailDate(message.date)}
                    </span>
                  </span>
                  <span className="truncate text-sm font-semibold">{message.subject}</span>
                  <span className="line-clamp-2 text-sm leading-6 text-steel">
                    {message.snippet}
                  </span>
                  <span className="flex flex-wrap items-center gap-1.5">
                    <span className="rounded-md border border-line bg-white px-2 py-0.5 text-xs text-steel">
                      {message.shape}
                    </span>
                    <span className="rounded-md border border-line bg-white px-2 py-0.5 text-xs text-steel">
                      {message.confidence.label}
                    </span>
                    {message.tags.slice(0, 4).map((tag) => (
                      <span
                        key={tag}
                        className="rounded-md border border-line bg-white px-2 py-0.5 text-xs text-steel"
                      >
                        {tag}
                      </span>
                    ))}
                  </span>
                </button>
              );
            })}
          </div>
        </section>

        <section className="min-w-0 bg-white">
          {selectedMessage ? (
            <article className="grid min-h-screen grid-rows-[auto_1fr_auto]">
              <header className="border-b border-line px-5 py-5">
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <div>
                    <p className="text-sm text-steel">{relativeAge(selectedMessage.date)}</p>
                    <h2 className="mt-1 text-2xl font-semibold leading-tight">
                      {selectedMessage.subject}
                    </h2>
                  </div>
                  <div className="flex flex-wrap items-center gap-2 print:hidden">
                    <button
                      type="button"
                      className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                      onClick={() => void copyValue(selectedMessage.body, "Body")}
                    >
                      <Clipboard aria-hidden="true" className="h-4 w-4" />
                      Body
                    </button>
                    <button
                      type="button"
                      className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                      onClick={() =>
                        downloadTextFile(
                          "localhuman-mail-messages.json",
                          messageListToJSON(currentMessages),
                          "application/json"
                        )
                      }
                    >
                      <FileJson aria-hidden="true" className="h-4 w-4" />
                      JSON
                    </button>
                    <button
                      type="button"
                      className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                      onClick={() =>
                        downloadTextFile(
                          "localhuman-mail-messages.csv",
                          messageListToCSV(currentMessages),
                          "text/csv"
                        )
                      }
                    >
                      <Download aria-hidden="true" className="h-4 w-4" />
                      CSV
                    </button>
                    <button
                      type="button"
                      className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                      onClick={() => void shareSnapshot()}
                    >
                      <Link aria-hidden="true" className="h-4 w-4" />
                      Share
                    </button>
                    <button
                      type="button"
                      className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                      onClick={() => window.print()}
                    >
                      <Printer aria-hidden="true" className="h-4 w-4" />
                      Print
                    </button>
                    <a
                      className="inline-flex items-center gap-2 rounded-lg border border-line px-3 py-2 text-sm font-semibold hover:bg-paper"
                      href={repoUrl}
                      target="_blank"
                      rel="noreferrer"
                    >
                      <ExternalLink aria-hidden="true" className="h-4 w-4" />
                      GitHub
                    </a>
                  </div>
                </div>
                <dl className="mt-4 grid gap-2 text-sm text-steel">
                  <div className="flex gap-2">
                    <dt className="w-20 text-ink">From</dt>
                    <dd>{selectedMessage.from}</dd>
                  </div>
                  <div className="flex gap-2">
                    <dt className="w-20 text-ink">To</dt>
                    <dd>{selectedMessage.to.join(", ") || "(none)"}</dd>
                  </div>
                  <div className="flex gap-2">
                    <dt className="w-20 text-ink">Shape</dt>
                    <dd>
                      {selectedMessage.shape} · {selectedMessage.confidence.label} confidence
                    </dd>
                  </div>
                </dl>
                {selectedMessage.attachments.length ? (
                  <div className="mt-4 flex flex-wrap gap-2 text-sm text-steel">
                    {selectedMessage.attachments.map((attachment) => (
                      <span
                        key={`${attachment.fileName}-${attachment.sizeBytes}`}
                        className="inline-flex items-center gap-2 rounded-lg border border-line px-2 py-1"
                      >
                        <Paperclip aria-hidden="true" className="h-4 w-4" />
                        {attachment.fileName}
                      </span>
                    ))}
                  </div>
                ) : null}
                {selectedMessage.warnings.length ? (
                  <div className="mt-4 grid gap-2">
                    {selectedMessage.warnings.map((warning) => (
                      <div
                        key={`${warning.field}-${warning.message}`}
                        className="rounded-lg border border-line bg-paper p-3 text-sm leading-6"
                      >
                        <div className="flex items-center gap-2 font-semibold">
                          <ShieldAlert aria-hidden="true" className="h-4 w-4 text-coral" />
                          {warning.message}
                        </div>
                        <p className="text-steel">{warning.nextStep}</p>
                      </div>
                    ))}
                  </div>
                ) : null}
              </header>

              <div className="overflow-auto px-5 py-6">
                <pre className="whitespace-pre-wrap font-sans text-base leading-7 text-ink">
                  {selectedMessage.body}
                </pre>
                {debugEnabled ? (
                  <pre className="mt-4 overflow-auto rounded-lg border border-line bg-paper p-3 text-xs leading-5 text-steel">
                    {JSON.stringify(selectedMessage, null, 2)}
                  </pre>
                ) : null}
              </div>

              <footer className="border-t border-line bg-paper px-5 py-4 print:bg-white">
                <div className="flex flex-wrap items-center gap-2 print:hidden">
                  {toneOptions.map((item) => (
                    <button
                      key={item}
                      type="button"
                      className={`inline-flex items-center gap-2 rounded-lg border px-3 py-2 text-sm font-semibold capitalize ${
                        uiState.tone === item
                          ? "border-fern bg-emerald-50 text-fern"
                          : "border-line bg-white"
                      }`}
                      onClick={() => patchUIState({ tone: item })}
                    >
                      {uiState.tone === item ? (
                        <Check aria-hidden="true" className="h-4 w-4" />
                      ) : null}
                      {item}
                    </button>
                  ))}
                  <button
                    type="button"
                    className="inline-flex items-center gap-2 rounded-lg bg-fern px-3 py-2 text-sm font-semibold text-white hover:bg-emerald-800 disabled:cursor-not-allowed disabled:opacity-60"
                    disabled={assistMutation.isPending}
                    onClick={() => assistMutation.mutate(selectedMessage)}
                  >
                    <Bot aria-hidden="true" className="h-4 w-4" />
                    Draft
                  </button>
                  {uiState.draft ? (
                    <button
                      type="button"
                      className="inline-flex items-center gap-2 rounded-lg border border-line bg-white px-3 py-2 text-sm font-semibold hover:bg-paper"
                      onClick={() => void copyValue(uiState.draft, "Draft")}
                    >
                      <Clipboard aria-hidden="true" className="h-4 w-4" />
                      Copy draft
                    </button>
                  ) : null}
                </div>
                {uiState.draft ? (
                  <textarea
                    aria-label="Generated reply draft"
                    className="mt-3 min-h-36 w-full resize-y rounded-lg border border-line bg-white p-3 text-sm leading-6 outline-none focus:border-fern print:border-0 print:p-0"
                    value={uiState.draft}
                    onChange={(event) => patchUIState({ draft: event.target.value })}
                  />
                ) : null}
              </footer>
            </article>
          ) : (
            <div className="grid min-h-screen place-items-center px-6 text-center text-steel">
              <p>No messages found.</p>
            </div>
          )}
        </section>
      </div>
      <Toast toast={toast} />
      <div className="sr-only" aria-live="polite">
        {versionQuery.isFetching ||
        messagesQuery.isFetching ||
        capabilitiesQuery.isFetching ||
        latestCommitQuery.isFetching ||
        isImporting
          ? "Refreshing"
          : "Idle"}
      </div>
    </main>
  );
}

function initialUIState() {
  if (typeof window === "undefined") {
    return defaultUIState(defaultBackendUrl);
  }
  const fallbackBackendUrl = window.localStorage.getItem(legacyBackendUrlKey) ?? defaultBackendUrl;
  return parseStoredUIState(window.localStorage.getItem(uiStateKey), fallbackBackendUrl);
}
