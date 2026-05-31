// constants.ts
import { css } from "@emotion/react"

export type ButtonVariant = "primary" | "danger" | "ghost"
export type ButtonSize = "sm" | "md" | "lg"

export const BUTTON_COLORS: Record<
  ButtonVariant,
  {
    bg: string
    hover: string
    text: string
  }
> = {
  primary: {
    bg: "#2563eb",
    hover: "#1d4ed8",
    text: "#ffffff",
  },
  danger: {
    bg: "#dc2626",
    hover: "#b91c1c",
    text: "#ffffff",
  },
  ghost: {
    bg: "transparent",
    hover: "#f3f4f6",
    text: "#111827",
  },
}

export const BUTTON_SIZE_STYLES: Record<ButtonSize, ReturnType<typeof css>> = {
  sm: css`
    height: 32px;
    padding: 0 12px;
    font-size: 13px;
  `,
  md: css`
    height: 40px;
    padding: 0 16px;
    font-size: 14px;
  `,
  lg: css`
    height: 48px;
    padding: 0 20px;
    font-size: 16px;
  `,
}
