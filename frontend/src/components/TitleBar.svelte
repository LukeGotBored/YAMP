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

    const {
        on_select,
        tabs = [],
    }: {
        on_select: (item_name: string) => void;
        tabs?: ItemTab[];
    } = $props();

    let isMaximized = $state(false);
    let selectedTab = $derived(tabs ? tabs[0].id : "");

    export function currentTab() {
        return selectedTab;
    }

    onMount(() => {
        Window.IsMaximised().then((state) => (isMaximized = state));

        const unsubMax = Events.On(
            "common:WindowMaximise",
            () => (isMaximized = true),
        );
        const unsubRestore = Events.On(
            "common:WindowRestore",
            () => (isMaximized = false),
        );

        return () => {
            unsubMax();
            unsubRestore();
        };
    });

    async function toggle() {
        if (isMaximized) {
            isMaximized = false;
            await Window.UnMaximise();
        } else {
            isMaximized = true;
            await Window.Maximise();
        }
    }
</script>

<div class="titlebar" style="--wails-draggable: drag;">
    <div id="tb-left">
        <div class="tb-title">
            <div class="tb-title-icon">
                <MusicNoteIcon weight="fill" />
            </div>
            YAMP
        </div>
        <div class="tb-separator"></div>
        <div class="tb-tabs">
            {#each tabs as tab}
                {@const isActive = tab.id === selectedTab}
                <button
                    class="tb-tab"
                    class:active={isActive}
                    onclick={() => {
                        on_select(tab.id);
                        selectedTab = tab.id;
                    }}
                >
                    {tab.label}
                </button>
            {/each}
        </div>
    </div>

    <div id="tb-right">
        <div class="wc">
            <button class="wc-btn wc-minimize" onclick={Window.Minimise}>
                <MinusIcon />
            </button>

            <button class="wc-btn wc-maximize" onclick={toggle}>
                {#if isMaximized}
                    <CardsIcon />
                {:else}
                    <SquareIcon />
                {/if}
            </button>

            <button class="wc-btn wc-close" onclick={Window.Close}>
                <XIcon />
            </button>
        </div>
    </div>
</div>

<style>
    .titlebar {
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

    #tb-left {
        display: flex;
        justify-content: flex-start;
        align-items: center;
        height: 100%;
    }

    #tb-right {
        display: flex;
        justify-content: flex-end;
        align-items: center;
        height: 100%;
    }

    .tb-title {
        padding: 0 20px;
        display: flex;
        justify-content: flex-start;
        align-items: center;
        line-height: 1;
        gap: 10px;
        font-weight: 600;
    }

    .tb-title-icon {
        background-color: var(--accent);
        padding: 2px;
        border-radius: 4px;
        font-size: 12px;
    }

    .tb-separator {
        height: 50%;
        border: var(--border) 1px solid;
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

    .wc-close:hover {
        background: var(--accent);
        color: #fff;
    }
</style>
