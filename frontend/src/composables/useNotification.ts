import { ref } from "vue"

export type NotificationType = "success" | "error" | "info"

const visible = ref(false)
const message = ref("")
const type = ref<NotificationType>("info")

let timeout: ReturnType<typeof setTimeout> | null = null

export function useNotification() {
  function show(msg: string, t: NotificationType = "info", duration = 4000) {
    if (timeout) clearTimeout(timeout)
    message.value = msg
    type.value = t
    visible.value = true
    timeout = setTimeout(() => { visible.value = false }, duration)
  }

  function close() {
    visible.value = false
    if (timeout) clearTimeout(timeout)
  }

  return { visible, message, type, show, close }
}
