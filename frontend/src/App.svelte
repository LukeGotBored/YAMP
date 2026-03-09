<script lang="ts">
    import { Events } from "@wailsio/runtime";
    import TitleBar from "./components/TitleBar.svelte";

    import HomeView from "./views/HomeView.svelte";
    import LibraryView from "./views/LibraryView.svelte";
    import ArtistsView from "./views/ArtistsView.svelte";
    import AlbumsView from "./views/AlbumsView.svelte";
    import PlaylistView from "./views/PlaylistView.svelte";
    import ErrorView from "./views/ErrorView.svelte";
    import SettingsView from "./views/SettingsView.svelte";

    import { ItemTab } from "./components/TitleBar";
    import { onMount } from "svelte";
    import { Playlist } from "phosphor-svelte";

    let titleBar: TitleBar;

    function select_tab(item_name: string) {}

    const tabs: ItemTab[] = [
        { id: "home", label: "Home", view: HomeView },
        { id: "library", label: "Library", view: LibraryView },
        { id: "albums", label: "Albums", view: AlbumsView },
        { id: "artists", label: "Artists", view: ArtistsView },
        { id: "playlists", label: "Playlists", view: PlaylistView },
        { id: "settings", label: "Settings", view: SettingsView },
    ];
</script>

<div id="main">
    <TitleBar bind:this={titleBar} on_select={select_tab} {tabs} />
    {#each tabs as itemTab}
        {#if titleBar?.currentTab() === itemTab.id}
            {@const ItemTabView = itemTab.view}
            <ItemTabView></ItemTabView>
        {/if}
    {/each}
</div>

<style>
    #main {
        height: 100vh;
    }
</style>
