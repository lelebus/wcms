<template lang="pug">
  .box
    h5.title.is-5 {{ wine.name }}
    h6.subtitle.is-6 {{ wine.winery || "Unknown" }}, {{ wine.year }}
    .field.is-grouped.is-grouped-multiline
      .control(
        v-for="{label, field, prefix} in tags"
        v-if="wine[field]"
        :key="field"
      )
        .tags.has-addons
          span.tag {{ label }}
          span.tag.is-primary {{ prefix }}{{ wine[field] }}
    .tags(v-if="wine.catalog.length")
      span.tag.is-info Catalogs
      span.tag(v-for="id in wine.catalog" :key="id")
        | {{ (catalogs.find(catalog => catalog.id === id) || {}).name }}
</template>

<script>
export default {
  name: "WineCard",

  props: {
    wine: {
      type: Object,
      default: () => ({})
    },

    catalogs: {
      type: Array,
      default: () => []
    }
  },

  data: () => ({
    tags: [
      { label: "Type", field: "type" },
      { label: "Size", field: "size" },
      { label: "Price", field: "price", prefix: "$" },
      { label: "Storage", field: "storage_area" },
      { label: "Territory", field: "territory" },
      { label: "Region", field: "region" },
      { label: "Country", field: "country" }
    ]
  })
};
</script>
