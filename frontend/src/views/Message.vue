<script setup lang="ts">
import { Ping, ReceiveMessage } from "../../wailsjs/go/main/App"
import BaseLayout from "../layout/BaseLayout.vue"
import { useMessageStore } from "../store/message"
import { userAccountStore } from "../store/account"
import { onMounted,ref } from "vue"
import { EventsOn } from "../../wailsjs/runtime/runtime"
import { model } from "../../wailsjs/go/models"

const onButtonClick = () => {
  Ping().then(res => console.log(res))
}

const userStore = userAccountStore()
const messageStore = useMessageStore()


const input = ref("")

const sendMessage = () => {
  if (input.value == "") {
    return
  }
  const message:model.Message = {
    from: userStore.user?.id as number,
    to: messageStore.receiver?.id as number,
    content: input.value,
    time: new Date().toISOString(),
    groupId: -1
  }
  ReceiveMessage(message)
  // SendMessage(message)
  input.value = ""
}

onMounted(() => {
  messageStore.mock()

  userStore.user = {
    id: 1,
    name: "Alice",
    address: "addr",
    port: 1234,
  }
  console.log(messageStore.messageList)
  console.log(messageStore.currentMessageList)
})

EventsOn("onMessage", (data: model.Message) => {
  console.log(data)
})

</script>

<template>
  <base-layout>
    <template #side-head>
      <input class="w-full rounded-[6px] bg-borderGrey-light h-[30px] text-sm px-3" type="text" />
    </template>
    <template #side-body>
      <button @click="onButtonClick">click</button>
    </template>
    <template #main>
      <div class="h-14 bg-fillColor-light-teritary opacity-60 flex items-center justify-center">{{ messageStore.receiver?.name }}</div>
      <div id="messageWindow" class="flex-grow flex flex-col p-2 overflow-auto">
        <div v-for="item in messageStore.currentMessageList">
          <div class="relative flex my-2 gap-2" :class="[item.from == userStore.user?.id ? 'flex-row-reverse ' : '']">
            <div class="bg-fillColor-light-secondary rounded-full w-8 h-8 flex items-center justify-center">
              <div class="text-xs font-light text-gray-500">{{ item.from == userStore.user?.id ? "me" : item.from }}</div>
            </div>
            <div
              class="py-2 px-2 text-sm rounded max-w-40"
              :class="[
                item.from == userStore.user?.id
                  ? 'bg-systemGreen-light text-systemWhite-light'
                  : 'bg-fillColor-light-teritary text-textBlack-light',
              ]"
            >
              {{ item.content }}
            </div>
          </div>
        </div>
        <div id="messageWindowBottom"></div>
      </div>
      <div class="h-12 flex items-center px-4 gap-4">
        <input
          type="text"
          v-model="input"
          class="flex-grow rounded-full border border-borderGrey-light px-2 py-0.5 text-sm placeholder:text-gray-400/80"
          placeholder="message"
          @keypress.enter="sendMessage"
        />
      </div>
    </template>
  </base-layout>
</template>
