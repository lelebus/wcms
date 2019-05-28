<template lang="pug">
  #export
    .button.is-primary(@click="save")
      span.icon
        i.fas.fa-download
      span Export
</template>

<script>
import { Workbook } from "exceljs";
import { saveAs } from "file-saver";

export default {
  name: "Export",

  data: function() {
    return {
      wines: [],
      catalogs: []
    };
  },

  mounted: function() {
    this.$http.get("/wines/").then(response => (this.wines = response.data));

    this.$http
      .get("/catalogs/")
      .then(response => (this.catalogs = response.data));
  },

  methods: {
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
        { key: "description", width: 26, style: { font } },
        { key: "year", width: 6, style: { font } },
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

      ws.addRow({
        price: 500,
        year: 2010,
        storage_area: "X 55",
        name: "test wine"
      });

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
            ws.lastRow.font = { size: 14, bold: true };
            break;
        }

        var childrens = this.catalogs.filter(c => c.parent === catalog.id);

        if (childrens.length) {
          childrens.forEach(addCatalog);
        } else if (catalog.wines.length) {
          ws.addRows(
            catalog.wines.map(w => ({
              storage_area: w.storage_area,
              name: w.name,
              description: w.description,
              year: w.year,
              winery: w.winery,
              price: w.price
            }))
          );
        }
      };

      this.catalogs
        .filter(c => c.level === 0)
        .sort(c => c.id)
        .forEach(addCatalog);

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

<style>
</style>
