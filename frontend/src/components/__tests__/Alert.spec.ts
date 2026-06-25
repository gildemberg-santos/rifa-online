import { describe, it, expect } from "vitest"
import { mount } from "@vue/test-utils"
import Alert from "../Alert.vue"

describe("Alert", () => {
  it("renders message text", () => {
    const wrapper = mount(Alert, {
      props: { message: "Test message" },
    })
    expect(wrapper.text()).toContain("Test message")
  })

  it("renders with default type (info styling)", () => {
    const wrapper = mount(Alert, {
      props: { message: "Info" },
    })
    const div = wrapper.find("div")
    expect(div.classes()).toContain("bg-blue-50")
    expect(div.classes()).toContain("text-blue-700")
  })

  it("renders with success type", () => {
    const wrapper = mount(Alert, {
      props: { message: "Success", type: "success" },
    })
    const div = wrapper.find("div")
    expect(div.classes()).toContain("bg-green-50")
    expect(div.classes()).toContain("text-green-700")
  })

  it("renders with error type", () => {
    const wrapper = mount(Alert, {
      props: { message: "Error", type: "error" },
    })
    const div = wrapper.find("div")
    expect(div.classes()).toContain("bg-red-50")
    expect(div.classes()).toContain("text-red-700")
  })

  it("renders with warning type", () => {
    const wrapper = mount(Alert, {
      props: { message: "Warning", type: "warning" },
    })
    const div = wrapper.find("div")
    expect(div.classes()).toContain("bg-yellow-50")
    expect(div.classes()).toContain("text-yellow-700")
  })
})
