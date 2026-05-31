// css.ts
import styled from "@emotion/styled"
import { css } from "@emotion/react"

import {
  BUTTON_COLORS,
  BUTTON_SIZE_STYLES,
  type ButtonSize,
  type ButtonVariant,
} from "./constants"

type ButtonProps = {
  variant?: ButtonVariant
  size?: ButtonSize
  fullWidth?: boolean
  disabled?: boolean
}

export const Button = styled.button<ButtonProps>`
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;

  border: none;
  border-radius: 10px;

  font-weight: 600;
  line-height: 1;

  cursor: pointer;
  transition:
    background-color 120ms ease,
    transform 120ms ease,
    opacity 120ms ease;

  &:active {
    transform: translateY(1px);
  }

  &:focus-visible {
    outline: 3px solid rgba(37, 99, 235, 0.35);
    outline-offset: 2px;
  }

  ${({ variant = "primary" }) => {
    const color = BUTTON_COLORS[variant]

    return css`
      background-color: ${color.bg};
      color: ${color.text};

      &:hover {
        background-color: ${color.hover};
      }
    `
  }}

  ${({ size = "md" }) => BUTTON_SIZE_STYLES[size]}

  ${({ fullWidth }) =>
    fullWidth &&
    css`
      width: 100%;
    `}

  ${({ disabled }) =>
    disabled &&
    css`
      opacity: 0.5;
      cursor: not-allowed;
      pointer-events: none;
    `}
`
