export async function getRect(selector = '', component: any = null) {
  return new Promise<UniNamespace.NodeInfo | UniNamespace.NodeInfo[]>(resolve => {
    uni.createSelectorQuery().in(component).select(selector).fields({ rect: true, size: true }, data => {
      resolve(data);
    }).exec();
  });
}

