import camelcaseKeys = require("camelcase-keys")
import { auth } from "./proto_gen/auth/auth_pb"

export namespace Coolcar{
    const serverAddr = 'http://localhost:8080'
    const AUTH_ERR = 'AUTH_ERR' // 错误常量，判断是不是Unauthorized


    const authData = {
        token: '',
        expiryMs: 0,
    }

    export interface RequestOption<REQ,RES>{
        method: 'GET'|'PUT'|'POST'|'DELETE'
        path: string
        data: REQ
        respMarshaller: (r: object) => RES
    }
    export interface AuthOption{
        attachAuthHeader: boolean // 是否要携带header，login不用
        retryOnAuthError:boolean
    }

    // 流程管理函数
    export async function sendRequestWithAuthRetry<REQ,RES>(o:RequestOption<REQ,RES>,a?:AuthOption) : Promise<RES> {
        const authOpt = a || {
            attachAuthHeader : true,
            retryOnAuthError : true
        }
        try {
            // 1、登录
            await login()
            // 2、token有效后发送业务请求
            return await sendRequest(o,authOpt)
            
        } catch (err) {
            if (err === AUTH_ERR && authOpt.retryOnAuthError) {
                // 清除token状态
                authData.token = ''
                authData.expiryMs = 0
                // 如果错误是Unauthorized,则重试
               return sendRequestWithAuthRetry(o,{
                    attachAuthHeader:authOpt.attachAuthHeader, // 原来携带还继续携带
                    retryOnAuthError:false,// 重试过就不用再重试
                })
            }else{
                throw err
            }
        }
    }

    
    
    
 
    export async function login() {
        // token 有效，不用登录，返回直接进行业务请求
      if (authData.token && authData.expiryMs >= Date.now()) {
          return
      }
      const wxResp = await wxLogin() 
      const reqTimeMs = Date.now()
      const start = new Date(reqTimeMs)
      console.log("当前时间",start)
      const resp = await sendRequest<auth.v1.ILoginRequest,auth.v1.ILoginResponse>({
          method: 'POST',
          path: '/v1/auth/login',
          data: {
              // 1、先获取code
              code: wxResp.code,
          },
          // 3、返回token
          respMarshaller:auth.v1.LoginResponse.fromObject
      },{
          // 登录不需要
          attachAuthHeader:false,
          retryOnAuthError:false,
      })
         // 4、把返回的token存起来
      authData.token = resp.accessToken!
      authData.expiryMs = reqTimeMs + resp.expiresIn! * 1000  
      const end = new Date(authData.expiryMs)
      console.log("过期时间",end)
    }
    
    function sendRequest<REQ,RES>(o:RequestOption<REQ,RES>,a:AuthOption) : Promise<RES>{

        return new Promise((resolve,reject) => {
            const header:Record<string,any> = {}
            if (a.attachAuthHeader){
                if (authData.token && authData.expiryMs >= Date.now()){
                     header.authorization = 'Bearer' + authData.token
                }else{
                    reject(AUTH_ERR)
                    return
                }
            }
                
            // 2、发送request
            wx.request({
                url:serverAddr + o.path,
                method:o.method,
                data:o.data,
                header,
                success: res =>{
                    // 401 代表Unauthorized（头部无效，token过期）
                    if (res.statusCode === 401) {
                        reject(AUTH_ERR)
                    }else if(res.statusCode >= 400){
                        reject(res)
                    }else{
                        
                        resolve(o.respMarshaller(camelcaseKeys(res.data as object,{
                            deep:true,
                        })))
                    }
                    
                },
                // 网不好没发出去，发出去报错仍在success里
                fail:reject,

            })
        })
    }

    // 把login 改成Promise的形式
    function wxLogin():Promise<WechatMiniprogram.LoginSuccessCallbackResult> {
        return new Promise((resolve,reject)=>{
            wx.login({
                success: resolve,
                fail:reject,
            })
        })
    }
}