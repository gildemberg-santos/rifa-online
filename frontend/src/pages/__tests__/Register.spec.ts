import { describe, it, expect, beforeEach, vi } from "vitest"
import { mount } from "@vue/test-utils"
import { setActivePinia, createPinia } from "pinia"
import Register from "../Register.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

const mockNotifyShow = vi.fn()
const mockRegister = vi.fn()
const mockVerifyEmail = vi.fn()
const mockResendCode = vi.fn()

vi.mock("../../stores/auth", () => ({
  useAuthStore: () => ({
    register: mockRegister,
    verifyEmail: mockVerifyEmail,
    resendCode: mockResendCode,
  }),
}))

vi.mock("../../composables/useNotification", () => ({
  useNotification: () => ({ show: mockNotifyShow }),
}))

vi.mock("vue-router", () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({}),
}))

describe("Register", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it("renders the registration form", () => {
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.find("form").exists()).toBe(true)
    expect(wrapper.text()).toContain("Criar Conta")
    expect(wrapper.text()).toContain("Cadastre-se para criar rifas")
  })

  it("has name, email and password inputs", () => {
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.find('input[type="text"]').exists()).toBe(true)
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
  })

  it("has submit button", () => {
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    const button = wrapper.find('button[type="submit"]')
    expect(button.exists()).toBe(true)
    expect(button.text()).toContain("Criar Conta")
  })

  it("has link to login page", () => {
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("J\u00e1 tem conta?")
    expect(wrapper.text()).toContain("Fazer login")
  })

  it("shows verification step after successful register", async () => {
    mockRegister.mockResolvedValueOnce(undefined)
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })

    await wrapper.find('input[type="text"]').setValue("John")
    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find('input[type="checkbox"]').setValue(true)
    await wrapper.find('form').trigger("submit.prevent")

    expect(mockRegister).toHaveBeenCalledWith("John", "john@test.com", "123456")
    expect(mockNotifyShow).toHaveBeenCalledWith(
      "Codigo de verificacao enviado para o email",
      "info",
    )

    expect(wrapper.text()).toContain("Verifique seu Email")
    expect(wrapper.text()).toContain("john@test.com")
  })

  it("shows verification code input after registration", async () => {
    mockRegister.mockResolvedValueOnce(undefined)
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })

    await wrapper.find('input[type="text"]').setValue("John")
    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find('input[type="checkbox"]').setValue(true)
    await wrapper.find('form').trigger("submit.prevent")

    const verifyInput = wrapper.find('input[placeholder="000000"]')
    expect(verifyInput.exists()).toBe(true)
    expect(wrapper.text()).toContain("Verificar C\u00f3digo")
    expect(wrapper.text()).toContain("Reenviar c\u00f3digo")
  })

  it("calls verifyEmail on verify button click", async () => {
    mockRegister.mockResolvedValueOnce(undefined)
    mockVerifyEmail.mockResolvedValueOnce(undefined)
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    await wrapper.find('input[type="text"]').setValue("John")
    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find('input[type="checkbox"]').setValue(true)
    await wrapper.find('form').trigger("submit.prevent")

    const codeInput = wrapper.find('input[placeholder="000000"]')
    await codeInput.setValue("123456")

    const verifyBtn = wrapper.findAll("button")[0]
    await verifyBtn.trigger("click")

    expect(mockVerifyEmail).toHaveBeenCalledWith("john@test.com", "123456")
  })

  it("calls resendCode on resend button click", async () => {
    mockRegister.mockResolvedValueOnce(undefined)
    mockResendCode.mockResolvedValueOnce(undefined)
    const wrapper = mount(Register, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    await wrapper.find('input[type="text"]').setValue("John")
    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find('input[type="checkbox"]').setValue(true)
    await wrapper.find('form').trigger("submit.prevent")

    const buttons = wrapper.findAll("button")
    const resendBtn = buttons[1]
    await resendBtn.trigger("click")

    expect(mockResendCode).toHaveBeenCalledWith("john@test.com")
  })
})
