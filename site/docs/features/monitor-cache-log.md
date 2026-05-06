---
title: 监控、缓存与日志
description: 系统监控、Redis 缓存管理和操作日志说明
---

# 监控、缓存与日志

## 系统监控

系统监控展示服务状态、运行信息和数据库状态，适合用于本地开发检查和部署后的基础验收。

<img src="/ShiyuAdmin/img/yi-biao-pan.png" alt="Shiyu Admin 运维态势全景" />

## 在线用户

在线用户页面展示当前登录会话、主机 IP、登录地点、操作系统和浏览器版本，并支持强退会话。

<img src="/ShiyuAdmin/img/zai-xian-yong-hu.png" alt="Shiyu Admin 在线用户" />

## Redis 缓存管理

缓存管理支持查看 Redis 0-15 号逻辑库，按 Key 表达式和数据类型查询缓存，并读取以下数据类型：

- String
- List
- Set
- ZSet
- Hash
- Stream

## 操作日志

操作日志记录新增、修改、删除等写操作，并记录登录成功、登录失败和参数错误等认证审计事件。

## 数据监控

数据监控页面可以查看数据表、字段注释和基础数据，适合开发排查、验收数据结构和了解系统内置表。

<img src="/ShiyuAdmin/img/shu-ju-jian-kong.png" alt="Shiyu Admin 数据监控" />
