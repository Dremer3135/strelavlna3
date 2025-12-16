<script lang="ts">
    import { pocketbase } from "$lib/pocketbase";

    let paperid = $state("");
    let paperfilter = $state("");
    async function genPapers() {
      await pocketbase.send("/api/papers", {query: {id: paperid, filter: paperfilter}})
    }

    let tfilter = $state("");
    let pfilter = $state("");
    async function setGameData() {
      if (confirm("resetovat data týmů")) return;
      await pocketbase.send("/api/setgamedata", {query: {teamf: tfilter, probf: pfilter}})
    }
    async function pushFreeProbs() {
      await pocketbase.send("/api/pushfreeprobs", {query: {teamf: tfilter, probf: pfilter}})
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
    <input type="text" bind:value={tfilter} placeholder="teams">
    <input type="text" bind:value={pfilter} placeholder="probs">
    <button onclick={setGameData}>SetGameData</button>
    <button onclick={pushFreeProbs}>PushFreeProbs</button>
    <br>
    <input type="text" bind:value={sql}>
    <button onclick={sendSql}>SendSql</button>
    <pre>
      select name, json_extract(inPersonData, '$.Money') as money from teams where json_array_length(inPersonData, '$.Bought') > 0 or json_array_length(inPersonData, '$.Sold') > 0 or json_array_length(inPersonData, '$.Solved') > 0 and contest = 'p2wd6fb1lcuyqbl' order by money desc
    </pre>
    <br>
    <input type="text" bind:value={rdbid}>
    <button onclick={genRdb}>GenRdb</button>
    <br>
    <button onclick={start}>Start</button>
    <br>
    <button onclick={end}>End</button>
</main>
