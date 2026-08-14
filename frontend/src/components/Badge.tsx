interface BadgeProps {
  label: string
  variant?: string
}

// variant はCSSクラスのサフィックスに使う (badge-approved 等)
export function Badge({ label, variant }: BadgeProps) {
  return <span className={`badge${variant ? ` badge-${variant}` : ''}`}>{label}</span>
}
