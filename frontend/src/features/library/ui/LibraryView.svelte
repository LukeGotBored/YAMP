<script lang="ts">
    import { onMount } from "svelte";
    import { ClockIcon, PlayCircleIcon } from "phosphor-svelte";
    import * as LibraryAPI from "../../../../bindings/changeme/internal/application/services/libraryservice";
    import type { Track } from "../../../../bindings/changeme/internal/domain/models/models";
    import { ShowLibraryItemMenu } from "../../../../bindings/changeme/internal/app/context";
    import { invokeCached } from "../../../shared/api/client";

    let tracks = $state.raw<Track[]>([]);
    let isLoading = $state(true);

    // Accessibility & Navigation State
    let cursorIndex = $state<number>(0);
    let trackTableContainer = $state<HTMLDivElement | null>(null);
    let isContainerFocused = $state(false);

    // Typeahead Search State
    let typeaheadBuffer = '';
    let typeaheadTimeout: ReturnType<typeof setTimeout> | null = null;

    // Sorting State
    type SortColumn = 'none' | 'title' | 'duration' | 'artist' | 'album';
    let sortColumn = $state<SortColumn>('title');
    let sortAsc = $state(true);

    function sortTracksList(column: SortColumn, silent = false) {
        if (!silent) {
            if (sortColumn === column) {
                sortAsc = !sortAsc;
            } else {
                sortColumn = column;
                sortAsc = true;
            }
        }

        const currentId = tracks[cursorIndex]?.id || (tracks[cursorIndex] as any)?.id;

        const sorted = [...tracks].sort((a, b) => {
            let valA, valB;
            switch(sortColumn) {
                case 'title':
                    valA = (a.title || (a as any).title || '').toLowerCase();
                    valB = (b.title || (b as any).title || '').toLowerCase();
                    break;
                case 'artist':
                    valA = (a.artist_name || (a as any).artist_name || '').toLowerCase();
                    valB = (b.artist_name || (b as any).artist_name || '').toLowerCase();
                    break;
                case 'album':
                    valA = (a.album_title || (a as any).album_title || '').toLowerCase();
                    valB = (b.album_title || (b as any).album_title || '').toLowerCase();
                    break;
                case 'duration':
                    valA = a.duration || (a as any).duration || 0;
                    valB = b.duration || (b as any).duration || 0;
                    break;
                default:
                    valA = a.id || (a as any).id;
                    valB = b.id || (b as any).id;
                    break;
            }

            if (valA < valB) return sortAsc ? -1 : 1;
            if (valA > valB) return sortAsc ? 1 : -1;
            return 0;
        });

        tracks = sorted;

        // Restore cursor native anchorage
        if (currentId) {
            const newIndex = tracks.findIndex(t => (t.id || (t as any).id) === currentId);
            if (newIndex !== -1) {
                cursorIndex = newIndex;
                scrollToCursor();
            }
        }
    }

    // Virtualization State
    let containerHeight = $state(0);
    let scrollY = $state(0);
    let pendingScrollY = 0;
    let isScrollFrameQueued = false;
    const ROW_HEIGHT = 32;

    let visibleRowCount = $derived(Math.max(1, Math.ceil(containerHeight / ROW_HEIGHT)));
    let bufferRows = $derived(Math.max(6, Math.min(14, Math.ceil(visibleRowCount * 0.5))));
    let startIndex = $derived(Math.max(0, Math.floor(scrollY / ROW_HEIGHT) - bufferRows));
    let endIndex = $derived(Math.min(tracks.length, Math.ceil((scrollY + containerHeight) / ROW_HEIGHT) + bufferRows));
    let visibleTracks = $derived(tracks.slice(startIndex, endIndex));

    function handleContainerScroll() {
        if (!trackTableContainer) return;

        pendingScrollY = trackTableContainer.scrollTop;
        if (isScrollFrameQueued) return;

        isScrollFrameQueued = true;
        requestAnimationFrame(() => {
            scrollY = pendingScrollY;
            isScrollFrameQueued = false;
        });
    }

    function getTrackSearchText(track: Track) {
        const title = (track.title || (track as any).title || "").toLowerCase();
        const artist = (track.artist_name || (track as any).artist_name || "").toLowerCase();
        const album = (track.album_title || (track as any).album_title || "").toLowerCase();
        return { title, artist, album };
    }

    function matchesTypeahead(track: Track, query: string) {
        const t = getTrackSearchText(track);
        return t.title.startsWith(query) || t.artist.startsWith(query) || t.album.startsWith(query);
    }

    function findNextTypeaheadMatch(query: string, fromIndex: number) {
        if (tracks.length === 0) return -1;

        const start = Math.max(0, fromIndex);
        for (let i = start; i < tracks.length; i += 1) {
            if (matchesTypeahead(tracks[i], query)) return i;
        }

        for (let i = 0; i < start; i += 1) {
            if (matchesTypeahead(tracks[i], query)) return i;
        }

        return -1;
    }

    function scrollToCursor() {
        if (!trackTableContainer) return;
        if (cursorIndex === -1) return;

        const top = cursorIndex * ROW_HEIGHT;
        const bottom = top + ROW_HEIGHT + 36; // Approx header height

        // Apple-style padding: keep 2 rows of context above/below when navigating with keyboard
        const padding = ROW_HEIGHT * 2;

        let newScrollY = scrollY;
        if (top < scrollY + padding) {
            newScrollY = Math.max(0, top - padding);
        } else if (bottom > scrollY + containerHeight - padding) {
            newScrollY = bottom - containerHeight + padding;
        }

        if (newScrollY !== scrollY) {
            trackTableContainer.scrollTop = newScrollY;
            scrollY = newScrollY; // Synchronous update eliminates virtual-dom jump flicker
        }
    }

    function handleWindowKeydown(e: KeyboardEvent) {
        if (tracks.length === 0) return;
        
        // Ignore input events if the user is typing in a native input field anywhere on the page
        const target = e.target as HTMLElement;
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) return;

        const isNavKey = ["ArrowDown", "ArrowUp", "Home", "End", "PageDown", "PageUp", "Enter", "ContextMenu"].includes(e.key) || (e.shiftKey && e.key === 'F10');
        const isTypeahead = e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey && !e.shiftKey;
        
        // If not focused, we only intercept if it's a known list navigation or typahead command
        if (!isContainerFocused && !isNavKey && !isTypeahead) return;

        // Auto-focus container to steal active window context back
        if (!isContainerFocused && (isNavKey || isTypeahead)) {
            trackTableContainer?.focus();
            if (cursorIndex === -1 && tracks.length > 0) {
                cursorIndex = 0;
            }
        }

        const pageJump = Math.max(1, Math.floor(containerHeight / ROW_HEIGHT) - 1);

        switch (e.key) {
            case "ArrowDown":
                e.preventDefault();
                cursorIndex = cursorIndex === -1 ? 0 : Math.min(cursorIndex + 1, tracks.length - 1);
                scrollToCursor();
                break;
            case "ArrowUp":
                e.preventDefault();
                cursorIndex = cursorIndex === -1 ? tracks.length - 1 : Math.max(cursorIndex - 1, 0);
                scrollToCursor();
                break;
            case "Home":
                e.preventDefault();
                cursorIndex = 0;
                scrollToCursor();
                break;
            case "End":
                e.preventDefault();
                cursorIndex = tracks.length - 1;
                scrollToCursor();
                break;
            case "PageDown":
                e.preventDefault();
                cursorIndex = cursorIndex === -1 ? pageJump : Math.min(cursorIndex + pageJump, tracks.length - 1);
                scrollToCursor();
                break;
            case "PageUp":
                e.preventDefault();
                cursorIndex = cursorIndex === -1 ? 0 : Math.max(cursorIndex - pageJump, 0);
                scrollToCursor();
                break;
            case "Enter":
                e.preventDefault();
                if (cursorIndex !== -1) {
                    const tId = tracks[cursorIndex].id || (tracks[cursorIndex] as any).id;
                    console.log('TODO: Play track from Enter', tId);
                }
                break;
            case "ContextMenu":
                e.preventDefault();
                if (cursorIndex !== -1) {
                    const trk = tracks[cursorIndex];
                    const trkId = (trk.id || (trk as any).id) as number;
                    handleTrackRightClick(e, cursorIndex, trkId);
                }
                break;
            default:
                if (e.shiftKey && e.key === 'F10') {
                    e.preventDefault();
                    if (cursorIndex !== -1) {
                        const trkShift = tracks[cursorIndex];
                        const trkIdShift = (trkShift.id || (trkShift as any).id) as number;
                        handleTrackRightClick(e, cursorIndex, trkIdShift);
                    }
                } else if (e.key.length === 1 && !e.ctrlKey && !e.metaKey && !e.altKey && !e.shiftKey) {
                    e.preventDefault();
                    if (typeaheadTimeout) clearTimeout(typeaheadTimeout);
                    
                    const char = e.key.toLowerCase();
                    const isCyclingSingleChar = typeaheadBuffer.length === 1 && typeaheadBuffer === char;
                    const searchBuffer = isCyclingSingleChar ? char : `${typeaheadBuffer}${char}`;
                    typeaheadBuffer = searchBuffer;

                    const searchStart = isCyclingSingleChar
                        ? (cursorIndex === -1 ? 0 : cursorIndex + 1)
                        : (typeaheadBuffer.length === 1 ? (cursorIndex === -1 ? 0 : cursorIndex + 1) : Math.max(0, cursorIndex));

                    typeaheadTimeout = setTimeout(() => {
                        typeaheadBuffer = '';
                    }, 400);

                    const foundIndex = findNextTypeaheadMatch(searchBuffer, searchStart);
                    if (foundIndex !== -1) {
                        cursorIndex = foundIndex;
                        scrollToCursor();
                    }
                }
                break;
        }
    }

    function handleHeaderKeydown(e: KeyboardEvent, column: SortColumn) {
        if (e.key !== "Enter" && e.key !== " ") return;
        e.preventDefault();
        sortTracksList(column);
    }
    async function loadTracks() {
        isLoading = true;
        
        const result = await invokeCached("library:tracks:list", LibraryAPI.ListAllTracks, []);
        if (result.error) {
             console.error("Failed to load tracks", result.error.message);
        } else {
             tracks = result.data || [];
             sortTracksList(sortColumn, true);
        }
        
        isLoading = false;
    }

    function handleTrackRightClick(event: MouseEvent | KeyboardEvent, index: number, trackId: number) {
        if (event && typeof event.preventDefault === "function") {
            event.preventDefault();
        }
        
        let x = 0;
        let y = 0;

        if ('clientX' in event) {
            x = event.clientX;
            y = event.clientY;
        } else if (trackTableContainer) {
            // Calculate coordinates for the selected row precisely
            const rect = trackTableContainer.getBoundingClientRect();
            const rowOffsets = (index * ROW_HEIGHT) - scrollY;
            x = rect.left + rect.width / 2;
            
            if (rowOffsets >= 0 && rowOffsets <= rect.height) {
                y = rect.top + rowOffsets + (ROW_HEIGHT / 2);
            } else {
                y = rect.top + rect.height / 2;
            }
        }

        // Display the native OS Context Menu using Wails v3
        ShowLibraryItemMenu(x, y, trackId.toString());
    }

    onMount(() => {
        loadTracks();

        return () => {
            if (typeaheadTimeout) clearTimeout(typeaheadTimeout);
        };
    });

    // Helper to format seconds into mm:ss or h:mm:ss
    function formatTime(seconds: number) {
        if (!seconds) return "0:00";
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = Math.floor(seconds % 60);
        if (h > 0) {
            return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
        }
        return `${m}:${s.toString().padStart(2, "0")}`;
    }
</script>

<svelte:window onkeydown={handleWindowKeydown} />

<div class="library-view">
    <div class="toolbar">
        <h2>Songs</h2>
        <span class="track-count">{tracks.length.toLocaleString()} items</span>
    </div>

    {#if isLoading}
        <div class="empty-state">
            <p>Loading library...</p>
        </div>
    {:else if tracks.length === 0}
        <div class="empty-state">
            <p>Your library is empty.</p>
        </div>
    {:else}
        <div class="table-header-container" class:scrolled={scrollY > 0}>
            <table class="track-table">
                <colgroup>
                    <col class="col-num" />
                    <col class="col-title" />
                    <col class="col-time" />
                    <col class="col-artist" />
                    <col class="col-album" />
                </colgroup>
                <thead>
                    <tr>
                        <th class="col-num">#</th>
                        <th
                            class="col-title"
                            tabindex="0"
                            aria-sort={sortColumn === 'title' ? (sortAsc ? 'ascending' : 'descending') : 'none'}
                            onclick={() => sortTracksList('title')}
                            onkeydown={(e) => handleHeaderKeydown(e, 'title')}
                        >
                            Name
                            {#if sortColumn === 'title'}
                                <span class="sort-arrow">{sortAsc ? '↑' : '↓'}</span>
                            {/if}
                        </th>
                        <th
                            class="col-time"
                            tabindex="0"
                            aria-sort={sortColumn === 'duration' ? (sortAsc ? 'ascending' : 'descending') : 'none'}
                            onclick={() => sortTracksList('duration')}
                            onkeydown={(e) => handleHeaderKeydown(e, 'duration')}
                        >
                            <ClockIcon size={14} />
                            {#if sortColumn === 'duration'}
                                <span class="sort-arrow">{sortAsc ? '↑' : '↓'}</span>
                            {/if}
                        </th>
                        <th
                            class="col-artist"
                            tabindex="0"
                            aria-sort={sortColumn === 'artist' ? (sortAsc ? 'ascending' : 'descending') : 'none'}
                            onclick={() => sortTracksList('artist')}
                            onkeydown={(e) => handleHeaderKeydown(e, 'artist')}
                        >
                            Artist
                            {#if sortColumn === 'artist'}
                                <span class="sort-arrow">{sortAsc ? '↑' : '↓'}</span>
                            {/if}
                        </th>
                        <th
                            class="col-album"
                            tabindex="0"
                            aria-sort={sortColumn === 'album' ? (sortAsc ? 'ascending' : 'descending') : 'none'}
                            onclick={() => sortTracksList('album')}
                            onkeydown={(e) => handleHeaderKeydown(e, 'album')}
                        >
                            Album
                            {#if sortColumn === 'album'}
                                <span class="sort-arrow">{sortAsc ? '↑' : '↓'}</span>
                            {/if}
                        </th>
                    </tr>
                </thead>
            </table>
        </div>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <div 
            class="track-table-container"
            bind:this={trackTableContainer}
            bind:clientHeight={containerHeight}
            onscroll={handleContainerScroll}
            tabindex="-1"
            role="listbox"
            onfocus={() => isContainerFocused = true}
            onblur={() => isContainerFocused = false}
            onclick={(e) => {
                // If the user clicks the empty space *below* the tracks
                if (e.target === trackTableContainer) {
                    cursorIndex = -1; // Unselect row (macOS/Windows style)
                    trackTableContainer?.focus();
                }
            }}
        >
            <table class="track-table">
                <colgroup>
                    <col class="col-num" />
                    <col class="col-title" />
                    <col class="col-time" />
                    <col class="col-artist" />
                    <col class="col-album" />
                </colgroup>
                <!-- svelte-ignore a11y_click_events_have_key_events -->
                <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
                <tbody
                    onclick={(e) => {
                        const tr = (e.target as HTMLElement).closest('tr[data-index]');
                        if (tr) {
                            cursorIndex = parseInt(tr.getAttribute('data-index') || '0', 10);
                            trackTableContainer?.focus();
                            
                            // If they specifically clicked the play button icon wrapper
                            const btn = (e.target as HTMLElement).closest('.play-btn-wrapper');
                            if (btn) {
                                const id = parseInt(tr.getAttribute('data-track-id') || '0', 10);
                                console.log('TODO: Play track immediately from icon', id);
                            }
                        }
                    }}
                    oncontextmenu={(e) => {
                        const tr = (e.target as HTMLElement).closest('tr[data-index]');
                        if (tr) {
                            const idx = parseInt(tr.getAttribute('data-index') || '0', 10);
                            const id = parseInt(tr.getAttribute('data-track-id') || '0', 10);
                            cursorIndex = idx;
                            trackTableContainer?.focus();
                            handleTrackRightClick(e, idx, id);
                        }
                    }}
                    ondblclick={(e) => {
                        const tr = (e.target as HTMLElement).closest('tr[data-index]');
                        if (tr) {
                            const id = parseInt(tr.getAttribute('data-track-id') || '0', 10);
                            console.log('TODO: Play track', id);
                        }
                    }}
                >
                    {#if startIndex > 0}
                        <tr class="virtual-spacer" style="height: {startIndex * ROW_HEIGHT}px;">
                            <td colspan="5"></td>
                        </tr>
                    {/if}
                    {#each visibleTracks as track, localIdx (track.id || (track as any).id)}
                        {@const i = startIndex + localIdx}
                        {@const trkId = track.id || (track as any).id}
                        {@const title = track.title || (track as any).title || "Unknown Track"}
                        {@const artist = track.artist_name || (track as any).artist_name || "Unknown Artist"}
                        {@const album = track.album_title || (track as any).album_title || "Unknown Album"}
                        <tr 
                            role="option"
                            aria-selected={i === cursorIndex}
                            class:stripe={i % 2 === 0}
                            class:cursor={i === cursorIndex}
                            data-index={i}
                            data-track-id={trkId}
                        >
                            <td class="col-num">
                                <button class="play-btn-wrapper" type="button" tabindex="-1" aria-label={`Play ${title}`}>
                                    <PlayCircleIcon size={16} weight="fill" />
                                </button>
                                <span class="track-number">{i + 1}</span>
                            </td>
                            <td class="col-title">
                                <span class="track-name" class:muted-text={title === "Unknown Track"}>{title}</span>
                            </td>
                            <td class="col-time">{formatTime(track.duration || (track as any).duration)}</td>
                            <td class="col-artist" class:muted-text={artist === "Unknown Artist"}>{artist}</td>
                            <td class="col-album" class:muted-text={album === "Unknown Album"}>{album}</td>
                        </tr>
                    {/each}
                    {#if endIndex < tracks.length}
                        <tr class="virtual-spacer" style="height: {(tracks.length - endIndex) * ROW_HEIGHT}px;">
                            <td colspan="5"></td>
                        </tr>
                    {/if}
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

    .table-header-container {
        /* Syncs physical header width precisely with the scrollbar on Windows */
        overflow-y: scroll;
        flex-shrink: 0;
        background: var(--bg-content);
        transition: box-shadow 0.2s ease;
        position: relative;
        z-index: 10;
    }
    
    .table-header-container.scrolled {
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    }
    
    .table-header-container::-webkit-scrollbar-thumb {
        background: transparent;
    }
    
    .table-header-container::-webkit-scrollbar-track {
        background: transparent;
    }
    
    .table-header-container::-webkit-scrollbar {
        background: transparent; /* Maintains space but prevents drawing */
    }

    .track-table-container {
        flex: 1;
        overflow-y: scroll; /* Force consistent body scrollbar styling */
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
        /* 1px border below the unstickied header */
        border-bottom: 1px solid var(--border);
        border-right: 1px solid var(--border);
        white-space: nowrap;
        cursor: default;
        user-select: none;
    }

    thead th:hover:not(.col-num) {
        background: var(--bg-hover, rgba(255, 255, 255, 0.05));
    }
    
    thead th:active:not(.col-num) {
        background: rgba(255, 255, 255, 0.08); /* Mimic pressed header state */
    }
    
    thead th:last-child {
        border-right: none;
    }

    .sort-arrow {
        font-size: 1.1em;
        display: inline-block;
        margin-left: 4px;
        color: var(--text-primary);
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
        height: 32px;
        box-sizing: border-box;
    }

    .muted-text {
        opacity: 0.4;
        font-style: italic;
    }

    tbody tr.stripe {
        background-color: var(--bg-stripe); /* Alternating row colors */
    }

    /* Selection state */
    /* Unfocused application / list state: Subtle gray highlighting */
    tbody tr.cursor {
        background-color: var(--bg-hover, rgba(255, 255, 255, 0.08));
    }
    
    /* Native OS focused application state: Bright primary highlight */
    .track-table-container:focus tbody tr.cursor {
        background-color: var(--accent); /* Highlight focused track */
        color: white;
    }
    
    .track-table-container:focus tbody tr.cursor td,
    .track-table-container:focus tbody tr.cursor .col-num,
    .track-table-container:focus tbody tr.cursor .col-time,
    .track-table-container:focus tbody tr.cursor .col-artist,
    .track-table-container:focus tbody tr.cursor .col-album {
        color: white;
    }

    .virtual-spacer, .virtual-spacer td {
        padding: 0 !important;
        margin: 0 !important;
        border: none !important;
        background: transparent !important;
    }

    tbody tr:hover:not(.virtual-spacer):not(.cursor) {
        background-color: var(--bg-hover);
        color: var(--text-primary);
    }

    tbody tr:hover:not(.virtual-spacer) td {
        color: var(--text-primary);
    }

    .col-num {
        width: 40px;
        text-align: right;
        padding-right: 12px;
        color: var(--text-muted);
        font-variant-numeric: tabular-nums;
    }

    .play-btn-wrapper {
        display: none;
    }

    .track-number {
        display: inline-block;
    }

    tbody tr:hover:not(.virtual-spacer) .play-btn-wrapper,
    .track-table-container:focus tbody tr.cursor .play-btn-wrapper {
        display: flex;
        justify-content: flex-end;
    }
    
    tbody tr:hover:not(.virtual-spacer) .track-number,
    .track-table-container:focus tbody tr.cursor .track-number {
        display: none;
    }
    
    tbody tr:hover:not(.virtual-spacer):not(.cursor) .col-num, 
    tbody tr:hover:not(.virtual-spacer):not(.cursor) .col-time {
        color: var(--text-secondary);
    }
    
    .col-time {
        width: 65px;
        text-align: right;
        color: var(--text-secondary);
        font-variant-numeric: tabular-nums;
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

</style>
