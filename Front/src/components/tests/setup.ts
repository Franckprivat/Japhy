// jest-dom is not available, so we skip importing it

globalThis.Audio = class {
  play() {}
  pause() {}
  currentTime = 0
} as any

globalThis.localStorage = window.localStorage
