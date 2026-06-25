import { describe, it, expect } from "vitest"
import { mount } from "@vue/test-utils"
import LoadingSpinner from "../LoadingSpinner.vue"

describe("LoadingSpinner", () => {
  it("renders by default (md size)", () => {
    const wrapper = mount(LoadingSpinner)
    const spinner = wrapper.find(".animate-spin")
    expect(spinner.exists()).toBe(true)
    expect(spinner.classes()).toContain("w-8")
    expect(spinner.classes()).toContain("h-8")
  })

  it("renders with sm size", () => {
    const wrapper = mount(LoadingSpinner, {
      props: { size: "sm" },
    })
    const spinner = wrapper.find(".animate-spin")
    expect(spinner.classes()).toContain("w-5")
    expect(spinner.classes()).toContain("h-5")
  })

  it("renders with lg size", () => {
    const wrapper = mount(LoadingSpinner, {
      props: { size: "lg" },
    })
    const spinner = wrapper.find(".animate-spin")
    expect(spinner.classes()).toContain("w-12")
    expect(spinner.classes()).toContain("h-12")
  })
})
