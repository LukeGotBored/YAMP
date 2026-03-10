<script lang="ts">
    import {
        MinusIcon,
        XIcon,
        SquareIcon,
        CardsIcon,
        MusicNoteIcon,
    } from "phosphor-svelte";
    import { Window, Events } from "@wailsio/runtime";
    import { onMount } from "svelte";
    import { ItemTab } from "./TitleBar";

    export interface DropdownOption {
        label?: string;
        action?: () => void;
        icon?: any;
        type?: "text" | "separator";
        shortcut?: string;
    }

    const {
        onSelect,
        tabs = [],
        dropdownOptions = [],
    }: {
        onSelect: (itemName: string) => void;
        tabs?: ItemTab[];
        dropdownOptions?: DropdownOption[];
    } = $props();

    let isMaximized = $state(false);
    let selectedTab = $state(tabs.length > 0 ? tabs[0].id : "");

    let isDropdownOpen = $state(false);
    let titleContainer: HTMLElement;

    export function currentTab() {
        return selectedTab;
    }

    onMount(() => {
        Window.IsMaximised().then((state) => (isMaximized = state));

        const unsubscribeMaximize = Events.On(
            "common:WindowMaximise",
            () => (isMaximized = true),
        );
        const unsubscribeRestore = Events.On(
            "common:WindowRestore",
            () => (isMaximized = false),
        );

        return () => {
            unsubscribeMaximize();
            unsubscribeRestore();
        };
    });

    async function toggleWindow() {
        if (isMaximized) {
            isMaximized = false;
            await Window.UnMaximise();
        } else {
            isMaximized = true;
            await Window.Maximise();
        }
    }

    function handleOutsideClick(event: MouseEvent) {
        if (
            isDropdownOpen &&
            titleContainer &&
            !titleContainer.contains(event.target as Node)
        ) {
            isDropdownOpen = false;
        }
    }

    function isShortcutMatch(event: KeyboardEvent, shortcut: string) {
        const parts = shortcut.split("+").map((p) => p.trim().toLowerCase());

        const ctrl = parts.includes("ctrl");
        const shift = parts.includes("shift");
        const alt = parts.includes("alt");
        const meta =
            parts.includes("cmd") ||
            parts.includes("meta") ||
            parts.includes("win");

        const key = parts.find(
            (p) => !["ctrl", "shift", "alt", "cmd", "meta", "win"].includes(p),
        );

        if (!key) return false;

        if (
            event.ctrlKey === ctrl &&
            event.shiftKey === shift &&
            event.altKey === alt &&
            event.metaKey === meta
        ) {
            return (
                event.key.toLowerCase() === key ||
                event.code.toLowerCase() === key ||
                event.code.toLowerCase() === `key${key}`
            );
        }

        return false;
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
                                        {option.shortcut}
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
            {/each}
        </div>
    </div>

    <div class="tb-right">
        <div class="wc">
            <button class="wc-btn wc-min" onclick={Window.Minimise}>
                <MinusIcon />
            </button>

            <button class="wc-btn wc-max" onclick={toggleWindow}>
                {#if isMaximized}
                    <CardsIcon />
                {:else}
                    <SquareIcon />
                {/if}
            </button>

            <button class="wc-btn wc-cls" onclick={Window.Close}>
                <XIcon />
            </button>
        </div>
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
        min-width: 250px;
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

    .wc {
        display: flex;
        align-items: stretch;
        -webkit-app-region: no-drag;
        pointer-events: all;
        margin-left: auto;
        height: 100%;
        --wails-draggable: no-drag;
    }

    .wc-btn {
        width: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--text-secondary);
        transition:
            background 0.1s,
            color 0.1s;
        -webkit-app-region: no-drag;
        cursor: pointer;
        pointer-events: all;
        background: none;
        border: none;
    }

    .wc-btn:hover {
        background: rgba(255, 255, 255, 0.07);
        color: var(--text-primary);
    }

    .wc-cls:hover {
        background: var(--accent);
        color: #fff;
    }
</style>
