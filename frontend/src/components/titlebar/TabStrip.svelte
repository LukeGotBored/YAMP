<script lang="ts">
    import type { ItemTab } from "../TitleBar";
    import { tabButtonId, tabPanelId } from "../TitleBar";

    let {
        tabs = [],
        selectedTab = $bindable(""),
    }: {
        tabs?: ItemTab[];
        selectedTab?: string;
    } = $props();

    let tabEls = $state<HTMLButtonElement[]>([]);

    let visibleTabs = $derived.by(() =>
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

<div class="tab-strip" role="tablist" aria-label="Library sections">
    {#each visibleTabs as tab, index}
        {@const isActive = tab.id === selectedTab}
        <button
            bind:this={tabEls[index]}
            id={tabButtonId(tab.id)}
            class="tab"
            class:active={isActive}
            type="button"
            role="tab"
            tabindex={isActive ? 0 : -1}
            aria-selected={isActive}
            aria-controls={tabPanelId(tab.id)}
            onclick={() => (selectedTab = tab.id)}
            onkeydown={(e) => onTabKeyDown(index, e)}
        >
            <span class="tab-label">{tab.label}</span>
        </button>
    {/each}
</div>

<style>
    .tab-strip {
        display: flex;
        align-items: stretch;
        height: 100%;
        min-width: 0;
        margin-left: 12px;
        --wails-draggable: no-drag;
    }

    .tab {
        opacity: 0.5;
        position: relative;
        min-width: 0;
        height: 100%;
        padding: 0 12.5px;
        display: inline-flex;
        align-items: center;
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
</style>
