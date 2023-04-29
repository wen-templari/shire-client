<script setup lang="ts">
import LayoutBase from "@/layout/LayoutBase.vue"
import InputBase from "@/components/Input/InputBase.vue"
import { ref } from "@vue/reactivity"
import { Login, Register } from "../../wailsjs/go/main/App"
import BaseLayout from "../layout/BaseLayout.vue"
import { userAccountStore } from "../store/account"
import router from "../router"

const id = ref<number>()
const password = ref<string>()

const name = ref("")
const isIDValid = ref("")
const isPasswordValid = ref("")

const store = userAccountStore()

const login = () => {
  if (id.value == undefined || password.value == undefined) {
    return
  }
  console.log(id.value, password.value)
  Login(id.value, password.value).then(res => {
    store.user = res
    router.push("/")
  })
}

const registerSwitch = ref(false)
const switchRegister = (b: boolean) => {
  registerSwitch.value = b
}

const register = () => {
  if (id.value == undefined || password.value == undefined) {
    return
  }
  Register(name.value, password.value).then(res => {
    store.user = res
    router.push("/")
  })
}
</script>
<template>
  <base-layout>
    <template #side-body>
      <div class="flex flex-col items-center h-full">
        <div class="h-1/4 flex flex-col justify-center">
          <div class="text-center text-3xl font-semibold tracking-wide mt-10">Shire Client</div>
          <div class="text-center text-2xl mt-6" v-if="!registerSwitch">登入来开始聊天</div>
          <div class="text-center text-2xl mt-6" v-if="registerSwitch">注册账号</div>
        </div>
        <form class="flex flex-col p-5 mt-15 w-full" v-if="!registerSwitch">
          <InputBase class="my-2 w-full" v-model:content.number="id" :passWarning="isIDValid" placeholder="ID" />
          <InputBase class="my-2 w-full" v-model:content="password" :passWarning="isPasswordValid" placeholder="密码" type="password" />
          <button
            class="mt-6 rounded text-[18px] px-[16px] py-[12px] text-systemBlue-light bg-systemWhite-light w-full"
            @click.prevent="login"
          >
            登入
          </button>
          <div class="mt-5 flex flex-col items-center text-center">
            <div class="w-20 text-indigo-700 text-sm font-semibold cursor-pointer" @click="switchRegister(true)">没有账号?</div>
          </div>
        </form>
        <form @submit.prevent="register" class="flex flex-col p-5 mt-15 w-full" v-if="registerSwitch">
          <InputBase class="my-2 w-full" v-model:content="name" placeholder="昵称" />
          <InputBase class="my-2 w-full" v-model:content="password" placeholder="密码" type="password" />
          <button class="mt-6 btn btnPrimary w-full" type="submit">注册</button>
          <div class="mt-5 flex flex-col items-center text-center">
            <div class="w-20 text-indigo-700 text-sm font-semibold cursor-pointer" @click="switchRegister(false)">有账号了</div>
          </div>
        </form>
        <div></div>
      </div>
    </template>
  </base-layout>
</template>
