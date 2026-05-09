import { z } from "zod";
import type { AssistTone, Message } from "@/api/client";

export const uiStateVersion = "localhuman.uiState.v1";
export const stateFileName = "localhuman-mail-state.json";
export const maxShareUrlLength = 20_000;

const confidenceSchema = z.object({
  score: z.number().min(0).max(1),
  label: z.enum(["low", "medium", "high"]),
  reasons: z.array(z.string())
});

const warningSchema = z.object({
  severity: z.string(),
  field: z.string(),
  message: z.string(),
  nextStep: z.string()
});

const attachmentSchema = z.object({
  fileName: z.string(),
  contentType: z.string(),
  sizeBytes: z.number()
});

const calendarSchema = z.object({
  summary: z.string(),
  location: z.string(),
  start: z.string(),
  end: z.string()
});

const provenanceSchema = z.object({
  sourceId: z.string(),
  sourceSha256: z.string(),
  parserVersion: z.string(),
  schemaVersion: z.string(),
  sizeBytes: z.number()
});

export const messageSchema = z.object({
  id: z.string(),
  sourceId: z.string(),
  subject: z.string(),
  from: z.string(),
  to: z.array(z.string()),
  date: z.string(),
  snippet: z.string(),
  body: z.string(),
  primaryBody: z.string(),
  shape: z.string(),
  tags: z.array(z.string()),
  confidence: confidenceSchema,
  fieldConfidence: z.record(confidenceSchema),
  warnings: z.array(warningSchema),
  attachments: z.array(attachmentSchema),
  calendar: calendarSchema.optional(),
  provenance: provenanceSchema
}) satisfies z.ZodType<Message>;

export const uiStateSchema = z.object({
  version: z.literal(uiStateVersion),
  backendUrl: z.string(),
  query: z.string(),
  selectedId: z.string(),
  tone: z.enum(["concise", "warm", "decisive"]),
  draft: z.string(),
  pasteText: z.string(),
  snapshotMessages: z.array(messageSchema).max(200)
});

export type UIState = z.infer<typeof uiStateSchema>;

export const stateSnapshotSchema = z.object({
  version: z.literal(uiStateVersion),
  exportedAt: z.string(),
  app: z.object({
    version: z.string(),
    commit: z.string()
  }),
  ui: uiStateSchema,
  messages: z.array(messageSchema).max(200)
});

export type StateSnapshot = z.infer<typeof stateSnapshotSchema>;

export type ImportStatus = "pending" | "imported" | "skipped" | "error";

export type FileImportResult = {
  id: string;
  name: string;
  size: number;
  status: ImportStatus;
  imported: number;
  skipped: number;
  error?: string;
};

export type BatchSummary = {
  imported: number;
  skipped: number;
  failed: number;
  total: number;
};

export function defaultUIState(backendUrl: string): UIState {
  return {
    version: uiStateVersion,
    backendUrl,
    query: "",
    selectedId: "",
    tone: "concise",
    draft: "",
    pasteText: "",
    snapshotMessages: []
  };
}

export function migrateUIState(value: unknown, fallbackBackendUrl: string): UIState {
  const parsed = uiStateSchema.safeParse(value);
  if (parsed.success) {
    return parsed.data;
  }
  const legacy = z
    .object({
      backendUrl: z.string().optional(),
      query: z.string().optional(),
      selectedId: z.string().optional(),
      tone: z.enum(["concise", "warm", "decisive"]).optional(),
      draft: z.string().optional()
    })
    .safeParse(value);
  if (!legacy.success) {
    return defaultUIState(fallbackBackendUrl);
  }
  return {
    ...defaultUIState(legacy.data.backendUrl ?? fallbackBackendUrl),
    query: legacy.data.query ?? "",
    selectedId: legacy.data.selectedId ?? "",
    tone: legacy.data.tone ?? "concise",
    draft: legacy.data.draft ?? ""
  };
}

export function parseStoredUIState(raw: string | null, fallbackBackendUrl: string): UIState {
  if (!raw) {
    return defaultUIState(fallbackBackendUrl);
  }
  try {
    return migrateUIState(JSON.parse(raw), fallbackBackendUrl);
  } catch {
    return defaultUIState(fallbackBackendUrl);
  }
}

export function serializeUIState(state: UIState) {
  return JSON.stringify(state);
}

export function createSnapshot(
  ui: UIState,
  messages: Message[],
  app: { version: string; commit: string },
  now = new Date()
): StateSnapshot {
  return {
    version: uiStateVersion,
    exportedAt: now.toISOString(),
    app,
    ui: {
      ...ui,
      snapshotMessages: []
    },
    messages: messages.slice(0, 200)
  };
}

export function parseSnapshotText(text: string): StateSnapshot {
  const parsed = stateSnapshotSchema.safeParse(JSON.parse(text));
  if (!parsed.success) {
    throw new Error("State file is not a valid localhuman-mail snapshot.");
  }
  return parsed.data;
}

export function snapshotToJSON(snapshot: StateSnapshot) {
  return JSON.stringify(snapshot, null, 2);
}

export function encodeSnapshotHash(snapshot: StateSnapshot) {
  const json = JSON.stringify(snapshot);
  const encoded = btoa(unescape(encodeURIComponent(json)))
    .replaceAll("+", "-")
    .replaceAll("/", "_")
    .replaceAll("=", "");
  return `state=${encoded}`;
}

export function decodeSnapshotHash(hash: string): StateSnapshot | null {
  const value = hash.replace(/^#/, "");
  if (!value.startsWith("state=")) {
    return null;
  }
  const encoded = value.slice("state=".length).replaceAll("-", "+").replaceAll("_", "/");
  const padded = encoded.padEnd(Math.ceil(encoded.length / 4) * 4, "=");
  try {
    return parseSnapshotText(decodeURIComponent(escape(atob(padded))));
  } catch {
    return null;
  }
}

export function messageListToJSON(messages: Message[]) {
  return JSON.stringify(
    {
      schemaVersion: "localhuman.messages.v1",
      exportedAt: new Date().toISOString(),
      messages
    },
    null,
    2
  );
}

export function messageListToCSV(messages: Message[]) {
  const header = [
    "id",
    "date",
    "from",
    "to",
    "subject",
    "shape",
    "confidence",
    "warnings",
    "snippet"
  ];
  const rows = messages.map((message) => [
    message.id,
    message.date,
    message.from,
    message.to.join("; "),
    message.subject,
    message.shape,
    message.confidence.label,
    message.warnings.map((warning) => warning.message).join("; "),
    message.snippet
  ]);
  return [header, ...rows].map((row) => row.map(csvCell).join(",")).join("\n") + "\n";
}

export function csvCell(value: string) {
  return `"${value.replaceAll('"', '""')}"`;
}

export function looksLikeEMLText(text: string) {
  return /(^|\n)(From|To|Subject|Date|Message-ID|Return-Path):/i.test(text);
}

export async function fileLooksLikeEML(file: File) {
  if (file.name.toLowerCase().endsWith(".eml") || file.type === "message/rfc822") {
    return true;
  }
  const sample = await file.slice(0, 4096).text();
  return looksLikeEMLText(sample);
}

export function textToEmlFile(text: string) {
  return new File([text], `pasted-${Date.now()}.eml`, { type: "message/rfc822" });
}

export function summarizeBatch(results: FileImportResult[]): BatchSummary {
  return results.reduce(
    (summary, result) => ({
      imported: summary.imported + result.imported,
      skipped: summary.skipped + result.skipped,
      failed: summary.failed + (result.status === "error" ? 1 : 0),
      total: summary.total + 1
    }),
    { imported: 0, skipped: 0, failed: 0, total: 0 }
  );
}

export function fileResultId(file: File, index: number) {
  return `${file.name}-${file.size}-${file.lastModified}-${index}`;
}

export function normalizeTone(value: string): AssistTone {
  const parsed = z.enum(["concise", "warm", "decisive"]).safeParse(value);
  return parsed.success ? parsed.data : "concise";
}
