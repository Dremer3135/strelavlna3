<script lang="ts">
    import { pocketbase } from "$lib/pocketbase";

    let paperid = $state("");
    let paperfilter = $state("");
    async function genPapers() {
      await pocketbase.send("/api/papers", {query: {id: paperid, filter: paperfilter}})
    }

    let sql = $state("");
    async function sendSql() {
      await pocketbase.send("/api/sql", {body: sql, method: "POST"})
    }

    let rdbid = $state("");
    async function genRdb() {
      await pocketbase.send("/api/rdb", {query: {id: rdbid} })
    }

    async function start() {
      await pocketbase.send("/api/start")
    }

    async function end() {
      await pocketbase.send("/api/end")
    }
</script>

<main>
    <h1>
        You are admin!!!!!!
    </h1>
    <input type="text" bind:value={paperid}>
    <input type="text" bind:value={paperfilter}>
    <button onclick={genPapers}>GenPapers</button>
    <br>
    <input type="text" bind:value={sql}>
    <button onclick={sendSql}>SendSql</button>
    <br>
    <input type="text" bind:value={rdbid}>
    <button onclick={genRdb}>GenRdb</button>
    <br>
    <button onclick={start}>Start</button>
    <br>
    <button onclick={end}>End</button>
</main>
