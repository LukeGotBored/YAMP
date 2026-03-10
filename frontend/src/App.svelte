<script lang="ts">
    import TitleBar, {
        type DropdownOption,
    } from "./components/TitleBar.svelte";
    import { Window, Events } from "@wailsio/runtime";
    import { onMount } from "svelte";
    import {
        BACKEND_FATAL_ERROR_EVENT,
        BACKEND_PANIC_EVENT,
        formatBackendError,
    } from "./lib/backendErrors";
    import { triggerBackendPanic } from "./lib/debugService";

    // View Imports
    import LibraryView from "./views/LibraryView.svelte";
    import AlbumsView from "./views/AlbumsView.svelte";
    import ArtistsView from "./views/ArtistsView.svelte";
    import PlaylistView from "./views/PlaylistView.svelte";
    import SettingsView from "./views/SettingsView.svelte";
    import ErrorView from "./views/ErrorView.svelte";
    import type { ItemTab } from "./components/TitleBar";

    import {
        ArrowsClockwiseIcon,
        InfoIcon,
        XIcon,
        BookIcon,
        BugIcon,
        TerminalIcon,
        ArrowClockwiseIcon,
        SkullIcon,
    } from "phosphor-svelte";

    let selectedTab = $state("library");
    let fatalError = $state<string | null>(null);

    function selectTab(itemName: string) {
        console.log(`Tab changed to: ${itemName}`);
    }

    onMount(() => {
        const unsubPanic = Events.On(BACKEND_PANIC_EVENT, (event) => {
            fatalError = formatBackendError("[PANIC]", event.data);
        });

        const unsubFatalError = Events.On(
            BACKEND_FATAL_ERROR_EVENT,
            (event) => {
                fatalError = formatBackendError("[FATAL]", event.data);
            },
        );

        return () => {
            unsubPanic();
            unsubFatalError();
        };
    });

    const tabs: ItemTab[] = [
        { id: "library", label: "Library", view: LibraryView },
        { id: "albums", label: "Albums", view: AlbumsView },
        { id: "artists", label: "Artists", view: ArtistsView },
        { id: "playlists", label: "Playlists", view: PlaylistView },
        { id: "settings", label: "Settings", view: SettingsView },
    ];

    const appMenu: DropdownOption[] = [
        {
            label: "About YAMP",
            icon: InfoIcon,
            action: () => alert("YAMP v1.0.0"),
        },
        {
            label: "Documentation",
            icon: BookIcon,
            action: () => alert("Documentation is not available yet."),
            shortcut: "F1",
        },
        {
            label: "Check for Updates...",
            icon: ArrowsClockwiseIcon,
            action: () => alert("No updates available."),
            type: "text",
        },
        {
            type: "separator",
        },
        {
            label: "Developer Tools",
            icon: TerminalIcon,
            action: () => Window.OpenDevTools(),
            shortcut: "F12",
        },
        {
            label: "Reload",
            icon: ArrowClockwiseIcon,
            action: () => Window.Reload(),
            shortcut: ["Cmd+R", "Ctrl+R"],
        },
        {
            label: "Force Backend Panic",
            icon: BugIcon,
            action: () => triggerBackendPanic(),
            shortcut: ["Cmd+Alt+B", "Ctrl+Alt+B"],
        },
        {
            type: "separator",
        },
        {
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
        <TitleBar
            bind:selectedTab
            onSelect={selectTab}
            {tabs}
            dropdownOptions={appMenu}
        />

        <div class="content-area">
            {#each tabs as itemTab}
                {#if selectedTab === itemTab.id}
                    {@const ItemTabView = itemTab.view}
                    <ItemTabView />
                {/if}
            {/each}
        </div>
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
    }
</style>
