<template>
  <div class="main-container">
    <div>
      <img alt="Wails logo" src="../assets/images/logo.png" class="logo">
    </div>
    <div class="app-title">HLAE Studio</div>
    <div class="version-area">
      {{ versionCode }} <span class="app-version">{{ appVersion }}</span>
    </div>
    <div class="p-status" >
      <a-progress type="circle" :percent="progress" width="18vw"
                  :strokeColor="{
                    '0%': '#108ee9',
                    '100%': '#87d068',
                    // '0%': '#9370DB',
                    // '100%': '#C71585',
                  //  自定主题后修改成功色，现在成功的时候很丑 @success-color: #52c41a;
                  }"
                  style="font-size: 4vw"
      />
<!--          <a-spin><a-icon slot="indicator" type="loading" style="font-size: 20vw;color: mediumpurple" spin/> </a-spin>-->
<!--      <a-icon type="check-circle" theme="twoTone" two-tone-color="#1ac41a" style="font-size: 12vw"/>-->
    </div>
    <div class="log-area" >
      {{ log }}
    </div>
    <div class="btn-control">
      <a-button class="btn" size="large" @click="tabSetting"><a-icon type="setting" /></a-button>
      <a-button class="btn" size="large" @click="openDirHLAE"><a-icon type="folder-open" /></a-button>
      <a-button class="btn" @click="launchHLAE" style="margin-right: 0;width: 36vw;font-size: 4.5vw;" size="large">打开HLAE</a-button>
    </div>
  </div>
</template>

<script>
import Wails from '@wailsapp/runtime';

export default {
  name: "Main",
  data() {
    return {
      versionCode: "Testify",
      appVersion: "v0.0.1",
      progress: 0,
      log: "",
    };
  },
  mounted() {
    //data多了之后改用解析Json形式赋值
    Wails.Events.On("SetProgess", (progress) => {
      this.progress = progress;
    });
    Wails.Events.On("SetLog", (log) => {
      this.log = log;
    });
    Wails.Events.On("SetVersionCode", (versionCode) => {
      this.versionCode = versionCode;
    });
    Wails.Events.On("SetAppVersion", (appVersion) => {
      this.appVersion = appVersion;
    });
    this.checkUpdate();
  },
  methods: {
    launchHLAE () {
      // console.log("启动HLAE");
      window.backend.App.LaunchHLAE();
    },
    tabSetting () {
      console.log("切换到设置Tab页"); //TODO
    },
    openDirHLAE () {
      // console.log("打开HLAE安装位置");
      window.backend.App.OpenHlaeDirectory();
      //发送wails信息->Go?
      // window.wails.Events.Emit("error", "这是一条错误信息！");
    },
    checkState () {
      // console.log("检查HLAE更新");
      window.backend.App.CheckState();
    },
    checkUpdate () {
      // console.log("检查HLAE更新");
      window.backend.App.CheckUpdate();
    }
  }
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
  text-shadow: 0.5vw 1vw 1.5vw rgba(0,0,0,0.2);
}

/*应用版本号*/
.version-area {
  text-align:center;
  font-size: 6.75vw;
  font-weight: lighter;
  margin-bottom: 4.5vw;
  text-shadow: 0.25vw 1vw 2vw rgba(0,0,0,0.2);
}

/*数字版本号*/
.app-version {
  text-align:center;
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

/*状态区域*/
.p-status {
  height: 18vw;
  margin: 2vw auto;
}

/*日志区域*/
.log-area {
  margin: 2vw auto 0;
  /*padding: 1vw;*/
  font-size: 3.5vw;
  width: 61.8vw;
  text-align:center;
  height: 5vw;
  /*background-color: gray;*/
  white-space: nowrap;   /*这是重点。文本换行*/
  /*text-overflow: clip;*/
  overflow: hidden;
  text-overflow:ellipsis;
}

/*按钮控制区域*/
.btn-control {
  margin: 5.25vw;
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