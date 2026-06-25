import { describe, it, expect } from "vitest"
import { mount } from "@vue/test-utils"
import Footer from "../Footer.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("Footer", () => {
  it("renders the footer text", () => {
    const wrapper = mount(Footer, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Rifa Online")
    expect(wrapper.text()).toContain("Plataforma para criação e gerenciamento de rifas online.")
  })

  it("contains router-links to home and login", () => {
    const wrapper = mount(Footer, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    const links = wrapper.findAll("a")
    const hrefs = links.map((l) => l.attributes("href"))
    expect(hrefs).toContain("/")
    expect(hrefs).toContain("/login")
  })

  it("shows current year in copyright", () => {
    const wrapper = mount(Footer, {
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    const year = new Date().getFullYear().toString()
    expect(wrapper.text()).toContain(year)
  })
})
