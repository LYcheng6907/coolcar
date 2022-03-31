import { routing } from "../../utils/routing"

Page({
    redirectURL: '',

    data: {
        licNo: '',
        name: '',
        genderIndex: 0,
        genders: ['未知', '男', '女', '其他'],
        birthDate: '1990-01-01',
        licImgURL: '',
        state: 'UNSUBMITTED' as 'UNSUBMITTED'|'PENDING'|'VERIFIED',
    },

    // onLoad(opt: Record<'redirect', string>) {
    //     const o: routing.RegisterOpts = opt
    //     if(o.redirect) {
    //         this.redirectURL = decodeURIComponent(o.redirect)
    //     }
    // },

    onLoad(opt: Record<'redirect', string>) {
            const o: routing.RegisterOpts = opt
            if(o.redirect) {
                this.redirectURL = decodeURIComponent(o.redirect)
            }
        },

    onUploadLic() {
        wx.chooseImage({
            success: res => {
                if (res.tempFilePaths.length > 0) {
                    this.setData({
                        licImgURL: res.tempFilePaths[0]
                    })
                    // TODO: upload image
                    setTimeout(() => {
                        this.setData({
                            licNo: '3252452345',
                            name: '张三',
                            genderIndex: 1,
                            birthDate: '1989-12-02',
                        })
                    }, 1000)
                }
            }
        })
    },

    onGenderChange(e: any) {
        this.setData({
            genderIndex: e.detail.value,
        })
    },

    onBirthDateChange(e: any) {
        this.setData({
            birthDate: e.detail.value,
        })
    },

    onSubmit() {
        this.setData({
            state: 'PENDING',
        })
        setTimeout(() => {
            this.onLicVerified()
        }, 3000);
    },

    onResubmit() {
        this.setData({
            state: 'UNSUBMITTED',
            licImgURL: '',
        })
    },

    onLicVerified() {
        this.setData({
            state: 'VERIFIED',
        })

        //  wx.redirectTo({
        //     url:'/pages/lock/lock' ,
        //  }) 

        // 判断用户是否注册过，注册过则跳转到租车页面
        if (this.redirectURL) {
            wx.redirectTo({
                url: this.redirectURL,
            })
        }
    }
})