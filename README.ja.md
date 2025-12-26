# ShiyuAdmin 通用管理システム（日本語・簡易版）

> 作者：王仕宇（JavaPub）  
> 公式サイト：https://javapub.net.cn/  
> リポジトリ：**ShiyuAdmin**（GitHub でオープンソース化、中国国内アクセス用に Gitee ミラーあり）
> - GitHub（メインリポジトリ）：`https://github.com/Rodert/ShiyuAdmin`  
> - Gitee（中国ミラー）：`https://gitee.com/rodert/ShiyuAdmin`

---

## 1. 概要

ShiyuAdmin は、Go バックエンド + React（Ant Design Pro）フロントエンドで構成された、
汎用的な管理画面向けスキャフォールドです。

主な用途：

- 自社システム向けの管理画面テンプレート
- Go + Gin + Gorm + React + Ant Design Pro + RBAC の学習用サンプル
- 個人 / チームの社内ツールを素早く立ち上げるためのベース

※ 詳しい説明は、中国語版の `README.md` を参照してください。

---

## 2. 技術スタック（概要）

- バックエンド
  - Go 1.23+
  - Gin, Gorm
  - PostgreSQL / MySQL / SQLite
  - Redis
  - JWT 認証 + RBAC 権限管理

- フロントエンド
  - React 19
  - Umi Max
  - Ant Design & Ant Design Pro Components

---

## 3. クイックスタート（Docker 推奨）

前提：

- Docker と Docker Compose がインストールされていること

手順：

```bash
# リポジトリをクローン（フォークした場合は実際の URL に置き換えてください）
git clone https://github.com/Rodert/ShiyuAdmin.git
cd ShiyuAdmin

docker-compose up -d
```

起動後：

- バックエンド：`http://localhost:8080`
- フロントエンド：`http://localhost:8000`
- デフォルト管理者：`admin` / `Admin@123`

より詳しい開発手順（ローカル開発、既存システムへの統合など）は、
中国語版 `README.md` を参照してください。

---

## 4. コミュニティ & サポート

- **Star をお願いします**：このプロジェクトが役に立った場合は GitHub で Star をお願いします  
- **バグ報告**：[Issue を作成](https://github.com/Rodert/ShiyuAdmin/issues/new?labels=bug)  
- **機能要望**：[Issue を作成](https://github.com/Rodert/ShiyuAdmin/issues/new?labels=enhancement)  
- **PR の提出**：コード改善、新機能、ドキュメント更新を歓迎します  
- **連絡先**：WeChat 公式アカウント「JavaPub」をフォローすると作者に連絡できます（中国語）
