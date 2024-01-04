'use client';

import type { ReadonlyURLSearchParams } from 'next/navigation';

export function mergeSearchParams(searchParams: ConstructorParameters<typeof URLSearchParams>[0] | ReadonlyURLSearchParams, params: Record<string, any>): URLSearchParams {
  const newSearchParams = new URLSearchParams(searchParams);
  Object.entries(params).forEach(([key, value]) => newSearchParams.set(key, value));

  return newSearchParams;
}

export * from './auth';
export * from './numerical';