<template lang="pug">
  .catalog(:class="{'indent': level}")
    .field
      label.label Name
      .control
        input.input(v-model="name")

    .field.parameters(
      v-if="!customized && (terminal || level === 2) && !childrens.length"
    )
      .field
        label.label Type
        .control
          multiselect(
            v-model="type"
            :options="parameters.types"
            :multiple="true"
            placeholder=""
          )

      .field
        label.label Size
        .control
          multiselect(
            v-model="size"
            :options="parameters.sizes"
            :multiple="true"
            placeholder=""
          )

      .field
        label.label Winery
        .control
          multiselect(
            v-model="winery"
            :options="parameters.wineries"
            :multiple="true"
            placeholder=""
          )

      .field
        label.label Territory
        .control
          multiselect(
            v-model="territory"
            :options="parameters.territories"
            :multiple="true"
            placeholder=""
          )

      .field
        label.label Region
        .control
          multiselect(
            v-model="region"
            :options="parameters.regions"
            :multiple="true"
            placeholder=""
          )

      .field
        label.label Country
        .control
          multiselect(
            v-model="country"
            :options="parameters.countries"
            :multiple="true"
            placeholder=""
          )

    .field.wines(
      v-if="customized && (terminal || level === 2) && !childrens.length"
    )
      .field
        label.label Wines
        .control
          .level
            .level-left(style="flex-grow:1;")
              .level-item(style="flex-grow:1;")
                multiselect(
                  v-model="wine"
                  :options="selectableWines"
                  :custom-label="customLabel"
                  placeholder=""
                )
              .level-item
                a.button.is-primary(@click="addWine") Add
      .field
        .box(v-for="wine in wines_")
          .level
            .level-left
              .level-item(style="display:block;")
                  h5.title.is-5(style="margin-bottom:1.2em;") {{ wine.name }}
                  h6.subtitle.is-6 {{ wine.winery || "Unknown" }}, {{ wine.year }}
            .level-right
              .level-item
                a.delete.is-large(@click="removeWine(wine.id)")

    p.buttons
      a.button.is-primary(
        v-if="customized && !terminal && !childrens.length && level < 2"
        @click="terminal = true"
      )
        span.icon
          i.fas.fa-wine-glass-alt
        span Add wines
      a.button.is-primary(
        v-if="!customized && !terminal && !childrens.length && level < 2"
        @click="terminal = true"
      )
        span.icon
          i.fas.fa-database
        span Add parameters
      a.button.is-info(
        v-if="!terminal && level < 2"
        @click="addChild"
      )
        span.icon
          i.fas.fa-file-alt
        span Add child
      a.button.is-danger(
        v-if="level > 0"
        @click="$emit('remove')"
      )
        span.icon
          i.fas.fa-trash-alt
        span Remove
      a.button.is-danger(
        v-if="terminal && level == 0 && id < 0"
        @click="reset"
      )
        span.icon
          i.fas.fa-redo
        span Reset

    Catalog(
      ref="childs"
      v-for="catalog in childrens"
      :key="catalog.id"
      :catalog="catalog"
      :catalogs="catalogs_"
      :wines="wines"
      :parameters="parameters"
      :customized="customized"
      @remove="removeChild(catalog.id)"
      @saved="$emit('saved')"
    )
</template>

<script>
import Multiselect from "vue-multiselect";
import { has, pick, clone, merge, isEmpty } from "lodash-es";

export default {
  name: "Catalog",

  components: { Multiselect },

  props: {
    catalog: {
      type: Object,
      default: () => ({})
    },

    catalogs: {
      type: Array,
      default: () => []
    },

    wines: {
      type: Array,
      default: () => []
    },

    parameters: {
      type: Object,
      default: () => []
    },

    customized: {
      type: Boolean,
      default: false
    }
  },

  data: () => ({
    id: undefined,
    level: undefined,
    parent: undefined,

    name: undefined,
    type: [],
    size: [],
    winery: [],
    territory: [],
    region: [],
    country: [],
    wines_: [],

    wine: undefined,

    created: [],
    removed: [],
    terminal: false
  }),

  computed: {
    catalogs_() {
      return this.catalogs
        .filter(c => !this.removed.includes(c.id))
        .concat(this.created);
    },

    childrens() {
      return this.catalogs_.filter(c => c.parent === this.id);
    },

    wineIds() {
      return this.wines_.map(w => w.id);
    },

    selectableWines() {
      return this.wines.filter(w => !this.wineIds.includes(w.id));
    },

    config: {
      get() {
        var fields = ["level", "parent", "name"];

        if (this.id > 0) {
          fields = fields.concat(["id"]);
        }

        if (this.childrens.length) {
          return pick(this, fields);
        }

        if (this.customized) {
          return merge(pick(this, fields), {
            wines: this.winesIds
          });
        } else {
          return pick(
            this,
            fields.concat([
              "type",
              "size",
              "winery",
              "territory",
              "region",
              "country"
            ])
          );
        }
      },

      set(config) {
        [
          "id",
          "level",
          "parent",
          "name",
          "type",
          "size",
          "winery",
          "territory",
          "region",
          "country"
        ].forEach(field => {
          if (has(config, field)) {
            this[field] = clone(config[field]);
          }
        });

        if (this.customized && config.wines) {
          this.wines_ = config.wines.map(
            id => this.wines.find(w => w.id === id) || {}
          );

          this.terminal = config.wines.length;
        }

        if (!this.customized) {
          this.terminal = !(
            isEmpty(config.type) &&
            isEmpty(config.size) &&
            isEmpty(config.winery) &&
            isEmpty(config.territory) &&
            isEmpty(config.region) &&
            isEmpty(config.country)
          );
        }
      }
    }
  },

  watch: {
    catalog(catalog) {
      this.config = catalog;
    }
  },

  mounted() {
    this.config = this.catalog;
  },

  methods: {
    reset() {
      this.config = {
        id: undefined,
        level: undefined,
        parent: undefined,
        name: undefined,
        type: [],
        size: [],
        winery: [],
        territory: [],
        region: [],
        country: [],
        wines: []
      };

      this.wine = undefined;
      this.created = [];
      this.removed = [];
      this.terminal = false;

      this.config = this.catalog;
    },

    addChild() {
      this.created.push({
        id: -1 - this.catalogs_.length,
        level: this.level + 1,
        parent: this.id
      });
    },

    removeChild(id) {
      if (id < 0) {
        var cid = this.created.findIndex(c => c.id === id);
        this.created.splice(cid, 1);
      } else {
        this.removed.push(id);
      }
    },

    customLabel(wine) {
      return `${wine.name}, ${wine.winery || "Unknown"}, ${wine.year}`;
    },

    addWine() {
      this.wines_.push(this.wine);
      this.wine = undefined;
    },

    removeWine(id) {
      this.wines_ = this.wines_.filter(w => w.id !== id);
    },

    save(parent) {
      this.$http
        .request({
          url: "/catalogs/",
          method: this.id < 0 ? "post" : "patch",
          params: { id: this.config.id },
          data: [merge(this.config, { parent: parent })]
        })
        .then(response => {
          this.removed.forEach(id =>
            this.$http.delete("/catalogs/", { params: { id: id } })
          );

          var parent = this.id < 0 ? response.data[0] : this.id;

          if (this.$refs.childs) {
            this.$refs.childs.forEach(child => child.save(parent));
          }

          this.$emit("saved");
        });
    }
  }
};
</script>

<style lang="stylus">
.catalog {
  margin-bottom: 8px;

  .multiselect {
    height: 2.25em;
  }

  .columns .column .box {
    cursor: pointer;
  }
}

.catalog.indent {
  margin-left: 15px;
  padding-left: 15px;
  border-left: 1px solid grey;
}
</style>
