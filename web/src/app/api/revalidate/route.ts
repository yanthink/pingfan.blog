import { NextRequest } from 'next/server'
import { revalidatePath, revalidateTag } from 'next/cache'

export async function POST(request: NextRequest) {
  const { paths = [], tags = [] } = await request.json();

  paths.forEach((path: string) => revalidatePath(path));
  tags.forEach((tag: string) => revalidateTag(tag));

  return Response.json({ revalidated: true, now: Date.now() })
}