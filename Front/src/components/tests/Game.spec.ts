import { mount } from '@vue/test-utils'
import Game from '../Game.vue'
import { describe, it, expect, vi } from 'vitest'
import { nextTick } from 'vue'

function flushPromises() {
  return new Promise(resolve => setTimeout(resolve))
}

describe('Game.vue', () => {
  it('affiche le message initial', () => {
    const wrapper = mount(Game)
    expect(wrapper.text()).toContain('Cliquez sur "Commencer" pour débuter la partie')
  })

  it('démarre une nouvelle partie au clic sur "Commencer"', async () => {
    const wrapper = mount(Game)

    global.fetch = vi.fn(() =>
      Promise.resolve({
        ok: true,
        text: () => Promise.resolve("1")
      })
    ) as any

    const button = wrapper.find('button')
    await button.trigger('click')

    await flushPromises()
    await nextTick()
    await new Promise(res => setTimeout(res, 100)) 
    expect(wrapper.text()).toContain('Regardez bien où est la balle')
  })
})
