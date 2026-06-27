import { describe, it, expect, beforeEach, vi } from "vitest"
import { mount } from "@vue/test-utils"
import { setActivePinia, createPinia } from "pinia"
import Login from "../Login.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

const mockNotifyShow = vi.fn()
const mockLogin = vi.fn()
const mockVerifyEmail = vi.fn()
const mockResendCode = vi.fn()

vi.mock("../../stores/auth", () => ({
  useAuthStore: () => ({
    login: mockLogin,
    verifyEmail: mockVerifyEmail,
    resendCode: mockResendCode,
  }),
}))

vi.mock("../../composables/useNotification", () => ({
  useNotification: () => ({ show: mockNotifyShow }),
}))

vi.mock("vue-router", () => ({
  useRouter: () => ({ push: vi.fn() }),
  useRoute: () => ({ query: {} }),
}))

describe("Login", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it("renders the login form", () => {
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.find("form").exists()).toBe(true)
    expect(wrapper.text()).toContain("Bem-vindo")
    expect(wrapper.text()).toContain("Fa\u00e7a login para acessar o painel")
  })

  it("has email and password inputs", () => {
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.find('input[type="email"]').exists()).toBe(true)
    expect(wrapper.find('input[type="password"]').exists()).toBe(true)
  })

  it("has submit button", () => {
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    const button = wrapper.find('button[type="submit"]')
    expect(button.exists()).toBe(true)
    expect(button.text()).toContain("Entrar")
  })

  it("has link to register page", () => {
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("N\u00e3o tem conta?")
    expect(wrapper.text()).toContain("Criar conta")
  })

  it("shows verification code form when email not verified", async () => {
    mockLogin.mockRejectedValueOnce(new Error("email not verified"))
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })

    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find("form").trigger("submit.prevent")

    expect(wrapper.text()).toContain("Codigo de verificacao")
    expect(wrapper.text()).toContain("Verificar e Entrar")
    expect(wrapper.text()).toContain("Reenviar codigo")
    expect(mockNotifyShow).toHaveBeenCalledWith(
      "Email nao verificado. Informe o codigo enviado no cadastro.",
      "info",
    )
  })

  it("calls verifyEmail when email not verified and code is entered", async () => {
    mockLogin.mockRejectedValueOnce(new Error("email not verified"))
    mockVerifyEmail.mockResolvedValueOnce(undefined)
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })

    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find("form").trigger("submit.prevent")

    const codeInput = wrapper.find('input[placeholder="000000"]')
    await codeInput.setValue("654321")

    const verifyBtn = wrapper.findAll("button")[0]
    await verifyBtn.trigger("click")

    expect(mockVerifyEmail).toHaveBeenCalledWith("john@test.com", "654321")
  })

  it("calls resendCode on resend button", async () => {
    mockLogin.mockRejectedValueOnce(new Error("email not verified"))
    mockResendCode.mockResolvedValueOnce(undefined)
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })

    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find("form").trigger("submit.prevent")

    const resendBtn = wrapper.findAll("button")[1]
    await resendBtn.trigger("click")

    expect(mockResendCode).toHaveBeenCalledWith("john@test.com")
  })

  it("calls api.post with correct credentials", async () => {
    mockLogin.mockResolvedValueOnce(undefined)
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })

    await wrapper.find('input[type="email"]').setValue("john@test.com")
    await wrapper.find('input[type="password"]').setValue("123456")
    await wrapper.find("form").trigger("submit.prevent")

    expect(mockLogin).toHaveBeenCalledWith("john@test.com", "123456")
  })
})
