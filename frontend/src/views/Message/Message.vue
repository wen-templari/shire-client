<script setup lang="ts">
import { GetUserById, Ping, ReceiveMessage, SendMessage } from "../../../wailsjs/go/main/App"
import { GetUsers } from "../../../wailsjs/go/main/App"
import BaseLayout from "../../layout/BaseLayout.vue"
import { type messageList, useMessageStore } from "../../store/message"
import { userAccountStore } from "../../store/account"
import { nextTick, onMounted, ref } from "vue"
import { EventsOn } from "../../../wailsjs/runtime/runtime"
import { model } from "../../../wailsjs/go/models"
import UserAvatar from "../../components/User/UserAvatar.vue"
import MessageView from "./MessageView.vue"
import GroupMessageStarterVue from "./GroupMessageStarter.vue"
import router from "../../router"

const userStore = userAccountStore()
const messageStore = useMessageStore()

const searchInput = ref("")
const searchResult = ref<model.User[]>([])
const onSearchInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  const lowercaseValue = target.value.toLowerCase()
  searchResult.value = userStore.userList.filter(user => user.name?.toLocaleLowerCase().indexOf(lowercaseValue) != -1)
  // messageStore.searchUser(target.value)
}
const onSelectUser = (user: model.User) => {
  startGroup.value = false
  searchInput.value = ""
  messageStore.selectContact(user)
  nextTick().then(() => {
    if (messageView.value != undefined) {
      messageView.value.scrollToBottom()
    }
  })
}

const startGroup = ref(false)
const onStartGroup = () => {
  startGroup.value = true
  userStore.updateUserList()
}
const onGroupStarted = (group: model.Group) => {
  startGroup.value = true
  messageStore.selectContact(group)
}

const onSelectContact = (e: messageList) => {
  console.log(e)

  startGroup.value = false
  if (e.user != undefined) {
    messageStore.selectContact(e.user)
  } else if (e.group != undefined) {
    messageStore.selectContact(e.group)
  }

  nextTick().then(() => {
    if (messageView.value != undefined) {
      messageView.value.scrollToBottom()
    }
  })
}

const messageView = ref()
onMounted(() => {
  if (userStore.user?.id == undefined) {
    router.push("/login")
  }
  messageStore.initMessageList()
  // userStore.updateUserList()
})

EventsOn("onMessage", (data: model.Message) => {
  messageStore.onReceiveMessage(data, GetUserById).then(() => {
    if (messageView.value != undefined) {
      messageView.value.scrollToBottom()
    }
  })
})
</script>

<template>
  <base-layout>
    <template #side-head>
      <div class="relative flex items-center gap-2">
        <input
          class="w-full rounded-[6px] bg-borderGrey-light h-[30px] text-sm px-3 placeholder:text-labelColor-light-secondary"
          type="text"
          placeholder="search users"
          v-model="searchInput"
          @input="onSearchInput"
        />
        <div
          class="hover:bg-borderGrey-light cursor-pointer h-[30px] w-[30px] flex items-center justify-center p-[2px] rounded-[6px]"
          @click="onStartGroup"
        >
          <svg width="28" height="15" viewBox="0 0 28 15" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M14 7.77783C15.7651 7.77783 17.2007 6.21777 17.2007 4.30615C17.2007 2.4165 15.7798 0.915039 14 0.915039C12.2349 0.915039 10.792 2.43848 10.792 4.31348C10.7993 6.2251 12.2275 7.77783 14 7.77783ZM5.50391 7.94629C7.04932 7.94629 8.29443 6.57666 8.29443 4.89209C8.29443 3.24414 7.04932 1.91113 5.50391 1.91113C3.98047 1.91113 2.71338 3.26611 2.7207 4.89941C2.7207 6.58398 3.97314 7.94629 5.50391 7.94629ZM22.4814 7.94629C24.0195 7.94629 25.2646 6.58398 25.272 4.89941C25.2793 3.26611 24.0122 1.91113 22.4814 1.91113C20.9434 1.91113 19.6909 3.24414 19.6909 4.89209C19.6909 6.57666 20.9434 7.94629 22.4814 7.94629ZM14 6.51807C12.9819 6.51807 12.125 5.55127 12.125 4.31348C12.1177 3.1123 12.9746 2.18213 14 2.18213C15.0254 2.18213 15.875 3.09766 15.875 4.30615C15.875 5.53662 15.0181 6.51807 14 6.51807ZM5.50391 6.70117C4.67627 6.70117 3.96582 5.90283 3.96582 4.89941C3.96582 3.93994 4.66895 3.15625 5.50391 3.15625C6.36084 3.15625 7.05664 3.92529 7.05664 4.89209C7.05664 5.90283 6.34619 6.70117 5.50391 6.70117ZM22.4814 6.70117C21.6392 6.70117 20.936 5.90283 20.936 4.89209C20.936 3.92529 21.6318 3.15625 22.4814 3.15625C23.3237 3.15625 24.0269 3.93994 24.0269 4.89941C24.0269 5.90283 23.3164 6.70117 22.4814 6.70117ZM1.47559 14.6772H7.1665C6.771 14.4502 6.5 13.9229 6.55127 13.4395H1.39502C1.25586 13.4395 1.19727 13.3809 1.19727 13.249C1.19727 11.5645 3.1748 9.95312 5.50391 9.95312C6.32422 9.95312 7.15186 10.1582 7.78906 10.5171C8.03076 10.1655 8.33838 9.86523 8.73389 9.60156C7.80371 9.02295 6.64648 8.71533 5.50391 8.71533C2.40576 8.71533 -0.0991211 10.9346 -0.0991211 13.3662C-0.0991211 14.2305 0.428223 14.6772 1.47559 14.6772ZM26.5171 14.6772C27.5645 14.6772 28.0918 14.2305 28.0918 13.3662C28.0918 10.9346 25.5796 8.71533 22.4888 8.71533C21.3462 8.71533 20.189 9.02295 19.2588 9.60156C19.6543 9.86523 19.9619 10.1655 20.2036 10.5171C20.8408 10.1582 21.6611 9.95312 22.4888 9.95312C24.8179 9.95312 26.7954 11.5645 26.7954 13.249C26.7954 13.3809 26.7368 13.4395 26.5977 13.4395H21.4414C21.4927 13.9229 21.2217 14.4502 20.8262 14.6772H26.5171ZM9.42969 14.6772H18.563C19.8228 14.6772 20.4307 14.2817 20.4307 13.4248C20.4307 11.418 17.9258 8.72266 13.9927 8.72266C10.0669 8.72266 7.55469 11.418 7.55469 13.4248C7.55469 14.2817 8.1626 14.6772 9.42969 14.6772ZM9.18799 13.4102C9.01221 13.4102 8.94629 13.3589 8.94629 13.2197C8.94629 12.0991 10.7627 9.98975 13.9927 9.98975C17.23 9.98975 19.0464 12.0991 19.0464 13.2197C19.0464 13.3589 18.9805 13.4102 18.8047 13.4102H9.18799Z"
              fill="#1C1C1E"
            />
          </svg>
        </div>

        <div class="absolute inset-x-0 h-32 -bottom-32 pr-[36px]" v-if="searchInput != ''">
          <div class="max-h-32 overflow-auto w-full mt-2 rounded-[6px] flex flex-col bg-systemWhite-light">
            <div class="group" v-for="user in searchResult" :key="user.id" @click="onSelectUser(user)">
              <div class="py-1 px-2 flex justify-between cursor-pointer hover:bg-systemBackground-lightSecondary">
                <span class="">{{ user.name }}</span>
                <span class="textDescription text-labelColor-light-secondary">({{ user.id }})</span>
              </div>
              <div class="h-[1px] bg-labelColor-light-tertiary w-full group-last:hidden"></div>
            </div>
          </div>
        </div>
        <!-- <div
          class="absolute max-h-40 overflow-auto w-full mt-2 rounded-[6px] flex flex-col bg-systemWhite-light"
          v-if="searchInput != ''"
        ></div> -->
      </div>
    </template>
    <template #side-body>
      <div class="flex flex-col justify-between h-full select-none cursor-default">
        <div>
          <div class="group mx-2" v-for="contact in messageStore.messageList" @click="onSelectContact(contact)">
            <div
              v-if="contact.user != undefined"
              class="flex items-center px-3 py-2 rounded-[6px]"
              :class="{ 'bg-systemBlue-light text-systemWhite-light': contact.user.id == messageStore.receiver?.id }"
            >
              <user-avatar class="shrink-0" :user="contact.user"></user-avatar>
              <div class="ml-3 text-sm flex flex-col items-start flex-grow h-12">
                <div class="flex justify-between w-full">
                  <div class="text-sm">
                    <span class="">{{ contact.user.name }}</span>
                    <span
                      :class="contact.user.id == messageStore.receiver?.id ? 'text-systemWhite-light' : 'text-labelColor-light-secondary'"
                    >
                      ({{ contact.user.id }})</span
                    >
                  </div>
                  <div
                    v-if="contact.messages && contact.messages.length > 0"
                    :class="contact.user.id == messageStore.receiver?.id ? 'text-systemWhite-light' : 'text-labelColor-light-secondary'"
                  >
                    {{ new Date(Date.parse(contact.messages[contact.messages.length - 1].time)).toLocaleDateString() }}
                  </div>
                </div>
                <div
                  v-if="contact.messages && contact.messages.length > 0"
                  class="text-xs overflow-hidden line-clamp-2 text-start"
                  :class="contact.user.id == messageStore.receiver?.id ? 'text-systemWhite-light' : 'text-[#6f6f6f]'"
                >
                  {{ contact.messages[contact.messages.length - 1].content }}
                </div>
              </div>
            </div>
            <div v-else class="flex px-3 py-2">
              <div class="flex -space-x-4 overflow-hidden">
                <user-avatar
                  class="border-2 border-labelColor-light-secondary"
                  v-for="user in contact.group?.users"
                  :user="user"
                ></user-avatar>
              </div>
            </div>
            <div class="h-[1px] ml-9 bg-labelColor-light-tertiary group-last:hidden"></div>
          </div>
        </div>
        <div class="flex-shrink-0 flex items-end justify-between pb-3 px-4">
          <div class="flex items-end p-1">
            <span class="font-semibold text-textBlack-light/80 ml-1">{{ userStore.user?.name }}</span>
            <span class="ml-1 textDescription text-labelColor-light-secondary">({{ userStore.user?.id }})</span>
          </div>
          <button
            class="text-center text-sm font-semibold rounded text-labelColor-light-secondary px-1 p-0.5 mb-0.5 hover:bg-gray-400/40"
            @click="userStore.logout"
          >
            登出
          </button>
        </div>
      </div>
    </template>
    <template #main>
      <message-view ref="messageView" v-if="startGroup == false && messageStore.isSelected"></message-view>
      <group-message-starter-vue v-else-if="startGroup" @started="onGroupStarted"></group-message-starter-vue>
    </template>
  </base-layout>
</template>
