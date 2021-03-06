<template lang="pug">
  .wine_form
    .field
      label.label Name
      .control
        input.input(
          v-model="name"
          :class="{'is-danger': errors.name}"
        )
      p.help.is-danger(v-if="errors.name") {{ errors.name }}

    .field
      label.label Type
      .control
        multiselect(
          v-model="type"
          :options="types"
          placeholder=""
        )
      p.help.is-danger(v-if="errors.type") {{ errors.type }}

    .field
      label.label Size
      multiselect(
        v-model="size"
        :options="sizes"
        placeholder=""
      )
      p.help.is-danger(v-if="errors.size") {{ errors.size }}

    .field
      label.label Year
      .control
        input.input(
          v-model="year"
          :class="{'is-danger': errors.year}"
          type="number"
        )
      p.help.is-danger(v-if="errors.year") {{ errors.year }}

    .field
      label.label Storage Area
      .control
        input.input(
          v-model="storage_area"
          :class="{'is-danger': errors.storage_area}"
        )
      p.help.is-danger(v-if="errors.storage_area") {{ errors.storage_area }}

    .field
      label.label Winery
      multiselect(
        v-model="winery"
        :options="wineries"
        :taggable="true"
        placeholder=""
        @tag="addTag('winery', $event)")
      p.help.is-danger(v-if="errors.winery") {{ errors.winery }}0

    .field
      label.label Territory
      multiselect(
        v-model="territory"
        :options="territories"
        :taggable="true"
        placeholder=""
        @tag="addTag('territory', $event)"
      )
      p.help.is-danger(v-if="errors.territories") {{ errors.territories }}

    .field
      label.label Region
      multiselect(
        v-model="region"
        :options="regions"
        :taggable="true"
        placeholder=""
        @tag="addTag('region', $event)"
      )
      p.help.is-danger(v-if="errors.region") {{ errors.region }}

    .field
      label.label Country
      multiselect(
        v-model="country"
        :options="countries"
        :taggable="true"
        placeholder=""
        @tag="addTag('country', $event)"
      )
      p.help.is-danger(v-if="errors.country") {{ errors.country }}

    .field
      label.label Price
      .control
        input.input(
          v-model="price"
          :class="{'is-danger': errors.price}"
          type="number"
        )
      p.help.is-danger(v-if="errors.price") {{ errors.price }}

    .field
      label.label Catalogs
      multiselect(
        v-model="catalog"
        :options="catalog_list"
        :multiple="true"
        label="name"
        track-by="id"
        placeholder=""
      )
      p.help.is-danger(v-if="errors.catalog") {{ errors.catalog }}

    .field
      label.label Details
      .control
        input.input(v-model="details")

    .field
      label.label Internal Notes
      .control
        textarea.textarea(v-model="internal_notes")

    .field.is-grouped
      .control
        button.button.is-primary(@click="$emit('save', config)") Save
      .control(v-if="id")
        button.button.is-danger(@click="$emit('delete')") Delete
</template>

<script>
import Multiselect from "vue-multiselect";
import { get, has, pick, merge } from "lodash-es";

export default {
  name: "WineEditor",

  components: { Multiselect },

  props: {
    parameters: {
      type: Object,
      default: () => ({})
    },

    catalogs: {
      type: Array,
      default: () => []
    },

    wine: {
      type: Object,
      default: () => ({})
    },

    errors: {
      type: Object,
      default: () => ({})
    }
  },

  data: () => ({
    id: undefined,

    name: undefined,
    type: undefined,
    size: undefined,
    year: undefined,
    storage_area: undefined,
    winery: undefined,
    territory: undefined,
    region: undefined,
    country: undefined,
    price: undefined,
    catalog: [],
    details: undefined,
    internal_notes: undefined,

    types: [],
    sizes: [],
    wineries: [],
    territories: [],
    regions: [],
    countries: []
  }),

  computed: {
    config: {
      get() {
        return merge(
          pick(this, [
            "id",
            "name",
            "type",
            "size",
            "year",
            "storage_area",
            "winery",
            "territory",
            "region",
            "country",
            "price",
            "details",
            "internal_notes"
          ]),
          {
            catalog: this.catalog.map(catalog => catalog.id)
          }
        );
      },

      set(config) {
        [
          "id",
          "name",
          "type",
          "size",
          "year",
          "storage_area",
          "winery",
          "territory",
          "region",
          "country",
          "price",
          "details",
          "internal_notes"
        ].forEach(field => {
          if (has(config, field)) {
            this[field] = config[field];
          }
        });

        if (config.catalog) {
          this.catalog = config.catalog.map(id =>
            this.catalog_list.find(c => c.id === id)
          );
        }
      }
    },

    catalog_list() {
      return this.catalogs.map(catalog => ({
        id: catalog.id,
        parent: catalog.parent,
        name: this.getCatalogPath(catalog.id)
      }));
    }
  },

  watch: {
    wine(wine) {
      this.config = merge(
        {
          id: undefined,
          name: undefined,
          type: undefined,
          size: undefined,
          year: undefined,
          storage_area: undefined,
          winery: undefined,
          territory: undefined,
          region: undefined,
          country: undefined,
          price: undefined,
          catalog: [],
          details: undefined,
          internal_notes: undefined
        },
        wine
      );
    },

    parameters(parameters) {
      [
        "types",
        "sizes",
        "wineries",
        "territories",
        "regions",
        "countries"
      ].forEach(field => {
        this[field] = get(parameters, field, []);
      });
    }
  },

  methods: {
    reset() {
      this.config = merge(
        {
          id: undefined,
          name: undefined,
          type: undefined,
          size: undefined,
          year: undefined,
          storage_area: undefined,
          winery: undefined,
          territory: undefined,
          region: undefined,
          country: undefined,
          price: undefined,
          catalog: [],
          details: undefined,
          internal_notes: undefined
        },
        this.wine
      );
    },

    addTag(source, value) {
      switch (source) {
        case "winery":
          this.wineries.push(value);
          this.winery = value;
          break;
        case "territory":
          this.territories.push(value);
          this.territory = value;
          break;
        case "region":
          this.regions.push(value);
          this.region = value;
          break;
        case "country":
          this.countries.push(value);
          this.country = value;
          break;
      }
    },

    getCatalogPath(id) {
      let catalog = this.catalogs.find(catalog => catalog.id === id);

      if (catalog) {
        if (catalog.parent) {
          return `${this.getCatalogPath(catalog.parent)} / ${catalog.name}`;
        } else {
          return catalog.name;
        }
      } else {
        return "Unknown";
      }
    }
  }
};
</script>

<style lang="stylus">
.multiselect {
  height: 2.25em;
}
</style>
