import { routing } from "../../utils/routing"

Page({
  isPageShowing: true,

  data: {
    avatarURL: '',
    setting: {
      skew: 0,
      rotate: 0,
      showLocation: true,
      showScale: true,
      subKey: '',
      layerStyle: -1,
      enableZoom: true,
      enableScroll: true,
      enableRotate: false,
      showCompass: false,
      enable3D: false,
      enableOverlooking: false,
      enableSatellite: false,
      enableTraffic: false,
    },
    location: {
      latitude: 23.099994,
      longitude: 113.32452,
    },
    scale: 10,
    markers: [
      {
        iconPath: "/resources/car.png",
        id: 0,
        latitude: 23.099994,
        longitude: 113.324520,
        width: 50,
        height: 50
      },
      {
        iconPath: "/resources/car.png",
        id: 1,
        latitude: 23.099994,
        longitude: 114.324520,
        width: 50,
        height: 50
      },
    ],
  },

  async onLoad() {
    const userInfo = await getApp<IAppOption>().globalData.userInfo
    this.setData({
      avatarURL: userInfo.avatarUrl,
    })
  },
  onMyLocationTap(){
    wx.getLocation({
      type:'gcj02',
      success: res => {
        this.setData({
          location:{
            latitude:res.latitude,
            longitude:res.longitude,
          }
        })
      },
      fail:()=>{
        wx.showToast({
          icon:'none',
          title:'请前往设置页授权',
        })
      }
    })
  },
  // onMyLocationTap() {
  //   wx.getLocation({
  //     type: 'gcj02',
  //     success: res => {
  //       this.setData({
  //         location: {
  //           latitude: res.latitude,
  //           longitude: res.longitude,
  //         },
  //       })
  //     }, 
  //     fail: () => {
  //       wx.showToast({
  //         icon: 'none',
  //         title: '请前往设置页授权',
  //       })
  //     }
  //   })
  // },
  OnScanTap(){
      wx.scanCode({
        success: () =>{
          wx.showModal({
            title: '身份认证',
            content: '需要身份验证才能租车',
            success: () => {
              // 假设扫码获得所租车辆的信息
              const carID = 'car123'
              // const redirectURL=`/pages/lock/lock?car_id=${carID}`
              const redirectURL = routing.lock({
                car_id: carID,
              })
              wx.navigateTo({
                // encodeURIComponent() 转义URL 防止路径出错
                // url:`/pages/register/register?redirect=${encodeURIComponent(redirectURL)}`,
                url: routing.register({
                  redirectURL:redirectURL
                })
              })
            }
          })
        },
        fail:console.error,
      })
  },
  // onScanTap() {
  //   wx.scanCode({
  //     success: async () => {
  //       await this.selectComponent('#licModal').showModal()
  //       // TODO: get car id from scan result
  //       const carID='car123'
  //       const redirectURL = routing.lock({
  //         car_id: carID,
  //       })
  //       wx.navigateTo({
  //         url: routing.register({
  //           redirectURL: redirectURL,
  //         })
  //       })
  //     },
  //     fail: console.error,
  //   })
  // },

  onShow() {
    this.isPageShowing = true;
  },
  
  onHide() {
    this.isPageShowing = false;
  },

  onMyTripsTap() {
    wx.navigateTo({
      url: routing.mytrips(),
      // url:'/pages/mytrips/mytrips'
    })
  },

  moveCars() {
    const map = wx.createMapContext("map")
    const dest = {
      latitude: 23.099994,
      longitude: 113.324520,
    }

    const moveCar = () => {
      dest.latitude += 0.1
      dest.longitude += 0.1
      map.translateMarker({
        destination: {
          latitude: dest.latitude,
          longitude: dest.longitude,
        },
        markerId: 0,
        autoRotate: true,
        rotate: 360,
        duration: 5000,
        animationEnd: () => {
          if (this.isPageShowing) {
            moveCar()
          }
        }
      })
    }

    moveCar()
  },
})
