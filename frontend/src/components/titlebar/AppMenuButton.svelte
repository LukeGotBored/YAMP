<script lang="ts">
    // I am genuinely so sorry for what you're about to witness.

    import { MusicNoteIcon } from "phosphor-svelte";
    import type { DropdownItem, DropdownOption } from "../TitleBar";
    import {
        parseShortcut,
        matchesEvent,
        isEditableTarget,
        type ShortcutMatcher,
    } from "./menuShortcuts";

    let { dropdownOptions = [] }: { dropdownOptions?: DropdownOption[] } =
        $props();

    let open = $state(false);
    let keyboard = $state(false);
    let activeIdx = $state(-1);
    let hoverIdx = $state(-1);
    let pendingFocus = $state<number | null>(null);

    let containerEl = $state<HTMLDivElement>();
    let triggerEl = $state<HTMLButtonElement>();
    let itemEls = $state<HTMLButtonElement[]>([]);

    type Entry =
        | { type: "separator" }
        | {
              type: "item";
              idx: number;
              option: DropdownItem;
              scLabel: string;
              scMatchers: ShortcutMatcher[];
          };
    type ItemEntry = Extract<Entry, { type: "item" }>;

    let entries = $derived.by<Entry[]>(() => {
        let idx = 0;
        return dropdownOptions.map((opt) => {
            if (opt.type === "separator") return { type: "separator" } as const;
            const sc = parseShortcut(opt.shortcut);
            return {
                type: "item",
                idx: idx++,
                option: opt,
                scLabel: sc.label,
                scMatchers: sc.matchers,
            } as const;
        });
    });

    let items = $derived(
        entries.filter((e): e is ItemEntry => e.type === "item"),
    );

    $effect(() => {
        if (pendingFocus !== null && itemEls.length > 0) {
            const idx = pendingFocus;
            pendingFocus = null;
            highlight(idx);
        }
    });

    //
    function wrap(i: number) {
        return ((i % items.length) + items.length) % items.length;
    }

    function highlight(index: number) {
        if (items.length === 0) return;
        keyboard = true;
        activeIdx = wrap(index);
        itemEls[activeIdx]?.focus();
    }

    function openMenu() {
        if (open || items.length === 0) return;
        open = true;
        keyboard = false;
        activeIdx = -1;
        hoverIdx = -1;
    }

    function closeMenu(restoreFocus = true) {
        if (!open) return;
        open = false;
        keyboard = false;
        activeIdx = -1;
        itemEls = [];
        if (restoreFocus) triggerEl?.focus();
    }

    function invoke(option: DropdownItem) {
        if (option.disabled) return;
        option.action();
        closeMenu();
    }

    //
    function origin() {
        return keyboard ? activeIdx : hoverIdx;
    }

    function moveDown() {
        if (!open) {
            openMenu();
            pendingFocus = 0;
            return;
        }
        highlight(origin() + 1);
    }

    function moveUp() {
        if (!open) {
            openMenu();
            pendingFocus = items.length - 1;
            return;
        }
        const o = origin();
        highlight(o >= 0 ? o - 1 : items.length - 1);
    }

    //
    function onPointerEnter(index: number) {
        hoverIdx = index;
        if (keyboard) {
            keyboard = false;
            activeIdx = -1;
        }
    }

    function onOutsidePointer(e: PointerEvent) {
        if (open && containerEl && !containerEl.contains(e.target as Node)) {
            closeMenu(false);
        }
    }

    function onGlobalKey(e: KeyboardEvent) {
        if (isEditableTarget(e.target)) return;
        for (const item of items) {
            if (!item.option.disabled && matchesEvent(e, item.scMatchers)) {
                e.preventDefault();
                item.option.action();
                closeMenu(false);
                return;
            }
        }
    }

    function onTriggerKey(e: KeyboardEvent) {
        switch (e.key) {
            case "ArrowDown":
                e.preventDefault();
                moveDown();
                break;
            case "ArrowUp":
                e.preventDefault();
                moveUp();
                break;
            case "Enter":
            case " ":
                e.preventDefault();
                open ? closeMenu() : openMenu();
                break;
            case "Escape":
                if (open) {
                    e.preventDefault();
                    closeMenu();
                }
                break;
        }
    }

    function onItemKey(e: KeyboardEvent) {
        switch (e.key) {
            case "ArrowDown":
                e.preventDefault();
                moveDown();
                break;
            case "ArrowUp":
                e.preventDefault();
                moveUp();
                break;
            case "Home":
                e.preventDefault();
                highlight(0);
                break;
            case "End":
                e.preventDefault();
                highlight(items.length - 1);
                break;
            case "Escape":
                e.preventDefault();
                closeMenu();
                break;
            case "Tab":
                closeMenu(false);
                break;
        }
    }
</script>

<svelte:window onpointerdown={onOutsidePointer} onkeydown={onGlobalKey} />

<div class="menu" bind:this={containerEl}>
    <button
        id="app-menu-trigger"
        bind:this={triggerEl}
        class="menu-trigger"
        type="button"
        aria-label="Open application menu"
        aria-haspopup="menu"
        aria-controls="app-menu"
        aria-expanded={open}
        onclick={() => (open ? closeMenu() : openMenu())}
        onkeydown={onTriggerKey}
    >
        <span class="menu-icon" aria-hidden="true">
            <MusicNoteIcon weight="fill" />
        </span>
        <span class="menu-label">YAMP</span>
    </button>

    {#if open}
        <div
            id="app-menu"
            class="menu-dropdown"
            class:keyboard
            role="menu"
            aria-label="Application menu"
            aria-labelledby="app-menu-trigger"
        >
            {#each entries as entry, i (i)}
                {#if entry.type === "separator"}
                    <div class="menu-separator" role="separator"></div>
                {:else}
                    <button
                        bind:this={itemEls[entry.idx]}
                        class="menu-option"
                        class:active={keyboard && activeIdx === entry.idx}
                        type="button"
                        role="menuitem"
                        disabled={entry.option.disabled}
                        aria-keyshortcuts={entry.scLabel || undefined}
                        onclick={() => invoke(entry.option)}
                        onpointerenter={() => onPointerEnter(entry.idx)}
                        onkeydown={onItemKey}
                    >
                        {#if entry.option.icon}
                            {@const Icon = entry.option.icon}
                            <span class="menu-option-icon" aria-hidden="true">
                                <Icon />
                            </span>
                        {/if}

                        <span class="menu-option-label"
                            >{entry.option.label}</span
                        >

                        {#if entry.scLabel}
                            <span
                                class="menu-option-shortcut"
                                aria-hidden="true"
                            >
                                {entry.scLabel}
                            </span>
                        {/if}
                    </button>
                {/if}
            {/each}
        </div>
    {/if}
</div>

<style>
    .menu {
        position: relative;
        display: flex;
        align-items: center;
        height: 100%;
        -webkit-app-region: no-drag;
        --wails-draggable: no-drag;
    }

    .menu-trigger {
        height: 100%;
        display: inline-flex;
        align-items: center;
        gap: 10px;
        padding: 0 20px;
        font-weight: 600;
        color: var(--text-primary);
        transition: var(--t-fast);
        border-radius: 0;
        --wails-draggable: no-drag;
    }

    .menu-trigger:hover,
    .menu-trigger[aria-expanded="true"] {
        background: rgba(255, 255, 255, 0.05);
    }

    .menu-trigger:focus-visible {
        outline-offset: -3px;
        border-radius: 6px;
    }

    .menu-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 2px;
        border-radius: 4px;
        background: var(--accent);
        font-size: 12px;
        opacity: 1;
    }

    .menu-label {
        letter-spacing: 0.01em;
    }

    .menu-dropdown {
        position: absolute;
        top: calc(100% + 5px);
        left: 5px;
        width: max-content;
        padding: 6px;
        display: flex;
        flex-direction: column;
        border: 1px solid var(--border);
        border-radius: 8px;
        background: var(--bg-tb);
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
        z-index: 300;
        -webkit-app-region: no-drag;
    }

    .menu-separator {
        height: 1px;
        margin: 5px 0;
        background: var(--border);
    }

    .menu-option {
        min-height: 0;
        padding: 8px 12px;
        display: grid;
        grid-template-columns: 16px 1fr auto;
        align-items: center;
        gap: 12px;
        color: var(--text-primary);
        border-radius: 4px;
        text-align: left;
        font-size: 13px;
        transition: var(--t-fast);
    }

    .menu-option:focus-visible {
        outline: none;
    }

    .menu-option:disabled {
        opacity: 0.45;
        cursor: default;
    }

    /* Highlighted state — pointer hover or keyboard .active */
    .menu-option:hover,
    .keyboard .menu-option.active {
        background: var(--accent);
        color: #fff;
    }

    /* In keyboard mode, suppress hover on non-active items */
    .keyboard .menu-option:not(.active):hover {
        background: transparent;
        color: var(--text-primary);
    }

    .menu-option-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 16px;
        opacity: 0.8;
    }

    .menu-option:hover .menu-option-icon,
    .keyboard .menu-option.active .menu-option-icon {
        opacity: 1;
    }

    .keyboard .menu-option:not(.active):hover .menu-option-icon {
        opacity: 0.8;
    }

    .menu-option-label {
        min-width: 0;
    }

    .menu-option-shortcut {
        padding-left: 25px;
        font-size: 11px;
        color: color-mix(in srgb, currentColor 60%, transparent);
        white-space: nowrap;
    }

    .menu-option:hover .menu-option-shortcut,
    .keyboard .menu-option.active .menu-option-shortcut {
        color: rgba(255, 255, 255, 0.8);
    }

    .keyboard .menu-option:not(.active):hover .menu-option-shortcut {
        color: color-mix(in srgb, currentColor 60%, transparent);
    }
</style>
