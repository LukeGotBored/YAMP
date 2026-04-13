<script lang="ts">
    import TitleBar from "./components/TitleBar.svelte";
    import { Window, Events } from "@wailsio/runtime";
    import { onMount } from "svelte";
    import {
        BACKEND_FATAL_ERROR_EVENT,
        BACKEND_PANIC_EVENT,
        formatBackendError,
    } from "./lib/backendErrors";
    import { openAboutDialog } from "./lib/sysPopupService";

    // View Imports
    import LibraryView from "./features/library/ui/LibraryView.svelte";
    import AlbumsView from "./views/AlbumsView.svelte";
    import ArtistsView from "./views/ArtistsView.svelte";
    import PlaylistView from "./views/PlaylistView.svelte";
    import SettingsView from "./views/SettingsView.svelte";
    import ErrorView from "./views/ErrorView.svelte";
    import DebugView from "./views/DebugView.svelte";
    import {
        tabButtonId,
        tabPanelId,
        type DropdownOption,
        type ItemTab,
    } from "./components/TitleBar";

    import {
        ArrowsClockwiseIcon,
        InfoIcon,
        XIcon,
        BookIcon,
        BugIcon,
        TerminalIcon,
        ArrowClockwiseIcon,
        GearIcon,
    } from "phosphor-svelte";

    let selectedTab = $state("library");
    let fatalError = $state<string | null>(null);

    onMount(() => {
        // Global Frontend Error Catching
        window.addEventListener("unhandledrejection", (e) => {
            console.error("Unhandled promise rejection:", e.reason);
            // We can decide to show an error view or log locally
        });
        window.addEventListener("error", (e) => {
            console.error("Global error:", e.error);
        });

        // Backend Events
        const unsubPanic = Events.On(BACKEND_PANIC_EVENT, (event) => {
            fatalError = formatBackendError("[PANIC]", event.data);
        });

        const unsubFatalError = Events.On(
            BACKEND_FATAL_ERROR_EVENT,
            (event) => {
                fatalError = formatBackendError("[FATAL]", event.data);
            },
        );

        const buildDiagnosticReport = (title: string, msg: string, stack: string = "") => {
            return `=========================================
${title}
=========================================
Time:      ${new Date().toISOString()}
Message:   ${msg}

--- Environment ---
UserAgent: ${navigator.userAgent}
Location:  ${window.location.href}
Viewport:  ${window.innerWidth}x${window.innerHeight}

--- Stack Trace ---
${stack || "No stack trace available."}
=========================================`;
        };

        const handleGlobalError = (event: ErrorEvent) => {
            const msg = event.message || "";
            // Ignore benign browser loops and dev-server/extension injection errors
            if (
                msg.includes("ResizeObserver") ||
                msg.includes("custom.js") ||
                msg.includes("Unexpected token '<'") // Often caused by Vite dev server caching issues or 404s
            ) {
                console.warn("Recovered from benign error/script failure:", msg);
                return;
            }
            
            console.error("Unhandled Frontend Error:", event.error || msg);
            fatalError = buildDiagnosticReport(
                "[FRONTEND CRASH]", 
                msg, 
                event.error?.stack || ""
            );
        };

        const handleUnhandledRejection = (event: PromiseRejectionEvent) => {
            console.error("Unhandled Promise Rejection:", event.reason);
            const reason = event.reason;
            const msg = reason instanceof Error ? reason.message : String(reason);
            const trace = reason instanceof Error ? reason.stack : String(reason);
            fatalError = buildDiagnosticReport(
                "[UNHANDLED PROMISE REJECTION]", 
                msg, 
                trace
            );
        };

        window.addEventListener("error", handleGlobalError);
        window.addEventListener("unhandledrejection", handleUnhandledRejection);

        return () => {
            unsubPanic();
            unsubFatalError();
            window.removeEventListener("error", handleGlobalError);
            window.removeEventListener("unhandledrejection", handleUnhandledRejection);
        };
    });

    const tabs: ItemTab[] = [
        { id: "library", label: "Library", view: LibraryView },
        { id: "albums", label: "Albums", view: AlbumsView },
        { id: "artists", label: "Artists", view: ArtistsView },
        { id: "playlists", label: "Playlists", view: PlaylistView },
        { id: "debug", label: "Developer", view: DebugView, hidden: true },
        { id: "settings", label: "Settings", view: SettingsView, hidden: true },
    ];

    let activeTab = $derived.by(
        () => tabs.find((tab) => tab.id === selectedTab) ?? tabs[0],
    );

    const appMenu: DropdownOption[] = [
        {
            type: "item",
            label: "About YAMP",
            icon: InfoIcon,
            action: () => openAboutDialog(),
        },
        {
            type: "item",
            label: "Help & Resources",
            icon: BookIcon,
            action: () => alert("Documentation is not available yet."),
            shortcut: "F1",
        },
        {
            type: "item",
            label: "Check for Updates...",
            icon: ArrowsClockwiseIcon,
            action: () => alert("No updates available."),
        },
        {
            type: "separator",
        },
        {
            type: "item",
            label: "WebView Tools",
            icon: TerminalIcon,
            action: () => Window.OpenDevTools(),
            shortcut: "F12",
        },
                {
            type: "item",
            label: "Developer Tools",
            icon: TerminalIcon,
            action: () => (selectedTab = "debug"),
            shortcut: ["Cmd+Alt+D", "Ctrl+Alt+D"],
        },
        {
            type: "item",
            label: "Reload",
            icon: ArrowClockwiseIcon,
            action: () => Window.Reload(),
            shortcut: ["Cmd+R", "Ctrl+R"],
        },
        {
            type: "separator",
        },
        {
            type: "item",
            label: "Settings",
            icon: GearIcon,
            action: () => (selectedTab = "settings"),
            shortcut: ["Cmd+,", "Ctrl+,"],
        },
        {
            type: "item",
            label: "Quit",
            icon: XIcon,
            action: () => Window.Close(),
            shortcut: ["Cmd+Q", "Ctrl+Q"],
        },
    ];
</script>

<div class="main">
    {#if fatalError}
        <ErrorView error={fatalError} />
    {:else}
        <TitleBar bind:selectedTab {tabs} dropdownOptions={appMenu} />

        <main class="content-area">
            {#each tabs as tab}
                {@const TabView = tab.view}
                <div
                    class="content-panel"
                    class:hidden={selectedTab !== tab.id}
                    role="tabpanel"
                    id={tabPanelId(tab.id)}
                    aria-labelledby={tabButtonId(tab.id)}
                >
                    <TabView />
                </div>
            {/each}
        </main>
    {/if}
</div>

<style>
    .main {
        height: 100vh;
        display: flex;
        flex-direction: column;
    }

    .content-area {
        flex: 1;
        overflow-y: auto;
        position: relative;
        background: var(--bg-content);
    }

    .content-panel {
        height: 100%;
    }
    
    .hidden {
        display: none !important;
    }
</style>
