<script setup lang="ts">
import { model } from "../../../wailsjs/go/models"
import { computed, onMounted, ref } from "vue"
import UserAvatar from "../../components/User/UserAvatar.vue"
import { userAccountStore } from "../../store/account"

const groupUserList = ref<model.User[]>([])
const userStore = userAccountStore()

const searchInput = ref("")
const searchResult = ref<model.User[]>([])
const onSearchInput = (e: Event) => {
  const target = e.target as HTMLInputElement
  const towerValue = target.value.toLowerCase()
  searchResult.value = userStore.userList.filter(
    user => user.name?.toLocaleLowerCase().indexOf(towerValue) != -1 && groupUserList.value.findIndex(u => u.id == user.id) == -1
  )
}

const onSelectUser = (user: model.User) => {
  searchInput.value = ""
  groupUserList.value.push(user)
}

const onDeleteUser = (user: model.User) => {
  const index = groupUserList.value.findIndex(u => u.id == user.id)
  groupUserList.value.splice(index, 1)
}
onMounted(() => {
  userStore.updateUserList()
  userStore.updateUserList()
  groupUserList.value.push(userStore.user as model.User)
})

const startGroup = () => {}
</script>
<template>
  <div class="flex flex-col items-center justify-center h-full">
    <div class="flex gap-3">
      <div v-for="user in groupUserList" :key="user.id">
        <UserAvatar :user="user"></UserAvatar>
      </div>
      <UserAvatar v-if="groupUserList.length < 3" v-for="item in 3 - groupUserList.length" :user="{}"></UserAvatar>
    </div>
    <div class="mt-2">在开始前，添加群聊中的成员</div>
    <div class="mt-4">
      <div class="relative">
        <input
          class="w-full rounded-[6px] border border-borderGrey-light h-[30px] text-sm px-3 placeholder:text-labelColor-light-secondary"
          type="text"
          placeholder="search users"
          v-model="searchInput"
          :disabled="groupUserList.length > 4"
          @input="onSearchInput"
        />
        <div class="absolute inset-x-0 h-32 -bottom-32" v-if="searchInput != ''">
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
      </div>
    </div>
    <button
      class="mt-40 rounded text-[18px] px-[12px] py-[6px] text-systemBlue-light bg-systemWhite-light w-52"
      @click.prevent="startGroup"
      :disabled="groupUserList.length < 3"
    >
      开始
    </button>

    <!-- :class="groupUserList.length < 3 ? 'bg-labelColor-light-tertiary text-labelColor-light-secondary' : 'cursor-pointer'" -->
  </div>
</template>
