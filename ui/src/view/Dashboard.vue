<template lang="pug">
  #dashboard.container
    .columns.is-multiline

      .column.is-one-third
        .box(@click="open({})")
          .columns.is-centered.is-vcentered.is-mobile
            span.icon.is-large.has-text-primary
              i.fas.fa-3x.fa-plus-circle

      .column.is-one-third(v-for="wine in wines" :key="wine.id")
        .box(@click="open(wine)")
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

    .modal(:class="{'is-active': is_modal_open}")
      .modal-background(@click="is_modal_open = false")
      .modal-content
        .box
          wine(
            :wine="wine"
            :parameters="params"
            :errors="errors"
            @save="save"
            @delete="remove"
          )
      button.modal-close.is-large(@click="is_modal_open = false")
</template>

<script>
import Wine from "../components/Wine";
import { find, merge } from "lodash-es";

export default {
  name: "Dashboard",

  components: { Wine },

  data: () => ({
    id: undefined,
    is_modal_open: false,

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
    catalogs: [],
    parameters: {},

    errors: {}
  }),

  computed: {
    wine() {
      return find(this.wines, ["id", this.id]);
    },

    params() {
      return merge(this.parameters, {
        catalogs: this.catalogs.filter(catalog => catalog.Customized)
      });
    }
  },

  mounted() {
    this.$http.get("/wines/").then(response => (this.wines = response.data));

    this.$http
      .get("/catalogs/parameters/")
      .then(response => (this.parameters = response.data));

    this.$http
      .get("/catalogs/")
      .then(response => (this.catalogs = response.data));
  },

  methods: {
    open(wine) {
      this.id = wine.id;
      this.is_modal_open = true;
    },

    save(wine) {
      this.$http
        .request({
          url: "/wines/",
          method: this.id ? "patch" : "post",
          params: { id: this.id },
          data: [wine]
        })
        .then(() => {
          this.is_modal_open = false;

          if (this.id) {
            let index = this.wines.findIndex(wine => wine.id === this.id);
            this.wines[index] = wine;
          } else {
            this.$http
              .get("/wines/")
              .then(response => (this.wines = response.data));
          }

          this.$http
            .get("/catalogs/parameters/")
            .then(response => (this.parameters = response.data));
        })
        .catch(error => {
          if (error.response.status === 422) {
            this.errors = error.response.data;
          }
        });
    },

    remove() {
      this.$http
        .request({ url: "/wines/", method: "delete", params: { id: this.id } })
        .then(() => {
          this.is_modal_open = false;
          this.wines.splice(
            this.wines.findIndex(wine => wine.id === this.id),
            1
          );
        });
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
