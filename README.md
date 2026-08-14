# Require Manager（要求管理ツール）

要件定義・要求管理・用語辞書・ドメインモデリングを一元的に行うための
Web アプリケーションです。`docker compose` で以下の3コンテナが起動します。

| レイヤー | 技術 | ポート |
| --- | --- | --- |
| フロントエンド | React + TypeScript + Vite | 5173 |
| バックエンド | Go + Gin + GORM | 8080 |
| データベース | PostgreSQL 16 | 5432 |

## 主な機能

- **プロジェクト管理** — 複数プロジェクトを作成し、切り替えて管理
- **要求管理（要件定義）** — 業務要求／機能要件／非機能要件／制約を
  親子階層で登録。優先度・ステータス・根拠・出典を管理し、キーワード／
  種別／状態で絞り込み
- **用語辞書** — 用語・読み・英語表記・定義・別名・分類をカード形式で管理し、
  横断検索
- **ドメインモデリング** — エンティティ（属性つき）と関連（多重度つき）を
  UML クラス図風に可視化。Mermaid クラス図テキストも自動生成
- **ダッシュボード** — 各要素の件数と要求ステータスの内訳を表示

## クイックスタート

```bash
# 環境変数を用意（そのままでも動きます）
cp .env.example .env

# 起動
docker compose up --build
```

起動後にアクセスします。

- フロントエンド: http://localhost:5173
- バックエンド API: http://localhost:8080/api
- ヘルスチェック: http://localhost:8080/health

初回起動時にサンプルデータ（ECサイトのプロジェクト）が自動投入されます。

停止・破棄:

```bash
docker compose down          # 停止
docker compose down -v       # DB のデータも削除
```

## Storybook（コンポーネント確認）

フロントエンドの UI コンポーネントを単体で確認できる Storybook を用意しています。
バックエンドや DB には依存しません。

```bash
docker compose -f docker-compose.storybook.yml up --build
```

起動後 http://localhost:6006 でアクセスします。

登録済みのストーリー:

- **Components/Badge** — ステータスバッジ（全バリアント）
- **Components/Modal** — モーダル（表示例・開閉のインタラクティブ例）
- **Design System/Overview** — バッジ／ボタン／チップ、統計カード、
  UML クラスカード、用語カードの視覚カタログ

コンテナを使わずローカルで起動する場合:

```bash
cd frontend
npm install
npm run storybook        # 開発サーバ（http://localhost:6006）
npm run build-storybook  # 静的ビルド（storybook-static/）
```

## ディレクトリ構成

```
require_manager/
├── docker-compose.yml        # 3サービスの構成
├── .env.example              # 環境変数のサンプル
├── db/
│   └── init.sql              # 拡張機能の有効化（スキーマは GORM が管理）
├── backend/                  # Go + Gin
│   ├── main.go
│   ├── Dockerfile
│   └── internal/
│       ├── config/           # 環境変数の読み込み
│       ├── database/         # 接続・マイグレーション
│       ├── models/           # ドメインモデル（GORM）
│       ├── handlers/         # HTTP ハンドラ
│       ├── router/           # ルーティング + CORS
│       └── seed/             # サンプルデータ投入
└── frontend/                 # React + TypeScript + Vite
    ├── Dockerfile
    └── src/
        ├── api.ts            # API クライアント（axios）
        ├── types.ts          # 型定義
        ├── ProjectContext.tsx
        ├── components/       # Modal / Badge
        └── pages/            # 各画面
```

## API エンドポイント

すべて `/api` プレフィックス配下です。

| メソッド | パス | 説明 |
| --- | --- | --- |
| GET/POST | `/projects` | プロジェクト一覧・作成 |
| GET/PUT/DELETE | `/projects/:id` | 取得・更新・削除 |
| GET | `/projects/:id/summary` | ダッシュボード用の集計 |
| GET/POST | `/requirements` | 要求一覧（`projectId`,`status`,`type`,`keyword` で絞り込み）・作成 |
| GET/PUT/DELETE | `/requirements/:id` | 取得・更新・削除 |
| GET/POST | `/terms` | 用語一覧（`projectId`,`keyword`）・作成 |
| GET/PUT/DELETE | `/terms/:id` | 取得・更新・削除 |
| GET/POST | `/entities` | エンティティ一覧・作成 |
| GET/PUT/DELETE | `/entities/:id` | 取得・更新・削除（属性も一括更新） |
| GET/POST | `/relationships` | 関連一覧・作成 |
| PUT/DELETE | `/relationships/:id` | 更新・削除 |

## ローカル開発（コンテナを使わない場合）

**バックエンド**（PostgreSQL が別途必要）:

```bash
cd backend
DB_HOST=localhost DB_PORT=5432 DB_USER=reqmgr DB_PASSWORD=reqmgr_pass \
  DB_NAME=reqmgr go run ./main.go
```

**フロントエンド**:

```bash
cd frontend
npm install
npm run dev
```

## 環境変数

`.env`（`.env.example` 参照）で設定します。

| 変数 | 既定値 | 用途 |
| --- | --- | --- |
| `POSTGRES_USER` / `POSTGRES_PASSWORD` / `POSTGRES_DB` | `reqmgr` / `reqmgr_pass` / `reqmgr` | DB 認証 |
| `BACKEND_PORT` | `8080` | バックエンド公開ポート |
| `FRONTEND_PORT` | `5173` | フロントエンド公開ポート |
| `CORS_ALLOW_ORIGINS` | `http://localhost:5173` | 許可するオリジン（カンマ区切り） |
| `VITE_API_BASE_URL` | `http://localhost:8080/api` | フロントから見た API のベース URL |

## データモデル

- **Project** — すべての情報を束ねる単位
- **Requirement** — 要求。`parentId` による自己参照で階層化
- **Term** — 用語辞書エントリ
- **DomainEntity** / **DomainAttribute** — エンティティと属性
- **DomainRelationship** — エンティティ間の関連（種類・多重度・ラベル）

スキーマはバックエンド起動時に GORM の AutoMigrate で自動生成されます。
