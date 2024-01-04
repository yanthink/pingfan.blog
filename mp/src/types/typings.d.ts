type UniPromiseOptions<T = Record<string, any>> = Omit<T, 'success' | 'fail' | 'complete'>;
