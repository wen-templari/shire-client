<script lang="ts" setup>
import HelloWorld from "./components/HelloWorld.vue"
import { WindowGetSize } from "../wailsjs/runtime"
import { ref } from "vue"
const size = ref()

const getSize = () => {
  WindowGetSize().then(res => {
    size.value = res
  })
}
</script>

<template>
  <div class="flex h-full">
    <div class="h-full flex flex-col w-80 resize-x overflow-auto">
      <div class="h-[92px] flex-shrink-0 px-2 pt-[50px]" style="--wails-draggable: drag">
        <input class="w-full rounded-[6px] bg-borderGrey-light  h-[30px] text-sm px-3" type="text" />
      </div>
      <div class="flex-grow flex flex-col"></div>
    </div>
    <div class="h-full bg-systemBackground-light flex-grow flex flex-col justify-between pb-1">
      <div></div>
      <div class="h-12 flex items-center px-4 gap-4">
        <input
          type="text"
          class="flex-grow rounded-full border border-borderGrey-light px-2 py-0.5 text-sm placeholder:text-gray-400/80"
          placeholder="message" />
      </div>
    </div>
  </div>
</template>

<style></style>
