// @ts-ignore
/* eslint-disable */
// 此文件中的验证码接口为示例演示，当前项目后端暂无实现。
// 如需短信验证码登录，可在后端增加对应接口后，再完善此处逻辑。
export async function getFakeCaptcha() {
  return Promise.resolve({ code: 200, status: 'ok' } as API.FakeCaptcha);
}
