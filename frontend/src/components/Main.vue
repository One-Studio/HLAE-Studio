<template>
  <div class="main-container">
    <div>
      <img alt="Wails logo" src="../assets/images/logo.png" class="logo">
    </div>
    <div class="app-title">HLAE Studio</div>
    <div class="app-version">
      Testify <span class="version">v1.0.1</span>
    </div>
    <a-progress
        type="circle"
        width="15vw"
        :stroke-color="{
          '0%': '#108ee9',
          '100%': '#87d068',
        }"
        :percent="66"
        status="active"
        style="margin-top: 0vw;font-size: 4vw;"
    />
    <br>
    <div class="log-area">
      正在检查hlae版本号...
    </div>
    <div class="btn-control">
      <a-button class="btn" size="large"><a-icon type="setting" /></a-button>
      <a-button class="btn" size="large"><a-icon type="folder-open" /></a-button>
      <a-button class="btn" style="margin-right: 0;width: 36vw;font-size: 4.5vw;" size="large">打开HLAE</a-button>
    </div>
  </div>
</template>

<script>
export default {
  name: "Main"
}
</script>

<style scoped>

/*Main组件容器*/
.main-container {
  margin: auto;
}

/*通用组件样式*/
.component {
  /*margin: 20px;*/
  /*padding: 5vw;*/
}

/*应用大标题*/
.app-title {
  font-size: 11vw;
  margin: auto;
  text-shadow: 0.5vw 1.75vw 2vw rgba(0,0,0,0.25);
}

/*应用版本号*/
.app-version {
  font-size: 6.75vw;
  font-weight: lighter;
  margin-bottom: 4.5vw;
  text-shadow: 0.25vw 1vw 2vw rgba(0,0,0,0.2);
}

/*数字版本号*/
.version {
  font-size: 5vw;
}

/*LOGO图标*/
.logo {
  width: 45vw;
  height: auto;
  margin-top: 24vw;
  margin-bottom: 2.5vw;
  mso-border-shadow: yes;
  filter:drop-shadow(-25vw 25vw 100vw rgba(26,58,70,0.7));
  /*box-shadow: 1vw 1vw 1vw gray;TODO*/
}

/*日志区域*/
.log-area {
  margin-top: 2.5vw;
  font-size: 3.5vw;
}

/*按钮控制区域*/
.btn-control {
  margin: 6.25vw;
  /*padding: 4vw;*/
}

.btn {
  margin-right: 2vw;
  width: 12.5vw;
  height: 11vw;
  font-size: 4vw;
  /*box-shadow: 0 0.25vw 2vw rgba(0,0,0,0.2);*/
  border-radius: 3vw;
}

</style>