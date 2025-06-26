<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

defineProps<{ msg?: string }>()

const shells = ref([
  { hasPearl: false, revealed: false, position: 0 },
  { hasPearl: false, revealed: false, position: 1 },
  { hasPearl: false, revealed: false, position: 2 }
])

const gameState = ref<'initial' | 'showing' | 'shuffling' | 'playing' | 'revealed'>('initial')
const message = ref('Cliquez sur "Commencer" pour débuter la partie !')
const isLoading = ref(false)
const isShuffling = ref(false)
const wins = ref(0)
const losses = ref(0)
const streak = ref(0)
const correctShellIndex = ref(-1)
const liftedIndexes = ref<number[]>([])
const shakeIndex = ref(-1)
let lastPearlPosition = -1

const messageClass = computed(() => {
  if (isLoading.value) return 'loading'
  if (message.value.includes('Bravo')) return 'success'
  if (message.value.includes('Dommage')) return 'error'
  return ''
})

const isWin = computed(() => message.value.includes('Bravo'))

onMounted(() => {
  wins.value = parseInt(localStorage.getItem('bonneteau-wins') || '0')
  losses.value = parseInt(localStorage.getItem('bonneteau-losses') || '0')
  streak.value = parseInt(localStorage.getItem('bonneteau-streak') || '0')
})

const saveScore = () => {
  localStorage.setItem('bonneteau-wins', wins.value.toString())
  localStorage.setItem('bonneteau-losses', losses.value.toString())
  localStorage.setItem('bonneteau-streak', streak.value.toString())
}

const resetShells = () => {
  shells.value.forEach((shell, idx) => {
    shell.hasPearl = false
    shell.revealed = false
    shell.position = idx
  })
}

const animateSwap = (idx1: number, idx2: number) => {
  const temp = shells.value[idx1].position
  shells.value[idx1].position = shells.value[idx2].position
  shells.value[idx2].position = temp
}

const getRandomPosition = async (): Promise<number> => {
  let newPosition = -1
  let tries = 0
  do {
    try {
      isLoading.value = true
      message.value = 'Récupération de la position aléatoire...'
      const response = await fetch('https://www.random.org/integers/?num=1&min=0&max=2&col=1&base=10&format=plain&rnd=new')
      const text = await response.text()
      newPosition = parseInt(text.trim())
    } catch {
      message.value = 'Erreur réseau, utilisation d\'un nombre local...'
      newPosition = Math.floor(Math.random() * 3)
    } finally {
      isLoading.value = false
    }
    tries++
  } while (newPosition === lastPearlPosition && tries < 5)
  lastPearlPosition = newPosition
  return newPosition
}

const startNewGame = async () => {
  shells.value.forEach((shell, idx) => (shell.position = idx))
  resetShells()
  gameState.value = 'showing'
  liftedIndexes.value = []
  shakeIndex.value = -1
  const newPosition = await getRandomPosition()
  shells.value[newPosition].hasPearl = true
  correctShellIndex.value = newPosition
  shells.value[newPosition].revealed = true
  liftedIndexes.value = [newPosition]
  message.value = `Regardez bien où est la balle 👀`

  setTimeout(() => {
    shells.value[newPosition].revealed = false
    liftedIndexes.value = []
    message.value = 'Mémorisez bien... mélange en cours 🌊'
    isShuffling.value = true
    gameState.value = 'shuffling'

    setTimeout(async () => {
      let swaps = Math.floor(Math.random() * 3) + 3
      for (let i = 0; i < swaps; i++) {
        const idx1 = Math.floor(Math.random() * 3)
        let idx2 = Math.floor(Math.random() * 3)
        while (idx1 === idx2) idx2 = Math.floor(Math.random() * 3)
        animateSwap(idx1, idx2)
        await new Promise(res => setTimeout(res, 250))
      }

      const logical = [...shells.value].sort((a, b) => a.position - b.position)
      shells.value.forEach((shell, i) => {
        shell.hasPearl = logical[i].hasPearl
      })
      correctShellIndex.value = shells.value.findIndex(s => s.hasPearl)

      shells.value.sort((a, b) => a.position - b.position)
      shells.value.forEach((shell, idx) => (shell.position = idx))

      isShuffling.value = false
      gameState.value = 'playing'
      message.value = 'À vous de jouer ! Cliquez sur un gobelet.'
    }, 1000)
  }, 2000)
}

const selectShell = (index: number) => {
  if (gameState.value !== 'playing' || isShuffling.value || isLoading.value) return
  gameState.value = 'revealed'
  liftedIndexes.value = [index]
  shells.value.forEach(shell => (shell.revealed = false))
  shells.value[index].revealed = true

  if (shells.value[index].hasPearl) {
    wins.value++
    streak.value++
    message.value = `🎉 Bravo ! Vous avez trouvé la balle ! Série de ${streak.value}`
  } else {
    losses.value++
    streak.value = 0
    shakeIndex.value = index
    message.value = '💔 Dommage ! La balle était ailleurs...'
    if (correctShellIndex.value !== index) {
      shells.value[correctShellIndex.value].revealed = true
      liftedIndexes.value = [index, correctShellIndex.value]
    }
  }

  setTimeout(() => {
    liftedIndexes.value = []
    shakeIndex.value = -1
    gameState.value = 'initial'
  }, 2000)

  saveScore()
}

const resetScore = () => {
  wins.value = 0
  losses.value = 0
  streak.value = 0
  saveScore()
  message.value = 'Scores remis à zéro !'
  gameState.value = 'initial'
}
</script>



<template>
  <div class="game-container">
    <h1>🎲 Jeu du Bonneteau</h1>
    <p :class="messageClass">{{ message }}</p>

    <div class="shells">
      <div
        v-for="(shell, index) in shells"
        :key="index"
        class="shell"
        :class="{
          revealed: shell.revealed,
          correct: isWin && liftedIndexes.includes(index) && shell.hasPearl && gameState === 'revealed',
          lifted: liftedIndexes.includes(index),
          shake: shakeIndex === index
        }"
        @click="gameState === 'playing' ? selectShell(index) : null"
        :style="{ transform: `translateX(${(shell.position - index) * 140}px)`, transition: isShuffling ? 'transform 0.25s' : 'transform 0.3s' }"
      >
        <svg width="120" height="160" viewBox="0 0 120 160" fill="none" xmlns="http://www.w3.org/2000/svg">
          <defs>
            <linearGradient id="cupBody" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#b0bec5"/>
              <stop offset="100%" stop-color="#607d8b"/>
            </linearGradient>
            <radialGradient id="cupShadow" cx="0.5" cy="0.8" r="0.5">
              <stop offset="0%" stop-color="#000" stop-opacity="0.18"/>
              <stop offset="100%" stop-color="#000" stop-opacity="0"/>
            </radialGradient>
          </defs>
          <ellipse cx="60" cy="150" rx="38" ry="10" fill="url(#cupShadow)"/>
          <path d="M25 30 Q60 10 95 30 L105 140 Q60 155 15 140 Z" fill="url(#cupBody)" stroke="none"/>
          <ellipse cx="60" cy="30" rx="35" ry="12" fill="#cfd8dc" stroke="none"/>
        </svg>
        <!-- Ici , on affiche la balle sous le cup levé qui a la perle, lors de la phase d'observation ou lors de la révélation -->
        <div v-if="shell.hasPearl && liftedIndexes.includes(index)" class="pearl">🔵</div>
      </div>
    </div>

    <div class="controls">
      <button @click="startNewGame" :disabled="isLoading || isShuffling || gameState === 'showing' || gameState === 'shuffling' || gameState === 'playing'">Commencer</button>
      <button @click="resetScore" :disabled="isShuffling || isLoading">Réinitialiser les scores</button>
    </div>

    <div class="scoreboard">
      ✅ Victoires : {{ wins }} | ❌ Défaites : {{ losses }} | 🔥 Série : {{ streak }}
    </div>
  </div>
</template>

<style scoped>
.game-container {
  max-width: 900px;
  margin: 5rem auto;
  text-align: center;
  padding: 4rem 2rem;
  border-radius: 20px;
  background: linear-gradient(to bottom right, #1c3b5a, #042c54);
  box-shadow: 0 15px 40px rgba(0, 0, 0, 0.25);
  font-family: 'Segoe UI', sans-serif;
}

h1 {
  font-size: 2.8rem;
  margin-bottom: 2rem;
  color: #00bfa5;
}

.message, p[class^="message"] {
  font-size: 1.4rem;
  font-weight: bold;
  margin: 2.5rem 0 3.5rem 0;
  color: #ffffff;
}

.success {
  color: #66bb6a;
}

.error {
  color: #ef5350;
}

.loading {
  color: #ffa726;
}

.shells {
  position: relative;
  width: 100%;
  max-width: 600px;
  min-width: 280px;
  height: 160px;
  margin: 4rem auto 3rem auto;
  display: flex;
  justify-content: center;
  align-items: flex-end;
  gap: 4rem;
}

@media (max-width: 700px) {
  .shells {
    gap: 1.2rem;
    max-width: 98vw;
    min-width: 0;
  }
  .shell {
    width: 80px;
    height: 110px;
  }
}

.shell {
  width: 120px;
  height: 160px;
  background: none;
  border-radius: 0;
  cursor: pointer;
  transition: transform 0.3s cubic-bezier(.4,2,.6,1), box-shadow 0.3s, filter 0.3s;
  box-shadow: 0 6px 18px rgba(0, 0, 0, 0.3);
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  overflow: visible;
}

.shell.revealed {
  transform: translateY(-30px) scale(1.1);
  z-index: 2;
}

.shell.correct {
  filter: drop-shadow(0 0 16px #00e676) brightness(1.15);
}

.shell.wrong {
  background: radial-gradient(ellipse at top, #ffcdd2, #e57373);
}

.shell.lifted {
  transform: translateY(-60px) scale(1.1) !important;
  z-index: 3;
  box-shadow: 0 20px 40px rgba(0,0,0,0.25);
}

.shell.shake {
  animation: shake 0.4s;
}
@keyframes shake {
  10%, 90% { transform: translateX(-2px); }
  20%, 80% { transform: translateX(4px); }
  30%, 50%, 70% { transform: translateX(-8px); }
  40%, 60% { transform: translateX(8px); }
}

.pearl {
  position: absolute;
  bottom: -20px;
  left: 50%;
  transform: translateX(-50%);
  font-size: 28px;
  animation: sparkle 1.5s infinite ease-in-out;
}

@keyframes sparkle {
  0%, 100% {
    transform: translateX(-50%) scale(1);
    opacity: 1;
  }
  50% {
    transform: translateX(-50%) scale(1.2);
    opacity: 0.8;
  }
}

.controls {
  display: flex;
  justify-content: center;
  gap: 3rem;
  margin-bottom: 3.5rem;
}

.controls button {
  padding: 14px 28px;
  border: none;
  border-radius: 30px;
  font-size: 1.1rem;
  font-weight: bold;
  background: linear-gradient(to right, #4db6ac, #00796b);
  color: white;
  cursor: pointer;
  box-shadow: 0 6px 15px rgba(0, 121, 107, 0.3);
  transition: all 0.3s ease;
  outline: none;
}
.controls button:hover, .controls button:focus {
  background: linear-gradient(to right, #00796b, #004d40);
  transform: scale(1.05);
  outline: 2px solid #00bfa5;
}
.controls button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.scoreboard {
  margin-top: 4rem;
  font-size: 1.2rem;
  background: #ffffff;
  color: #222;
  border-radius: 15px;
  padding: 2rem 2.5rem;
  font-weight: 600;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.2);
}
</style>
