<template>
  <div id="app">
    <Main/>
    <Setting/>
  </div>
</template>

<script>
import Main from "@/components/Main";
import Setting from "@/components/Setting";
import "./assets/css/main.css";

export default {
  name: "app",
  components: {
    Main,
    Setting
  }
};
</script>
