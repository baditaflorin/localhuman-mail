import createClient, { type Client } from "openapi-fetch";
import type { components, paths } from "./schema";

export type ApiClient = Client<paths>;
export type Message = components["schemas"]["Message"];
export type Capability = components["schemas"]["Capability"];
export type AssistTone = components["schemas"]["AssistReplyRequest"]["tone"];
export type AssistReplyResponse = components["schemas"]["AssistReplyResponse"];

export function normalizeBackendUrl(value: string) {
  return value.trim().replace(/\/+$/, "");
}

export function createApiClient(baseUrl: string): ApiClient {
  return createClient<paths>({
    baseUrl: normalizeBackendUrl(baseUrl)
  });
}

export function errorMessage(error: unknown, fallback: string) {
  if (error && typeof error === "object") {
    const payload = error as Partial<components["schemas"]["ErrorResponse"]>;
    if (payload.what && payload.why && payload.nowWhat) {
      return `${payload.what}. ${payload.why} ${payload.nowWhat}`;
    }
    if (typeof payload.error === "string") {
      return payload.error;
    }
  }
  if (error instanceof Error) {
    return error.message;
  }
  return fallback;
}
