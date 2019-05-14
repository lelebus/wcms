<template>
  <div id="dashboard">
    <table class="table is-fullwidth is-striped is-hoverable">
      <thead>
        <tr>
          <th>ID</th>
          <th>Storage</th>
          <th>Name</th>
          <th>Size</th>
          <th>Year</th>
          <th>Winery</th>
          <th>Region</th>
          <th>Country</th>
          <th>Price</th>
          <th>Type</th>
        </tr>
      </thead>
      <tfoot>
        <tr>
          <th>ID</th>
          <th>Storage</th>
          <th>Name</th>
          <th>Size</th>
          <th>Year</th>
          <th>Winery</th>
          <th>Region</th>
          <th>Country</th>
          <th>Price</th>
          <th>Type</th>
        </tr>
      </tfoot>
      <tbody>
        <router-link
          v-for="wine in wines"
          :key="wine.id"
          :to="'/wine/' + wine.id"
          tag="tr"
          class="wine-entry"
        >
          <td>{{ wine.id }}</td>
          <td>{{ wine['storage-area'] }}</td>
          <td>{{ wine.name }}</td>
          <td>{{ wine.size }}</td>
          <td>{{ wine.year }}</td>
          <td>{{ wine.winery }}</td>
          <td>{{ wine.region }}</td>
          <td>{{ wine.country }}</td>
          <td>{{ wine.price }}</td>
          <td>{{ wine.type }}</td>
        </router-link>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: "Dashboard",

  data() {
    return {
      wines: []
    };
  },

  mounted() {
    this.$http.get("wines/").then(response => (this.wines = response.data));
  }
};
</script>

<style>
.wine-entry {
  cursor: pointer;
}
</style>
