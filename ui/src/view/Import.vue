<template lang="pug">
  #import
    .columns
      .column.is-narrow
        .file.has-name.is-primary
          label.file-label
            input.file-input(
              type="file"
              accept=".xls,.xlsx"
              @change="read($event.target.files[0])"
            )
            span.file-cta
              span.file-icon
                i.fas.fa-upload
              span.file-label Choose a file...
            span.file-name(v-if="file_name") {{ file_name }}
      .column(v-if="wines.length")
        button.button.is-primary(@click="upload") Upload

    .columns.is-multiline(v-if="wines.length" style="margin-top:0;")
      .column.is-one-third(v-for="(wine, index) in wines" :key="index")
        Card(
          :wine="wine"
          :catalogs="catalogs"
          :class="{'is-danger': size(wine.errors)}"
          @click.native="open(index)"
        )

    .modal(:class="{'is-active': is_modal_open}")
      .modal-background(@click="is_modal_open = false")
      .modal-content
        .box
          Editor(
            ref="editor"
            :wine="wine"
            :parameters="parameters"
            :catalogs="catalogs"
            :errors="errors"
            @save="save"
            @delete="remove"
          )
      button.modal-close.is-large(@click="is_modal_open = false")
</template>

<script>
import { Workbook } from "exceljs";
import { omit, includes, last, size, clone } from "lodash-es";
import Card from "../components/WineCard";
import Editor from "../components/WineEditor";

export default {
  name: "Import",

  components: { Card, Editor },

  data: () => ({
    last_wine: 0,
    last_catalog: 0,
    open_id: undefined,
    file_name: undefined,
    is_modal_open: false,

    wines: [],
    catalogs: [],
    parameters: {},

    size: size
  }),

  computed: {
    wine() {
      return this.wines[this.open_id];
    },

    errors() {
      if (this.wine) {
        return this.wine.errors;
      } else {
        return {};
      }
    }
  },

  methods: {
    read(file) {
      this.file_name = file.name;
      let file_reader = new FileReader();
      file_reader.onload = event => this.parse(event.target.result);
      file_reader.readAsArrayBuffer(file);
    },

    parse(buffer) {
      this.wines = [];

      this.$http.get("/catalogs/").then(response => {
        this.catalogs = response.data.filter(catalog => catalog.Customized);

        this.$http.get("/catalogs/parameters").then(response => {
          this.parameters = response.data;

          new Workbook().xlsx.load(buffer).then(wb => {
            wb.getWorksheet(1).eachRow((row, index) => {
              // header row
              if (index === 1) {
                return;
              } else {
                var wine = {
                  id: --this.last_wine,
                  name: row.values[4],
                  type: row.values[11],
                  size: (row.values[3] || "").toString(),
                  year: (row.values[6] || "").toString(),
                  storage_area: row.values[2],
                  winery: row.values[5],
                  territory: row.values[8],
                  region: row.values[9],
                  country: row.values[10],
                  price: (row.values[7] || "").toString(),
                  catalog: row.values[1],
                  details: row.values[12],
                  internal_notes: row.values[13],
                  errors: {}
                };

                this.updateParams(wine);

                if (wine.catalog) {
                  var parents = [0];
                  wine.catalog
                    .split("/")
                    .map(catalog => catalog.trim())
                    .forEach(name => {
                      var catalog = this.catalogs
                        .filter(c => c.parent === last(parents))
                        .find(c => c.name === name);

                      if (catalog === undefined) {
                        catalog = {
                          id: --this.last_catalog,
                          name: name,
                          level: parents.length - 1,
                          parent: last(parents),
                          wines: [],
                          Customized: true
                        };
                        this.catalogs.push(catalog);
                      }

                      parents.push(catalog.id);
                      wine.catalog = [catalog.id];
                    });
                } else {
                  wine.catalog = [];
                }

                this.wines.push(wine);
              }
            });
          });
        });
      });
    },

    open(index) {
      this.open_id = index;
      this.$refs.editor.reset();
      this.$nextTick(() => (this.is_modal_open = true));
    },

    save(wine) {
      this.updateParams(wine);
      this.is_modal_open = false;
      this.wines[this.open_id] = wine;
    },

    remove() {
      this.is_modal_open = false;
      this.wines.splice(this.open_id, 1);
    },

    upload() {
      this.upload_r(this.catalogs.filter(c => c.id < 0), [0], { 0: 0 });
    },

    upload_r(catalogs, parents, idmap) {
      if (catalogs.length) {
        var payload = catalogs
          .filter(c => includes(parents, c.parent))
          .map(c => {
            c.parent = idmap[c.parent];
            return c;
          });

        this.$http
          .request({
            url: "/catalogs/",
            method: "post",
            data: payload
          })
          .then(response => {
            response.data.forEach((id, index) => {
              var cid = payload[index].id;
              var aid = this.catalogs.findIndex(c => c.id === cid);

              idmap[cid] = id;

              parents.push(cid);
              parents.push(idmap[cid]);

              this.catalogs[aid].id = id;
            });

            this.upload_r(
              catalogs.filter(c => !includes(parents, c.id)),
              parents,
              idmap
            );
          });
      } else {
        this.wines.forEach((wine, index) => {
          this.wines[index].catalog = wine.catalog.map(id =>
            id < 0 ? idmap[id] : id
          );
        });

        clone(this.wines).forEach(wine => {
          this.$http
            .request({
              url: "/wines/",
              method: "post",
              data: [omit(wine, ["id", "errors"])]
            })
            .then(() => {
              var idx = this.wines.findIndex(w => w.id === wine.id);
              this.wines.splice(idx, 1);
            })
            .catch(error => {
              var idx = this.wines.findIndex(w => w.id === wine.id);

              if (error.response.status == 422) {
                this.wines[idx].errors = error.response.data;
              }
            });
        });
      }
    },

    updateParams(wine) {
      if (!includes(this.parameters.wineries, wine.winery)) {
        this.parameters.wineries.push(wine.winery);
      }

      if (
        wine.territory &&
        !includes(this.parameters.territories, wine.territory)
      ) {
        this.parameters.territories.push(wine.territory);
      }

      if (wine.region && !includes(this.parameters.regions, wine.region)) {
        this.parameters.regions.push(wine.region);
      }

      if (!includes(this.parameters.countries, wine.country)) {
        this.parameters.countries.push(wine.country);
      }
    }
  }
};
</script>