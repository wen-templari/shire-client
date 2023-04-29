import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"
import { ref } from "vue"
import router from "../router"
import { GetUsers } from "../../wailsjs/go/main/App"

export const userAccountStore = defineStore("userAccount", () => {
  const user = ref<model.User>()
  const userList = ref<model.User[]>([])
  const token = ref<string>()

  const logout = () => {
    user.value = {}
    router.push("/login")
    router.push("/login")
  }

  const updateUserList = async () => {
    userList.value = await GetUsers()
  }

  return {
    user,
    token,
    logout,
    userList,
    updateUserList,
  }
})
