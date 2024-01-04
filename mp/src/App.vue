<script setup lang="ts">
import { onShow } from '@dcloudio/uni-app';
import { useWebsocket } from '@/hooks';
import { ReadyState } from "@/stores";

onShow(() => {
  const { readyState, connect } = useWebsocket();

  if (readyState.value === ReadyState.Closed) {
    connect();
  }
})
</script>
<style lang="scss">
@import "./global.scss";
</style>
