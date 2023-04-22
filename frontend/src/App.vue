<script lang="ts" setup>
import { reactive, ref } from "vue"
import { EventsOn, WindowGetSize } from "../wailsjs/runtime"
const size = ref()

const getSize = () => {
  WindowGetSize().then(res => {
    size.value = res
  })
}
EventsOn("ping", e => console.log(e))
</script>

<template>
  <router-view></router-view>
</template>

<style></style>
