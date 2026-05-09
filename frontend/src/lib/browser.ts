export async function copyText(value: string) {
  if (!navigator.clipboard?.writeText) {
    throw new Error("Clipboard write is not available in this browser.");
  }
  await navigator.clipboard.writeText(value);
}

export async function readClipboardText() {
  if (!navigator.clipboard?.readText) {
    throw new Error("Clipboard read is not available in this browser.");
  }
  return navigator.clipboard.readText();
}

export function downloadTextFile(filename: string, contents: string, type: string) {
  const blob = new Blob([contents], { type });
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = url;
  anchor.download = filename;
  anchor.click();
  URL.revokeObjectURL(url);
}
