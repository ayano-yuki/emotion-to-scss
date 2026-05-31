Go で CLI ツール emotion-to-scss を実装してください。

## 目的:
TypeScript / JavaScript ファイル内の Emotion `css` tagged template literal を検出し、静的 CSS を SCSS ファイルへ変換する。
また、変換前 Emotion CSS と、生成 SCSS をコンパイルした CSS を比較し、正規化後の CSS AST が等価であることを検証する。

## 必須コマンド:
```
1. convert <input> --out <dir>
2. verify <input> --out <dir>
3. check <input>
```

## 対象拡張子:
```
.ts, .tsx, .js, .jsx
```

## 条件:
- css tagged template literal のみ対応
- style 変数名を SCSS class name として使う

変換例:
```
import styled from '@emotion/styled'

const buttonStyle = styled.button`
  color: red;

  &:hover {
    color: blue;
  }
`
```

生成 SCSS:
```
.buttonStyle {
  color: red;

  &:hover {
    color: blue;
  }
}
```

## CLI オプション:
```
--out <dir>
--dry-run
--overwrite
--fail-on-unsupported
--report <path>
```

## verify の仕様:
- Emotion CSS を抽出する
- SCSS を生成する
- SCSS を CSS にコンパイルする
- 変換元 CSS と変換後 CSS を CSS AST としてパースする
- コメント、空白、末尾セミコロン差分は無視する
- セレクタ、at-rule、宣言プロパティ、宣言値は保持して比較する
- 同一セレクタ内の宣言順序は保持する
- 等価なら exit code 0
- 差分があれば exit code 1

## exit code:
0 success
1 verification failed or unsupported syntax with --fail-on-unsupported
2 invalid CLI arguments
3 file read/write error
4 parse error

## 出力レポート:
--report が指定された場合 JSON を出力する。
summary と files 配列を含める。
各 style には name, status, line, className, verification, reason を含める。

## 実装要件:
- Go module として構成する
- オニオンアーキテクチャを意識した実装にする
- cobra などの CLI ライブラリを使ってよい
- 単体テストを必ず書く
- parser, converter, verifier, report を分離する
- 主要な正常系・異常系テストを含める
- README に使い方を書く