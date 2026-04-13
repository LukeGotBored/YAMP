<!-- THIS SCREEN IS CURRENTLY WORK IN PROGRESS, NOTHING -->
<!-- and I can't stress this enough -->
<!-- NOTHING -->
<!-- is final -->

<script lang="ts">
    import { onMount } from "svelte";
    import { ClockIcon, PlayCircleIcon } from "phosphor-svelte";
    import * as DB from "../../bindings/changeme/internal/db/queries";
    import type { Track } from "../../bindings/changeme/internal/db/models";
    import { ShowLibraryItemMenu } from "../../bindings/changeme/internal/app/context";

    let tracks = $state<Track[]>([]);
    let isLoading = $state(true);
    let hoveredTrack = $state<number | null>(null);

    async function loadTracks() {
        isLoading = true;
        try {
            tracks = (await DB.ListTracks()) || [];
        } catch (e) {
            console.error("Failed to load tracks", e);
        } finally {
            isLoading = false;
        }
    }

    function handleTrackRightClick(event: MouseEvent, index: number, trackId: number) {
        event.preventDefault();
        // Display the native OS Context Menu using Wails v3
        ShowLibraryItemMenu(event.clientX, event.clientY, trackId.toString());
    }

    onMount(() => {
        loadTracks();
    });

    // Helper to format seconds into mm:ss
    function formatTime(seconds: number) {
        if (!seconds) return "0:00";
        const m = Math.floor(seconds / 60);
        const s = Math.floor(seconds % 60);
        return `${m}:${s.toString().padStart(2, "0")}`;
    }
</script>

<div class="library-view">
    <div class="toolbar">
        <h2>Songs</h2>
        <span class="track-count">{tracks.length} items</span>
    </div>

    {#if isLoading}
        <div class="empty-state">
            <p>Loading library...</p>
        </div>
    {:else if tracks.length === 0}
        <div class="empty-state">
            <p>Your library is empty. Go to Developer tool to generate dummies!</p>
        </div>
    {:else}
        <div class="track-table-container">
            <table class="track-table">
                <thead>
                    <tr>
                        <th class="col-num">#</th>
                        <th class="col-title">Name <span class="sort-arrow">↑</span></th>
                        <th class="col-time"><ClockIcon size={14} /></th>
                        <th class="col-artist">Artist</th>
                        <th class="col-album">Album</th>
                    </tr>
                </thead>
                <tbody>
                    {#each tracks as track, i}
                        {@const trkId = track.id || (track as any).id}
                        <!-- svelte-ignore a11y_mouse_events_have_key_events -->
                        <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
                        <tr 
                            oncontextmenu={(e) => handleTrackRightClick(e, i, trkId)}
                            onmouseover={() => hoveredTrack = trkId}
                            onmouseleave={() => hoveredTrack = null}
                        >
                            <td class="col-num">
                                {#if hoveredTrack === trkId}
                                    <div class="play-btn">
                                        <PlayCircleIcon size={16} weight="fill" />
                                    </div>
                                {:else}
                                    {i + 1}
                                {/if}
                            </td>
                            <td class="col-title">
                                <span class="track-name">{track.title || (track as any).title}</span>
                            </td>
                            <td class="col-time">{formatTime(track.duration || (track as any).duration)}</td>
                            <td class="col-artist">Artist {track.artist_id || (track as any).artist_id}</td>
                            <td class="col-album">Album {track.album_id || (track as any).album_id}</td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</div>

<style>
    .library-view {
        color: var(--text-primary);
        height: 100%;
        display: flex;
        flex-direction: column;
        background: var(--bg-content);
        font-size: 13px; /* Much more compact typography */
        user-select: none; /* Native desktop feel */
    }

    .toolbar {
        padding: 10px 20px;
        background: var(--bg-base);
        border-bottom: 1px solid var(--border);
        display: flex;
        align-items: baseline;
        gap: 12px;
        flex-shrink: 0;
    }

    .toolbar h2 {
        margin: 0;
        font-size: 16px;
        font-weight: 500;
        color: var(--text-primary);
    }

    .track-count {
        color: var(--text-secondary);
        font-size: 12px;
    }

    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        flex: 1;
        color: var(--text-secondary);
        font-size: 13px;
    }

    .track-table-container {
        flex: 1;
        overflow-y: auto;
    }

    .track-table {
        width: 100%;
        border-collapse: collapse;
        table-layout: fixed;
    }

    thead th {
        font-size: 12px;
        font-weight: 400;
        text-align: left;
        color: var(--text-secondary);
        padding: 6px 12px;
        background: var(--bg-content);
        border-bottom: 1px solid var(--border);
        border-right: 1px solid var(--border);
        position: sticky;
        top: 0;
        z-index: 5;
        white-space: nowrap;
    }
    
    thead th:last-child {
        border-right: none;
    }

    .sort-arrow {
        font-size: 1.1em;
    }

    tbody td {
        padding: 4px 12px; /* Tighter rows */
        color: var(--text-primary);
        vertical-align: middle;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    tbody tr {
        background-color: transparent;
    }

    tbody tr:nth-child(even) {
        background-color: var(--bg-stripe); /* Alternating row colors */
    }

    tbody tr:hover {
        background-color: var(--bg-hover);
        color: var(--text-primary);
    }

    tbody tr:hover td {
        color: var(--text-primary);
    }

    .col-num {
        width: 40px;
        text-align: right;
        padding-right: 12px;
        color: var(--text-muted);
    }
    
    tbody tr:hover .col-num, tbody tr:hover .col-time {
        color: var(--text-secondary);
    }
    
    .col-time {
        width: 65px;
        text-align: right;
        color: var(--text-secondary);
    }
    
    .col-title {
        width: 35%;
        color: var(--text-primary);
    }
    
    .col-artist {
        width: 25%;
        color: var(--text-secondary);
    }
    
    .col-album {
        width: 25%;
        color: var(--text-secondary);
    }

    .track-name {
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .play-btn {
        display: flex;
        align-items: center;
        justify-content: flex-end;
        color: var(--text-primary);
    }

</style>
