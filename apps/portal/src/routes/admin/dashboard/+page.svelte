<script lang="ts">
    import { pocketbase } from "$lib/pocketbase";

    let paperid = $state("");
    async function genPapers() {
      await pocketbase.send("/api/papers", {query: {id: paperid}})
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
    <button onclick={End}>End</button>
</main>
