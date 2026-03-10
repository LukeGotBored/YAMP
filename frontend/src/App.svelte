<script lang="ts">
    import TitleBar, {
        type DropdownOption,
    } from "./components/TitleBar.svelte";
    import { Window, Events } from "@wailsio/runtime";
    import LibraryView from "./views/LibraryView.svelte";
    import AlbumsView from "./views/AlbumsView.svelte";
    import ArtistsView from "./views/ArtistsView.svelte";
    import PlaylistView from "./views/PlaylistView.svelte";
    import SettingsView from "./views/SettingsView.svelte";
    import { ItemTab } from "./components/TitleBar";

    import {
        ArrowsClockwiseIcon,
        InfoIcon,
        XIcon,
        BookIcon,
    } from "phosphor-svelte";

    let titleBar: TitleBar;

    function select_tab(item_name: string) {
        console.log(`Tab changed to: ${item_name}`);
    }

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
            action: () =>
                alert(
                    "No updates available. (not even sure where to check, tbh)",
                ),
            type: "text",
        },
        {
            type: "separator",
        },
        {
            label: "Quit",
            icon: XIcon,
            action: () => Window.Close(), // TODO: maybe add a confirmation dialog here?
        },
    ];
</script>

<div class="main">
    <TitleBar
        bind:this={titleBar}
        onSelect={select_tab}
        {tabs}
        dropdownOptions={appMenu}
    />

    {#each tabs as itemTab}
        {#if titleBar?.currentTab() === itemTab.id}
            {@const ItemTabView = itemTab.view}
            <ItemTabView></ItemTabView>
        {/if}
    {/each}
</div>

<style>
    .main {
        height: 100vh;
        display: flex;
        flex-direction: column;
    }
</style>
