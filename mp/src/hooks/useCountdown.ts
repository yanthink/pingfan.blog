import { computed, readonly, ref } from 'vue';

export interface Formatted {
  days: number;
  hours: number;
  minutes: number;
  seconds: number;
}

export interface CountdownOption {
  leftSeconds?: number;
  timeup?: () => void;
}

function parse(seconds: number): Formatted {
  return {
    days: Math.floor(seconds / 86400),
    hours: Math.floor(seconds / 3600) % 24,
    minutes: Math.floor(seconds / 60) % 60,
    seconds: seconds % 60,
  };
}

export function useCountdown(option?: CountdownOption) {
  const leftSeconds = ref(option?.leftSeconds ?? 0);
  let timer: NodeJS.Timeout;

  function setLeftSeconds(seconds: number) {
    leftSeconds.value = seconds;
  }

  function start() {
    clearInterval(timer);

    if (leftSeconds.value <= 0) {
      option?.timeup?.();
      return;
    }

    timer = setInterval(() => {
      leftSeconds.value -= 1;

      if (leftSeconds.value <= 0) {
        clearInterval(timer);
        option?.timeup?.();
      }
    }, 1000);
  }

  function stop() {
    clearInterval(timer);
  }

  const formatted = computed(() => parse(leftSeconds.value));

  return {
    leftSeconds: readonly(leftSeconds),
    formatted,
    setLeftSeconds,
    start,
    stop,
  };
}
