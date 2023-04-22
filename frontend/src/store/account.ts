import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"
import { ref } from "vue"
import router from "../router"

export const userAccountStore = defineStore("userAccount", () => {
  const user = ref<model.User>()
  const token = ref<string>()

  const logout = () => {
    user.value = {}
    router.push("/login")
    router.push("/login")
  }

  return {
    user,
    token,
    logout,
  }
})
