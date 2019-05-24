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
        :options="catalogs"
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

    .field
      label.checkbox.is-block
        input(v-model="is_active" type="checkbox")
        |  Is active

    button.button.is-primary(@click="save") Save
</template>

<script>
import Multiselect from "vue-multiselect";
import { get, has, pick, reduce } from "lodash-es";

export default {
  name: "Wine",

  components: { Multiselect },

  props: {
    wine: {
      type: Number,
      default: undefined
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
    is_active: true,

    types: [],
    sizes: [],
    wineries: [],
    territories: [],
    regions: [],
    countries: [],
    catalogs: [],

    errors: {}
  }),

  computed: {
    config: {
      get() {
        return pick(this, [
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
          "catalog",
          "details",
          "internal_notes",
          "is_active"
        ]);
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
          "catalog",
          "details",
          "internal_notes",
          "is_active"
        ].forEach(field => {
          if (has(config, field)) {
            this[field] = config[field];
          }
        });
      }
    }
  },

  methods: {
    reset() {
      this.config = {
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
        is_active: true
      };

      this.errors = {};

      this.$http
        .get("/wines/", { params: { id: this.wine } })
        .then(response => (this.config = response.data[0]));

      this.$http.get("/catalogs/parameters/").then(response => {
        [
          "types",
          "sizes",
          "wineries",
          "territories",
          "regions",
          "countries"
        ].forEach(field => {
          this[field] = get(response.data, field, []);
        });
      });

      this.$http
        .get("/catalogs/")
        .then(
          response =>
            (this.catalogs = response.data.filter(
              catalog => !catalog.Customized
            ))
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

    save() {
      this.$http
        .request({
          url: "/wines/",
          method: this.id ? "patch" : "post",
          params: { id: this.id },
          data: [this.config]
        })
        .then(() => {
          this.$parent.is_active = false;
        })
        .catch(error => {
          this.errors = reduce(
            error.response.data,
            (errors, error) => {
              errors[error.id] = error.message;
              return errors;
            },
            {}
          );
        });
    }
  }
};
</script>

<style lang="stylus">
.multiselect {
  height: 2.25em;
}
</style>
