<template>
  <div
    class="ai-assistant-bubble"
    :class="containerClass"
    @mouseenter="onMouseEnter"
    @mouseleave="onMouseLeave"
  >
    <button
      type="button"
      class="ai-float-btn"
      :class="{
        'ai-float-btn--open': openClassEnabled && (open || visibleHint),
        'ai-float-btn--pulse': pulse
      }"
      @click="$emit('toggle')"
      :aria-label="ariaLabel"
      :aria-expanded="open"
    >
      <span class="ai-float-icon">✨</span>
    </button>

    <Transition name="ai-panel">
      <div v-if="open" class="ai-assistant-panel-shell">
        <slot />
      </div>
    </Transition>

    <Transition v-if="useTransition" name="ai-hint">
      <div v-if="visibleHint" class="ai-float-hint" :class="hintClass">
        <p class="ai-float-hint-text">{{ hintText }}</p>
        <span class="ai-float-hint-arrow"></span>
      </div>
    </Transition>
    <div v-else-if="visibleHint" class="ai-float-hint" :class="hintClass">
      <p class="ai-float-hint-text">{{ hintText }}</p>
      <span class="ai-float-hint-arrow"></span>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  hintText: { type: String, default: 'Where to go? Ask me!' },
  showHint: { type: Boolean, default: false },
  open: { type: Boolean, default: false },
  pulse: { type: Boolean, default: false },
  ariaLabel: { type: String, default: 'AI travel assistant' },
  openClassEnabled: { type: Boolean, default: true },
  useTransition: { type: Boolean, default: true },
  containerClass: { type: [String, Array, Object], default: '' },
  hintClass: { type: [String, Array, Object], default: '' },
  emitHover: { type: Boolean, default: true }
})

const emit = defineEmits(['toggle', 'hover-change'])
const visibleHint = computed(() => props.showHint && !props.open)

function onMouseEnter() {
  if (props.emitHover) emit('hover-change', true)
}

function onMouseLeave() {
  if (props.emitHover) emit('hover-change', false)
}
</script>
