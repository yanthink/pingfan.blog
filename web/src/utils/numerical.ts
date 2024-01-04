export function prettyNumber(n: number, p = 1): string {
  const v = Math.pow(10, p);

  if (n > 1000) {
    return `${Math.floor(n * v / 1000) / v}k`;
  }

  return String(n);
}