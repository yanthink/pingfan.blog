interface Uni {
  request(options: RequestConfig): Promise<Parameters<Exclude<UniNamespace.RequestOptions['success'], undefined>>[0]>;

  updateManager(): any;
}
