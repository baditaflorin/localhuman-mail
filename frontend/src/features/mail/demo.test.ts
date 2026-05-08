import { describe, expect, it } from "vitest";
import { demoMessages, fallbackDraft, filterMessages } from "./demo";

describe("mail demo helpers", () => {
  it("filters messages by subject, sender, body, and tags", () => {
    expect(filterMessages(demoMessages, "liability")).toHaveLength(1);
    expect(filterMessages(demoMessages, "launch")).toHaveLength(1);
    expect(filterMessages(demoMessages, "missing")).toHaveLength(0);
  });

  it("builds a fallback reply without remote services", () => {
    const draft = fallbackDraft(demoMessages[0], "warm");

    expect(draft).toContain("Thanks for the clear context");
    expect(draft).toContain(demoMessages[0].subject);
  });
});
