import { reactive, ref, computed } from "vue"
import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"

const MessageStore = reactive({})

export { MessageStore }

export const useMessageStore = defineStore("counter", () => {
  const receiver = ref<model.User>()
  const group = ref<number>()
  const messageList = ref(new Map<number, model.Message[]>())

  const currentMessageList = computed(() => {
    if (receiver.value) {
      console.log(receiver.value)
      return messageList.value.get(receiver.value.id)
    } else if (group.value) {
      return messageList.value.get(group.value)
    }
  })

  const mock = () => {
    receiver.value = {
      id: 7,
      name: "bob",
      address: "addr",
      port: 1234,
    }

    messageList.value.set(7, [
      {
        from: 1,
        to: 7,
        content: "hello",
        time: "123",
        groupId: -1,
      },
      {
        from: 7,
        to: 1,
        content: "hello again",
        time: "123",
        groupId: -1,
      },
    ])
  }

  return {
    receiver,
    group,
    messageList,
    currentMessageList,
    mock,
  }
})
