<template>
  <div id="wine">
    <div style="padding: 0 20px;">
      <div style="margin-bottom: 15px;">
        <button
          class="button is-primary"
          @click="onsave"
        >
          <span>SAVE</span>
        </button>

        <button
          class="button is-danger"
          @click="ondelete"
        >
          <span>DELETE</span>
        </button>
      </div>

      <div class="columns">
        <div class="column">
          <div class="field">
            <label class="label">ID</label>
            <div class="control">
              <input
                :value="data.id"
                class="input"
                type="text"
                readonly
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Type</label>
            <div class="control">
              <input
                v-model="data.type"
                class="input"
                type="text"
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Catalog</label>
            <div class="control">
              <input
                :value="data.catalog"
                class="input"
                type="text"
                readonly
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Storage</label>
            <div class="control">
              <input
                v-model="data['storage-area']"
                class="input"
                type="text"
              >
            </div>
          </div>
        </div>

        <div class="column">
          <div class="field">
            <label class="label">Name</label>
            <div class="control">
              <input
                v-model="data.name"
                class="input"
                type="text"
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Year of Production</label>
            <div class="control">
              <input
                v-model="data.year"
                class="input"
                type="text"
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Winery</label>
            <div class="control">
              <input
                v-model="data.winery"
                class="input"
                type="text"
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Region</label>
            <div class="control">
              <input
                v-model="data.region"
                class="input"
                type="text"
              >
            </div>
          </div>
        </div>

        <div class="column">
          <div class="field">
            <label class="label">Details</label>
            <div class="control">
              <input
                v-model="data.details"
                class="input"
                type="text"
              >
            </div>
          </div>

          <div class="field">
            <label class="label">Price</label>
            <div class="control">
              <input
                v-model="data.price"
                class="input"
                type="text"
              >
            </div>
          </div>
        </div>
      </div>

      <div class="columns is-variable is-8">
        <div class="column">
          <div class="field">
            <label class="label">Internal Notes</label>
            <div class="control">
              <textarea
                v-model="data['internal-notes']"
                class="textarea"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="column">
          <table class="table is-striped is-hoverable">
            <thead>
              <th>Data</th>
              <th>Supplier</th>
              <th>Price / Bottle</th>
              <th>Qty</th>
            </thead>
            <tfoot>
              <th>Data</th>
              <th>Supplier</th>
              <th>Price / Bottle</th>
              <th>Qty</th>
            </tfoot>
            <tbody>
              <tr>
                <td>+ Add new</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "Wine",

  data() {
    return {
      data: {}
    };
  },

  mounted() {
    this.$http
      .get("/wines/" + this.$route.params.id)
      .then(response => (this.data = response.data[0]));
  },

  methods: {
    onsave() {
      this.$http.patch("/wines/" + this.$route.params.id, this.data);
    },

    ondelete() {
      this.$http.delete("/wines/" + this.$route.params.id).then(response => {
        if (response.status === 200) {
          this.$router.replace("/");
        }
      });
    }
  }
};
</script>

<style>
</style>
