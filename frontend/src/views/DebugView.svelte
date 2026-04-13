<script lang="ts">
    import { invoke } from "../shared/api/client";
    import * as LibraryService from "../../bindings/changeme/internal/application/services/libraryservice";
    import * as DBService from "../../bindings/changeme/internal/db/service";

    let logs = $state<string[]>([]);
    let isWorking = $state(false);

    function addLog(msg: string) {
        logs = [...logs, `[${new Date().toLocaleTimeString()}] ${msg}`];
    }

    async function populateBigDummyData() {
        isWorking = true;
        addLog("Generating 50,000 dummy tracks for benchmarking...");
        const t0 = performance.now();
        const res = await invoke(LibraryService.GenerateDummyData, 50000);
        const t1 = performance.now();

        if (!res.error) {
            addLog(`Generated 50k tracks in ${(t1 - t0).toFixed(2)}ms.`);
        } else {
            addLog(`Error generating tracks: ${res.error.message}`);
        }
        isWorking = false;
    }

    async function populateSmallDummyData() {
        isWorking = true;
        addLog("Generating 1,000 dummy tracks for benchmarking...");
        const t0 = performance.now();
        const res = await invoke(LibraryService.GenerateDummyData, 1000);
        const t1 = performance.now();

        if (!res.error) {
            addLog(`Generated 1k tracks in ${(t1 - t0).toFixed(2)}ms.`);
        } else {
            addLog(`Error generating tracks: ${res.error.message}`);
        }
        isWorking = false;
    }

    async function testListTracks() {
        isWorking = true;
        addLog("Benchmarking ListAllTracks() query scale...");
        const t0 = performance.now();
        const res = await invoke(LibraryService.ListAllTracks);
        const t1 = performance.now();

        if (!res.error) {
            const tracks = res.data || [];
            addLog(
                `Fetched ${tracks.length} tracks to JS in ${(t1 - t0).toFixed(2)}ms.`,
            );
        } else {
            addLog(`Error fetching tracks: ${res.error.message}`);
        }
        isWorking = false;
    }

    function clearLogs() {
        logs = [];
    }

    async function resetDatabase() {
        isWorking = true;
        addLog("Purging all database records and sequences...");
        const res = await invoke(DBService.ResetDatabase);

        if (!res.error) {
            addLog("Database successfully reset.");
        } else {
            addLog(`Reset failed: ${res.error.message}`);
        }
        isWorking = false;
    }
</script>

<div class="debug-view">
    <div class="header">
        <h2>Developer Tools</h2>
        <div class="actions">
            <button class="bg-btn secondary" onclick={clearLogs}
                >Clear Logs</button
            >
        </div>
    </div>
    <div class="sandbox-info">
        <p>
            This view is intended for development and testing purposes. Use
            these features only if you know what you're doing!
        </p>
        <p><b>You might accidentally delete your data!</b></p>
    </div>

    <div class="dev-layout">
        <div class="dev-section db-section">
            <h3>Database Throughput</h3>
            <div class="actions" style="margin-bottom: 10px;">
                <button
                    class="bg-btn secondary"
                    onclick={populateSmallDummyData}
                    disabled={isWorking}>Generate 1k Tracks</button
                >
                <button
                    class="bg-btn secondary"
                    onclick={populateBigDummyData}
                    disabled={isWorking}>Generate 50k Tracks</button
                >
                <button
                    class="bg-btn secondary"
                    onclick={testListTracks}
                    disabled={isWorking}>Clock ListAllTracks()</button
                >
                <button
                    class="bg-btn custom-danger"
                    onclick={resetDatabase}
                    disabled={isWorking}>Hard Reset System</button
                >
            </div>

            <div class="log-console">
                <div id="log-header">Benchmark Output</div>
                <div class="log-body">
                    {#each logs as log}
                        <div class="log-entry">{log}</div>
                    {/each}
                    {#if logs.length === 0}
                        <div class="log-empty">
                            Waiting for a benchmark to run...
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .debug-view {
        padding: 20px;
        color: var(--text-primary);
        height: 100%;
        display: flex;
        flex-direction: column;
        overflow-y: auto;
    }
    .header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 20px;
    }
    .actions {
        display: flex;
        gap: 10px;
        flex-wrap: wrap;
    }
    h2 {
        margin: 0;
        font-weight: 500;
    }
    h3 {
        margin: 0 0 15px 0;
        font-weight: 500;
        font-size: 1.1em;
        color: var(--text-secondary);
    }
    .sandbox-info {
        background: rgba(142, 56, 56, 0.1);
        border: 1px solid rgba(142, 56, 56, 0.3);
        padding: 10px 15px;
        margin-bottom: 20px;
        border-radius: 4px;
    }
    .sandbox-info p {
        margin: 0;
        font-size: 0.9em;
        color: rgba(255, 255, 255, 0.8);
    }
    .dev-layout {
        display: flex;
        flex-direction: column;
        gap: 20px;
        flex: 1;
    }
    .dev-section {
        background: rgba(255, 255, 255, 0.02);
        border: 1px solid rgba(255, 255, 255, 0.05);
        padding: 15px;
        border-radius: 6px;
    }
    .db-section {
        display: flex;
        flex-direction: column;
    }
    .bg-btn {
        background: var(--accent);
        color: white;
        border: none;
        padding: 6px 12px;
        border-radius: 4px;
        cursor: pointer;
    }
    .bg-btn:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }
    .bg-btn:hover:not(:disabled) {
        background: var(--accent-hover);
    }
    .bg-btn.secondary {
        background: rgba(255, 255, 255, 0.1);
    }
    .bg-btn.secondary:hover {
        background: rgba(255, 255, 255, 0.2);
    }
    .custom-danger {
        background: #e53935;
    }
    .custom-danger:hover:not(:disabled) {
        background: #d32f2f;
    }
    .log-console {
        flex: 1;
        background: #111;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 4px;
        display: flex;
        flex-direction: column;
        font-family: monospace;
        font-size: 0.9em;
        min-height: 200px;
    }
    #log-header {
        font-weight: 500;
        padding: 8px 10px;
        background: rgba(255, 255, 255, 0.02);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        color: var(--text-secondary);
    }
    .log-body {
        padding: 10px;
        overflow-y: auto;
        flex: 1;
    }
    .log-entry {
        margin-bottom: 5px;
        color: #a5d6ff;
    }
    .log-empty {
        color: rgba(255, 255, 255, 0.3);
        font-style: italic;
        text-align: center;
        margin-top: 10px;
    }
</style>
