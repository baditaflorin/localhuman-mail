import { Component, type ErrorInfo, type ReactNode } from "react";

type Props = {
  children: ReactNode;
};

type State = {
  failed: boolean;
};

export class ErrorBoundary extends Component<Props, State> {
  state: State = { failed: false };

  static getDerivedStateFromError() {
    return { failed: true };
  }

  componentDidCatch(_error: Error, _errorInfo: ErrorInfo) {
    if (import.meta.env.DEV) {
      // Development-only signal; production builds keep mailbox-safe silence.
      console.error(_error, _errorInfo);
    }
  }

  render() {
    if (this.state.failed) {
      return (
        <main className="grid min-h-screen place-items-center bg-paper p-6 text-ink">
          <section className="max-w-md rounded-lg border border-line bg-white p-6 shadow-pane">
            <h1 className="text-xl font-semibold">Something broke</h1>
            <p className="mt-3 text-sm leading-6 text-steel">
              Refresh the app. If the problem repeats, check the backend URL and browser storage.
            </p>
          </section>
        </main>
      );
    }

    return this.props.children;
  }
}
