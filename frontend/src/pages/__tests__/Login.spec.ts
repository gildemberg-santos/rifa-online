import { describe, it, expect, beforeEach } from "vitest"
import { mount } from "@vue/test-utils"
import { setActivePinia, createPinia } from "pinia"
import Login from "../Login.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("Login", () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
  })

  it("renders the login form", () => {
    const wrapper = mount(Login, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.find("form").exists()).toBe(true)
    expect(wrapper.text()).toContain("Bem-vindo")
    expect(wrapper.text()).toContain("Faça login para acessar o painel")
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
    expect(wrapper.text()).toContain("Não tem conta?")
    expect(wrapper.text()).toContain("Criar conta")
  })
})
