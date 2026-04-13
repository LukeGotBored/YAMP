<script lang="ts">
    import {
        MinusIcon,
        XIcon,
        SquareIcon,
        CardsIcon,
    } from "phosphor-svelte";
    import { Window, Events } from "@wailsio/runtime";
    import { onMount } from "svelte";

    let isMaximized = $state(false);

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
</script>

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

<style>
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
