import clsx from 'clsx';
import Heading from '@theme/Heading';
import Link from '@docusaurus/Link';
import Layout from '@theme/Layout';
import styles from './index.module.css';

const features = [
  {
    title: '开箱即用',
    description: 'Docker Compose 一条命令启动前端、后端、PostgreSQL 和 Redis，适合本地体验与二次开发。',
  },
  {
    title: '权限完整',
    description: '内置用户、角色、菜单、部门、JWT 登录认证、RBAC、动态菜单和接口权限控制。',
  },
  {
    title: '运维可见',
    description: '提供系统监控、Redis 缓存管理、数据管理和操作日志审计，便于排查与演示。',
  },
  {
    title: '部署友好',
    description: '覆盖 Docker、GitHub Pages、Render、Fly.io 和免费部署方案，适合个人项目上线。',
  },
];

const stack = [
  'Go',
  'Gin',
  'Gorm',
  'Viper',
  'JWT',
  'React',
  'Umi Max',
  'Ant Design Pro',
  'TypeScript',
  'ECharts',
  'PostgreSQL',
  'MySQL',
  'SQLite',
  'Redis',
  'Docker Compose',
];

function HomepageHeader() {
  return (
    <header className={styles.hero}>
      <div className={styles.heroInner}>
        <section className={styles.heroCopy}>
          <img className={styles.logo} src="/ShiyuAdmin/img/logo.png" alt="Shiyu Admin Logo" />
          <p className={styles.eyebrow}>开源通用后台管理系统</p>
          <Heading as="h1" className={styles.heroTitle}>
            Shiyu Admin 仕宇通用管理后台
          </Heading>
          <p className={styles.heroSubtitle}>
            基于 Go、Gin、Gorm、React、Ant Design Pro 和 RBAC 的前后端分离后台脚手架，
            适合快速搭建中后台、学习权限模型，或作为新业务系统的基础工程。
          </p>
          <div className={styles.actions}>
            <Link className="button button--primary button--lg" to="/docs/getting-started/quick-start">
              快速开始
            </Link>
            <Link className="button button--secondary button--lg" to="/docs/resources/cloud-offers">
              云服务器优惠
            </Link>
            <Link className="button button--outline button--lg" href="https://github.com/Rodert/ShiyuAdmin">
              GitHub
            </Link>
          </div>
        </section>
        <section className={styles.heroVisual} aria-label="Shiyu Admin 后台首页截图">
          <img src="/ShiyuAdmin/img/home-img.png" alt="Shiyu Admin 后台首页截图" />
        </section>
      </div>
    </header>
  );
}

function FeatureGrid() {
  return (
    <section className={styles.section}>
      <div className="container">
        <div className={styles.sectionHead}>
          <p className={styles.eyebrow}>Why Shiyu Admin</p>
          <Heading as="h2">适合直接开箱，也适合二次开发</Heading>
        </div>
        <div className={styles.featureGrid}>
          {features.map((feature, index) => (
            <article className={styles.featureCard} key={feature.title}>
              <span>{String(index + 1).padStart(2, '0')}</span>
              <Heading as="h3">{feature.title}</Heading>
              <p>{feature.description}</p>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}

function QuickStart() {
  return (
    <section className={clsx(styles.section, styles.quickStart)}>
      <div className="container">
        <div className={styles.quickStartGrid}>
          <div>
            <p className={styles.eyebrow}>Local Quick Start</p>
            <Heading as="h2">3 条命令启动完整环境</Heading>
            <p>
              本地只需要安装 Docker 和 Docker Compose，即可启动前端后台、后端接口、PostgreSQL 和 Redis。
            </p>
          </div>
          <pre>
            <code>{`git clone https://github.com/Rodert/ShiyuAdmin.git
cd ShiyuAdmin
docker compose up -d`}</code>
          </pre>
        </div>
      </div>
    </section>
  );
}

function Preview() {
  return (
    <section className={styles.section}>
      <div className="container">
        <div className={styles.sectionHead}>
          <p className={styles.eyebrow}>Preview</p>
          <Heading as="h2">演示图</Heading>
        </div>
        <div className={styles.screens}>
          <figure>
            <img src="/ShiyuAdmin/img/login-img.png" alt="Shiyu Admin 登录页截图" />
            <figcaption>登录页</figcaption>
          </figure>
          <figure>
            <img src="/ShiyuAdmin/img/home-img.png" alt="Shiyu Admin 后台首页截图" />
            <figcaption>后台首页</figcaption>
          </figure>
          <figure>
            <img src="/ShiyuAdmin/img/dashboard.png" alt="Shiyu Admin 数据仪表盘截图" />
            <figcaption>数据仪表盘</figcaption>
          </figure>
        </div>
      </div>
    </section>
  );
}

function Stack() {
  return (
    <section className={styles.section}>
      <div className="container">
        <div className={styles.stackPanel}>
          <div>
            <p className={styles.eyebrow}>Tech Stack</p>
            <Heading as="h2">组件官网引用已整理进文档</Heading>
            <p>
              文档内集中维护前后端、数据库、缓存、部署平台和文档站相关组件官网链接，便于读者继续学习。
            </p>
          </div>
          <div className={styles.tags}>
            {stack.map((item) => (
              <span key={item}>{item}</span>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}

export default function Home() {
  return (
    <Layout
      title="仕宇通用管理后台"
      description="Shiyu Admin 是一个 Go + React + RBAC 的开源通用后台管理系统。">
      <HomepageHeader />
      <main>
        <FeatureGrid />
        <QuickStart />
        <Preview />
        <Stack />
      </main>
    </Layout>
  );
}
