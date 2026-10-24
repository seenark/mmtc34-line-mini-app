import type { ErrorResponse } from "./types.gen";

/** `${BASE_URL}api` where BASE_URL is Vite's `base`, i.e. "<BASE_PATH>/". */
export const API = `${import.meta.env.BASE_URL}api`;

/** Thrown for every non-2xx. message is the server's Thai `error`. */
export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message);
  }
}

/** JSON fetch with the session cookie. A FormData body is sent as multipart (uploads). */
export async function api<T>(
  method: string,
  path: string,
  body?: unknown,
): Promise<T> {
  const form = body instanceof FormData;
  const res = await fetch(`${API}${path}`, {
    method,
    credentials: "include",
    headers:
      body === undefined || form
        ? undefined
        : { "Content-Type": "application/json" },
    body: body === undefined ? undefined : form ? body : JSON.stringify(body),
  });
  if (!res.ok) {
    const err: Partial<ErrorResponse> = await res.json().catch(() => ({}));
    throw new ApiError(res.status, err.error ?? res.statusText);
  }
  return res.status === 204 ? (undefined as T) : res.json();
}
