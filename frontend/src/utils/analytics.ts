declare global {
  interface Window {
    gtag: (command: string, ...args: any[]) => void
    dataLayer: any[]
  }
}

export function sendEvent(name: string, params?: Record<string, unknown>) {
  if (typeof window !== "undefined" && typeof window.gtag === "function") {
    window.gtag("event", name, params)
  }
}

export function sendUserAction(action: string, label?: string, value?: number) {
  sendEvent(action, {
    event_category: "user_action",
    event_label: label,
    value,
  })
}

export function sendRaffleAction(action: string, raffleId: string, label?: string) {
  sendEvent(`raffle_${action}`, {
    raffle_id: raffleId,
    event_label: label,
  })
}

export function sendPaymentAction(action: string, amount?: number, label?: string) {
  sendEvent(`payment_${action}`, {
    amount,
    event_label: label,
  })
}

export function sendAdminAction(action: string, targetId: string, label?: string) {
  sendEvent(`admin_${action}`, {
    target_id: targetId,
    event_label: label,
  })
}
