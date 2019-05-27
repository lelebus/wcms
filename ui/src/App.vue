<template lang="pug">
  #app
    nav.navbar.is-fixed-top
      .navbar-brand
        router-link.navbar-item(
          to="/"
          @click.native="navbar_expanded = false"
        )
          img(src="../static/logo.jpg" style="height:50px;max-height:initial;")
        a.navbar-burger.burger(
          :class="{'is-active': navbar_expanded}"
          @click="navbar_expanded = !navbar_expanded"
        )
          span
          span
          span
      .navbar-menu(:class="{'is-active': navbar_expanded}")
        .navbar-end
          router-link.navbar-item(
            to="/import"
            exact-active-class="is-active"
            @click.native="navbar_expanded = false"
          ) IMPORT
          router-link.navbar-item(
            to="/export"
            exact-active-class="is-active"
            @click.native="navbar_expanded = false"
          ) EXPORT
          router-link.navbar-item(
            to="/catalog/dashboard"
            exact-active-class="is-active"
            @click.native="navbar_expanded = false"
          ) CATALOGS
    router-view(style="margin-top:25px;")
  </div>
</template>

<script>
export default {
  name: "App",

  data() {
    return {
      navbar_expanded: false
    };
  }
};
</script>

<style>
</style>
