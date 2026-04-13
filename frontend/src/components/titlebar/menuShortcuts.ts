// considering wether to replace this with a standardised system or not
// this may be a bit overkill lol


export interface ShortcutMatcher {
  ctrl: boolean;
  shift: boolean;
  alt: boolean;
  meta: boolean;
  key: string;
}

const MODIFIERS = new Set(["ctrl", "shift", "alt", "cmd", "meta", "win"]);

const LABELS: Record<string, string> = {
  ctrl: "Ctrl",
  cmd: "Ctrl",
  meta: "Ctrl",
  win: "Win",
  shift: "Shift",
  alt: "Alt",
};

function norm(token: string) {
  const t = token.trim().toLowerCase();
  if (t === " " || t === "spacebar") return "space";
  if (t === "esc") return "escape";
  return t;
}

function cap(s: string) {
  return s.charAt(0).toUpperCase() + s.slice(1).toLowerCase();
}

function parse(combo: string): ShortcutMatcher | null {
  const parts = combo.split("+").map(norm).filter(Boolean);
  const key = parts.find((p) => !MODIFIERS.has(p));
  if (!key) return null;

  return {
    ctrl: parts.includes("ctrl"),
    shift: parts.includes("shift"),
    alt: parts.includes("alt"),
    meta:
      parts.includes("cmd") ||
      parts.includes("meta") ||
      parts.includes("win"),
    key,
  };
}

export function parseShortcut(shortcut?: string | string[]): {
  label: string;
  matchers: ShortcutMatcher[];
} {
  if (!shortcut) return { label: "", matchers: [] };

  const combos = Array.isArray(shortcut) ? shortcut : [shortcut];
  return {
    label: combos[0]
      .split("+")
      .map(norm)
      .map((p) => LABELS[p] ?? cap(p))
      .join("+"),
    matchers: combos
      .map(parse)
      .filter((m): m is ShortcutMatcher => m !== null),
  };
}

export function matchesEvent(
  event: KeyboardEvent,
  matchers: ShortcutMatcher[],
) {
  const key = norm(event.key);
  const code = event.code.toLowerCase();

  return matchers.some(
    (m) =>
      event.ctrlKey === m.ctrl &&
      event.shiftKey === m.shift &&
      event.altKey === m.alt &&
      event.metaKey === m.meta &&
      (key === m.key || code === m.key || code === `key${m.key}`),
  );
}

// prevent shortcuts from triggering when typing in inputs, textareas, selects, contenteditable elements, or role="textbox" elements
// this will probably be useful for shortcuts such as space for play/pause, or any shortcuts that involve letters/numbers
export function isEditableTarget(target: EventTarget | null) {
  return !!(target as HTMLElement | null)?.closest(
    'input, textarea, select, [contenteditable="true"], [role="textbox"]',
  );
}
