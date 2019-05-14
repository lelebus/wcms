<template>
  <div id="app">
    <nav class="navbar is-fixed-top">
      <div class="navbar-brand">
        <router-link
          to="/"
          class="navbar-item"
        >
          <img src="http://via.placeholder.com/120x30">
        </router-link>

        <a
          :class="['navbar-burger', 'burger', {'is-active': expanded}]"
          @click="expanded = !expanded"
        >
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
          <span aria-hidden="true"></span>
        </a>
      </div>

      <div :class="['navbar-menu', {'is-active': expanded}]">
        <div class="navbar-end">
          <router-link
            to="/import"
            class="navbar-item"
            exact-active-class="is-active"
          >
            <span>IMPORT</span>
          </router-link>

          <router-link
            to="/export"
            class="navbar-item"
            exact-active-class="is-active"
          >
            <span>EXPORT</span>
          </router-link>

          <router-link
            to="/catalog/new"
            class="navbar-item"
            exact-active-class="is-active"
          >
            <span>CREATE CATALOG</span>
          </router-link>
        </div>
      </div>
    </nav>

    <router-view style="margin-top: 25px;"/>
  </div>
</template>

<script>
export default {
  name: "App",

  data() {
    return {
      expanded: false
    };
  }
};
</script>

<style>
</style>
