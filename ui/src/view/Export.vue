<template lang="pug">
  #export
    .button.is-primary(@click="save")
      span.icon
        i.fas.fa-download
      span Export
    SortableList(v-model="mcatalog")
        Catalog(v-for="(item, index) in mcatalog", :index="index", :key="index", :item="item")
</template>

<script>
import { Workbook } from "exceljs";
import { saveAs } from "file-saver";
import SortableList from "../components/SortableList";
import SortableItem from "../components/SortableItem";
import Catalog from "../components/CatalogList";

export default {
  name: "Export",

  components: { SortableList, SortableItem, Catalog },

  data: function() {
    return {
      wines: [],
      catalogs: [],
      mcatalog: []
    };
  },

  watch: {
    mcatalog: {
      deep: true,
      handler(mcatalog) {
        const traverse = list =>
          list.forEach((c, idx) => {
            this.setPosition(c.id, true, idx);
            traverse(c.child);
          });

        traverse(mcatalog);
      }
    }
  },

  mounted: function() {
    this.$http.get("/wines/").then(response => {
      this.wines = response.data;

      this.$http.get("/catalogs/").then(response => {
        this.catalogs = response.data;
        const inflate = c => {
          c.child = this.catalogs
            .filter(a => a.parent === c.id)
            .sort(a => a.position);

          c.child.forEach(inflate);

          if (c.wines) {
            c.wines = c.wines
              .map(w => this.wines.find(e => e.id === w))
              .sort(w => w.name + ", " + w.winery + ", " + w.year);
          }

          return c;
        };

        this.mcatalog = this.catalogs
          .filter(c => c.level === 0)
          .sort(c => c.position)
          .map(inflate);
      });
    });
  },

  methods: {
    setPosition(id, type, pos) {
      const idx = this[type ? "catalogs" : "wines"].findIndex(e => e.id === id);

      if (idx != -1) {
        var obj = this[type ? "catalogs" : "wines"][idx];
        obj.position = pos;
        this[type ? "catalogs" : "wines"][idx] = obj;
      }

      this.$http.request({
        url: "/export/",
        method: "post",
        data: {
          id: id,
          type: type ? "catalog" : "wine",
          position: pos
        }
      });
    },

    save() {
      const wb = new Workbook();
      const ws = wb.addWorksheet("Wines");
      const font = { name: "Monotype Corsiva", size: 12 };

      ws.pageSetup = {
        margins: {
          top: 0.35,
          left: 0.35,
          bottom: 0.35,
          right: 0.35,
          header: 0.45,
          footer: 0.45
        },
        paperSize: 9
      };

      ws.properties.defaultRowHeight = 18;

      ws.columns = [
        { key: "storage_area", width: 8, style: { font } },
        { key: "name", width: 28, style: { font } },
        { key: "description", width: 25, style: { font } },
        {
          key: "year",
          width: 7,
          style: { font, alignment: { horizontal: "left" } }
        },
        { key: "winery", width: 20, style: { font } },
        {
          key: "price",
          width: 10,
          style: {
            font,
            numFmt: "[$€-410] 0;-[$€-410] 0",
            alignment: { horizontal: "left" }
          }
        }
      ];

      const addCatalog = catalog => {
        ws.addRow([catalog.name]);

        switch (catalog.level) {
          case 0:
            ws.lastRow.font = { size: 20, underline: true };
            break;
          case 1:
            ws.lastRow.font = { size: 16, bold: true };
            break;
          case 2:
            ws.lastRow.font = { bold: true };
            break;
        }

        catalog.child.forEach(addCatalog);

        if (catalog.wines) {
          ws.addRows(
            catalog.wines.map(w => {
              if (w.year) {
                return {
                  storage_area: w.storage_area,
                  name: w.name,
                  description: w.description,
                  year: Number(w.year),
                  winery: w.winery,
                  price: Number(w.price)
                };
              } else {
                return {
                  storage_area: w.storage_area,
                  name: w.name,
                  description: w.description,
                  winery: w.winery,
                  price: Number(w.price)
                };
              }
            })
          );
        }
      };

      this.mcatalog.forEach(addCatalog);

      wb.xlsx.writeBuffer().then(buffer => {
        var blob = new Blob([buffer.buffer], {
          type: "application/vnd.ms-excel"
        });

        saveAs(blob, "wines.xls");
      });
    }
  }
};
</script>

<style lang="scss">
.root {
  display: flex;
  height: 100%;
  box-sizing: border-box;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

a {
  text-decoration: none;
  &:hover {
    text-decoration: underline;
  }
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  padding: 10px;
}

pre {
  width: 80%;
  max-width: 500px;
  border-radius: 20px;
  padding: 20px;
  background: #fefefe;
}

.list {
  width: 80%;
  max-height: 80vh;
  max-width: 500px;
  margin: 0 auto;
  padding: 0;
  overflow: auto;
  background-color: #f3f3f3;
  border: 1px solid #efefef;
  border-radius: 3;
}

.list-item {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 20px;
  background-color: #fff;
  border-bottom: 1px solid #efefef;
  box-sizing: border-box;
  user-select: none;

  color: #333;
  font-weight: 400;
}

.stylizedList {
  position: relative;
  z-index: 0;
  background-color: #f3f3f3;
  border: 1px solid #efefef;
  border-radius: 3px;
  outline: none;
}
.stylizedItem {
  display: flex;
  align-items: center;
  width: 100%;
  padding: 0 20px;
  background-color: #fff;
  border-bottom: 1px solid #efefef;
  box-sizing: border-box;
  user-select: none;

  color: #333;
  font-weight: 400;
}

.handle {
  display: block;
  width: 18px;
  height: 18px;
  background-image: url('data:image/svg+xml;charset=utf-8,<svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 50 50"><path d="M 0 7.5 L 0 12.5 L 50 12.5 L 50 7.5 L 0 7.5 z M 0 22.5 L 0 27.5 L 50 27.5 L 50 22.5 L 0 22.5 z M 0 37.5 L 0 42.5 L 50 42.5 L 50 37.5 L 0 37.5 z" color="black"></path></svg>');
  background-size: contain;
  background-repeat: no-repeat;
  opacity: 0.25;
  margin-right: 20px;
  cursor: row-resize;
}

.horizontalList {
  display: flex;
  width: 600px;
  height: 300px;
  white-space: nowrap;
}

.horizontalItem {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 200px;
  border-right: 1px solid #efefef;
  border-bottom: 0;
}

.grid {
  display: block;
  width: 130 * 4px;
  height: 350px;
  white-space: nowrap;
  border: 0;
  background-color: transparent;
}

.gridItem {
  float: left;
  width: 130px;
  padding: 8px;
  background: transparent;
  border: 0;

  .wrapper {
    display: flex;
    align-items: center;
    justify-content: center;

    width: 100%;
    height: 100%;
    background: #fff;
    border: 1px solid #efefef;

    font-size: 28px;

    span {
      display: none;
    }
  }
}

.category {
  height: auto;

  .categoryHeader {
    display: flex;
    flex-flow: row nowrap;
    align-items: center;
    padding: 10px 14px;
    background: #f9f9f9;
    border-bottom: 1px solid #efefef;
  }

  .categoryList {
    height: auto;
  }
}

.divider {
  padding: 10px 20px;
  background: #f9f9f9;
  border-bottom: 1px solid #efefef;
  text-transform: uppercase;
  font-size: 14px;
  color: #333;
}

.helper {
  box-shadow: 0 5px 5px -5px rgba(0, 0, 0, 0.2),
    0 -5px 5px -5px rgba(0, 0, 0, 0.2);
}
.stylizedHelper {
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
  background-color: rgba(31, 136, 255, 0.8);
  color: #fff;
  cursor: row-resize;
  border: 1px solid white;

  &.horizontalItem {
    cursor: col-resize;
  }
  &.gridItem {
    background-color: transparent;
    white-space: nowrap;
    box-shadow: none;

    .wrapper {
      background-color: rgba(255, 255, 255, 0.8);
      box-shadow: 0 0 7px rgba(0, 0, 0, 0.15);
    }
  }
}

.shrinkedHelper {
  height: 20px !important;
}
</style>
