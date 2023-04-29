import { reactive, ref, computed, useCssModule } from "vue"
import { defineStore } from "pinia"
import { model } from "../../wailsjs/go/models"
import { GetGroupById, GetUserById, GetMessages } from "../../wailsjs/go/main/App"
import { userAccountStore } from "./account"
import { updateLanguageServiceSourceFile } from "typescript"

const MessageStore = reactive({})

export { MessageStore }

export type messageList = {
  user?: model.User
  group?: model.Group
  messages: model.Message[]
}
// (model.User | { id: number }) & { messages: model.Message[] }
export type contactList = messageList[]

export const useMessageStore = defineStore("counter", () => {
  const receiver = ref<model.User>()
  const group = ref<model.Group>()
  const messageList = ref<contactList>([])
  const userStore = userAccountStore()

  const userOrGroup = (
    contact: model.User | model.Group,
    userCallBack: (u: model.User) => void,
    groupCallBack: (u: model.Group) => void
  ) => {
    if ("name" in contact) {
      userCallBack(contact)
    } else {
      groupCallBack(contact as model.Group)
    }
  }

  const selectContact = (contact: model.User | model.Group) => {
    userOrGroup(
      contact,
      u => (receiver.value = u),
      g => (group.value = g)
    )

    appendMessage(contact, [])
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
      return {
        messages: [],
      }
    }
  })

  const initMessageList = async () => {
    const messages = await GetMessages()
    await userStore.updateUserList()
    const userDataSource = async (userId: number) => {
      return userStore.userList.find(v => v.id == userId) as model.User
    }
    for (const message of messages) {
      await onReceiveMessage(message, userDataSource)
    }
  }

  const onReceiveMessage = async (message: model.Message, userDataSource: (userId: number) => Promise<model.User>) => {
    if (message.groupId > 0) {
      // a group message
      console.log("group message")
      const group = await GetGroupById(message.groupId)
      appendMessage(group, [message])
    } else {
      const userId = message.from == userStore.user?.id ? message.to : message.from
      const user = await userDataSource(userId)
      appendMessage(user, [message])
    }
  }

  const appendMessage = (contact: model.User | model.Group, messages: model.Message[]) => {
    const index = messageList.value.findIndex(v => {
      if (v.user != undefined && v.group == undefined) {
        return v.user.id == contact.id
      } else if (v.group != undefined) {
        return v.group.id == contact.id
      }
    })
    console.log(index)
    if (index == -1) {
      messageList.value.push({
        user: "name" in contact ? contact : undefined,
        group: "name" in contact ? undefined : (contact as model.Group),
        messages: [...messages],
      })
    } else {
      messageList.value[index].messages.push(...messages)
    }
  }

  const mock = () => {
    // messageList.value = []
    // messageList.value.push({
    //   user: {
    //     id: 7,
    //     name: "Alice",
    //     address: "addr",
    //     port: 1234,
    //   },
    //   messages: [
    //     {
    //       from: 1,
    //       to: 7,
    //       content: "hello",
    //       time: "123",
    //       groupId: -1,
    //     },
    //     {
    //       from: 7,
    //       to: 1,
    //       content: "hello again",
    //       time: "123",
    //       groupId: -1,
    //     },
    //   ],
    // })
    // messageList.value.push({
    //   user: {
    //     id: 9,
    //     name: "Mando",
    //     address: "addr",
    //     port: 1234,
    //   },
    //   messages: [
    //     {
    //       from: 1,
    //       to: 9,
    //       content: "Greeting!",
    //       time: "123",
    //       groupId: -1,
    //     },
    //     {
    //       from: 9,
    //       to: 1,
    //       content: "This is the way.",
    //       time: "123",
    //       groupId: -1,
    //     },
    //   ],
    // })
  }

  return {
    receiver,
    group,
    initMessageList,
    selectContact,
    onReceiveMessage,
    isSelected,
    messageList,
    currentMessageList,
    mock,
  }
})
