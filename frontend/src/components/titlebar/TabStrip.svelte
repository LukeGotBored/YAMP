<script lang="ts">
    import type { ItemTab } from "../TitleBar";
    import { tabButtonId, tabPanelId } from "../TitleBar";

    let {
        tabs = [],
        selectedTab = $bindable(""),
        side = "left",
    }: {
        tabs?: ItemTab[];
        selectedTab?: string;
        side?: "left" | "right";
    } = $props();

    let tabEls = $state<HTMLButtonElement[]>([]);

    let visibleTabs = $derived.by(() =>
        tabs.filter((tab) => !tab.hidden || tab.id === selectedTab),
    );

    function focusTab(index: number) {
        if (visibleTabs.length === 0) return;
        const i =
            ((index % visibleTabs.length) + visibleTabs.length) %
            visibleTabs.length;
        selectedTab = visibleTabs[i].id;
        tabEls[i]?.focus();
    }

    function onTabKeyDown(index: number, e: KeyboardEvent) {
        switch (e.key) {
            case "ArrowRight":
            case "ArrowDown":
                e.preventDefault();
                focusTab(index + 1);
                break;
            case "ArrowLeft":
            case "ArrowUp":
                e.preventDefault();
                focusTab(index - 1);
                break;
            case "Home":
                e.preventDefault();
                focusTab(0);
                break;
            case "End":
                e.preventDefault();
                focusTab(visibleTabs.length - 1);
                break;
        }
    }
</script>

<div
    class="tab-strip"
    class:align-right={side === "right"}
    role="tablist"
    aria-label={side === "right" ? "Secondary sections" : "Primary sections"}
>
    {#each visibleTabs as tab, index}
        {@const isActive = tab.id === selectedTab}
        {@const Icon = tab.icon}
        {@const iconWeight = isActive
            ? (tab.activeIconWeight ?? "fill")
            : (tab.iconWeight ?? "regular")}
        <button
            bind:this={tabEls[index]}
            id={tabButtonId(tab.id)}
            class="tab"
            class:active={isActive}
            type="button"
            role="tab"
            tabindex={isActive ? 0 : -1}
            aria-selected={isActive}
            aria-label={tab.ariaLabel ?? tab.label ?? tab.id}
            aria-controls={tabPanelId(tab.id)}
            onclick={() => (selectedTab = tab.id)}
            onkeydown={(e) => onTabKeyDown(index, e)}
        >
            {#if Icon}
                <span class="tab-icon" aria-hidden="true">
                    <Icon size={tab.iconSize ?? 15} weight={iconWeight} />
                </span>
                <!-- TODO: not sure if I should keep this -->
            {/if}
            {#if tab.label}
                <span class="tab-label">{tab.label}</span>
            {/if}
        </button>
    {/each}
</div>

<style>
    .tab-strip {
        display: flex;
        align-items: stretch;
        height: 100%;
        min-width: 0;
        gap: 2px;
        --wails-draggable: no-drag;
    }

    .tab-strip.align-right {
        justify-content: flex-end;
    }

    .tab {
        opacity: 0.5;
        position: relative;
        min-width: 0;
        height: 100%;
        padding: 0 10px;
        display: inline-flex;
        align-items: center;
        gap: 7px;
        color: inherit;
        outline-color: var(--accent);
        outline-offset: -2px;
        outline-width: 2px;
        border-radius: var(--r-sm);
        transition: var(--t-fast);
    }

    .tab::after {
        content: "";
        position: absolute;
        left: 50%;
        width: 50%;
        bottom: 0;
        height: 2px;
        background: var(--accent);
        opacity: 0;
        transform: translateX(-50%) scaleX(0.5);
        transform-origin: center;
        transition: var(--t-smooth);
    }

    .tab:hover {
        opacity: 0.78;
    }

    .tab.active {
        opacity: 1;
    }

    .tab.active::after {
        opacity: 1;
        transform: translateX(-50%) scaleX(1);
    }

    .tab:focus-visible {
        outline-offset: -2px;
        border-radius: var(--r-sm);
    }

    .tab-label {
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .tab-icon {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        flex: 0 0 auto;
        line-height: 1;
    }
</style>
