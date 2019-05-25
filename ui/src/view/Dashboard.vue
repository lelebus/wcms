<template lang="pug">
  #dashboard.container
    .columns.is-multiline

      .column.is-one-third
        .box(@click="open({})")
          .columns.is-centered.is-vcentered.is-mobile
            span.icon.is-large.has-text-primary
              i.fas.fa-3x.fa-plus-circle

      .column.is-one-third(v-for="wine in wines" :key="wine.id")
        Card(
          :wine="wine"
          :catalogs="catalogs"
          @click.native="open(wine)"
        )

    .modal(:class="{'is-active': is_modal_open}")
      .modal-background(@click="is_modal_open = false")
      .modal-content
        .box
          Editor(
            ref="editor"
            :wine="wine"
            :parameters="parameters"
            :catalogs="catalogs.filter(catalog => catalog.Customized)"
            :errors="errors"
            @save="save"
            @delete="remove"
          )
      button.modal-close.is-large(@click="is_modal_open = false")
</template>

<script>
import Card from "../components/WineCard";
import Editor from "../components/WineEditor";
import { find, merge } from "lodash-es";

export default {
  name: "Dashboard",

  components: { Card, Editor },

  data: () => ({
    id: undefined,
    is_modal_open: false,

    wines: [],
    catalogs: [],
    parameters: {},

    errors: {}
  }),

  computed: {
    wine() {
      return find(this.wines, ["id", this.id]);
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
      this.$refs.editor.reset();
      this.$nextTick(() => (this.is_modal_open = true));
    },

    save(wine) {
      this.$http
        .request({
          url: "/wines/",
          method: this.id ? "patch" : "post",
          params: { id: this.id },
          data: [wine]
        })
        .then(response => {
          this.is_modal_open = false;

          if (this.id) {
            let index = this.wines.findIndex(wine => wine.id === this.id);
            this.wines[index] = wine;
          } else {
            wine.id = response.data[0];
            this.wines.push(wine);
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

  & .columns {
    height: 100%;
    min-height: 160px;
  }
}
</style>
