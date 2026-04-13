export type ErrorCode = 
  | 'NOT_FOUND' 
  | 'INVALID_INPUT' 
  | 'INTERNAL_ERROR' 
  | 'UNAUTHORIZED' 
  | 'DB_ERROR'
  | 'UNKNOWN';

export interface AppError {
  code: ErrorCode;
  message: string;
  details?: string;
}

export type Result<T> = 
  | { data: T; error: null }
  | { data: null; error: AppError };

type CacheEntry = {
  expiresAt: number;
  value: unknown;
};

const DEFAULT_CACHE_TTL_MS = 30_000;
const responseCache = new Map<string, CacheEntry>();

export function invalidateInvokeCache(keyPrefix?: string): void {
  if (!keyPrefix) {
    responseCache.clear();
    return;
  }

  for (const key of responseCache.keys()) {
    if (key.startsWith(keyPrefix)) {
      responseCache.delete(key);
    }
  }
}

function readCache<T>(key: string): Result<T> | null {
  const hit = responseCache.get(key);
  if (!hit) return null;

  if (Date.now() >= hit.expiresAt) {
    responseCache.delete(key);
    return null;
  }

  return { data: hit.value as T, error: null };
}

function writeCache<T>(key: string, data: T, ttlMs: number): void {
  responseCache.set(key, {
    value: data,
    expiresAt: Date.now() + ttlMs,
  });
}

/**
 * Safely invokes a Wails TS binding function, catching Go panics or rejected promises
 * and converting them into a structured AppError.
 */
export async function invoke<T, Args extends any[]>(
  bindFn: (...args: Args) => Promise<T>,
  ...args: Args
): Promise<Result<T>> {
  try {
    const data = await bindFn(...args);
    return { data, error: null };
  } catch (err: any) {
    // Determine if the error is already our structured AppError format
    // (In Wails, if Go returns an error, it comes back as a string, unless
    // specifically wrapped and parsed. Here we handle the most common formats)
    let appError: AppError = {
      code: 'UNKNOWN',
      message: 'An unexpected error occurred.',
    };

    if (typeof err === 'string') {
      try {
        const parsed = JSON.parse(err);
        if (parsed.code && parsed.message) {
          appError = parsed as AppError;
        } else {
          appError.message = err;
        }
      } catch {
        appError.message = err;
      }
    } else if (err instanceof Error) {
      appError.message = err.message;
    } else if (typeof err === 'object' && err !== null) {
       if (err.code && err.message) {
         appError = err as AppError;
       } else {
         appError.message = JSON.stringify(err);
       }
    }

    return { data: null, error: appError };
  }
}

export async function invokeCached<T, Args extends any[]>(
  cacheKey: string,
  bindFn: (...args: Args) => Promise<T>,
  args: Args,
  ttlMs: number = DEFAULT_CACHE_TTL_MS,
): Promise<Result<T>> {
  const cached = readCache<T>(cacheKey);
  if (cached) return cached;

  const result = await invoke(bindFn, ...args);
  if (!result.error) {
    writeCache(cacheKey, result.data, ttlMs);
  }

  return result;
}
