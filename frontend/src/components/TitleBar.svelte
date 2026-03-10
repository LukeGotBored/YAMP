<script lang="ts">
    import { MusicNoteIcon } from "phosphor-svelte";
    import WindowControls from "./WindowControls.svelte";
    import type { ItemTab } from "./TitleBar";

    export interface DropdownOption {
        label?: string;
        action?: () => void;
        icon?: any;
        type?: "text" | "separator";
        shortcut?: string | string[];
    }

    let {
        onSelect,
        tabs = [],
        dropdownOptions = [],
        selectedTab = $bindable(""),
    }: {
        onSelect: (itemName: string) => void;
        tabs?: ItemTab[];
        dropdownOptions?: DropdownOption[];
        selectedTab?: string;
    } = $props();

    let isDropdownOpen = $state(false);
    let titleContainer: HTMLElement;

    export function currentTab() {
        return selectedTab;
    }

    $effect(() => {
        if (!selectedTab && tabs.length > 0) {
            selectedTab = tabs[0].id;
        }
    });

    function handleOutsideClick(event: MouseEvent) {
        if (
            isDropdownOpen &&
            titleContainer &&
            !titleContainer.contains(event.target as Node)
        ) {
            isDropdownOpen = false;
        }
    }

    function formatShortcut(shortcut?: string | string[]) {
        if (!shortcut) return "";

        const primary = Array.isArray(shortcut) ? shortcut[0] : shortcut;

        return primary
            .replace(/cmd|meta/gi, "Ctrl")
            .replace(/win/gi, "Win")
            .split("+")
            .map((p) => p.charAt(0).toUpperCase() + p.slice(1).toLowerCase())
            .join("+");
    }

    function checkSingleShortcut(event: KeyboardEvent, shortcutStr: string) {
        const parts = shortcutStr.split("+").map((p) => p.trim().toLowerCase());

        const ctrl = parts.includes("ctrl");
        const shift = parts.includes("shift");
        const alt = parts.includes("alt");
        const meta =
            parts.includes("cmd") ||
            parts.includes("meta") ||
            parts.includes("win");

        const keys = parts.filter(
            (p) => !["ctrl", "shift", "alt", "cmd", "meta", "win"].includes(p),
        );

        const key = keys[0];
        if (!key) return false;

        if (
            event.ctrlKey === ctrl &&
            event.shiftKey === shift &&
            event.altKey === alt &&
            event.metaKey === meta
        ) {
            const pressedKey = event.key.toLowerCase();
            const pressedCode = event.code.toLowerCase();

            return (
                pressedKey === key ||
                pressedCode === key ||
                pressedCode === `key${key}` ||
                (key === "space" && pressedKey === " ")
            );
        }

        return false;
    }

    function isShortcutMatch(
        event: KeyboardEvent,
        shortcut: string | string[],
    ) {
        const combinations = Array.isArray(shortcut) ? shortcut : [shortcut];
        return combinations.some((combo) => checkSingleShortcut(event, combo));
    }

    function handleKeyDown(event: KeyboardEvent) {
        for (const option of dropdownOptions) {
            if (
                option.type !== "separator" &&
                option.shortcut &&
                option.action
            ) {
                if (isShortcutMatch(event, option.shortcut)) {
                    event.preventDefault();
                    option.action();
                    isDropdownOpen = false;
                    return;
                }
            }
        }
    }
</script>

<svelte:window onclick={handleOutsideClick} onkeydown={handleKeyDown} />

<div class="tb" style="--wails-draggable: drag;">
    <div class="tb-left">
        <div class="tb-title" bind:this={titleContainer}>
            <button
                class="tb-title-btn"
                onclick={() => (isDropdownOpen = !isDropdownOpen)}
            >
                <div class="tb-title-icon">
                    <MusicNoteIcon weight="fill" />
                </div>
                YAMP
            </button>

            {#if isDropdownOpen}
                <div class="tb-dd">
                    {#each dropdownOptions as option}
                        {#if option.type === "separator"}
                            <div class="tb-dd-sep"></div>
                        {:else}
                            <button
                                class="tb-dd-opt"
                                onclick={() => {
                                    option.action?.();
                                    isDropdownOpen = false;
                                }}
                            >
                                {#if option.icon}
                                    {@const Icon = option.icon}
                                    <div class="tb-dd-opt-icon">
                                        <Icon />
                                    </div>
                                {/if}
                                {option.label}

                                {#if option.shortcut}
                                    <span class="tb-dd-opt-sc">
                                        {formatShortcut(option.shortcut)}
                                    </span>
                                {/if}
                            </button>
                        {/if}
                    {/each}
                </div>
            {/if}
        </div>

        <div class="tb-sep"></div>

        <div class="tb-tabs">
            {#each tabs as tab}
                {@const isActive = tab.id === selectedTab}
                {#if !tab.hidden || isActive}
                    <button
                        class="tb-tab"
                        class:active={isActive}
                        onclick={() => {
                            onSelect(tab.id);
                            selectedTab = tab.id;
                        }}
                    >
                        {tab.label}
                    </button>
                {/if}
            {/each}
        </div>
    </div>

    <div class="tb-right">
        <WindowControls />
    </div>
</div>

<style>
    .tb {
        contain: layout style;
        height: var(--tb-h);
        flex-shrink: 0;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        background: var(--bg-tb);
        border-bottom: 1px solid var(--border);
        user-select: none;
        z-index: 200;
        box-sizing: border-box;
    }

    .tb-left,
    .tb-right {
        display: flex;
        align-items: center;
        height: 100%;
    }

    .tb-left {
        justify-content: flex-start;
    }

    .tb-right {
        justify-content: flex-end;
    }

    .tb-title {
        position: relative;
        display: flex;
        align-items: center;
        height: 100%;
        --wails-draggable: no-drag;
    }

    .tb-title-btn {
        padding: 0 20px;
        display: flex;
        align-items: center;
        gap: 10px;
        font-weight: 600;
        background: none;
        border: none;
        color: inherit;
        cursor: pointer;
        -webkit-app-region: no-drag;
        height: 100%;
        transition: background 0.1s;
    }

    .tb-title-btn:hover {
        background: rgba(255, 255, 255, 0.05);
    }

    .tb-title-icon {
        background-color: var(--accent);
        padding: 2px;
        border-radius: 4px;
        font-size: 12px;
        display: flex;
    }

    .tb-dd {
        position: absolute;
        top: calc(100% + 5px);
        left: 5px;
        background-color: var(--bg-tb);
        border: 1px solid var(--border);
        border-radius: 8px;
        display: flex;
        flex-direction: column;
        width: max-content;
        padding: 6px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
        z-index: 300;
        -webkit-app-region: no-drag;
    }

    .tb-dd-sep {
        height: 1px;
        background-color: var(--border);
        margin: 5px 0;
    }

    .tb-dd-opt {
        background: none;
        border: none;
        color: var(--text-primary);
        padding: 8px 12px;
        text-align: left;
        cursor: pointer;
        font-size: 13px;
        border-radius: 4px;
        display: flex;
        align-items: center;
        flex-direction: row;
        gap: 12px;
        transition:
            background 0.1s,
            color 0.1s;
    }

    .tb-dd-opt-sc {
        margin-left: auto;
        opacity: 0.5;
        font-size: 11px;
        font-family: inherit;
        letter-spacing: 0.5px;
        padding-left: 25px;
    }

    .tb-dd-opt:hover {
        background-color: var(--accent);
        color: #fff;
    }

    .tb-dd-opt:hover .tb-dd-opt-sc {
        opacity: 0.8;
    }

    .tb-dd-opt-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 16px;
        opacity: 0.8;
    }

    .tb-dd-opt:hover .tb-dd-opt-icon {
        opacity: 1;
    }

    .tb-sep {
        height: 50%;
        border-left: 1px solid var(--border);
    }

    .tb-tabs {
        --wails-draggable: no-drag;
        display: flex;
        flex-direction: row;
        height: 100%;
        margin-left: 12.5px;
    }

    .tb-tab {
        opacity: 0.5;
        position: relative;
        transition: var(--t-fast);
        padding: 0 12.5px;
        background: none;
        border: none;
        color: inherit;
        cursor: pointer;
    }

    .tb-tab::after {
        content: "";
        position: absolute;
        bottom: 0;
        left: 50%;
        width: 50%;
        height: 2px;
        background-color: var(--accent);
        transform: translateX(-50%) scaleX(0.5);
        opacity: 0;
        transition: var(--t-fast);
        transform-origin: center;
    }

    .tb-tab.active {
        opacity: 1;
    }

    .tb-tab.active::after {
        transform: translateX(-50%) scaleX(1);
        opacity: 1;
    }

</style>
