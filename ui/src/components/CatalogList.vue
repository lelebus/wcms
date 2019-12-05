<template>
  <li class="list-item">
    <h3>{{ item.name }}</h3>
    <SortableList
      v-if="item.child.length"
      v-model="item.child"
      class="shortList"
    >
      <Catalog
        v-for="(c, idx) in item.child"
        :key="idx"
        collection="lists"
        :index="idx"
        :item="c"
      />
    </SortableList>
  </li>
</template>

<script>
import { ElementMixin } from "vue-slicksort";
import SortableList from "../components/SortableList";
export default {
  name: "Catalog",
  components: { SortableList },
  mixins: [ElementMixin],
  props: ["item"]
};
</script>

<style>
</style>