# emotion-to-scss

`emotion-to-scss` は TypeScript / JavaScript ファイル内の Emotion style template を検出し、変数名を class name として SCSS に変換する Go 製 CLI です。

対象拡張子:

- `.ts`
- `.tsx`
- `.js`
- `.jsx`

## インストール

```sh
go build ./cmd/emotion-to-scss
```

## 使い方

```sh
emotion-to-scss convert <input> --out <dir>
emotion-to-scss verify <input> --out <dir>
emotion-to-scss check <input>
```

`<input>` にはファイルまたはディレクトリを指定できます。ディレクトリ指定時は対象拡張子のファイルを再帰的に処理し、出力先にも相対パスを保持します。

## コマンド

### convert

Emotion style template を抽出して SCSS ファイルへ変換します。

```sh
emotion-to-scss convert src --out styles --overwrite
```

### verify

SCSS 生成後、変換元 style と生成 SCSS を CSS AST として比較します。コメント、空白、末尾セミコロンの差分は無視し、セレクタ、at-rule、宣言プロパティ、宣言値、同一セレクタ内の宣言順を比較します。

```sh
emotion-to-scss verify src/Button.tsx --out styles --report report.json
```

### check

対象ファイルの style template を抽出し、変換後 SCSS を CSS として検証します。ファイル出力は行いません。

```sh
emotion-to-scss check src
```

## オプション

- `--out <dir>`: SCSS の出力先。`convert` と `verify` で必須です。
- `--dry-run`: ファイルを書き込まず、変換結果を標準出力します。
- `--overwrite`: 既存の出力ファイルを上書きします。
- `--fail-on-unsupported`: unsupported syntax を失敗として扱うためのオプションです。現在、動的補間は placeholder 化して処理対象に含めます。
- `--report <path>`: JSON レポートを出力します。

## 対応する style template

変数宣言に代入された `css` tagged template literal と `styled.*` tagged template literal を対象にします。

```tsx
const buttonStyle = css`
  color: red;

  &:hover {
    color: blue;
  }
`
```

生成される SCSS:

```scss
.buttonStyle {
  color: red;

  &:hover {
    color: blue;
  }
}
```

動的補間は、検証可能な安定値にするため placeholder に置換します。

```tsx
const boxStyle = css`
  color: ${color};
`
```

生成される SCSS:

```scss
.boxStyle {
  color: __emotion_to_scss_dynamic_1__;
}
```

単独の動的補間は、検証可能な custom property placeholder に置換します。

```tsx
const boxStyle = css`
  ${baseStyle}
`
```

生成される SCSS:

```scss
.boxStyle {
  --emotion-to-scss-dynamic-1: __emotion_to_scss_dynamic_1__;
}
```

## Exit Code

- `0`: success
- `1`: verification failed or unsupported syntax with `--fail-on-unsupported`
- `2`: invalid CLI arguments
- `3`: file read/write error
- `4`: parse error
