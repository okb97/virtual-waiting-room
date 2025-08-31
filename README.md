# virtual-waiting-room

学習用の仮想待合室システムです。  
ユーザーが順番待ちに入り、待ち時間や進捗を動的に確認できます。

---

## 技術スタック

- **フロントエンド**: Next.js / React  
  待合室ページUI、カウントダウン、進捗バー表示

- **デプロイ/ホスティング**: Vercel (無料枠)  
  フロントページのCDN配信、サーバーレス関数としてAPIも動作可能

- **バックエンド（待合室API）**: Go (Serverless Function)  
  キュー管理、順番・予想待ち時間計算、入場トークン発行

- **キュー管理**: Redis (Upstash無料枠)  
  順番待ちリスト、TTL管理、FIFOキュー操作

- **通信方式**: HTTP REST + ポーリング  
  フロントからAPIへ数秒ごとに順番・進捗情報取得

- **認証**: なし（学習用）  
  トークン管理は最小限。ユーザー認証は他システムに依存せずモックでOK

---

## フォルダ構成
```text
/virtual-waiting-room
├─ frontend
│   ├─ components
│   │   ├─ Countdown.tsx        ← カウントダウン表示コンポーネント
│   │   ├─ ProgressBar.tsx      ← 待ち順の進捗バー表示
│   │   └─ TicketButton.tsx     ← チケットサイトにアクセスするボタン
│   └─ pages
│       ├─ index.tsx            ← 待合室トップページ
│       ├─ status.tsx           ← ポーリングで順番・進捗表示
│       └─ _app.tsx             ← 共通レイアウト・環境変数設定
├─ api
│   ├─ queue.go                 ← join, status, enter エンドポイント
│   └─ redis.go                 ← Redis操作ラッパー
├─ utils
│   ├─ timeUtils.go             ← 推定待ち時間計算関数
│   └─ queueUtils.go            ← 順番・進捗計算関数
└─ .env.local                   ← 接続先URLやRedis情報など環境変数
```

## 通信・動作フロー

1. ユーザーが待合室ページにアクセス → VercelのCDNから高速表示  
2. `/api/queue/join` を呼び出してキューに追加 → `ticketId` 発行  
3. フロントが数秒ごとに `/api/queue/status/:ticketId` をポーリング  
4. APIが Redis から順番・進捗を計算して返却  
5. フロントで残り時間と進捗バーを更新  
6. 順番が来たら `/api/queue/enter/:ticketId` で入場許可

## 無料枠内アクセス制限ß

- VercelやUpstashの無料枠を超えると自動課金される可能性があるため、アクセス制限を導入
- API側でキュー件数の上限をチェック
  - 上限超過時はHTTP 429を返して新規ユーザーを拒否
- フロント側で制御
  - 429返却時は「アクセス集中のため一時的に待合室に入れません」の画面を表示
- 上限例:
```go
MAX_QUEUE = 5000  # 無料枠想定の最大同時待機ユーザー数
```
