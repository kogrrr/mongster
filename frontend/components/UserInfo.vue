<template>
  <div v-if="!loggedIn">
    <a href="/auth/login">
      <button class="btn btn-primary"><i class="fab fa-google"></i> Sign in with Google</button>
    </a>
  </div>

  <b-nav-item-dropdown right v-else-if="loggedIn">
    <!-- Using 'button-content' slot -->
    <template v-slot:button-content>
      <img id="user_icon" width="35" class="rounded-circle" alt="user icon" v-bind:src="icon">
    </template>
    <div class="userinfo-header">User Info:</div>
    <b-dropdown-item><em>{{ username }}</em></b-dropdown-item>
    <b-dropdown-item href="/auth/logout">Sign Out</b-dropdown-item>
  </b-nav-item-dropdown>
</template>

<script>
  import $ from 'jquery'

  export default {
    name: 'UserInfo',
    data() {
      return {
        loggedIn: false,
        icon: '',
        username: 'placeholder'
      }
    },
    mounted() {
      console.log("userinfo component mounted");
      fetch('/auth/self')
        .then((response) => {
          if (response.status == 200) {
                return response.json();
          } else {
            throw new Error("Userinfo API responded with " + response.status);
          }
        })
        .catch((error) => {
          console.error("error requesting userinfo: ", error);
        })
        .then((userdata) => {
          if ($.isEmptyObject(userdata)) {
            this.loggedIn = false;
            this.icon = "";
            this.username = "anonymous";
          } else {
            this.loggedIn = true;
            this.icon = userdata.icon;
            this.username = userdata.name;
          }
        })
        .catch((error) => {
          console.error("error parsing userinfo: ", error);
        });

    }
  }
</script>

<style scoped>
.userinfo-header {
  font-weight: bold;
  padding: 0.25rem 1.5rem;
}
</style>
