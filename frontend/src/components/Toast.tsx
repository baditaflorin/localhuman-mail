import { AlertCircle, CheckCircle2 } from "lucide-react";

export type ToastState = {
  kind: "success" | "error";
  message: string;
} | null;

type Props = {
  toast: ToastState;
};

export function Toast({ toast }: Props) {
  if (!toast) {
    return null;
  }

  const Icon = toast.kind === "success" ? CheckCircle2 : AlertCircle;
  const color = toast.kind === "success" ? "text-fern" : "text-rose-700";

  return (
    <div
      role="status"
      className="pointer-events-none fixed bottom-5 right-5 z-50 flex max-w-sm items-start gap-3 rounded-lg border border-line bg-white px-4 py-3 text-sm shadow-pane"
    >
      <Icon aria-hidden="true" className={`mt-0.5 h-4 w-4 ${color}`} />
      <span className="leading-5 text-ink">{toast.message}</span>
    </div>
  );
}
