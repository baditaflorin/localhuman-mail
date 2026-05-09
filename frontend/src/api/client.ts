import createClient, { type Client } from "openapi-fetch";
import { z } from "zod";
import type { components, paths } from "./schema";

export type ApiClient = Client<paths>;
export type Message = components["schemas"]["Message"];
export type Capability = components["schemas"]["Capability"];
export type AssistTone = components["schemas"]["AssistReplyRequest"]["tone"];
export type AssistReplyResponse = components["schemas"]["AssistReplyResponse"];
export type ImportResponse = components["schemas"]["ImportResponse"];

const errorResponseSchema = z.object({
  error: z.string(),
  kind: z.string().optional(),
  what: z.string().optional(),
  why: z.string().optional(),
  nowWhat: z.string().optional()
});

const importResponseSchema = z.object({
  imported: z.number(),
  skipped: z.number(),
  total: z.number()
});

export class UserFacingError extends Error {
  constructor(
    message: string,
    readonly detail?: components["schemas"]["ErrorResponse"]
  ) {
    super(message);
    this.name = "UserFacingError";
  }
}

export function normalizeBackendUrl(value: string) {
  return value.trim().replace(/\/+$/, "");
}

export function createApiClient(baseUrl: string): ApiClient {
  return createClient<paths>({
    baseUrl: normalizeBackendUrl(baseUrl)
  });
}

export async function uploadEmlFile(baseUrl: string, file: File): Promise<ImportResponse> {
  const body = new FormData();
  body.append("file", file);
  const response = await fetch(`${normalizeBackendUrl(baseUrl)}/api/v1/import/eml`, {
    method: "POST",
    body
  });
  const payload = await parseJSONResponse(response);
  if (!response.ok) {
    throw new UserFacingError(
      errorMessage(payload, "EML import failed"),
      parseErrorResponse(payload)
    );
  }
  return importResponseSchema.parse(payload);
}

export function errorMessage(error: unknown, fallback: string) {
  const parsed = errorResponseSchema.safeParse(error);
  if (parsed.success) {
    const payload = parsed.data;
    if (payload.what && payload.why && payload.nowWhat) {
      return `${payload.what}. ${payload.why} ${payload.nowWhat}`;
    }
    return payload.error;
  }
  if (error instanceof UserFacingError) {
    return error.message;
  }
  if (error instanceof Error) {
    return error.message;
  }
  return fallback;
}

async function parseJSONResponse(response: Response): Promise<unknown> {
  const text = await response.text();
  if (!text) {
    return {};
  }
  try {
    return JSON.parse(text);
  } catch {
    return { error: text };
  }
}

function parseErrorResponse(value: unknown) {
  const parsed = errorResponseSchema.safeParse(value);
  return parsed.success ? parsed.data : undefined;
}
