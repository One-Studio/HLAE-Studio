<template>
  <div class="main-container">
    <div class="component">
      <img alt="Wails logo" src="../assets/images/logo-alpha.png" class="logo">
    </div>
    <div class="app-title">HLAE Studio</div>
    <div class="app-version">Testify v1.0.1</div>
    <a-progress
        type="circle"
        :width="80"
        :stroke-color="{
          '0%': '#108ee9',
          '100%': '#87d068',
        }"
        :percent="66.57"
        status="active"
        style="margin-top: 10px;"
    />
    <br>
    <div class="log-area">
      正在检查hlae版本号...
    </div>
    <div class="btn-control">
      <a-button style="margin-right: 10px" size="large">设置</a-button>
      <a-button style="margin-right: 10px" size="large">位置</a-button>
      <a-button style="margin-right: 0" size="large">打开HLAE</a-button>
    </div>
  </div>
</template>

<script>
export default {
  name: "Main"
}
</script>

<style scoped>
.main-container {
  margin: auto;
}

.component {
  /*margin: 20px;*/
  padding: 20px;
}

.app-title {
  font-size: 46px;
  margin: auto;
  text-shadow: 2px 4px 6px rgba(0,0,0,0.3);
}

.app-version {
  font-size: 24px;
  font-weight: lighter;
  mso-ansi-font-weight: lighter;
  mso-bidi-font-weight: lighter;
  margin-bottom: 10px;
}

.logo {
  width: 40%;
  height: auto;
  margin-top: 60px;
}

.log-area {
  margin-top: 10px;
}

.btn-control {
  margin: 10px;
  padding: 20px;
}

</style>