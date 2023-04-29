<script setup lang="ts">
import { computed } from "vue"
import { model } from "../../../wailsjs/go/models"

const props = defineProps<{
  user?: model.User
}>()

const userNameDisplay = computed(() => {
  if (props.user == undefined || props.user.name == undefined) {
    return "??"
  } else {
    const endIndex = props.user.name.length > 1 ? 2 : props.user.name.length
    return props.user.name.slice(0, endIndex).toUpperCase()
  }
})
</script>
<template>
  <div v-if="userNameDisplay == '??'" class="w-9 h-9 rounded-full border-4 border-dotted border-labelColor-light-tertiary flex items-center justify-center text-labelColor-light-tertiary font-bold">?</div>
  <div
    v-else
    class="bg-gradient-to-b from-labelColor-light-secondary/80 to-labelColor-light-secondary rounded-full w-9 h-9 flex items-center justify-center font-light text-systemBackground-light select-none cursor-default"
  >
    {{ userNameDisplay }}
  </div>
</template>
