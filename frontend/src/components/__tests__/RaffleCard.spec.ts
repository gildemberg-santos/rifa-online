import { describe, it, expect } from "vitest"
import { mount } from "@vue/test-utils"
import RaffleCard from "../RaffleCard.vue"

const RouterLinkStub = {
  name: "RouterLink",
  template: '<a :href="to"><slot /></a>',
  props: { to: { type: [String, Object], required: true } },
}

describe("RaffleCard", () => {
  const raffle = {
    id: "abc123",
    title: "Test Raffle",
    description: "A test raffle description",
    ticketPrice: 500,
    maxNumbers: 100,
    drawDate: "2026-12-25T00:00:00Z",
    status: "ACTIVE",
  }

  it("renders title and description", () => {
    const wrapper = mount(RaffleCard, {
      props: { raffle },
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("Test Raffle")
    expect(wrapper.text()).toContain("A test raffle description")
  })

  it("formats price correctly", () => {
    const wrapper = mount(RaffleCard, {
      props: { raffle },
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("R$ 5.00")
  })

  it("shows max numbers count", () => {
    const wrapper = mount(RaffleCard, {
      props: { raffle },
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    expect(wrapper.text()).toContain("100")
  })

  it("links to raffle detail", () => {
    const wrapper = mount(RaffleCard, {
      props: { raffle },
      global: { stubs: { RouterLink: RouterLinkStub } },
    })
    const link = wrapper.find("a")
    expect(link.attributes("href")).toBe("/raffles/abc123")
  })
})
