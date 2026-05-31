# LOG

## basic-ok fixture の確認

`test-css\basic-ok` が等価ケースとして成功するよう、単純な Emotion 補間の正規化を修正しました。

実装内容の詳細:
- Emotion CSS の `${COLOR}` のような単純な識別子補間を `var(--COLOR)` として扱うようにしました。
- 同名 SCSS 側の `var(--COLOR)` と正規化後 CSS AST が一致するようにしました。
- 識別子以外の複雑な補間は従来どおり placeholder に正規化します。
- parser の単体テストに、単純識別子補間が CSS variable になるケースを追加しました。

変更ファイル:
- `internal/parser/parser.go`
- `internal/parser/parser_test.go`
- `docs/log.md`

検証:
- `go test ./...` を実行し、成功しました。
- `go run .\cmd\emotion-to-scss check test-css\basic-ok` を実行し、`OK test-css\basic-ok\test.ts` を確認しました。

## CLIの実装

Emotion CSS と同一ディレクトリ・同一ファイル名の SCSS を比較する検証専用 CLI を実装しました。

実装内容の詳細:
- `check <input>` コマンドのみを持つ CLI エントリポイントを追加しました。
- `.ts`, `.tsx`, `.js`, `.jsx` の入力に対して、同じディレクトリにある同名 `.scss` を自動検出するようにしました。
- Emotion の `css` / `styled.*` tagged template literal を抽出し、変数名またはコンポーネント名を class selector として比較用 SCSS に変換する処理を追加しました。
- SCSS の簡易コンパイル処理を追加し、ネストセレクタと at-rule を CSS へ展開できるようにしました。
- CSS AST パーサを追加し、コメント、空白、末尾セミコロン差分を正規化して、セレクタ、at-rule、宣言プロパティ、宣言値、宣言順を比較するようにしました。
- 一致・不一致・不正引数の単体テストを追加しました。

変更ファイル:
- `cmd/emotion-to-scss/main.go`
- `internal/app/app.go`
- `internal/app/app_test.go`
- `internal/cssast/cssast.go`
- `internal/domain/domain.go`
- `internal/parser/parser.go`
- `internal/parser/parser_test.go`
- `internal/scss/compiler.go`
- `internal/scss/compiler_test.go`
- `internal/verifier/verifier.go`
- `internal/verifier/verifier_test.go`
- `docs/log.md`

検証:
- `go test ./...` を実行し、成功しました。
- `go build ./cmd/emotion-to-scss` と `go build -o emotion-to-scss.exe .\cmd\emotion-to-scss` は、Go の一時 exe 作成時に Windows Defender の `virus or potentially unwanted software` 判定でブロックされたため完了できませんでした。
