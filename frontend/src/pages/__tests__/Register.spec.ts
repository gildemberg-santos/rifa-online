import { describe, it, expect, beforeEach } from "vitest"
import { mount } from "@vue/test-utils"
import { setActivePinia, createPinia } from "pinia"
import Register from "../Register.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("Register", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
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
    expect(wrapper.text()).toContain("Já tem conta?")
    expect(wrapper.text()).toContain("Fazer login")
  })
})
