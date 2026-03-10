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

    import { ArrowsClockwiseIcon, InfoIcon, XIcon } from "phosphor-svelte";

    let titleBar: TitleBar;

    // Optional: Sync selection if you need to access it outside of the views
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
            label: "Check for Updates...",
            icon: ArrowsClockwiseIcon,
            action: () => console.log("Updating..."),
            type: "text",
        },
        {
            label: "About YAMP",
            icon: InfoIcon,
            action: () => alert("YAMP v1.0.0"),
        },
        {
            type: "separator",
        },
        {
            label: "Quit",
            icon: XIcon,
            action: () => Window.Close(),
        },
    ];
</script>

<div class="main">
    <TitleBar
        bind:this={titleBar}
        on_select={select_tab}
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
