<script lang="ts">
    import WindowControls from "./WindowControls.svelte";
    import AppMenuButton from "./titlebar/AppMenuButton.svelte";
    import TabStrip from "./titlebar/TabStrip.svelte";
    import type { DropdownOption, ItemTab } from "./TitleBar";

    let {
        tabs = [],
        dropdownOptions = [],
        selectedTab = $bindable(""),
    }: {
        tabs?: ItemTab[];
        dropdownOptions?: DropdownOption[];
        selectedTab?: string;
    } = $props();

    const leftTabs = $derived.by(() =>
        tabs.filter((tab) => (tab.placement ?? "left") === "left"),
    );
    const rightTabs = $derived.by(() =>
        tabs.filter((tab) => tab.placement === "right"),
    );

    const visibleTabs = $derived.by(() =>
        tabs.filter((tab) => !tab.hidden || tab.id === selectedTab),
    );

    $effect(() => {
        if (
            visibleTabs.length > 0 &&
            !visibleTabs.some((tab) => tab.id === selectedTab)
        ) {
            selectedTab = visibleTabs[0].id;
        }
    });
</script>

<header class="tb" style="--wails-draggable: drag;" aria-label="Application title bar">
    <div class="tb-left">
        <AppMenuButton {dropdownOptions} />

        <div class="tb-sep" style="margin-left: 0; transform: translateX(-1px)" aria-hidden="true"></div>

        <TabStrip bind:selectedTab tabs={leftTabs} side="left" />
    </div>

    <div class="tb-right">
        {#if rightTabs.length > 0}
            <TabStrip bind:selectedTab tabs={rightTabs} side="right" />
            <div class="tb-sep" aria-hidden="true"></div>
        {/if}
        <WindowControls />
    </div>
</header>

<style>
    .tb {
        contain: layout style;
        height: var(--tb-h);
        flex-shrink: 0;
        display: flex;
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
        min-width: 0;
    }

    .tb-left {
        flex: 1 1 auto;
        gap: 0;
        justify-content: flex-start;
    }

    .tb-right {
        flex: 0 0 auto;
        gap: 0;
        justify-content: flex-end;
    }

    .tb-sep {
        height: 50%;
        border-left: 1px solid var(--border);
        margin: 0 12px;
    }
</style>
