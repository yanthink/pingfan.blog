export function useLockFn<P extends any[] = any[], V extends any = any>(fn: (...args: P) => Promise<V>) {
  let lock = false;

  return async (...args: P) => {
    if (lock) return Promise.reject();

    lock = true;

    try {
      return await fn(...args);
    } finally {
      lock = false;
    }
  };
}