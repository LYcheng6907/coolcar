import { getSetting, getUserInfo } from "./utils/wxapi"
import { Coolcar } from "./service/request"
let resolveUserInfo:(value: WechatMiniprogram.UserInfo | PromiseLike<WechatMiniprogram.UserInfo>) => void
let rejectUserInfo:(reason?: any) => void
// app.ts
App<IAppOption>({
  globalData: {
    userInfo:new Promise((resolve, reject) => {
      resolveUserInfo = resolve,
      rejectUserInfo = reject
    })
   
  }, 
  async onLaunch() {
    // // 展示本地存储能力
    // const logs = wx.getStorageSync('logs') || []
    // logs.unshift(Date.now())
    // wx.setStorageSync('logs', logs)
    
    // 登录
    Coolcar.login()

    try {
      const setting = await getSetting()
      if(setting.authSetting['scope.userInfo']){
        const  userInfoRes = await getUserInfo()
        resolveUserInfo(userInfoRes.userInfo)
      }
    } catch (err) {
      rejectUserInfo(err)
    }
    // 使用promise替代callback，改写getsetting
    // userInfo: new Promise((resolve, reject) => {
    //   getSetting().then(res => {
    //     if(res.authSetting['scope.userInfo']){
    //       return getUserInfo()
    //     }
    //     return Promise.resolve(undefined)
    //   }).then(res =>{
    //     if(!res){
    //       return
    //     }

    //     //通知页面我获得了用户信息
    //     resolveUserInfo(res.userInfo)    
    //   }).catch(rejectUserInfo)
    // })
    // // 获取用户信息
    // wx.getSetting({
    //   success: res => {
    //     if (res.authSetting['scope.userInfo']) {
    //       // 已经授权，可以直接调用 getUserInfo 获取头像昵称，不会弹框
    //       wx.getUserInfo({
    //         success: res => {
    //           // 可以将 res 发送给后台解码出 unionId
    //           this.globalData.userInfo = res.userInfo

    //           // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
    //           // 所以此处加入 callback 以防止这种情况
    //           if (this.userInfoReadyCallback) {
    //             this.userInfoReadyCallback(res)
    //           }
    //         },
    //       })
    //     }
    //   },
    // })
  },
  resolveUserInfo(userInfo:WechatMiniprogram.UserInfo){
    resolveUserInfo(userInfo)
  }
})