import { useUser } from "@/hooks/useUser";
import { watch } from "vue";

export function useUserPromise() {
  const { user } = useUser();

  let userResolve: (u: typeof user) => void;
  const userPromise = new Promise<typeof user>(resolve => userResolve = resolve);

  watch(user, newUser => {
    if (newUser.id) {
      userResolve?.(user);
    }
  }, { immediate: true });

  return userPromise;
}