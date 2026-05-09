import { describe, expect, it } from "vitest";
import { demoMessages } from "./demo";
import {
  createSnapshot,
  csvCell,
  decodeSnapshotHash,
  defaultUIState,
  encodeSnapshotHash,
  fileLooksLikeEML,
  looksLikeEMLText,
  messageListToCSV,
  migrateUIState,
  parseSnapshotText,
  snapshotToJSON,
  summarizeBatch,
  textToEmlFile
} from "./workbench";

describe("mail workbench helpers", () => {
  it("detects raw EML text and file-like inputs", async () => {
    const text = "From: Maya <maya@example.com>\nSubject: Hi\n\nBody";

    expect(looksLikeEMLText(text)).toBe(true);
    expect(looksLikeEMLText("just a note")).toBe(false);
    expect(await fileLooksLikeEML(textToEmlFile(text))).toBe(true);
  });

  it("round-trips a versioned state snapshot", () => {
    const state = defaultUIState("http://localhost:8080");
    const snapshot = createSnapshot(
      { ...state, query: "launch", selectedId: demoMessages[1].id },
      demoMessages,
      { version: "0.3.0", commit: "test" },
      new Date("2026-05-09T12:00:00Z")
    );

    const parsed = parseSnapshotText(snapshotToJSON(snapshot));

    expect(parsed.ui.query).toBe("launch");
    expect(parsed.messages).toHaveLength(demoMessages.length);
  });

  it("encodes and decodes share hashes", () => {
    const snapshot = createSnapshot(defaultUIState("http://localhost:8080"), [demoMessages[0]], {
      version: "0.3.0",
      commit: "test"
    });
    const decoded = decodeSnapshotHash(`#${encodeSnapshotHash(snapshot)}`);

    expect(decoded?.messages[0].subject).toBe(demoMessages[0].subject);
  });

  it("exports deterministic CSV rows", () => {
    expect(csvCell('a "quoted" value')).toBe('"a ""quoted"" value"');
    expect(messageListToCSV([demoMessages[0]])).toContain("personal_reply");
  });

  it("migrates legacy local state and summarizes batches", () => {
    const migrated = migrateUIState({ backendUrl: "http://api", tone: "warm" }, "fallback");

    expect(migrated.backendUrl).toBe("http://api");
    expect(migrated.tone).toBe("warm");
    expect(
      summarizeBatch([
        { id: "1", name: "a.eml", size: 1, status: "imported", imported: 1, skipped: 0 },
        { id: "2", name: "b.eml", size: 1, status: "error", imported: 0, skipped: 0 }
      ])
    ).toEqual({ imported: 1, skipped: 0, failed: 1, total: 2 });
  });
});
