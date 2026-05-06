---
title: 云服务器优惠
description: 云服务器、AI 编程工具优惠与开源使用成本说明
---

# 云资源优惠

如果你准备把 Shiyu Admin 部署到公网，云服务器是最直接的方式。下面整理了适合个人演示、小团队试用和长期部署的云资源优惠入口。

:::info 说明
部分链接可能包含推广或优惠标识。实际价格、库存、地域、带宽、续费规则和优惠资格以各平台活动页为准。
:::

## 云服务器

<div className="cloud-grid">
  <div className="cloud-card">
    <h3>阿里云 2 核 2G 3M</h3>
    <p>99 元云服务器，适合 Shiyu Admin 个人演示、轻量部署和学习环境。购买后建议马上续费一次，续费价同样为 99 元时更适合长期持有。</p>
    <a href="https://www.aliyun.com/minisite/goods?userCode=lw7epbel" target="_blank" rel="noreferrer">查看阿里云 99 元活动</a>
  </div>
  <div className="cloud-card">
    <h3>百度智能云 2 核 2G</h3>
    <p>低至 59 元，有 3 年可选，带宽可以按需选择。适合预算敏感的演示站、个人项目和轻量后台服务。</p>
    <a href="https://cloud.baidu.com/campaign/ambassador-product/index.html?ambassadorId=9a1b569b368440f6bd582fdaa3ec2702#knowledge-bcc" target="_blank" rel="noreferrer">查看百度智能云活动</a>
  </div>
  <div className="cloud-card">
    <h3>京东云 2 核 2G 3M</h3>
    <p>低至 58 元起，有 3 年可选。适合使用 Docker Compose 部署 Shiyu Admin、PostgreSQL 和 Redis 的入门环境。</p>
    <a href="https://3.cn/2-7uWUqd" target="_blank" rel="noreferrer">查看京东云活动</a>
  </div>
  <div className="cloud-card">
    <h3>腾讯云轻量 2 核 2G 4M</h3>
    <p>轻量应用服务器低至 79 元，适合快速部署后台演示、官网文档站和个人项目服务。</p>
    <a href="https://curl.qcloud.com/qqVt9Tp3" target="_blank" rel="noreferrer">查看腾讯云优惠合集</a>
  </div>
</div>

## AI 编程工具

<div className="cloud-grid">
  <div className="cloud-card">
    <h3>智谱 Coding Plan</h3>
    <p>国内编程大模型服务，适配 20+ 主流编程工具，适合代码生成、重构、解释和项目开发辅助。</p>
    <a href="https://www.bigmodel.cn/glm-coding?ic=CEO8NJ2ABY" target="_blank" rel="noreferrer">查看智谱 Coding Plan</a>
  </div>
  <div className="cloud-card">
    <h3>MiniMax Coding Plan</h3>
    <p>订阅 Coding Plan 可享 9 折优惠，适合需要长期使用 AI 编程能力的开发者。</p>
    <a href="https://platform.minimaxi.com/subscribe/coding-plan?code=EEOEFBbvhp&source=link" target="_blank" rel="noreferrer">查看 MiniMax Coding Plan 优惠</a>
  </div>
</div>

## Shiyu Admin 推荐配置

| 用途 | 推荐配置 |
| --- | --- |
| 个人演示 | 2 核 2G，Docker Compose，SQLite 或云数据库 |
| 小团队试用 | 2 核 4G，PostgreSQL，Redis |
| 长期运行 | 4 核 8G 起，独立数据库、备份、日志和 HTTPS |

## 开源使用成本

Shiyu Admin 本身是开源项目，仓库代码按 Apache-2.0 协议开放使用，项目源码不收费。真正需要预算的是运行和推广时的基础资源。

| 项目 | 是否必须 | 说明 |
| --- | --- | --- |
| Shiyu Admin 源码 | 否 | 开源免费，可本地运行、学习和二次开发 |
| 本地 Docker 环境 | 否 | 本地体验不需要云服务器，只需要自己的电脑安装 Docker |
| 云服务器 | 是，公网部署时 | 个人演示通常 2 核 2G 起步即可 |
| 域名 | 否 | 如果需要正式访问地址、备案和品牌展示，建议购买 |
| HTTPS 证书 | 否 | 可使用免费证书，生产环境建议开启 HTTPS |
| 数据库 / Redis | 视情况 | 可跟应用同机 Docker 部署，也可以购买托管服务 |
| 对象存储 / CDN | 否 | 有文件上传、图片分发或公网加速需求时再考虑 |
| AI 编程工具 | 否 | 只影响开发效率，不影响 Shiyu Admin 运行 |

低成本部署建议：

- 本地学习：0 元，直接使用 Docker Compose。
- 公网演示：选择 2 核 2G 云服务器，使用 Docker Compose 同机部署。
- 长期运行：云服务器 + 独立数据库 + 备份 + HTTPS。

## 页面维护建议

- 价格类文案只写“低至”或“活动价”，不要承诺长期有效。
- 推荐使用可追踪链接，方便统计转化。
- 明确标注推广或优惠属性，保持透明。
- 将部署教程和优惠页互相引用，提高转化路径清晰度。

> 价格、库存、地域和优惠规则以各云厂商官网活动页为准。
