<script setup lang="ts">
import { ref } from "vue"
import { SendMessage } from "../../../wailsjs/go/main/App"
import { model } from "../../../wailsjs/go/models"
import { userAccountStore } from "../../store/account"
import { useMessageStore } from "../../store/message"

const userStore = userAccountStore()
const messageStore = useMessageStore()

const input = ref("")
const sendMessage = () => {
  if (input.value == "") {
    return
  }
  const message: model.Message = {
    from: userStore.user?.id as number,
    to: messageStore.receiver?.id as number,
    content: input.value,
    time: new Date().toISOString(),
    groupId: -1,
  }
  // ReceiveMessage(message)
  SendMessage(message)
  input.value = ""
}
</script>
<template>
  <div class="h-14 bg-fillColor-light-teritary opacity-60 flex items-center px-4 border-b border-labelColor-light-tertiary text-sm">
    <div>收件人：</div>
    <span class="font-semibold ml-1">{{ messageStore.receiver?.name }}</span>
    <span class="ml-1 textDescription text-labelColor-light-secondary">({{ messageStore.receiver?.id }})</span>
  </div>
  <div id="messageWindow" class="flex-grow flex flex-col px-4 py-3 overflow-auto">
    <div>
      <div v-for="item in messageStore.currentMessageList?.messages">
        <div class="relative flex my-2 gap-2" :class="[item?.from == userStore.user?.id ? 'flex-row-reverse ' : '']">
          <user-avatar :user="item.from == userStore.user?.id ? userStore.user : messageStore.receiver"></user-avatar>
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
