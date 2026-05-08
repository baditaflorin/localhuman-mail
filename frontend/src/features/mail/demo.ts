import type { Message } from "@/api/client";

export const demoMessages: Message[] = [
  {
    id: "demo-1",
    subject: "Re: Contract redlines before the investor call",
    from: "maya@northstar.example",
    to: ["you@localhuman.dev"],
    date: new Date(Date.now() - 18 * 60_000).toISOString(),
    snippet: "Can you send the condensed version with only the liability deltas highlighted?",
    body: [
      "Can you send the condensed version with only the liability deltas highlighted?",
      "",
      "The long memo is useful, but the partner group is going to skim. I mostly need the parts that affect downside exposure, timeline, and approval authority."
    ].join("\n"),
    tags: ["legal", "priority"]
  },
  {
    id: "demo-2",
    subject: "Thursday launch checklist",
    from: "ops@atelier.example",
    to: ["you@localhuman.dev", "shipping@atelier.example"],
    date: new Date(Date.now() - 2 * 60 * 60_000).toISOString(),
    snippet: "The only unresolved items are the DNS cutover and the fallback status page copy.",
    body: [
      "The only unresolved items are the DNS cutover and the fallback status page copy.",
      "",
      "Everything else is green. If you can confirm owner and timing by noon, we can keep Thursday without compressing QA."
    ].join("\n"),
    tags: ["launch"]
  },
  {
    id: "demo-3",
    subject: "Receipt for local model hardware",
    from: "receipts@compute.example",
    to: ["you@localhuman.dev"],
    date: new Date(Date.now() - 22 * 60 * 60_000).toISOString(),
    snippet: "Your order for the workstation memory kit and storage upgrade has shipped.",
    body: [
      "Your order for the workstation memory kit and storage upgrade has shipped.",
      "",
      "Tracking will update after carrier intake. The invoice is attached to this message in the original mailbox."
    ].join("\n"),
    tags: ["finance"]
  }
];

export function filterMessages(messages: Message[], query: string) {
  const normalized = query.trim().toLowerCase();
  if (!normalized) {
    return messages;
  }

  return messages.filter((message) =>
    [message.subject, message.from, message.snippet, message.body, ...message.tags]
      .join(" ")
      .toLowerCase()
      .includes(normalized)
  );
}

export function fallbackDraft(message: Message, tone: string) {
  const opener =
    tone === "warm"
      ? "Thanks for the clear context."
      : tone === "decisive"
        ? "I can take this from here."
        : "Got it.";

  return `${opener}\n\nI reviewed "${message.subject}" and will send the focused version with the decision points, risks, and next action separated. I will keep it tight and flag anything that needs approval.\n\nBest,\n`;
}

