<template>
  <li
    class="list-item"
    :style="{height: '50px'}"
  >
    <span
      v-if="showHandle"
      v-handle
      class="handle"
    ></span>
    {{ item.name }}
  </li>
</template>

<script>
import { ElementMixin, HandleDirective } from "vue-slicksort";
export default {
  directives: { handle: HandleDirective },
  mixins: [ElementMixin],
  props: ["item", "showHandle"]
};
</script>