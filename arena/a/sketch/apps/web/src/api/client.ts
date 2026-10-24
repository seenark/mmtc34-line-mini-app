import type { ErrorResponse } from "./types.gen";

/** `${BASE_URL}api` — BASE_URL is Vite's `base`, i.e. "<BASE_PATH>/". */
export const API = `${import.meta.env.BASE_URL}api`;

/** Thrown for every non-2xx; message is the server's Thai `error`. */
export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message);
  }
}

/** JSON fetch with the session cookie. body FormData is sent as multipart (uploads). */
export async function api<T>(method: string, path: string, body?: unknown): Promise<T> {
  throw new Error("not implemented");
}

export type { ErrorResponse };
