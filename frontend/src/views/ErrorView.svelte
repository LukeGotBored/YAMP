<script lang="ts">
    import { BugIcon } from "phosphor-svelte";
</script>

<main>
    <div id="icon">
        <BugIcon />
    </div>

    <div id="main">Whoops!</div>
    <div id="description">Something went wrong, somehow.</div>
</main>

<style>
    main {
        height: calc(100% - var(--tb-h));
        width: 100vw;
        display: flex;
        justify-content: center;
        align-items: center;
        flex-direction: column;
    }

    #icon {
        font-size: 5em;
        opacity: 1;
        display: flex;
        justify-content: center;
        align-items: center;
    }

    #main {
        font-weight: 600;
        font-size: 3em;
    }

    #description {
        font-size: 1.25em;
        opacity: 0.5;
    }
</style>
