<script lang="ts">
    import {
        BugIcon,
        CopyIcon,
        ArrowsClockwiseIcon,
        XIcon,
    } from "phosphor-svelte";
    import WindowControls from "../components/WindowControls.svelte";
    import { Window } from "@wailsio/runtime";

    let {
        error = "An unknown fatal error occurred.",
    }: {
        error?: string;
    } = $props();

    async function copyError() {
        await navigator.clipboard.writeText(error);
        alert("Error trace copied to clipboard!");
    }

    function reloadApp() {
        Window.Reload();
    }

    function quitApp() {
        Window.Close();
    }
</script>

<main class="panic-screen" style="--wails-draggable: drag;">
    <div class="window-controls-wrap">
        <WindowControls />
    </div>

    <div id="icon">
        <BugIcon weight="fill" />
    </div>

    <div id="main">Oh no!</div>

    <div id="description">
        Something went (very) wrong and YAMP cannot continue. Please copy the
        details below and report this issue, thanks!
    </div>

    <div class="error-trace" style="--wails-draggable: no-drag;">
        <code>{error}</code>
    </div>

    <div class="actions" style="--wails-draggable: no-drag;">
        <button class="btn secondary" onclick={copyError}>
            <CopyIcon /> Copy Trace
        </button>
        <button class="btn secondary" onclick={reloadApp}>
            <ArrowsClockwiseIcon /> Reload Window
        </button>
        <button class="btn danger" onclick={quitApp}>
            <XIcon /> Quit App
        </button>
    </div>
</main>

<style>
    .panic-screen {
        height: 100vh;
        width: 100vw;
        position: relative;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
        gap: 15px;
        background-color: #1a0f0f;
        color: #fff;
        padding: 40px;
        box-sizing: border-box;
    }

    .window-controls-wrap {
        position: absolute;
        top: 0;
        right: 0;
        height: var(--tb-h);
        z-index: 2;
    }

    #icon {
        font-size: 4em;
        color: #ff5f56;
    }

    #main {
        font-weight: 700;
        font-size: 2.5em;
        margin-bottom: -5px;
    }

    #description {
        font-size: 1.1em;
        opacity: 0.8;
        text-align: center;
        max-width: 600px;
        text-wrap: balance;
    }

    .error-trace {
        background: rgba(0, 0, 0, 0.5);
        border: 1px solid rgba(255, 95, 86, 0.3);
        border-radius: 8px;
        padding: 15px;
        width: 100%;
        max-width: 800px;
        max-height: 300px;
        overflow-y: auto;
        text-align: left;
        margin: 20px 0;
        user-select: text;
    }

    code {
        font-family: monospace;
        font-size: 0.9em;
        white-space: pre-wrap;
        word-wrap: break-word;
        color: #ffb8b8;
    }

    .actions {
        display: flex;
        gap: 15px;
    }

    .btn {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 10px 20px;
        font-size: 1em;
        font-weight: 600;
        border: none;
        border-radius: 6px;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn.secondary {
        background-color: rgba(255, 255, 255, 0.1);
        color: #fff;
    }

    .btn.secondary:hover {
        background-color: rgba(255, 255, 255, 0.2);
    }

    .btn.danger {
        background-color: #ff5f56;
        color: #fff;
    }

    .btn.danger:hover {
        background-color: #e0443e;
    }
</style>
