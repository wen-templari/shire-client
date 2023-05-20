<script setup lang="ts">
import { ref } from "vue"
import { SendMessage } from "../../../wailsjs/go/main/App"
import { model } from "../../../wailsjs/go/models"
import { userAccountStore } from "../../store/account"
import { useMessageStore } from "../../store/message"
import UserAvatar from "../../components/User/UserAvatar.vue"

const userStore = userAccountStore()
const messageStore = useMessageStore()

const input = ref("")
const errorFlag = ref(false)
const sendMessage = async () => {
  if (input.value == "") {
    return
  }
  const message: model.Message = {
    from: userStore.user?.id as number,
    to: messageStore.receiver?.id as number,
    groupId: messageStore.group == undefined ? -1 : (messageStore.group.id as number),
    content: input.value,
    time: new Date().toISOString(),
  }
  // ReceiveMessage(message)
  try {
    await SendMessage(message)
    input.value = ""
    errorFlag.value = false
  } catch (error) {
    errorFlag.value = true
    setTimeout(() => {
      errorFlag.value = false
    }, 3000)
    console.log(error)
  }
}

const scrollToBottom = () => {
  var bottom = document.getElementById("messageWindowBottom")
  bottom?.scrollIntoView({ behavior: "smooth" })
}
defineExpose({ scrollToBottom })
</script>
<template>
  <div
    class="h-14 shrink-0 bg-fillColor-light-teritary opacity-60 flex items-center px-4 border-b border-labelColor-light-tertiary text-sm"
  >
    <div>收件人：</div>
    <div v-if="messageStore.receiver != undefined">
      <span class="font-semibold ml-1">{{ messageStore.receiver?.name }}</span>
      <span class="ml-1 textDescription text-labelColor-light-secondary">({{ messageStore.receiver?.id }})</span>
    </div>
    <div v-if="messageStore.group != undefined" class="flex">
      <div v-for="user in messageStore.group.users" class="">
        <span class="font-semibold ml-1">{{ user.name }}</span>
        <span class="ml-1 textDescription text-labelColor-light-secondary">({{ user.id }})</span>
      </div>
    </div>
  </div>
  <div id="messageyWindow" class="flex-grow flex flex-col px-4 py-3 overflow-auto">
    <div>
      <div v-for="item in messageStore.currentMessageList?.messages">
        <div class="relative flex mb-4 gap-2" :class="[item?.from == userStore.user?.id ? 'flex-row-reverse ' : '']">
          <user-avatar :user="userStore.userList.find(u => u.id == item.from)" class="mt-4"></user-avatar>
          <div>
            <div class="text-xs text-labelColor-light-secondary text-left pl-2">
              {{ userStore.userList.find(u => u.id == item.from)?.name }}
            </div>
            <div
              class="py-2 px-2 text-sm rounded text-start"
              style="max-width: 260px"
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
    </div>
    <div id="messageWindowBottom"></div>
  </div>
  <div class="h-12 flex items-center px-4 gap-4 relative flex-shrink-0">
    <transition>
      <div class="absolute left-0 -top-3 px-4 text-sm text-systemRed-light" v-if="errorFlag">发送失败！对方可能不在线。</div>
    </transition>
    <input
      type="text"
      v-model="input"
      class="flex-grow rounded-full border border-borderGrey-light px-2 py-0.5 text-sm placeholder:text-gray-400/80"
      placeholder="message"
      @keypress.enter="sendMessage"
    />
  </div>
  <div></div>
</template>

<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
