import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"
import { ref } from "vue"
export const userAccountStore = defineStore("userAccount", () => {
  const user = ref<model.User>()
  const token = ref<string>()

  return {
    user,
    token,
  }
})
