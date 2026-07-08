<script setup>
defineProps({
  categories: { type: Array, required: true },
  modelValue: { type: String, required: true }
})

defineEmits(['update:modelValue'])
</script>

<template>
  <div class="flex gap-2 overflow-x-auto pb-3 mb-6 custom-scrollbar scrollbar-none">
    <button
        v-for="cat in categories"
        :key="cat"
        @click="$emit('update:modelValue', cat)"
        :class="modelValue === cat
        ? 'bg-indigo-600 text-white shadow-sm border-indigo-600'
        : 'bg-white text-gray-600 hover:bg-gray-100 border-gray-200'"
        class="px-5 py-2 rounded-full text-sm font-bold border whitespace-nowrap transition"
    >
      {{ cat }}
    </button>
  </div>
</template>

