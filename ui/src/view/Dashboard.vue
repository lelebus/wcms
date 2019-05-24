<template lang="pug">
  #dashboard.container
    .columns.is-multiline

      .column.is-one-third
        .box(@click="open(undefined)")
          .columns.is-centered.is-vcentered.is-mobile
            span.icon.is-large.has-text-primary
              i.fas.fa-3x.fa-plus-circle

      .column.is-one-third(v-for="wine in wines" :key="wine.id")
        .box(@click="open(wine.id)")
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

    .modal(:class="{'is-active': is_active}")
      .modal-background(@click="is_active = false")
      .modal-content
        .box
          wine(ref="wine" :wine="id")
      button.modal-close.is-large(@click="is_active = false")
</template>

<script>
import Wine from "../components/Wine";

export default {
  name: "Dashboard",

  components: { Wine },

  data: () => ({
    id: undefined,
    is_active: false,

    tags: [
      { label: "Type", field: "type" },
      { label: "Size", field: "size" },
      { label: "Price", field: "price", prefix: "$" },
      { label: "Storage", field: "storage_area" },
      { label: "Territory", field: "territory" },
      { label: "Region", field: "region" },
      { label: "Country", field: "country" }
    ],

    wines: [],
    catalogs: []
  }),

  mounted() {
    this.$http.get("/wines/").then(response => (this.wines = response.data));
    this.$http
      .get("/catalogs/")
      .then(response => (this.catalogs = response.data));
  },

  methods: {
    open(id) {
      this.id = id;
      this.is_active = true;
      this.$nextTick(() => this.$refs.wine.reset());
    }
  }
};
</script>

<style lang="stylus">
#dashboard .columns .column .box {
  height: 100%;
  cursor: pointer;

  & .columns {
    height: 100%;
    min-height: 160px;
  }
}
</style>
