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

    let isMaximized = false;

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
    <div class="tb-title">
        <div class="tb-title-icon">
            <MusicNoteIcon weight="fill" />
        </div>
        YAMP
    </div>
    <div class="tb-separator"></div>
    <div class="tb-tabs">
        <button class="tb-tab">Home</button>
        <button class="tb-tab">Library</button>
        <button class="tb-tab">Albums</button>
        <button class="tb-tab">Artists</button>
        <button class="tb-tab">Playlists</button>
    </div>
    <div class="wc">
        <button class="wc-btn wc-minimize" on:click={Window.Minimise}>
            <MinusIcon />
        </button>

        <button class="wc-btn wc-maximize" on:click={toggle}>
            {#if isMaximized}
                <CardsIcon />
            {:else}
                <SquareIcon />
            {/if}
        </button>

        <button class="wc-btn wc-close" on:click={Window.Close}>
            <XIcon />
        </button>
    </div>
</div>

<style>
    .titlebar {
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
        display: flex;
        flex-direction: row;
        height: 100%;
        gap: 25px;
        margin-left: 25px;
    }

    .wc {
        display: flex;
        align-items: center;
        -webkit-app-region: no-drag;
        pointer-events: all;
        margin-left: auto;

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
