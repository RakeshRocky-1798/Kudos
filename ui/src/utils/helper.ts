export const filterOption = (
  input: string,
  option?: { label: string; value: string; aEmoji: string },
): boolean =>
  (option?.label.toLowerCase() ?? '').startsWith(input.toLowerCase());
