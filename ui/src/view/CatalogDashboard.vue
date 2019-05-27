<template lang="pug">
  #catalog-dashboard
    .tabs
      ul
        li(
          :class="{'is-active': !customized}"
          @click="customized = false"
        )
          a Automatic
        li(
          :class="{'is-active': customized}"
          @click="customized = true"
        )
          a Customized

    .columns.is-multiline
      .column.is-one-third
        .box(@click="open({id: -1, parent: 0, level: 0, name: undefined})")
          .columns.is-centered.is-vcentered.is-mobile
            span.icon.is-large.has-text-primary
              i.fas.fa-3x.fa-plus-circle
      .column.is-one-third(
        v-for="catalog in rootCatalogs.filter(c => c.Customized === customized)"
        :key="catalog.id"
      )
        .box(@click="open(catalog)")
          .title.is-5 {{ catalog.name }}

    .modal(:class="{'is-active': is_modal_open}")
      .modal-background(@click="is_modal_open = false")
      .modal-content
        .box
          Catalog(
            ref="editor"
            :catalog="catalog"
            :catalogs="catalogs"
            :wines="wines"
            :parameters="parameters"
            :customized="customized"
            @saved="reload"
          )
          .field.is-grouped(style="margin-top:1.5em;")
            .control
              button.button.is-success(@click="save") Save
            .control(v-if="catalog.id > 0")
              button.button.is-danger(@click="remove") Delete
      button.modal-close.is-large(@click="is_modal_open = false")
</template>

<script>
import { debounce } from "lodash-es";
import Catalog from "../components/Catalog";

export default {
  name: "CatalogDashboard",

  components: { Catalog },

  data: () => ({
    wines: [],
    catalogs: [],
    parameters: {},
    customized: false,
    catalog: {},
    is_modal_open: false
  }),

  computed: {
    rootCatalogs() {
      return this.catalogs.filter(c => c.parent === 0);
    },

    reload() {
      return debounce(
        () =>
          this.$http.get("/catalogs").then(response => {
            this.catalogs = response.data;
          }),
        100
      );
    }
  },

  mounted() {
    this.$http.get("/wines").then(response => {
      this.wines = response.data;
    });

    this.$http.get("/catalogs").then(response => {
      this.catalogs = response.data;
    });

    this.$http.get("/catalogs/parameters").then(response => {
      this.parameters = response.data;
    });
  },

  methods: {
    open(catalog) {
      this.catalog = catalog;
      this.$refs.editor.reset();
      this.$nextTick(() => (this.is_modal_open = true));
    },

    save() {
      this.$refs.editor.save(0);
      this.is_modal_open = false;
    },

    remove() {
      this.$http
        .delete("/catalogs", { params: { id: this.catalog.id } })
        .then(() => {
          this.is_modal_open = false;
          this.catalogs = this.catalogs.filter(c => c.id !== this.catalog.id);
        });
    }
  }
};
</script>

<style lang="stylus">
#catalog-dashboard .columns .column .box {
  height: 100%;
  cursor: pointer;
  min-height: 160px;

  & .columns {
    height: 100%;
    min-height: 160px;
  }
}
</style>
