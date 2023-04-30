import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"
import { ref } from "vue"
import router from "../router"
import { GetUsers, Logout } from "../../wailsjs/go/main/App"

export const userAccountStore = defineStore("userAccount", () => {
  const user = ref<model.User>()
  const userList = ref<model.User[]>([])
  const token = ref<string>()

  const setToken = (t: string) => {
    localStorage.setItem("token", t)
    token.value = t
  }
  const getToken = () => {
    if (token.value == undefined) {
      token.value = localStorage.getItem("token") as string | undefined
    }
    return token.value
  }

  const logout = async () => {
    user.value = {}
    await Logout()
    router.push("/login")
    router.push("/login")
  }

  const updateUserList = async () => {
    userList.value = await GetUsers()
  }

  return {
    user,
    setToken,
    getToken,
    logout,
    userList,
    updateUserList,
  }
})
