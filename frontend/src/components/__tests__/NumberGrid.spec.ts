import { describe, it, expect } from "vitest"
import { mount } from "@vue/test-utils"
import NumberGrid from "../NumberGrid.vue"

describe("NumberGrid", () => {
  const tickets = [
    { number: 1, status: "AVAILABLE" },
    { number: 2, status: "RESERVED" },
    { number: 3, status: "PAID" },
  ]

  it("renders all tickets", () => {
    const wrapper = mount(NumberGrid, {
      props: { tickets, selectedNumbers: [] },
    })
    const buttons = wrapper.findAll("button")
    expect(buttons).toHaveLength(3)
  })

  it("shows correct number text", () => {
    const wrapper = mount(NumberGrid, {
      props: { tickets, selectedNumbers: [] },
    })
    expect(wrapper.text()).toContain("1")
    expect(wrapper.text()).toContain("2")
    expect(wrapper.text()).toContain("3")
  })

  it("disables non-available tickets", () => {
    const wrapper = mount(NumberGrid, {
      props: { tickets, selectedNumbers: [] },
    })
    const buttons = wrapper.findAll("button")
    expect(buttons[0].attributes("disabled")).toBeUndefined()
    expect(buttons[1].attributes("disabled")).toBeDefined()
    expect(buttons[2].attributes("disabled")).toBeDefined()
  })

  it("emits toggle event when clicking available ticket", async () => {
    const wrapper = mount(NumberGrid, {
      props: { tickets, selectedNumbers: [] },
    })
    await wrapper.findAll("button")[0].trigger("click")
    expect(wrapper.emitted("toggle")).toBeTruthy()
    expect(wrapper.emitted("toggle")![0]).toEqual([1])
  })

  it("highlights selected numbers", () => {
    const wrapper = mount(NumberGrid, {
      props: { tickets, selectedNumbers: [1] },
    })
    const button = wrapper.findAll("button")[0]
    expect(button.classes()).toContain("bg-indigo-600")
  })
})
