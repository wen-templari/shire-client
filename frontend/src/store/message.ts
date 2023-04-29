import { reactive, ref, computed } from "vue"
import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"

const MessageStore = reactive({})

export { MessageStore }

export type messageList = {
  user?: model.User
  group?: { id: number }
  messages?: model.Message[]
}
// (model.User | { id: number }) & { messages: model.Message[] }
export type contactList = messageList[]

export const useMessageStore = defineStore("counter", () => {
  const receiver = ref<model.User>()
  const group = ref<number>()
  const messageList = ref<contactList>([])
  const selectContact = (contact: messageList) => {
    receiver.value = contact.user
    group.value = contact.group?.id
  }
  const isSelected = computed(() => {
    return receiver.value != undefined || group.value != undefined
  })
  const currentMessageList = computed(() => {
    if (receiver.value != undefined && receiver.value.id != undefined) {
      return messageList.value.find(v => {
        if (v.user != undefined && v.group == undefined) {
          return v.user.id == receiver.value?.id
        }
      })
    } else if (group.value) {
      return messageList.value.find(v => {
        if (v.group != undefined) {
          return v.group.id == group.value
        }
      })
    } else {
      return [] as messageList
    }
  })

  const mock = () => {
    messageList.value = []
    messageList.value.push({
      user: {
        id: 7,
        name: "Alice",
        address: "addr",
        port: 1234,
      },
      messages: [
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
      ],
    })
    messageList.value.push({
      user: {
        id: 9,
        name: "Mando",
        address: "addr",
        port: 1234,
      },
      messages: [
        {
          from: 1,
          to: 9,
          content: "Greeting!",
          time: "123",
          groupId: -1,
        },
        {
          from: 9,
          to: 1,
          content: "This is the way.",
          time: "123",
          groupId: -1,
        },
      ],
    })
  }

  return {
    receiver,
    group,
    selectContact,
    isSelected,
    messageList,
    currentMessageList,
    mock,
  }
})
